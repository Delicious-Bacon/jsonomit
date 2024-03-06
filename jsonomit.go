// jsonomitpackage provides JSON marshal funcs to omit empty (time)
// structs and null fields from MarshalJSON custom implementations
package jsonomit

import (
	"bytes"
	"encoding/json"
	"regexp"
)

var (
	emptyTimeRGX   = regexp.MustCompile(`"0001-01-01T00:00:00Z",?`)
	nullFieldRGX   = regexp.MustCompile(`:null,?`)
	emptyStructRGX = regexp.MustCompile(`{},?`)

	cleanupRgxs = []*regexp.Regexp{
		emptyTimeRGX,
		nullFieldRGX,
		emptyStructRGX,
	}
)

// Marshal returns the JSON encoding of v clean of empty values
// (zero time, null fields and empty structs).
// Reference the standard json package for JSON encoding information:
// https://pkg.go.dev/encoding/json.
func Marshal(v any) ([]byte, error) {
	// Do the standard JSON marshal.
	b, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}

	// Cleanup the JSON from empty structs.
	for _, rgx := range cleanupRgxs {

		if matches := rgx.FindAllIndex(b, -1); len(matches) > 0 {

			for i := len(matches) - 1; i >= 0; i-- {
				var c int
				for j := matches[i][0] - 1; j >= 0; j-- {
					if b[j] == '"' {
						c++
						if c == 2 {
							b = append(b[:j], b[matches[i][1]:]...)
						}
					}
				}
			}
			b = bytes.Replace(b, []byte(`,}`), []byte(`}`), -1)
		}
	}

	return b, nil
}
