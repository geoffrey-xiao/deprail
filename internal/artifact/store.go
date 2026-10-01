package artifact

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type ErrorCode string

const (
	ErrInvalidInput ErrorCode = "CONFIG_INVALID"
	ErrWriteFailed  ErrorCode = "ARTIFACT_WRITE_FAILED"
	ErrNotFound     ErrorCode = "ARTIFACT_NOT_FOUND"
	ErrIntegrity    ErrorCode = "ARTIFACT_INTEGRITY_FAILED"
	ErrOutputLimit  ErrorCode = "SCANNER_OUTPUT_LIMIT"
)

type Error struct {
	Code    ErrorCode
	Message string
}

func (e *Error) Error() string { return fmt.Sprintf("%s: %s", e.Code, e.Message) }
func IsCode(err error, code ErrorCode) bool {
	var target *Error
	return errors.As(err, &target) && target.Code == code
}

type Store struct {
	Root     string
	MaxBytes int64
}
type Artifact struct {
	Digest string
	Size   int64
	Data   []byte
}

func (s Store) Put(data []byte) (Artifact, error) {
	if s.Root == "" || s.MaxBytes <= 0 {
		return Artifact{}, &Error{Code: ErrInvalidInput, Message: "artifact root and positive size limit are required"}
	}
	if int64(len(data)) > s.MaxBytes {
		return Artifact{}, &Error{Code: ErrOutputLimit, Message: "artifact exceeds configured size limit"}
	}
	digestBytes := sha256.Sum256(data)
	digest := hex.EncodeToString(digestBytes[:])
	dir := filepath.Join(s.Root, digest[:2])
	path := filepath.Join(dir, digest)
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return Artifact{}, &Error{Code: ErrWriteFailed, Message: "unable to create artifact directory"}
	}
	if existing, err := os.ReadFile(path); err == nil {
		if !equalDigest(existing, digest) {
			return Artifact{}, &Error{Code: ErrIntegrity, Message: "existing artifact digest mismatch"}
		}
		return Artifact{Digest: digest, Size: int64(len(existing)), Data: append([]byte(nil), existing...)}, nil
	} else if !errors.Is(err, os.ErrNotExist) {
		return Artifact{}, &Error{Code: ErrWriteFailed, Message: "unable to inspect artifact"}
	}
	tmp, err := os.CreateTemp(dir, ".artifact-*")
	if err != nil {
		return Artifact{}, &Error{Code: ErrWriteFailed, Message: "unable to create temporary artifact"}
	}
	tmpName := tmp.Name()
	defer os.Remove(tmpName)
	if err = tmp.Chmod(0o600); err == nil {
		_, err = tmp.Write(data)
	}
	if closeErr := tmp.Close(); err == nil {
		err = closeErr
	}
	if err != nil {
		return Artifact{}, &Error{Code: ErrWriteFailed, Message: "unable to write artifact"}
	}
	if err := os.Rename(tmpName, path); err != nil && !errors.Is(err, os.ErrExist) {
		return Artifact{}, &Error{Code: ErrWriteFailed, Message: "unable to publish artifact"}
	}
	return Artifact{Digest: digest, Size: int64(len(data)), Data: append([]byte(nil), data...)}, nil
}

func (s Store) Get(digest string) ([]byte, error) {
	if len(digest) != sha256.Size*2 || !isHex(digest) {
		return nil, &Error{Code: ErrInvalidInput, Message: "artifact digest is invalid"}
	}
	data, err := os.ReadFile(filepath.Join(s.Root, digest[:2], digest))
	if errors.Is(err, os.ErrNotExist) {
		return nil, &Error{Code: ErrNotFound, Message: "artifact is not available"}
	}
	if err != nil {
		return nil, &Error{Code: ErrWriteFailed, Message: "unable to read artifact"}
	}
	if !equalDigest(data, digest) {
		return nil, &Error{Code: ErrIntegrity, Message: "artifact digest verification failed"}
	}
	return data, nil
}

