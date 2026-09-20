package presenter

import (
	"fmt"
	"io"

	"github.com/geoffrey-xiao/deprail/internal/app"
)

// ProgressSink renders concise lifecycle events to an already-approved stream.
type ProgressSink struct {
	Writer io.Writer
}

func (s ProgressSink) Emit(event app.Event) error {
	label := string(event.Type)
	if event.WorkspacePath != "" {
		label += " " + event.WorkspacePath
	}
	_, err := fmt.Fprintf(s.Writer, "Progress: %s\n", SanitizeLabel(label))
	return err
}
