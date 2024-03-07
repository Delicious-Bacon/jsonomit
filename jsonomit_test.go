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
	Num         float64
	NumEmpty    float64
	Custom      customStruct
	CustomEmpty customStruct
	StructEmpty struct{}
}

type testStruct struct {
	T           time.Time
	TEmpty      time.Time
	StringTrap  string `json:",omitempty"`
	Num         float64
	NumEmpty    float64 // missing omitempty
	Custom      customStruct
	CustomEmpty customStruct
	Nested      nestedStruct
	NestedEmpty nestedStruct
	StructEmpty struct{}
}

type icon struct {
	Pixels  []uint16 // missing omitempty
	ImgSize point    // missing omitempty
}
type point struct {
	X, Y int // missing omitempty
}

var (
	ts = time.Unix(0, 0).UTC()

	allTogether = testStruct{
		T:          ts,
		TEmpty:     time.Time{},
		StringTrap: `"Time":"0001-01-01T00:00:00Z"`,
		Num:        0.1,
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
				Value: "value",
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

	// Time
	b, err := Marshal(
		struct {
			T time.Time
		}{},
	)
	if err != nil {
		t.Fatal(err)
	}
	want := `{}`
	if string(b) != want {
		t.Fatalf("Failed 'time' want!\nWanted:\n%s\nGot:\n%s", want, string(b))
	} else {
		t.Log("Clean time: OK!")
	}

	// Custom marshal
	b, err = Marshal(customStruct{Value: "value", valid: false})
	if err != nil {
		t.Fatal(err)
	}
	want = `null`
	if string(b) != want {
		t.Fatalf("Failed 'custom marshal' want!\nWanted:\n%s\nGot:\n%s", want, string(b))
	} else {
		t.Logf("Custom marshal: OK!")
	}

	// Custom in a struct
	b, err = Marshal(struct{ Custom customStruct }{Custom: customStruct{Value: "value", valid: false}})
	if err != nil {
		t.Fatal(err)
	}
	want = `{}`
	if string(b) != want {
		t.Fatalf("Failed 'custom marshal in a struct' want!\nWanted:\n%s\nGot:\n%s", want, string(b))
	} else {
		t.Logf("Custom marshal: OK!")
	}

	// Nested empty.
	b, err = Marshal(struct{ _ struct{ _ struct{} } }{struct{ _ struct{} }{struct{}{}}})
	if err != nil {
		t.Fatal(err)
	}
	want = `{}`
	if string(b) != want {
		t.Fatalf("Failed 'nested empty' want!\nWanted:\n%s\nGot:\n%s", want, string(b))
	} else {
		t.Logf("Nested empty: OK!")
	}

	// Slice of empty structs.
	b, err = Marshal([]struct{}{{}, {}, {}})
	if err != nil {
		t.Fatal(err)
	}
	want = `[{},{},{}]`
	if string(b) != want {
		t.Fatalf("Failed 'slice of empty structs' want!\n%s\n%s", want, string(b))
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
		t.Fatalf("Failed 'map of empty structs' want!\n%s\n%s", want, string(b))
	} else {
		t.Logf("Map of empty struct: OK!")
	}

	// All together
	b, err = Marshal(allTogether)
	if err != nil {
		t.Fatal(err)
	}
	want = `{"T":"1970-01-01T00:00:00Z","StringTrap":"\"Time\":\"0001-01-01T00:00:00Z\"","Num":0.1,"Custom":"value","Nested":{"T":"1970-01-01T00:00:00Z","StringTrap":"\"MyStruct\":null","Custom":"value"}}`
	if string(b) != want {
		t.Fatalf("Failed 'all together' want!\n%s\n%s", want, string(b))
	} else {
		t.Logf("All together: OK!")
	}

	b, err = Marshal(icon{})
	if err != nil {
		t.Fatal(err)
	}

	want = `{}`
	if string(b) != want {
		t.Fatalf("Failed 'all together' want!\n%s\n%s", want, string(b))
	} else {
		t.Logf("All together: OK!")
	}
}

func TestMarshalIndent(t *testing.T) {
	b, err := MarshalIndent(allTogether, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	want := `{
  "T": "1970-01-01T00:00:00Z",
  "StringTrap": "\"Time\":\"0001-01-01T00:00:00Z\"",
  "Num": 0.1,
  "Custom": "value",
  "Nested": {
    "T": "1970-01-01T00:00:00Z",
    "StringTrap": "\"MyStruct\":null",
    "Custom": "value"
  }
}`
	if string(b) != want {
		t.Fatalf("Failed 'all together' want!\n%s\n%s", want, string(b))
	}
}

func TestMarshalCustom(t *testing.T) {

	// With values
	b, err := MarshalCustom(allTogether, OptionTime, OptionZeroNum, OptionStruct)
	if err != nil {
		t.Fatal(err)
	}
	want := `{"T":"1970-01-01T00:00:00Z","StringTrap":"\"Time\":\"0001-01-01T00:00:00Z\"","Num":0.1,"Custom":"value","CustomEmpty":null,"Nested":{"T":"1970-01-01T00:00:00Z","StringTrap":"\"MyStruct\":null","Custom":"value","CustomEmpty":null},"NestedEmpty":{"Custom":null,"CustomEmpty":null}}`
	if string(b) != want {
		t.Fatalf("Failed 'all together' want!\nWanted:\n%s\nGot:\n%s", want, string(b))
	} else {
		t.Log("With values: OK!")
	}

	// Empty testStruct
	b, err = MarshalCustom(testStruct{}, OptionTime, OptionZeroNum, OptionStruct)
	if err != nil {
		t.Fatal(err)
	}
	want = `{"Custom":null,"CustomEmpty":null,"Nested":{"Custom":null,"CustomEmpty":null},"NestedEmpty":{"Custom":null,"CustomEmpty":null}}`
	if string(b) != want {
		t.Fatalf("Failed 'empty testStruct' want!\nWanted:\n%s\nGot:\n%s", want, string(b))
	} else {
		t.Log("Empty testStruct: OK!")
	}

	// Map of time structs.
	b, err = MarshalCustom(map[string]struct{ T time.Time }{"a": {ts}, "b": {}, "c": {}}, OptionTime)
	if err != nil {
		t.Fatal(err)
	}
	want = `{"a":{"T":"1970-01-01T00:00:00Z"},"b":{},"c":{}}`
	if string(b) != want {
		t.Fatalf("Failed 'map of time structs' want!\n%s\n%s", want, string(b))
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
		t.Fatalf("Failed 'map of time structs #2' want!\n%s\n%s", want, string(b))
	} else {
		t.Logf("Map of time struct #2: OK!")
	}
}
