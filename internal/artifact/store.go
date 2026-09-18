package artifact

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
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
