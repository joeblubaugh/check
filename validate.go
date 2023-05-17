package validate

import (
	"reflect"
	"strings"
)

type Data interface {
	IsValid() error
}

type MultiError interface {
	Errors() []error
}

type ValidationError struct {
	key    string
	errors []error
}

func (v *ValidationError) Errors() []error {
	return v.errors
}

func (v *ValidationError) Add(m ...error) *ValidationError {
	v.errors = append(v.errors, m...)
	return v
}

func Validate(key string, d Data) error {
	ve := &ValidationError{}
	err := d.IsValid()
	if err == nil {
		return nil
	}

	if e, ok := err.(*ValidationError); ok {
		ve = e
	} else {
		ve.Add(err)
	}

	ve.key = key
	if len(key) == 0 {
		v := reflect.ValueOf(d)
		t := v.Type()
		for t.Kind() == reflect.Pointer {
			t = t.Elem()
		}
		n := t.String()
		ve.key = n
	}

	return ve
}

func (v *ValidationError) Error() string {
	// Iterate over errors. Prints the key, and if there is only one message, adds message to the same line.
	// Otherwise, indents message errors by \t
	str := v.key + ":\n"
	for _, m := range v.errors {
		// Not quite right. Need to pass indentation level into a nested function, avoid "\n" where not needed.
		str += align("\t", m.Error())
		str += "\n"
	}

	return str
}

// Prepends "indent" to every line of string
func align(indent string, source string) string {
	reader := strings.NewReader(source)
	builder := &strings.Builder{}
	builder.WriteString(indent)
	reset := false
	for reader.Len() > 0 {
		if reset {
			builder.WriteString(indent)
			reset = false
		}
		rune, _, _ := reader.ReadRune()
		builder.WriteRune(rune)
		if rune == '\n' {
			reset = true
		}
	}

	return builder.String()
}
