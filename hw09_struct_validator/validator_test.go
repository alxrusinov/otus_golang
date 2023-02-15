package hw09structvalidator

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int      `validate:"min:18|max:50"`
		Email  string   `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole `validate:"in:admin,stuff"`
		Phones []string `validate:"len:11"`
		meta   json.RawMessage
	}

	App struct {
		Version string `validate:"len:5"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}
)

var tests = []struct {
	in             interface{}
	errProgramming error
	errValidation  ValidationErrors
}{
	{
		in:             12,
		errProgramming: ProgrammingError{Msg: getTypeErrorMsg()},
		errValidation:  nil,
	},
	{
		in:             nil,
		errProgramming: ProgrammingError{Msg: getTypeErrorMsg()},
		errValidation:  nil,
	},
	{
		in: struct {
			Foo string `validate:"len-10"`
		}{
			Foo: "clown",
		},
		errProgramming: ProgrammingError{Msg: getInvalidValidatorErrorMsg()},
		errValidation:  nil,
	},
	{
		in: struct {
			Foo string `validate:"regexp:+++"`
		}{
			Foo: "foo",
		},
		errProgramming: ProgrammingError{Msg: getInvalidValidatorErrorMsg()},
		errValidation:  nil,
	},
	{
		in: struct {
			Foo string `validate:"len:len"`
		}{
			Foo: "foo",
		},
		errProgramming: ProgrammingError{Msg: getInvalidValidatorErrorMsg()},
		errValidation:  nil,
	},
	{
		in: struct {
			Foo int `validate:"min:min"`
		}{
			Foo: 1,
		},
		errProgramming: ProgrammingError{Msg: getInvalidValidatorErrorMsg()},
		errValidation:  nil,
	},
	{
		in: struct {
			Foo int `validate:"max:min"`
		}{
			Foo: 1,
		},
		errProgramming: ProgrammingError{Msg: getInvalidValidatorErrorMsg()},
		errValidation:  nil,
	},
	{
		in: struct {
			Foo int `validate:"in:12,foo"`
		}{
			Foo: 1,
		},
		errProgramming: ProgrammingError{Msg: getInvalidValidatorErrorMsg()},
		errValidation:  nil,
	},
	{
		in: User{
			ID:     "123456789123456789123456789123456789",
			Name:   "qwerty",
			Age:    21,
			Email:  "foo@bar.baz",
			Role:   "admin",
			Phones: []string{"88002223311", "88002223311"},
			meta:   []byte("123"),
		},
		errProgramming: nil,
		errValidation:  nil,
	},
	{
		in: User{
			ID:     "123456789123456789123456789123456789",
			Name:   "qwerty",
			Age:    60,
			Email:  "foo@bar.baz",
			Role:   "admin",
			Phones: []string{"88002223311", "88002223311"},
			meta:   []byte("123"),
		},
		errProgramming: nil,
		errValidation: ValidationErrors{
			ValidationError{Field: "Age", Err: errors.New(getValidationErrorMsg())},
		},
	},
	{
		in: User{
			ID:     "12",
			Name:   "qwerty",
			Age:    17,
			Email:  "foo@bar.baz.qwqwq@.",
			Role:   "clown",
			Phones: []string{"1", "88002223311"},
			meta:   []byte("123"),
		},
		errProgramming: nil,
		errValidation: ValidationErrors{
			ValidationError{Field: "ID", Err: errors.New(getValidationErrorMsg())},
			ValidationError{Field: "Age", Err: errors.New(getValidationErrorMsg())},
			ValidationError{Field: "Email", Err: errors.New(getValidationErrorMsg())},
			ValidationError{Field: "Role", Err: errors.New(getValidationErrorMsg())},
			ValidationError{Field: "Phones[0]", Err: errors.New(getValidationErrorMsg())},
		},
	},
	{
		in:             App{Version: "ww"},
		errProgramming: nil,
		errValidation:  ValidationErrors{ValidationError{Field: "Version", Err: errors.New(getValidationErrorMsg())}},
	},
	{
		in: Token{
			Header:    []byte("123"),
			Payload:   []byte("123"),
			Signature: []byte("123"),
		},
		errProgramming: nil,
		errValidation:  nil,
	},
	{
		in: Response{
			Code: 200,
			Body: "123",
		},
		errProgramming: nil,
		errValidation:  nil,
	},
	{
		in: Response{
			Code: 400,
			Body: "123",
		},
		errProgramming: nil,
		errValidation:  ValidationErrors{ValidationError{Field: "Code", Err: errors.New(getValidationErrorMsg())}},
	},
}

func TestValidate(t *testing.T) {
	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			err := Validate(tt.in)

			if tt.errValidation == nil && tt.errProgramming == nil {
				require.Nil(t, err)
				return
			}

			if tt.errProgramming != nil {
				require.ErrorAs(t, err, &tt.errProgramming)
				return
			}

			require.ErrorAs(t, err, &tt.errValidation)

			// err type-casting for iterating on slice of errors
			//nolint:errorlint
			for i, errValidation := range err.(ValidationErrors) {
				require.Equal(t, tt.errValidation[i].Field, errValidation.Field)
				require.ErrorIs(t, tt.errValidation[i].Err, errValidation.Err)
			}
		})
	}
}
