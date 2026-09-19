package v0alpha1_test

import (
	"embed"
	"encoding/json"
	"testing"

	"github.com/santhosh-tekuri/jsonschema/v6"
)

//go:embed remediation-plan.schema.json examples/*.json
var schemaFiles embed.FS

const schemaURL = "https://deprail.dev/schemas/v0alpha1/remediation-plan.schema.json"

func compiledSchema(t *testing.T) *jsonschema.Schema {
	t.Helper()
	bytes, err := schemaFiles.ReadFile("remediation-plan.schema.json")
	if err != nil {
		t.Fatal(err)
	}
	var document any
	if err := json.Unmarshal(bytes, &document); err != nil {
		t.Fatal(err)
	}
	compiler := jsonschema.NewCompiler()
	if err := compiler.AddResource(schemaURL, document); err != nil {
		t.Fatal(err)
	}
	compiled, err := compiler.Compile(schemaURL)
	if err != nil {
		t.Fatal(err)
	}
	return compiled
}

func TestExamplesValidateAgainstSchema(t *testing.T) {
	compiled := compiledSchema(t)
	for _, name := range []string{
		"recommended", "multiple-candidates", "no-recommendation", "unknown-candidate",
		"rejected-candidate", "stale-input", "incomplete-plan",
	} {
		bytes, err := schemaFiles.ReadFile("examples/" + name + ".json")
		if err != nil {
			t.Fatal(err)
		}
		var document any
		if err := json.Unmarshal(bytes, &document); err != nil {
			t.Fatalf("%s: invalid JSON: %v", name, err)
		}
		if err := compiled.Validate(document); err != nil {
			t.Errorf("%s: schema validation failed: %v", name, err)
		}
	}
}

func TestSchemaRejectsMissingProvenanceAndUnsafeCommandPath(t *testing.T) {
	compiled := compiledSchema(t)
	bytes, err := schemaFiles.ReadFile("examples/recommended.json")
	if err != nil {
		t.Fatal(err)
	}
	var document map[string]any
	if err := json.Unmarshal(bytes, &document); err != nil {
		t.Fatal(err)
	}
	delete(document, "created_from")
	if err := compiled.Validate(document); err == nil {
		t.Fatal("missing provenance unexpectedly validated")
	}

	if err := json.Unmarshal(bytes, &document); err != nil {
		t.Fatal(err)
	}
	commands := document["commands"].([]any)
	commands[0].(map[string]any)["working_directory"] = "/outside"
	if err := compiled.Validate(document); err == nil {
		t.Fatal("absolute command working directory unexpectedly validated")
	}
}

func TestSchemaRejectsUnsupportedVersionAndMissingWorkingDirectory(t *testing.T) {
	compiled := compiledSchema(t)
	bytes, err := schemaFiles.ReadFile("examples/recommended.json")
	if err != nil {
		t.Fatal(err)
	}
	var document map[string]any
	if err := json.Unmarshal(bytes, &document); err != nil {
		t.Fatal(err)
	}
	document["schema_version"] = "v9"
	if err := compiled.Validate(document); err == nil {
		t.Fatal("unsupported schema version unexpectedly validated")
	}

	if err := json.Unmarshal(bytes, &document); err != nil {
		t.Fatal(err)
	}
	commands := document["commands"].([]any)
	delete(commands[0].(map[string]any), "working_directory")
	if err := compiled.Validate(document); err == nil {
		t.Fatal("missing working directory unexpectedly validated")
	}
}
