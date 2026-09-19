package v1alpha_test

import (
	"embed"
	"encoding/json"
	"testing"

	"github.com/santhosh-tekuri/jsonschema/v6"
)

//go:embed deprail.schema.json examples/*.json
var schemaFiles embed.FS

func TestSchemaAndExamplesAreValidJSON(t *testing.T) {
	for _, name := range []string{
		"deprail.schema.json",
		"examples/project.json",
		"examples/scan.json",
		"examples/partial-scan.json",
		"examples/failed-scan.json",
		"examples/invalid-scan.json",
	} {
		data, err := schemaFiles.ReadFile(name)
		if err != nil {
			t.Fatalf("read %s: %v", name, err)
		}
		if !json.Valid(data) {
			t.Errorf("%s is not valid JSON", name)
		}
	}
}

func TestExamplesValidateAgainstSchema(t *testing.T) {
	schemaBytes, err := schemaFiles.ReadFile("deprail.schema.json")
	if err != nil {
		t.Fatal(err)
	}
	var schemaDocument any
	if err := json.Unmarshal(schemaBytes, &schemaDocument); err != nil {
		t.Fatal(err)
	}

	const schemaURL = "https://deprail.dev/schemas/v1alpha/deprail.schema.json"
	compiler := jsonschema.NewCompiler()
	if err := compiler.AddResource(schemaURL, schemaDocument); err != nil {
		t.Fatal(err)
	}
	compiled, err := compiler.Compile(schemaURL)
	if err != nil {
		t.Fatalf("compile schema: %v", err)
	}

	for _, name := range []string{"examples/project.json", "examples/scan.json", "examples/partial-scan.json", "examples/failed-scan.json"} {
		data, err := schemaFiles.ReadFile(name)
		if err != nil {
			t.Fatal(err)
		}
		var document any
		if err := json.Unmarshal(data, &document); err != nil {
			t.Fatal(err)
		}
		if err := compiled.Validate(document); err != nil {
			t.Errorf("%s should validate: %v", name, err)
		}
	}

	invalidBytes, err := schemaFiles.ReadFile("examples/invalid-scan.json")
	if err != nil {
		t.Fatal(err)
	}
	var invalidDocument any
	if err := json.Unmarshal(invalidBytes, &invalidDocument); err != nil {
		t.Fatal(err)
	}
	if err := compiled.Validate(invalidDocument); err == nil {
		t.Error("invalid-scan.json unexpectedly validates")
	}
}

func TestValidExamplesHaveVersionedDocumentTypes(t *testing.T) {
	for _, name := range []string{"examples/project.json", "examples/scan.json", "examples/partial-scan.json", "examples/failed-scan.json"} {
		data, err := schemaFiles.ReadFile(name)
		if err != nil {
			t.Fatalf("read %s: %v", name, err)
		}
		var document struct {
			SchemaVersion string `json:"schema_version"`
			DocumentType  string `json:"document_type"`
		}
		if err := json.Unmarshal(data, &document); err != nil {
			t.Fatalf("decode %s: %v", name, err)
		}
		if document.SchemaVersion != "v1alpha" || document.DocumentType == "" {
			t.Errorf("%s must declare schema_version v1alpha and document_type", name)
		}
	}
}

func TestInvalidExampleDemonstratesContractFailure(t *testing.T) {
	data, err := schemaFiles.ReadFile("examples/invalid-scan.json")
	if err != nil {
		t.Fatal(err)
	}
	var document struct {
		Status string `json:"status"`
		Errors []any  `json:"errors"`
	}
	if err := json.Unmarshal(data, &document); err != nil {
		t.Fatal(err)
	}
	if document.Status != "safe" || document.Errors != nil {
		t.Fatal("invalid fixture no longer exercises invalid status and missing errors")
	}
}
