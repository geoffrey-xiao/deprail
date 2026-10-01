package app

import (
	"context"
	"errors"
	"path/filepath"
	"testing"
)

type cancellationCaptureSink struct {
	started chan struct{}
	calls   int
}

func (s *cancellationCaptureSink) Capture(ctx context.Context, _ CaptureInput) (SaveResult, error) {
	s.calls++
	close(s.started)
	<-ctx.Done()
	return SaveResult{}, ctx.Err()
}

func TestCaptureCanBeCancelledBySecondTermination(t *testing.T) {
	captureCancel, cancelCapture := context.WithCancel(context.Background())
	defer cancelCapture()
	sink := &cancellationCaptureSink{started: make(chan struct{})}
	done := make(chan error, 1)
	go func() {
		_, err := Scan(context.Background(), filepath.Join(t.TempDir(), "missing"), ScanOptions{
			Capture: sink, CaptureCancel: captureCancel,
		})
		done <- err
	}()
	<-sink.started
	cancelCapture()
	err := <-done
	var captureErr *ScanCaptureError
	if !errors.As(err, &captureErr) || captureErr.PersistenceErr == nil {
		t.Fatalf("capture cancellation error = %v", err)
	}
	if sink.calls != 1 || !errors.Is(captureErr.PersistenceErr, context.Canceled) {
		t.Fatalf("capture call/cancellation = %d/%v", sink.calls, captureErr.PersistenceErr)
	}
}
