package orderedmap_test

import (
	"testing"

	"github.com/schoren/orderedmap"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var sample = []struct {
	key   string
	value string
}{
	{"first item", "a"},
	{"this is the second item", "b"},
	{"3rd item", "c"},
}

var sampleUnordered = map[string]string{
	"first item":              "a",
	"this is the second item": "b",
	"3rd item":                "c",
}

var sampleJSON = `[
	{"Key":"first item","Value":"a"},
	{"Key":"this is the second item","Value":"b"},
	{"Key":"3rd item","Value":"c"}
]`

func newOMFromSample(t *testing.T) orderedmap.OrderedMap[string, string] {
	om := orderedmap.New[string, string]()
	for _, s := range sample {
		var err error

		om, err = om.Set(s.key, s.value)
		require.NoError(t, err)
	}

	return om
}

func TestOrderIsMaintained(t *testing.T) {
	t.Parallel()

	om := newOMFromSample(t)

	assert.Equal(t, len(sample), om.Len())

	assertSampleOrder(t, om)

}

func assertSampleOrder(t *testing.T, om orderedmap.OrderedMap[string, string]) {
	t.Helper()

	i := 0
	foreachErr := om.ForEach(func(key string, val string) error {
		assert.Equal(t, sample[i].key, key, `expected index %d to have key "%s", got "%s"`, i, sample[i].key, key)
		assert.Equal(t, sample[i].value, val, `expected index %d to have value "%s", got "%s"`, i, sample[i].value, val)
		i++
		return nil
	})
	require.NoError(t, foreachErr)
}

func TestUnordered(t *testing.T) {
	t.Parallel()

	om := newOMFromSample(t)

	assert.Equal(t, sampleUnordered, om.Unordered())
}

func TestForeachError(t *testing.T) {
	t.Parallel()

	om := newOMFromSample(t)

	err := om.ForEach(func(key string, val string) error {
		if key == "this is the second item" {
			return assert.AnError
		}
		return nil
	})

	assert.ErrorIs(t, err, assert.AnError)
}

func TestSetReturnsErrorIfKeyAlreadyExists(t *testing.T) {
	t.Parallel()

	om := orderedmap.New[string, string]()
	om, err := om.Set("key", "value")
	require.NoError(t, err)

	_, err = om.Set("key", "value")
	assert.ErrorIs(t, err, orderedmap.ErrKeyAlreadyExists)

	errorMsg := `key "key" already exists`
	assert.EqualError(t, err, errorMsg)
}

func TestMustSetPanic(t *testing.T) {
	t.Parallel()

	om := orderedmap.New[string, string]()
	om = om.MustSet("key", "value")

	assert.Panics(t, func() {
		om.MustSet("key", "value")
	})
}

func TestSetOnUninitializedMap(t *testing.T) {
	t.Parallel()

	var om orderedmap.OrderedMap[string, string]
	om, err := om.Set("key", "value")
	require.NoError(t, err)

	assert.Equal(t, 1, om.Len())
}

func TestGet(t *testing.T) {
	t.Parallel()

	om := orderedmap.New[string, string]().
		MustSet("key", "value")

	assert.Equal(t, "value", om.Get("key"))
	assert.Equal(t, "", om.Get("NotExists"))
}

func TestContains(t *testing.T) {
	t.Parallel()

	om := orderedmap.New[string, string]().
		MustSet("key", "value")

	assert.True(t, om.Contains("key"))
	assert.False(t, om.Contains("NotExists"))
}

func TestDelete(t *testing.T) {
	t.Parallel()

	om := newOMFromSample(t)
	om = om.Delete("this is the second item")

	assert.Equal(t, 2, om.Len())
	assert.True(t, om.Contains("first item"))
	assert.False(t, om.Contains("this is the second item"))
	assert.True(t, om.Contains("3rd item"))

}

func TestJSON(t *testing.T) {
	t.Parallel()

	om := newOMFromSample(t)

	data, err := om.MarshalJSON()
	require.NoError(t, err)

	assert.JSONEq(t, sampleJSON, string(data))

	var unmarshalled orderedmap.OrderedMap[string, string]
	err = unmarshalled.UnmarshalJSON(data)
	require.NoError(t, err)

	assertSampleOrder(t, om)
}

func TestNonUniqueJson(t *testing.T) {
	t.Parallel()

	nonUniqueJSON := `[{"Key":"1","Value":"a"},{"Key":"1","Value":"b"}]`

	var om orderedmap.OrderedMap[string, string]
	err := om.UnmarshalJSON([]byte(nonUniqueJSON))
	assert.ErrorIs(t, err, orderedmap.ErrKeyAlreadyExists)
}
