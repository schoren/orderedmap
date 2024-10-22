package orderedmap

import (
	"encoding/json"
	"errors"
	"fmt"
)

// OrderedMap is a map that maintains the order of insertion.
type OrderedMap[K comparable, V any] struct {
	list        []V
	keyPosition map[K]int
	positionKey map[int]K
}

// New creates a new OrderedMap.
func New[K comparable, V any]() OrderedMap[K, V] {
	return OrderedMap[K, V]{
		list:        []V{},
		keyPosition: make(map[K]int),
		positionKey: make(map[int]K),
	}
}

// MustSet is like Set, but panics if an error occurs.
// It simplifies initialization enabling chaining.
func (om OrderedMap[K, V]) MustSet(key K, asserts V) OrderedMap[K, V] {
	def, err := om.Set(key, asserts)
	if err != nil {
		panic(err)
	}
	return def
}

type keyAlreadyExistsError struct {
	key any
}

func (e keyAlreadyExistsError) Unwrap() error {
	return ErrKeyAlreadyExists
}

func (e keyAlreadyExistsError) Error() string {
	return fmt.Sprintf(`key "%v" already exists`, e.key)
}

// ErrKeyAlreadyExists is returned when trying to add a key that already exists.
var ErrKeyAlreadyExists = errors.New("key already exists")

// Set adds a new key-value pair to the map.
// If the key already exists, an error is returned.
func (om OrderedMap[K, V]) Set(key K, asserts V) (OrderedMap[K, V], error) {
	if om.keyPosition == nil {
		om.keyPosition = make(map[K]int)
	}
	if om.positionKey == nil {
		om.positionKey = make(map[int]K)
	}

	if _, exists := om.keyPosition[key]; exists {
		return OrderedMap[K, V]{}, keyAlreadyExistsError{key}
	}

	om.list = append(om.list, asserts)
	ix := len(om.list) - 1
	om.keyPosition[key] = ix
	om.positionKey[ix] = key

	return om, nil
}

// Delete removes a key from the map.
// If the key does not exist, the map is returned unchanged.
func (om OrderedMap[K, V]) Delete(key K) OrderedMap[K, V] {
	ix, exists := om.keyPosition[key]
	if !exists {
		return om
	}

	delete(om.keyPosition, key)
	delete(om.positionKey, ix)

	om.list = append(om.list[:ix], om.list[ix+1:]...)
	for i := ix; i < len(om.list); i++ {
		k := om.positionKey[i+1]
		om.keyPosition[k] = i
		om.positionKey[i] = k
	}

	return om
}

// Len returns the number of elements in the map.
func (om OrderedMap[K, V]) Len() int {
	return len(om.list)
}

// Contains returns true if the key exists in the map.
func (om OrderedMap[K, V]) Contains(key K) bool {
	_, exists := om.keyPosition[key]
	return exists
}

// Get returns the value associated with the key.
// If the key does not exist, the zero value of the value type is returned.
func (om OrderedMap[K, V]) Get(key K) V {
	ix, exists := om.keyPosition[key]
	if !exists {
		var result V
		return result
	}

	return om.list[ix]
}

// ForEach iterates over the map, calling the function for each key-value pair.
// If the function returns an error, the iteration stops and the error is returned.
func (om *OrderedMap[K, V]) ForEach(fn func(key K, val V) error) error {
	for ix, asserts := range om.list {
		K := om.positionKey[ix]
		err := fn(K, asserts)
		if err != nil {
			return err
		}
	}

	return nil
}

// Unordered returns a map with the same key-value pairs, but in an unordered map.
func (om OrderedMap[K, V]) Unordered() map[K]V {
	m := map[K]V{}
	om.ForEach(func(key K, val V) error {
		m[key] = val
		return nil
	})

	return m
}

func (om *OrderedMap[K, V]) replace(om2 *OrderedMap[K, V]) {
	*om = *om2
}

type jsonOrderedMapEntry[K comparable, V any] struct {
	Key   K
	Value V
}

func (om OrderedMap[K, V]) MarshalJSON() ([]byte, error) {
	j := []jsonOrderedMapEntry[K, V]{}
	om.ForEach(func(key K, asserts V) error {
		j = append(j, jsonOrderedMapEntry[K, V]{key, asserts})
		return nil
	})

	return json.Marshal(j)
}

func (om *OrderedMap[K, V]) UnmarshalJSON(data []byte) error {
	aux := []jsonOrderedMapEntry[K, V]{}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	newMap := OrderedMap[K, V]{}
	var err error
	for _, s := range aux {
		newMap, err = newMap.Set(s.Key, s.Value)
		if err != nil {
			return err
		}
	}

	om.replace(&newMap)

	return nil
}