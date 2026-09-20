package presenter

import (
	"io"
	"os"
)

type OutputMode string

const (
	OutputTerminal OutputMode = "terminal"
	OutputJSON     OutputMode = "json"
)

// CapabilityOptions supplies stream and environment facts to capability selection.
type CapabilityOptions struct {
	Mode       OutputMode
	Stdout     io.Writer
	Stderr     io.Writer
	CI         bool
	Color      bool
	Width      int
	CanCancel  bool
	CanRestore bool
	IsTerminal func(io.Writer) bool
}

// Capabilities is the renderer's explicit, deterministic terminal contract.
type Capabilities struct {
	Mode              OutputMode
	Interactive       bool
	Color             bool
	Width             int
	Animation         bool
	Cancellation      bool
	RestoreOnExit     bool
	MachineOutputOnly bool
}

// SelectCapabilities chooses safe output behavior from approved stream facts.
func SelectCapabilities(options CapabilityOptions) Capabilities {
	mode := options.Mode
	if mode == "" {
		mode = OutputTerminal
	}
	isTerminal := options.IsTerminal
	if isTerminal == nil {
		isTerminal = IsTerminal
	}
	interactive := mode == OutputTerminal && !options.CI && isTerminal(options.Stderr)
	width := options.Width
	if width <= 0 {
		width = 80
	}
	machine := mode == OutputJSON
	return Capabilities{
		Mode:              mode,
		Interactive:       interactive,
		Color:             !machine && options.Color && interactive,
		Width:             width,
		Animation:         interactive,
		Cancellation:      options.CanCancel,
		RestoreOnExit:     options.CanRestore && interactive,
		MachineOutputOnly: machine,
	}
}

// IsTerminal reports whether a writer is backed by a character device.
func IsTerminal(writer io.Writer) bool {
	file, ok := writer.(*os.File)
	if !ok {
		return false
	}
	info, err := file.Stat()
	return err == nil && info.Mode()&os.ModeCharDevice != 0
}
