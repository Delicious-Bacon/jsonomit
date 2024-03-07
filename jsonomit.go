// Package jsonomit provides JSON marshal functions that can omit empty structs
// and null fields.
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
	zeroNumRGX     = regexp.MustCompile(`"\w+":0(?:,|("|{|}|\[|\]))`)

	cleanupRgxs = []*regexp.Regexp{
		emptyTimeRGX,
		nullFieldRGX,
		zeroNumRGX,
	}

	opt_Rgx = map[option]*regexp.Regexp{
		OptionTime:    emptyTimeRGX,
		OptionNull:    nullFieldRGX,
		OptionZeroNum: zeroNumRGX,
	}

	// Cleans up empty time.Time fields.
	/* "field":"0001-01-01T00:00:00Z" */
	OptionTime = option{1}
	// Cleans up null fields.
	/* "field":null */
	OptionNull = option{2}
	// Cleans up empty structs.
	/* "field":{} */
	OptionStruct = option{3}
	// Cleans up zero number fields. Useful when external package
	// does not omitempty on zero number fields.
	/* "field":0 */
	OptionZeroNum = option{4}
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

		b = rgx.ReplaceAll(b, []byte("$1"))
	}
	b = bytes.Replace(b, []byte(`,}`), []byte(`}`), -1)

	// Clean the JSON from empty structs.
	for emptyStructRGX.Match(b) {

		b = emptyStructRGX.ReplaceAll(b, []byte(""))
		b = bytes.Replace(b, []byte(`,}`), []byte(`}`), -1)
	}

	return b, nil
}

// MarshalIndent is like Marshal but applies Indent to format the output.
// Each JSON element in the output will begin on a new line beginning with prefix
// followed by one or more copies of indent according to the indentation nesting.
func MarshalIndent(v any, prefix, indent string) ([]byte, error) {
	b, err := Marshal(v)
	if err != nil {
		return nil, err
	}

	b2 := bytes.NewBuffer([]byte{})
	err = json.Indent(b2, b, prefix, indent)
	if err != nil {
		return nil, err
	}
	return b2.Bytes(), nil
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

	count := len(opts)
	var cleanStructs bool
	for _, opt := range opts {
		if opt == OptionStruct {
			cleanStructs = true
			count--
			continue
		}

		b = opt_Rgx[opt].ReplaceAll(b, []byte("$1"))
	}
	if count > 0 {
		b = bytes.Replace(b, []byte(`,}`), []byte(`}`), -1)
	}

	if cleanStructs {
		// Clean the JSON from empty structs.
		for emptyStructRGX.Match(b) {

			b = emptyStructRGX.ReplaceAll(b, []byte(""))

			b = bytes.Replace(b, []byte(`,}`), []byte(`}`), -1)
		}
	}

	return b, nil
}

// MarshalCustomIndent is like MarshalCustom but applies Indent to format the output.
// Each JSON element in the output will begin on a new line beginning with prefix
// followed by one or more copies of indent according to the indentation nesting.
func MarshalCustomIndent(v any, prefix, indent string, opts ...option) ([]byte, error) {
	b, err := MarshalCustom(v, opts...)
	if err != nil {
		return nil, err
	}
	b2 := bytes.NewBuffer([]byte{})
	err = json.Indent(b2, b, prefix, indent)
	if err != nil {
		return nil, err
	}
	return b2.Bytes(), nil
}
