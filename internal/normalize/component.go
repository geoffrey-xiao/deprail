package normalize

import (
	"fmt"
	"net/url"
	"strings"
)

type ErrorCode string

const ErrIncompleteIdentity ErrorCode = "NORMALIZATION_FAILED"

type Error struct {
	Code    ErrorCode
	Message string
}

func (e *Error) Error() string { return fmt.Sprintf("%s: %s", e.Code, e.Message) }

type ComponentInput struct {
	Name        string
	Version     string
	Ecosystem   string
	WorkspaceID string
}

type Component struct {
	Name        string
	Version     string
	Ecosystem   string
	PURL        string
	WorkspaceID string
}

func NormalizeComponent(input ComponentInput) (Component, error) {
	name := strings.TrimSpace(input.Name)
	version := strings.TrimSpace(input.Version)
	ecosystem := strings.ToLower(strings.TrimSpace(input.Ecosystem))
	if name == "" || version == "" || ecosystem == "" {
		return Component{}, &Error{Code: ErrIncompleteIdentity, Message: "component name, version, and ecosystem are required"}
	}
	purl, err := makePURL(ecosystem, name, version)
	if err != nil {
		return Component{}, err
	}
	return Component{Name: name, Version: version, Ecosystem: ecosystem, PURL: purl, WorkspaceID: input.WorkspaceID}, nil
}

func makePURL(ecosystem, name, version string) (string, error) {
	encode := func(value string) string { return strings.ReplaceAll(url.PathEscape(value), "@", "%40") }
	switch ecosystem {
	case "npm", "pnpm", "yarn":
		return "pkg:npm/" + encode(name) + "@" + encode(version), nil
	case "pip", "uv", "poetry":
		return "pkg:pypi/" + encode(strings.ToLower(name)) + "@" + encode(version), nil
	case "maven", "gradle":
		parts := strings.SplitN(name, ":", 2)
		if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
			return "", &Error{Code: ErrIncompleteIdentity, Message: "Maven components require group and artifact"}
		}
		return "pkg:maven/" + encode(parts[0]) + "/" + encode(parts[1]) + "@" + encode(version), nil
	default:
		return "", &Error{Code: ErrIncompleteIdentity, Message: "component ecosystem is unsupported"}
	}
}
