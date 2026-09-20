package presenter

import (
	"strings"
	"testing"

	"github.com/geoffrey-xiao/deprail/internal/app"
)

func TestProgressSinkRendersStructuredEventsToOneStream(t *testing.T) {
	var output strings.Builder
	sink := ProgressSink{Writer: &output}
	if err := sink.Emit(app.Event{Type: app.EventWorkspaceScanStarted, WorkspacePath: "apps/web"}); err != nil {
		t.Fatal(err)
	}
	if got := output.String(); got != "Progress: WorkspaceScanStarted apps/web\n" {
		t.Fatalf("progress = %q", got)
	}
}
