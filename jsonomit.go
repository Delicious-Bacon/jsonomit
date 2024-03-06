// Package jsonomit provides JSON marshal functions to omit empty structs
// and null fields. By default, the functions omit empty structs.
//
// Provided functions can omit zero value time.Time fields, or null fields that
// result from custom MarshalJSON implementations.
package jsonomit

import (
	"bytes"
	"encoding/json"
	"regexp"
)

var (
	emptyTimeRGX   = regexp.MustCompile(`"\w+":"0001-01-01T00:00:00Z",?`)
	nullFieldRGX   = regexp.MustCompile(`"\w+":null,?`)
	emptyStructRGX = regexp.MustCompile(`"\w+":{},?`)

	cleanupRgxs = []*regexp.Regexp{
		emptyTimeRGX,
		nullFieldRGX,
	}

	opt_Rgx = map[option]*regexp.Regexp{
		OptionTime: emptyTimeRGX,
		OptionNull: nullFieldRGX,
	}

	// Cleans up empty time.Time fields.
	/* "field":"0001-01-01T00:00:00Z" */
	OptionTime = option{1}
	// Cleans up null fields.
	/* "field":null */
	OptionNull = option{2}
)

// option is used to specify which fields to clean up from the JSON
// encoding output.
type option struct {
	int
}

// Marshal returns the JSON encoding of v clean of empty values
// (zero time, null fields and empty structs).
//
// Reference the standard json package for JSON encoding information:
// https://pkg.go.dev/encoding/json.
func Marshal(v any) ([]byte, error) {
	// Do the standard JSON marshal.
	b, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}

	// Clean the JSON from empty values.
	for _, rgx := range cleanupRgxs {

		b = rgx.ReplaceAll(b, []byte(""))
	}
	b = bytes.Replace(b, []byte(`,}`), []byte(`}`), -1)

	// Clean the JSON from empty structs.
	for emptyStructRGX.Match(b) {

		b = emptyStructRGX.ReplaceAll(b, []byte(""))
		b = bytes.Replace(b, []byte(`,}`), []byte(`}`), -1)
	}

	return b, nil
}

// MarshalCustom returns the JSON encoding of v clean of
// empty values for the given options.
//
// This version can be faster than Marshal if you don't
// need to clean all the empty values.
//
// Reference the standard json package for JSON encoding information:
// https://pkg.go.dev/encoding/json.
func MarshalCustom(v any, opts ...option) ([]byte, error) {
	// Do the standard JSON marshal.
	b, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}

	// Clean the JSON from empty values of the given options.
	for _, opt := range opts {

		b = opt_Rgx[opt].ReplaceAll(b, []byte(""))
	}
	b = bytes.Replace(b, []byte(`,}`), []byte(`}`), -1)

	// Clean the JSON from empty structs.
	for emptyStructRGX.Match(b) {

		b = emptyStructRGX.ReplaceAll(b, []byte(""))

		b = bytes.Replace(b, []byte(`,}`), []byte(`}`), -1)
	}

	return b, nil
}
