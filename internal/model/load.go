package model

import (
	"bytes"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Load reads and parses an intent.yaml file from disk.
func Load(path string) (*Root, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	r, err := Parse(b)
	if err != nil {
		return nil, fmt.Errorf("parsing %s: %w", path, err)
	}
	return r, nil
}

// Parse turns raw YAML bytes into a Root.
func Parse(b []byte) (*Root, error) {
	var r Root
	if err := yaml.Unmarshal(b, &r); err != nil {
		return nil, err
	}
	return &r, nil
}

// Serialize renders a Root back to canonical YAML.
//
// The CLI is the sole writer and always emits this canonical form, so two
// serializations of the same model are byte-identical and diffs stay clean
// (DESIGN §11). Indentation is fixed at two spaces.
func (r *Root) Serialize() ([]byte, error) {
	var buf bytes.Buffer
	enc := yaml.NewEncoder(&buf)
	enc.SetIndent(2)
	if err := enc.Encode(r); err != nil {
		return nil, err
	}
	if err := enc.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
