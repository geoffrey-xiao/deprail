package historyv1

import (
	"bytes"
	"embed"
	"encoding/json"
	"errors"
	"io"
	"regexp"
	"strings"
	"sync"

	"github.com/santhosh-tekuri/jsonschema/v6"
)

//go:embed history.schema.json
var schemaFS embed.FS

const schemaURL = "https://deprail.dev/schemas/history-v1/history.schema.json"
const historyPathPattern = `^(?![\s\S]*[\u0000-\u001f])(?:\.|(?!/)(?![A-Za-z]:)(?!.*(?:^|/)\.\.(?:/|$))(?!.*(?:^|/)\.(?:/|$))(?!.*\\)[^/]+(?:/[^/]+)*)$`

var (
	compiledSchemaOnce sync.Once
	compiledSchema     *jsonschema.Schema
	compiledSchemaErr  error
)

func readHistorySchema() ([]byte, error) {
	return schemaFS.ReadFile("history.schema.json")
}

// ValidateJSON checks a serialized history projection against the committed
// history-v1 JSON Schema, including its reviewed repository-relative path rule.
func ValidateJSON(data []byte) error {
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()
	var document any
	if err := decoder.Decode(&document); err != nil {
		return errors.New("history projection is not valid JSON")
	}
	var trailing any
	if err := decoder.Decode(&trailing); !errors.Is(err, io.EOF) {
		return errors.New("history projection has trailing JSON data")
	}
	schema, err := compiledHistorySchema()
	if err != nil {
		return errors.New("history schema could not be loaded")
	}
	if err := schema.Validate(document); err != nil {
		return errors.New("history projection does not satisfy history-v1 schema")
	}
	return nil
}

func compiledHistorySchema() (*jsonschema.Schema, error) {
	compiledSchemaOnce.Do(func() {
		schemaBytes, err := readHistorySchema()
		if err != nil {
			compiledSchemaErr = err
			return
		}
		var schemaDocument any
		if err := json.Unmarshal(schemaBytes, &schemaDocument); err != nil {
			compiledSchemaErr = err
			return
		}
		compiler := newHistorySchemaCompiler()
		if err := compiler.AddResource(schemaURL, schemaDocument); err != nil {
			compiledSchemaErr = err
			return
		}
		compiledSchema, compiledSchemaErr = compiler.Compile(schemaURL)
	})
	return compiledSchema, compiledSchemaErr
}

// The schema path pattern uses negative lookaheads unsupported by Go regexp.
// This matcher preserves the committed pattern's traversal and platform rules.
type historySchemaRegexp struct {
	source string
	re     *regexp.Regexp
	path   bool
}

func (r historySchemaRegexp) String() string { return r.source }
func (r historySchemaRegexp) MatchString(value string) bool {
	if !r.re.MatchString(value) {
		return false
	}
	if !r.path {
		return true
	}
	if value == "." {
		return true
	}
	if value == "" || strings.HasPrefix(value, "/") || strings.Contains(value, `\`) {
		return false
	}
	if len(value) >= 2 && ((value[0] >= 'A' && value[0] <= 'Z') || (value[0] >= 'a' && value[0] <= 'z')) && value[1] == ':' {
		return false
	}
	for _, char := range value {
		if char < 0x20 {
			return false
		}
	}
	for _, part := range strings.Split(value, "/") {
		if part == "" || part == "." || part == ".." {
			return false
		}
	}
	return true
}

func historyRegexpEngine(pattern string) (jsonschema.Regexp, error) {
	if pattern == historyPathPattern {
		return historySchemaRegexp{source: pattern, re: regexp.MustCompile(`^(?:\.|[^/]+(?:/[^/]+)*)$`), path: true}, nil
	}
	compiled, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}
	return historySchemaRegexp{source: pattern, re: compiled}, nil
}

func newHistorySchemaCompiler() *jsonschema.Compiler {
	compiler := jsonschema.NewCompiler()
	compiler.UseRegexpEngine(historyRegexpEngine)
	return compiler
}
