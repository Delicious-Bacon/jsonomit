# jsonomit
Package jsonomit provides JSON marshal functions that can omit empty structs
and null fields.

Provided functions can omit zero value time.Time fields, or null fields that
result from custom MarshalJSON implementations.

`go get github.com/Delicious-Bacon/jsonomit`

## Example

Working example on goplay: https://goplay.tools/snippet/Az0AigVRcQd

```go
package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/Delicious-Bacon/jsonomit"
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

func main() {
    // ============================
    // Clean all empty values.
    // ============================
    b, _ := jsonomit.MarshalIndent(
        testStruct{
            T:          time.Unix(0, 0).UTC(),
            TEmpty:     time.Time{}, // will be omitted
            StringTrap: `"Time":"0001-01-01T00:00:00Z"`,
            Custom: customStruct{
                Value: "value",
                valid: true,
            },
            CustomEmpty: customStruct{ // will be omitted
                Value: "empty",
                valid: false,
            },
            Nested: nestedStruct{
                T: time.Unix(0, 0).UTC(),
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
            NestedEmpty: nestedStruct{ // will be omitted
                TEmpty: time.Time{},
                CustomEmpty: customStruct{
                    Value: "empty",
                    valid: false,
                },
            },
        },
        "", // Prefix.
        "    ", // Indent.
    )

    fmt.Println(string(b))
    // Output:
    // {
    //     "T": "1970-01-01T00:00:00Z",
    //     "StringTrap": "\"Time\":\"0001-01-01T00:00:00Z\"",
    //     "Custom": "value",
    //     "Nested": {
    //       "T": "1970-01-01T00:00:00Z",
    //       "StringTrap": "\"MyStruct\":null",
    //       "Custom": "test"
    //     }
    // }

    // ============================
    // Customized marshal cleaning.
    // ============================
    b, _ = jsonomit.MarshalCustomIndent(
        map[string]struct{
            T time.Time
        }{
            "a": {time.Unix(0, 0).UTC()},
            "b": {},
            "c": {},
        },
        "", // Prefix.
        "    ", // Indent.
        jsonomit.OptionTime, // Clean zero time.Time structs.
        // jsonomit.OptionNull, -> Would clean null fields.
        // jsonomit.OptionStruct, -> Would remove keys with empty structs.
    )

    fmt.Println(string(b))
    // Output:
    // {
    //     "a": {
    //       "T": "1970-01-01T00:00:00Z"
    //     },
    //     "b": {},
    //     "c": {}
    // }
}
```
