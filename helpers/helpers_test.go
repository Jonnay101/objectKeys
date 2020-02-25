package helpers

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Thing is a struct type
type Thing struct {
	Key1 string     `json:"key1"`
	Key2 int        `json:"key2"`
	Key3 []int      `json:"key3"`
	Key4 InnerThing `json:"key4"`
}

// InnerThing - is a struct type
type InnerThing struct {
	Key5 float32         `json:"key5"`
	Key6 InnerInnerThing `json:"key6"`
}

// InnerInnerThing - is a struct in type inside the inside struct
type InnerInnerThing struct {
	Key7 string `json:"key7"`
	Key8 string `json:"key8"`
}

type UnExportedField struct {
	key1 string
	Key2 string
}

type PointerField struct {
	Key1 *InnerInnerThing
	Key2 string
}

type SliceThing []string

func TestObjectKeys(t *testing.T) {
	table := []struct {
		Obj     interface{}
		Expect  interface{}
		Errored error
		Msg     string
	}{
		{Thing{},
			[]string{"Key1", "Key2", "Key3", "Key4"},
			nil,
			"this is a plain struct type",
		},
		{&Thing{},
			[]string{"Key1", "Key2", "Key3", "Key4"},
			nil,
			"this is a plain struct type",
		},
		{UnExportedField{"this", "that"},
			[]string{"key1", "Key2"},
			nil,
			"this is with unexported fields",
		},
		{PointerField{&InnerInnerThing{}, "three"},
			[]string{"Key1", "Key2"},
			nil,
			"this is with one field having a pointer type",
		},
		{SliceThing{},
			[]string(nil),
			errors.New("this function only accepts struct types or pointer-to-struct types"),
			"non struct type passed to the function",
		},
	}

	for _, test := range table {
		res, err := ObjectKeys(test.Obj)

		assert.Equal(t, test.Errored, err, test.Msg)
		assert.Equal(t, test.Expect, res, test.Msg)
	}
}

func TestObjectKeysFlatten(t *testing.T) {
	table := []struct {
		Obj     interface{}
		Expect  interface{}
		Errored error
		Msg     string
	}{
		{Thing{},
			[]string{"Key1", "Key2", "Key3", "Key4", "Key5", "Key6", "Key7", "Key8"},
			nil,
			"this is a plain struct type",
		},
		{&Thing{},
			[]string{"Key1", "Key2", "Key3", "Key4", "Key5", "Key6", "Key7", "Key8"},
			nil,
			"this is a plain struct type",
		},
		{UnExportedField{"this", "that"},
			[]string{"key1", "Key2"},
			nil,
			"this is with unexported fields",
		},
		{PointerField{&InnerInnerThing{}, "three"},
			[]string{"Key1", "Key7", "Key8", "Key2"},
			nil,
			"this is with one field having a pointer type",
		},
		{SliceThing{},
			[]string(nil),
			errors.New("this function only accepts struct types or pointer-to-struct types"),
			"non struct type passed to the function",
		},
	}

	for _, test := range table {
		res, err := ObjectKeysFlatten(test.Obj)

		assert.Equal(t, test.Errored, err, test.Msg)
		assert.Equal(t, test.Expect, res, test.Msg)
	}
}

type GetFieldStruct struct {
	Key1  string
	Inner InnerGet
}

type InnerGet struct {
	Number int
}

func TestGet(t *testing.T) {
	table := []struct {
		Obj     interface{}
		Key     string
		Expect  interface{}
		Errored error
		Msg     string
	}{
		{GetFieldStruct{"result1", InnerGet{3}}, "Key1", "result1", nil, "should pass"},
		{GetFieldStruct{"result1", InnerGet{3}}, "Inner", InnerGet{3}, nil, "should pass"},
		{GetFieldStruct{"result1", InnerGet{3}}, "Foo", nil, errors.New("Sorry, your key doesn't match any fields in the provided struct type"), "should error"},
	}

	for _, test := range table {
		res, err := Get(test.Obj, test.Key)

		assert.Equal(t, test.Expect, res, test.Msg)
		assert.Equal(t, test.Errored, err, test.Msg)
	}
}

func TestGetVals(t *testing.T) {
	table := []struct {
		Obj     interface{}
		Expect  interface{}
		Errored error
		Msg     string
	}{
		{
			GetFieldStruct{
				Key1:  "result1",
				Inner: InnerGet{3},
			},
			[]interface{}{"result1", InnerGet{3}},
			nil,
			"should pass",
		},
		{
			Thing{
				Key1: "string",
				Key2: 3,
				Key3: []int{3, 2, 1},
				Key4: InnerThing{
					Key5: 3.1,
					Key6: InnerInnerThing{"this", "that"},
				},
			},
			[]interface{}{"string", 3, []int{3, 2, 1}, InnerThing{3.1, InnerInnerThing{"this", "that"}}},
			nil,
			"this should pass becouse it's",
		},
		{
			[]string{"this", "is", "not", "a", "struct"},
			[]interface{}{},
			errors.New("this function only accepts struct types or pointer-to-struct types"),
			"this should pass too",
		},
	}

	for _, test := range table {
		t.Run(test.Msg, func(t *testing.T) {
			res, err := GetVals(test.Obj)
			assert.Equal(t, test.Expect, res, test.Msg)
			assert.Equal(t, test.Errored, err, test.Msg)
		})
	}
}
