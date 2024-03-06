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
	StringTrap  string `json:",omitempty"`
	Custom      customStruct
	CustomEmpty customStruct
	StructEmpty struct{}
}

type testStruct struct {
	T           time.Time
	TEmpty      time.Time
	StringTrap  string `json:",omitempty"`
	Custom      customStruct
	CustomEmpty customStruct
	Nested      nestedStruct
	NestedEmpty nestedStruct
	StructEmpty struct{}
}

var (
	ts = time.Unix(0, 0).UTC()

	withVals = testStruct{
		T:          ts,
		TEmpty:     time.Time{},
		StringTrap: `"Time":"0001-01-01T00:00:00Z"`,
		Custom: customStruct{
			Value: "value",
			valid: true,
		},
		CustomEmpty: customStruct{
			Value: "empty",
			valid: false,
		},
		Nested: nestedStruct{
			T:          ts,
			TEmpty:     time.Time{},
			StringTrap: `"MyStruct":null`,
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
)

func TestMarshal(t *testing.T) {

	// Populated testStruct
	b, err := Marshal(withVals)
	if err != nil {
		t.Fatal(err)
	}
	want := `{"T":"1970-01-01T00:00:00Z","StringTrap":"\"Time\":\"0001-01-01T00:00:00Z\"","Custom":"value","Nested":{"T":"1970-01-01T00:00:00Z","StringTrap":"\"MyStruct\":null","Custom":"test"}}`
	if string(b) != want {
		t.Fatalf("Failed want!\nWanted:\n%s\nGot:\n%s", want, string(b))
	} else {
		t.Log("Populated testStruct: OK!")
	}

	// Zero value testStruct
	b, err = Marshal(testStruct{})
	if err != nil {
		t.Fatal(err)
	}
	want = `{}`
	if string(b) != want {
		t.Fatalf("Failed want!\nWanted:\n%s\nGot:\n%s", want, string(b))
	} else {
		t.Logf("Zero value testStruct: OK!")
	}

	// Empty struct.
	b, err = Marshal(struct{}{})
	if err != nil {
		t.Fatal(err)
	}
	want = `{}`
	if string(b) != want {
		t.Fatalf("Failed want!\nWanted:\n%s\nGot:\n%s", want, string(b))
	} else {
		t.Logf("Empty struct: OK!")
	}

	// Slice of empty structs.
	b, err = Marshal([]struct{}{{}, {}, {}})
	if err != nil {
		t.Fatal(err)
	}
	want = `[{},{},{}]`
	if string(b) != want {
		t.Fatalf("Want, got:\n%s\n%s", want, string(b))
	} else {
		t.Logf("Slice of empty struct: OK!")
	}

	// Map of empty structs.
	b, err = Marshal(map[string]struct{}{"a": {}, "b": {}, "c": {}})
	if err != nil {
		t.Fatal(err)
	}
	want = `{}`
	if string(b) != want {
		t.Fatalf("Want, got:\n%s\n%s", want, string(b))
	} else {
		t.Logf("Map of empty struct: OK!")
	}
}

func TestMarshalIndent(t *testing.T) {
	b, err := MarshalIndent(withVals, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	want := `{
  "T": "1970-01-01T00:00:00Z",
  "StringTrap": "\"Time\":\"0001-01-01T00:00:00Z\"",
  "Custom": "value",
  "Nested": {
    "T": "1970-01-01T00:00:00Z",
    "StringTrap": "\"MyStruct\":null",
    "Custom": "test"
  }
}`
	if string(b) != want {
		t.Fatalf("Want, got:\n%s\n%s", want, string(b))
	}
}

func TestMarshalCustom(t *testing.T) {

	// With values
	b, err := MarshalCustom(withVals, OptionTime, OptionStruct)
	if err != nil {
		t.Fatal(err)
	}
	want := `{"T":"1970-01-01T00:00:00Z","StringTrap":"\"Time\":\"0001-01-01T00:00:00Z\"","Custom":"value","CustomEmpty":null,"Nested":{"T":"1970-01-01T00:00:00Z","StringTrap":"\"MyStruct\":null","Custom":"test","CustomEmpty":null},"NestedEmpty":{"Custom":null,"CustomEmpty":null}}`
	if string(b) != want {
		t.Fatalf("Failed want!\nWanted:\n%s\nGot:\n%s", want, string(b))
	} else {
		t.Log("With values: OK!")
	}

	// With no values
	b, err = MarshalCustom(testStruct{}, OptionTime, OptionStruct)
	if err != nil {
		t.Fatal(err)
	}
	want = `{"Custom":null,"CustomEmpty":null,"Nested":{"Custom":null,"CustomEmpty":null},"NestedEmpty":{"Custom":null,"CustomEmpty":null}}`
	if string(b) != want {
		t.Fatalf("Failed want!\nWanted:\n%s\nGot:\n%s", want, string(b))
	} else {
		t.Log("With no values: OK!")
	}

	// Map of time structs.
	b, err = MarshalCustom(map[string]struct{ T time.Time }{"a": {ts}, "b": {}, "c": {}}, OptionTime)
	if err != nil {
		t.Fatal(err)
	}
	want = `{"a":{"T":"1970-01-01T00:00:00Z"},"b":{},"c":{}}`
	if string(b) != want {
		t.Fatalf("Want, got:\n%s\n%s", want, string(b))
	} else {
		t.Logf("Map of time struct: OK!")
	}

	// Map of time structs.
	b, err = MarshalCustom(map[string]struct{ T time.Time }{"a": {ts}, "b": {}, "c": {}}, OptionTime, OptionStruct)
	if err != nil {
		t.Fatal(err)
	}
	want = `{"a":{"T":"1970-01-01T00:00:00Z"}}`
	if string(b) != want {
		t.Fatalf("Want, got:\n%s\n%s", want, string(b))
	} else {
		t.Logf("Map of time struct #2: OK!")
	}
}
