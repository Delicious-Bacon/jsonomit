package jsonomit

import (
	"encoding/json"
	"testing"
	"time"
)

type customStruct struct {
	Value any
	valid bool
}

func (s customStruct) MarshalJSON() ([]byte, error) {
	if !s.valid { // logic to set the struct to null.
		return []byte("null"), nil
	}

	return json.Marshal(s.Value)
}

type nestedStruct struct {
	T           time.Time
	TEmpty      time.Time
	Level       string `json:",omitempty"`
	Custom      customStruct
	CustomEmpty customStruct
}

type testStruct struct {
	T           time.Time
	TEmpty      time.Time
	Level       string `json:",omitempty"`
	Custom      customStruct
	CustomEmpty customStruct
	Nested      nestedStruct
	NestedEmpty nestedStruct
}

func TestMarshal(t *testing.T) {
	ts := time.Unix(0, 0).UTC()
	withVals := testStruct{
		T:      ts,
		TEmpty: time.Time{},
		Level:  "one",
		Custom: customStruct{
			Value: "value",
			valid: true,
		},
		CustomEmpty: customStruct{
			Value: "empty",
			valid: false,
		},
		Nested: nestedStruct{
			T:      ts,
			TEmpty: time.Time{},
			Level:  "two",
			Custom: customStruct{
				Value: "test",
				valid: true,
			},
			CustomEmpty: customStruct{
				Value: "empty",
				valid: false,
			},
		},
		NestedEmpty: nestedStruct{
			TEmpty: time.Time{},
			CustomEmpty: customStruct{
				Value: "empty",
				valid: false,
			},
		},
	}

	b, err := Marshal(withVals)
	if err != nil {
		t.Fatal(err)
	}
	if string(b) != `{"T":"1970-01-01T00:00:00Z","Level":"one","Custom":"value","Nested":{"T":"1970-01-01T00:00:00Z","Level":"two","Custom":"test"}}` {
		t.Fatal(string(b))
	} else {
		t.Log("With values: OK!")
		t.Log(string(b))
	}

	b, err = Marshal(testStruct{})
	if err != nil {
		t.Fatal(err)
	}
	if string(b) != `{}` {
		t.Fatal(string(b))
	} else {
		t.Log("With no values: OK!")
		t.Log(string(b))
	}
}
