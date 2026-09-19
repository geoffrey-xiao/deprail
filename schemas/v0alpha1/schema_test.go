package v0alpha1_test

import (
	"embed"
	"encoding/json"
	"testing"

	"github.com/geoffrey-xiao/deprail/internal/remediation"
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
		var plan remediation.Plan
		if err := json.Unmarshal(bytes, &plan); err != nil {
			t.Fatalf("%s: domain decode failed: %v", name, err)
		}
		if err := plan.Validate(); err != nil {
			t.Errorf("%s: domain validation failed: %v", name, err)
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

func TestSchemaRejectsTraversalUnsafeRecommendationAndBadDigest(t *testing.T) {
	compiled := compiledSchema(t)
	bytes, err := schemaFiles.ReadFile("examples/recommended.json")
	if err != nil {
		t.Fatal(err)
	}
	for _, directory := range []string{"../outside", "foo/../../outside", "./service"} {
		var document map[string]any
		if err := json.Unmarshal(bytes, &document); err != nil {
			t.Fatal(err)
		}
		document["commands"].([]any)[0].(map[string]any)["working_directory"] = directory
		if err := compiled.Validate(document); err == nil {
			t.Fatalf("traversal path %q unexpectedly validated", directory)
		}
	}

	var document map[string]any
	if err := json.Unmarshal(bytes, &document); err != nil {
		t.Fatal(err)
	}
	candidate := document["candidates"].([]any)[0].(map[string]any)
	candidate["evidence"].(map[string]any)["peer_compatible"] = false
	if err := compiled.Validate(document); err == nil {
		t.Fatal("unsafe recommended candidate unexpectedly validated")
	}

	if err := json.Unmarshal(bytes, &document); err != nil {
		t.Fatal(err)
	}
	document["provenance"].(map[string]any)["artifact_digests"] = []any{"not-a-digest"}
	if err := compiled.Validate(document); err == nil {
		t.Fatal("invalid artifact digest unexpectedly validated")
	}

	if err := json.Unmarshal(bytes, &document); err != nil {
		t.Fatal(err)
	}
	digest := document["provenance"].(map[string]any)["artifact_digests"].([]any)[0]
	document["provenance"].(map[string]any)["artifact_digests"] = []any{digest, digest}
	if err := compiled.Validate(document); err == nil {
		t.Fatal("duplicate artifact digest unexpectedly validated")
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
