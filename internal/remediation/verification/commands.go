package verification

import (
	"errors"
	"path/filepath"
	"sort"
	"strings"
)

type CommandKind string

const (
	Test      CommandKind = "test"
	Build     CommandKind = "build"
	TypeCheck CommandKind = "typecheck"
)

type Command struct {
	ID               string      `json:"id"`
	Kind             CommandKind `json:"kind"`
	Path             string      `json:"path"`
	Args             []string    `json:"args"`
	WorkingDirectory string      `json:"working_directory"`
	Reason           string      `json:"reason,omitempty"`
	Enabled          bool        `json:"enabled"`
}

func SelectCommands(candidates []Command) ([]Command, error) {
	seen := make(map[string]struct{}, len(candidates))
	selected := make([]Command, 0, len(candidates))
	for _, command := range candidates {
		if command.ID == "" || command.Path == "" || !filepath.IsAbs(command.Path) || !command.Enabled {
			continue
		}
		if command.Kind != Test && command.Kind != Build && command.Kind != TypeCheck {
			continue
		}
		working := filepath.Clean(command.WorkingDirectory)
		if command.WorkingDirectory == "" || filepath.IsAbs(command.WorkingDirectory) || working == ".." || strings.HasPrefix(working, ".."+string(filepath.Separator)) {
			return nil, errors.New("verification working directory must be repository-relative")
		}
		if isShell(command.Path, command.Args) {
			return nil, errors.New("shell verification commands are not allowed")
		}
		if _, exists := seen[command.ID]; exists {
			return nil, errors.New("duplicate verification command ID")
		}
		seen[command.ID] = struct{}{}
		command.Args = append([]string(nil), command.Args...)
		selected = append(selected, command)
	}
	sort.Slice(selected, func(i, j int) bool {
		if selected[i].Kind != selected[j].Kind {
			return selected[i].Kind < selected[j].Kind
		}
		return selected[i].ID < selected[j].ID
	})
	return selected, nil
}

func isShell(path string, args []string) bool {
	switch strings.ToLower(filepath.Base(path)) {
	case "sh", "bash", "zsh", "fish", "cmd", "cmd.exe", "powershell", "powershell.exe", "pwsh", "pwsh.exe":
		return true
	}
	for _, arg := range args {
		if arg == "-c" || arg == "/c" || arg == "-Command" {
			return true
		}
	}
	return false
}
