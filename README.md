# jsonomit
Package jsonomit provides JSON marshal functions to omit empty structs
and null fields.

Provided functions can omit zero value time.Time fields, or null fields that
result from custom MarshalJSON implementations.

`go get github.com/Delicious-Bacon/jsonomit`

## Example

```go
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
		TEmpty:     time.Time{}, // will be omitted
		StringTrap: `"Time":"0001-01-01T00:00:00Z"`, // won't be omitted
		Custom: customStruct{
			Value: "value",
			valid: true,
		},
		CustomEmpty: customStruct{
			Value: "empty",
			valid: false, // will be omitted
		},
		Nested: nestedStruct{
			T:          ts,
			TEmpty:     time.Time{},
			StringTrap: `"MyStruct":null`, // won't be omitted
			Custom: customStruct{
				Value: "test",
				valid: true,
			},
			CustomEmpty: customStruct{
				Value: "empty",
				valid: false,
			},
		},
		NestedEmpty: nestedStruct{ // will be omitted
			TEmpty: time.Time{},
			CustomEmpty: customStruct{
				Value: "empty",
				valid: false,
			},
		},
	}
)

func main() {
	// Populated testStruct
	b, _ := Marshal(withVals)

	// Output: {"T":"1970-01-01T00:00:00Z","StringTrap":"\"Time\":\"0001-01-01T00:00:00Z\"","Custom":"value","Nested":{"T":"1970-01-01T00:00:00Z","StringTrap":"\"MyStruct\":null","Custom":"test"}}

    // Customized marshal cleaning
	b, _ = MarshalCustom(
        map[string]struct{
            T time.Time
        }{
            "a": {ts},
            "b": {},
            "c": {},
        },
        OptionTime, // Clean zero time.Time structs.
        // OptionStruct, -> would remove keys with empty structs.
    )

    // Preserves map keys:
	// Output: {"a":{"T":"1970-01-01T00:00:00Z"},"b":{},"c":{}}
}
