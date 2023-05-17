package validate

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type simple struct {
	value int
}

func (s *simple) IsValid() error {
	if s.value > 5 {
		return (&ValidationError{}).Add(errors.New("Value must be 5 or less"))
	}
	return nil
}

// Examples of use of the validate package
func ExampleValidate_simple() {
	d := &simple{
		value: 10,
	}

	fmt.Println(Validate("simple type", d))
	// Output:
	//simple type:
	//	Value must be 5 or less
}

type Complex struct {
	value int
	child *Complex
}

func (c *Complex) IsValid() error {
	var ve *ValidationError
	add := func(err error) {
		if ve == nil {
			ve = &ValidationError{}
		}
		ve.Add(err)
	}

	if c.value > 5 {
		add(errors.New("Value must be 5 or less"))
	}

	if c.child != nil {
		err := Validate("child", c.child)
		if err != nil {
			add(err)
		}
	}

	// Without this line, the example crashes. WHY?????
	if ve == nil {
		return nil
	}

	return ve
}

func TestErrorWeirdness(t *testing.T) {
	var err error
	get := func() error {
		var e *ValidationError = nil
		return e
	}

	err = get()
	if u, ok := err.(*ValidationError); ok {
		assert.Nil(t, u)
	}

}

func ExampleValidate_complex() {
	d := &Complex{
		value: 1,
		child: &Complex{
			value: 10,
			child: &Complex{
				value: 5,
			},
		},
	}

	fmt.Println(Validate("complex type", d))
	// Output:
	//complex type:
	//	child:
	//		Value must be 5 or less
	//		child:
	//			Value must be 5 or less
}
