// Copyright 2020 Julien Cretel (jub0bs). All rights reserved.
// Use of this source code is governed by a BSD 3-clause
// license that can be found in the LICENSE file.

// Package bimap provides a bidirectional map of some comparable
// key and value types.
package bimap

import "fmt"

// A Bimap is a bidirectional map, i.e. an associative data
// structure in which key-value pairs form a one-to-one
// correspondence.
// Both keys and values must be comparable.
// Loads, stores, and deletes run in amortized constant time.
//
// The zero value for Bimap is empty and ready for use.
// A Bimap must not be copied after first use.
type Bimap[K, V comparable] struct {
	forward map[K]V
	inverse map[V]K
}

// New returns a new, empty Bimap.
func New[K, V comparable]() *Bimap[K, V] {
	return &Bimap[K, V]{}
}

// Store creates a key-value pair and returns whether or not the
// operation was successful. Pre-existing key-value pairs (if any)
// that involve the given key and/or the given value are silently
// removed from the Bimap. Keys and values for which equality is
// not reflexive are disallowed.
func (bi *Bimap[K, V]) Store(key K, value V) bool {
	if !isEqualityReflexive(key) || !isEqualityReflexive(value) {
		return false
	}
	k, exists := bi.inverse[value]
	if exists { // value is already associated with k
		delete(bi.forward, k)
	}
	v, exists := bi.forward[key]
	if exists { // key is already associated with v
		delete(bi.inverse, v)
	}
	if bi.forward == nil { // bi hasn't been initialised yet
		bi.forward = make(map[K]V)
		bi.inverse = make(map[V]K)
	}
	bi.forward[key] = value
	bi.inverse[value] = key
	return true
}

func isEqualityReflexive[T comparable](t T) bool {
	return t == t
}

// LoadValue returns the value stored in the Bimap for a key,
// or the zero value of the K type if no value is present.
// The ok result indicates whether the key was found in the map.
func (bi *Bimap[K, V]) LoadValue(k K) (V, bool) {
	v, ok := bi.forward[k]
	return v, ok
}

// LoadKey returns the key stored in the Bimap for a key,
// or the zero value of the V type if no key is present.
// The ok result indicates whether the value was found in the
// map.
func (bi *Bimap[K, V]) LoadKey(v V) (K, bool) {
	k, ok := bi.inverse[v]
	return k, ok
}

// DeleteByKey deletes the key-value pair involving the given
// key.
func (bi *Bimap[K, V]) DeleteByKey(k K) {
	v := bi.forward[k]
	delete(bi.forward, k)
	delete(bi.inverse, v)
}

// DeleteByValue deletes the key-value pair involving the given
// value.
func (bi *Bimap[K, V]) DeleteByValue(v V) {
	k := bi.inverse[v]
	delete(bi.inverse, v)
	delete(bi.forward, k)
}

// Size returns the number of key-value pairs in the Bimap.
// The complexity is O(1).
func (bi *Bimap[K, V]) Size() int {
	return len(bi.forward)
}

// Keys returns a slice of the keys in the Bimap.
func (bi *Bimap[K, V]) Keys() []K {
	var keys []K
	for k := range bi.forward {
		keys = append(keys, k)
	}
	return keys
}

// Values returns a slice of the values in the Bimap.
func (bi *Bimap[K, V]) Values() []V {
	var values []V
	for v := range bi.inverse {
		values = append(values, v)
	}
	return values
}

// String returns a string representing the Bimap. That string
// representation is similar to the string representation of a
// built-in map.
func (bi *Bimap[K, V]) String() string {
	return fmt.Sprintf("Bi%v", bi.forward)
}
