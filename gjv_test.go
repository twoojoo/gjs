package gjs_test

import (
	"gjs"
	"os"
	"testing"
)

func TestValidate(t *testing.T) {
	type TestData struct {
		Name  string `json:"name" jsonschema:"maxLength=2"`
		Value int    `json:"value" jsonschema:"maximum=4"`
	}

	data := TestData{Name: "Example", Value: 42}
	schema := gjs.NewSchema[TestData]()
	result, err := schema.Validate(&data)
	if err != nil {
		t.Fatalf("Validation failed: %v", err)
	}

	t.Logf("Validation result: %+v", result.Errors())

	if result.Valid() {
		t.Errorf("Expected validation to fail, but it passed")
	}

	if len(result.Errors()) == 0 {
		t.Errorf("Expected validation errors, but got none")
	}

	if len(result.Errors()) != 2 {
		t.Errorf("Expected 2 validation errors, got %d", len(result.Errors()))
	}
}

func TestSchemString(t *testing.T) {
	type TestData struct {
		Name  string `json:"name" jsonschema:"maxLength=2"`
		Value int    `json:"value" jsonschema:"maximum=4"`
	}

	schema := gjs.NewSchema[TestData]()
	str, err := schema.String(gjs.WithIndent(" "))
	if err != nil {
		t.Errorf("Schema extraction failed: %v", err)
	}

	t.Logf("%v", str)
}

func TestStore(t *testing.T) {
	type TestData struct {
		Name  string `json:"name" jsonschema:"maxLength=2"`
		Value int    `json:"value" jsonschema:"maximum=4"`
	}

	val := gjs.NewSchema[TestData]()
	err := val.Store("./schema.json", gjs.WithIndent("\t"))
	if err != nil {
		t.Fatalf("Schema to file failed: %v", err)
	}

	file, err := os.ReadFile("./schema.json")
	if err != nil {
		t.Fatalf("Failed to open schema file: %v", err)
	}

	t.Log(string(file))
}
