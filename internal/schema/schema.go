// Package schema validates raw intent.yaml bytes against the embedded JSON
// Schema contract.
//
// The schema checks shape only — required fields, lowercase identifier keys, no
// unknown fields, at-least-one edges. Meaning-level rules (does a reference
// resolve? is a record's affects domain-pure? does a key match this project's
// configured casing?) belong to the linter, not here.
package schema

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"

	"github.com/santhosh-tekuri/jsonschema/v6"
	"gopkg.in/yaml.v3"
)

//go:embed intent.schema.json
var schemaJSON []byte

// compiled is built once, on first use.
var compiled *jsonschema.Schema

func load() (*jsonschema.Schema, error) {
	if compiled != nil {
		return compiled, nil
	}
	doc, err := jsonschema.UnmarshalJSON(bytes.NewReader(schemaJSON))
	if err != nil {
		return nil, fmt.Errorf("reading embedded schema: %w", err)
	}
	c := jsonschema.NewCompiler()
	if err := c.AddResource("intent.schema.json", doc); err != nil {
		return nil, fmt.Errorf("registering schema: %w", err)
	}
	s, err := c.Compile("intent.schema.json")
	if err != nil {
		return nil, fmt.Errorf("compiling schema: %w", err)
	}
	compiled = s
	return compiled, nil
}

// Validate checks that YAML bytes satisfy the structural schema. It returns nil
// when the shape is valid.
func Validate(yamlBytes []byte) error {
	s, err := load()
	if err != nil {
		return err
	}
	// YAML -> generic value -> JSON -> JSON-native value, so the validator sees
	// the exact types it expects (map[string]any, []any, string, ...).
	var raw interface{}
	if err := yaml.Unmarshal(yamlBytes, &raw); err != nil {
		return fmt.Errorf("parsing yaml: %w", err)
	}
	jsonBytes, err := json.Marshal(raw)
	if err != nil {
		return fmt.Errorf("converting to json: %w", err)
	}
	inst, err := jsonschema.UnmarshalJSON(bytes.NewReader(jsonBytes))
	if err != nil {
		return fmt.Errorf("reading instance: %w", err)
	}
	if err := s.Validate(inst); err != nil {
		return err
	}
	return nil
}
