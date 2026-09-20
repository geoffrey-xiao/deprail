package presenter

import "io"

type MessageColor string

const (
	ColorNone   MessageColor = ""
	ColorCyan   MessageColor = "cyan"
	ColorGreen  MessageColor = "green"
	ColorYellow MessageColor = "yellow"
	ColorRed    MessageColor = "red"
	ColorBlue   MessageColor = "blue"
)

const (
	ansiReset  = "\x1b[0m"
	ansiCyan   = "\x1b[36m"
	ansiGreen  = "\x1b[32m"
	ansiYellow = "\x1b[33m"
	ansiRed    = "\x1b[31m"
	ansiBlue   = "\x1b[34m"
)

func WriteColoredHeader(w io.Writer, message string, color MessageColor, enabled bool) error {
	if !enabled || color == ColorNone {
		return WriteHeader(w, message)
	}
	_, err := io.WriteString(w, colorPrefix(color)+SanitizeLabel(message)+ansiReset+"\n")
	return err
}

func colorPrefix(color MessageColor) string {
	switch color {
	case ColorCyan:
		return ansiCyan
	case ColorGreen:
		return ansiGreen
	case ColorYellow:
		return ansiYellow
	case ColorRed:
		return ansiRed
	case ColorBlue:
		return ansiBlue
	default:
		return ""
	}
}
