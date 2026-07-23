package model

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

// OrderedMap is a string-keyed map that remembers the order keys were inserted.
//
// The intent model is full of dicts whose order is meaningful: presentation
// order equals insertion order (DESIGN §3). A plain Go map loses that order and
// yaml.v3 would re-emit its keys sorted, so we use this type wherever the file
// holds a dict of elements (jobs, outcomes, components, records, ...).
type OrderedMap[V any] struct {
	keys []string
	m    map[string]V
}

// Len reports how many entries the map holds.
func (o *OrderedMap[V]) Len() int { return len(o.keys) }

// Keys returns the keys in insertion order. The slice is a copy; callers may
// modify it freely.
func (o *OrderedMap[V]) Keys() []string {
	out := make([]string, len(o.keys))
	copy(out, o.keys)
	return out
}

// Get looks up a value by key.
func (o *OrderedMap[V]) Get(k string) (V, bool) {
	v, ok := o.m[k]
	return v, ok
}

// Set inserts or updates a key. New keys are appended to the end; existing keys
// keep their position.
func (o *OrderedMap[V]) Set(k string, v V) {
	if o.m == nil {
		o.m = map[string]V{}
	}
	if _, ok := o.m[k]; !ok {
		o.keys = append(o.keys, k)
	}
	o.m[k] = v
}

// Delete removes a key, preserving the order of the remaining keys. It reports
// whether the key was present. The write path (rm) is the caller (DESIGN §11).
func (o *OrderedMap[V]) Delete(k string) bool {
	if _, ok := o.m[k]; !ok {
		return false
	}
	delete(o.m, k)
	for i, kk := range o.keys {
		if kk == k {
			o.keys = append(o.keys[:i], o.keys[i+1:]...)
			break
		}
	}
	return true
}

// Rename changes key `old` to `new`, keeping the entry's value and its position
// in the order. It reports false if `old` is absent or `new` already exists (a
// rename must never collide with or reorder around a sibling). The write path
// (mv --rename) is the caller.
func (o *OrderedMap[V]) Rename(old, new string) bool {
	if _, ok := o.m[old]; !ok {
		return false
	}
	if _, exists := o.m[new]; exists {
		return false
	}
	o.m[new] = o.m[old]
	delete(o.m, old)
	for i, k := range o.keys {
		if k == old {
			o.keys[i] = new
			break
		}
	}
	return true
}

// Reorder moves an existing key to sit immediately before or after ref (also
// existing). It errors if either key is absent; moving a key relative to itself
// is a no-op. The write path (mv --before/--after) is the caller.
func (o *OrderedMap[V]) Reorder(key, ref string, after bool) error {
	if _, ok := o.m[key]; !ok {
		return fmt.Errorf("no such key %q", key)
	}
	if _, ok := o.m[ref]; !ok {
		return fmt.Errorf("no such key %q", ref)
	}
	if key == ref {
		return nil
	}
	ks := make([]string, 0, len(o.keys))
	for _, k := range o.keys {
		if k != key {
			ks = append(ks, k)
		}
	}
	pos := 0
	for i, k := range ks {
		if k == ref {
			pos = i
			if after {
				pos = i + 1
			}
			break
		}
	}
	ks = append(ks[:pos:pos], append([]string{key}, ks[pos:]...)...)
	o.keys = ks
	return nil
}

// IsZero lets yaml.v3's `omitempty` drop an empty map. Without this, a struct
// with only unexported fields always looks "zero" to yaml.v3, which would omit
// even a populated map — so we spell the rule out explicitly.
func (o OrderedMap[V]) IsZero() bool { return len(o.keys) == 0 }

// MarshalYAML emits a mapping node with entries in insertion order.
func (o OrderedMap[V]) MarshalYAML() (interface{}, error) {
	n := &yaml.Node{Kind: yaml.MappingNode, Tag: "!!map"}
	for _, k := range o.keys {
		var kn yaml.Node
		if err := kn.Encode(k); err != nil {
			return nil, err
		}
		var vn yaml.Node
		if err := vn.Encode(o.m[k]); err != nil {
			return nil, err
		}
		n.Content = append(n.Content, &kn, &vn)
	}
	return n, nil
}

// UnmarshalYAML reads a mapping node, preserving the order keys appear in the
// file.
func (o *OrderedMap[V]) UnmarshalYAML(node *yaml.Node) error {
	if node.Kind != yaml.MappingNode {
		return fmt.Errorf("expected a mapping, got yaml kind %d", node.Kind)
	}
	o.m = map[string]V{}
	o.keys = o.keys[:0]
	for i := 0; i+1 < len(node.Content); i += 2 {
		k := node.Content[i].Value
		var v V
		if err := node.Content[i+1].Decode(&v); err != nil {
			return fmt.Errorf("decoding %q: %w", k, err)
		}
		o.keys = append(o.keys, k)
		o.m[k] = v
	}
	return nil
}