// Verify streams a content-addressed artifact and never returns its bytes.
func (s Store) Verify(ctx context.Context, digest string) error {
	if s.Root == "" || s.MaxBytes <= 0 || len(digest) != sha256.Size*2 || !isHex(digest) {
		return &Error{Code: ErrInvalidInput, Message: "artifact verification input is invalid"}
	}
	if err := ctx.Err(); err != nil {
		return &Error{Code: ErrWriteFailed, Message: "artifact verification was interrupted"}
	}
	rootInfo, err := os.Lstat(s.Root)
	if errors.Is(err, os.ErrNotExist) {
		return &Error{Code: ErrNotFound, Message: "artifact is not available"}
	}
	if err != nil {
		return &Error{Code: ErrWriteFailed, Message: "unable to inspect artifact store"}
	}
	if !rootInfo.IsDir() || rootInfo.Mode()&os.ModeSymlink != 0 {
		return &Error{Code: ErrIntegrity, Message: "artifact store boundary is invalid"}
	}
	directory := filepath.Join(s.Root, digest[:2])
	directoryInfo, err := os.Lstat(directory)
	if errors.Is(err, os.ErrNotExist) {
		return &Error{Code: ErrNotFound, Message: "artifact is not available"}
	}
	if err != nil {
		return &Error{Code: ErrWriteFailed, Message: "unable to inspect artifact"}
	}
	if !directoryInfo.IsDir() || directoryInfo.Mode()&os.ModeSymlink != 0 {
		return &Error{Code: ErrIntegrity, Message: "artifact path is invalid"}
	}
	path := filepath.Join(directory, digest)
	pathInfo, err := os.Lstat(path)
	if errors.Is(err, os.ErrNotExist) {
		return &Error{Code: ErrNotFound, Message: "artifact is not available"}
	}
	if err != nil {
		return &Error{Code: ErrWriteFailed, Message: "unable to inspect artifact"}
	}
	if !pathInfo.Mode().IsRegular() || pathInfo.Size() < 0 || pathInfo.Size() > s.MaxBytes {
		return &Error{Code: ErrIntegrity, Message: "artifact file is invalid"}
	}
	file, err := os.Open(path)
	if errors.Is(err, os.ErrNotExist) {
		return &Error{Code: ErrNotFound, Message: "artifact is not available"}
	}
	if err != nil {
		return &Error{Code: ErrWriteFailed, Message: "unable to read artifact"}
	}
	defer file.Close()
	openedInfo, err := file.Stat()
	if err != nil || !openedInfo.Mode().IsRegular() || !os.SameFile(pathInfo, openedInfo) {
		return &Error{Code: ErrIntegrity, Message: "artifact file changed during verification"}
	}
	hasher := sha256.New()
	var total int64
	buffer := make([]byte, 32<<10)
	for {
		if err := ctx.Err(); err != nil {
			return &Error{Code: ErrWriteFailed, Message: "artifact verification was interrupted"}
		}
		n, readErr := file.Read(buffer)
		total += int64(n)
		if total > int64(s.MaxBytes) {
			return &Error{Code: ErrIntegrity, Message: "artifact file exceeds the configured limit"}
		}
		_, _ = hasher.Write(buffer[:n])
		if errors.Is(readErr, io.EOF) {
			break
		}
		if readErr != nil {
			return &Error{Code: ErrWriteFailed, Message: "unable to verify artifact"}
		}
	}
	if hex.EncodeToString(hasher.Sum(nil)) != digest {
		return &Error{Code: ErrIntegrity, Message: "artifact digest verification failed"}
	}
	return nil
}
func equalDigest(data []byte, expected string) bool {
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:]) == expected
}
func isHex(value string) bool {
	for _, r := range value {
		if !((r >= '0' && r <= '9') || (r >= 'a' && r <= 'f')) {
			return false
		}
	}
	return true
}
