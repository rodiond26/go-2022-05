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

	SimpleTestStruct struct {
		Value string
	}

	SliceTestStruct struct {
		Value []string
	}

	LenTestStruct struct {
		Value string `validate:"len:6"`
	}

	RegExpTestStruct struct {
		Value string `validate:"regexp:^\\d+$"`
	}

	StringInTagTestStruct struct {
		Value string `validate:"in:foo,bar"`
	}

	MinIntTagTestStruct struct {
		Value int `validate:"min:18"`
	}

	MaxIntTagTestStruct struct {
		Value int `validate:"max:10"`
	}

	IntInTagTestStruct struct {
		Value int `validate:"in:200,404,500"`
	}
	MinMaxIntTagTestStruct struct {
		Value int `validate:"min:18|max:50"`
	}

	NotValidTagTestStruct struct {
		StringValue string `validate:"default:common"`
		IntValue    int    `validate:"default:10"`
	}
)

func TestValidateWithoutErrors(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			in:          nil,
			expectedErr: nil,
		},
		{
			in: SimpleTestStruct{
				Value: "qwerty",
			},
			expectedErr: nil,
		},
		{
			in: SliceTestStruct{
				Value: []string{"foo", "bar"},
			},
			expectedErr: nil,
		},
		{
			in: LenTestStruct{
				Value: "qwerty",
			},
			expectedErr: nil,
		},
		{
			in: RegExpTestStruct{
				Value: "1234",
			},
			expectedErr: nil,
		},
		{
			in: StringInTagTestStruct{
				Value: "foo",
			},
			expectedErr: nil,
		},
		{
			in: MinIntTagTestStruct{
				Value: 100,
			},
			expectedErr: nil,
		},
		{
			in: MaxIntTagTestStruct{
				Value: 0,
			},
			expectedErr: nil,
		},
		{
			in: IntInTagTestStruct{
				Value: 200,
			},
			expectedErr: nil,
		},
		{
			in: User{
				ID:     "012345678901234567890123456789012345",
				Name:   "Name",
				Age:    33,
				Email:  "mail@email.com",
				Role:   UserRole("admin"),
				Phones: []string{"12345678901", "01234567890"},
				meta:   nil,
			},
			expectedErr: nil,
		},
		{
			in: App{
				Version: "12345",
			},
			expectedErr: nil,
		},
		{
			in: Token{
				Header:    []byte("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9"),
				Payload:   []byte("eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ"),
				Signature: []byte("Eg0sUd_l7s_PW_4b0iRoA-tqPkkM6NgdZt7mzXj2fmM"),
			},
			expectedErr: nil,
		},
		{
			in: Response{
				Code: 404,
				Body: "",
			},
			expectedErr: nil,
		},
	}

	for i, testCase := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			t.Parallel()
			actualErrors := Validate(testCase.in)
			if testCase.expectedErr == nil {
				require.NoError(t, actualErrors)
			} else {
				var validationErrors ValidationErrors
				if errors.As(actualErrors, &validationErrors) {
					var expectedErrors ValidationErrors
					require.ErrorAs(t, testCase.expectedErr, &expectedErrors)
					for j, err := range validationErrors {
						require.ErrorIs(t, err, validationErrors[j])
					}
				} else {
					require.ErrorIs(t, actualErrors, testCase.expectedErr)
				}
			}
		})
	}
}

func TestValidateWithOneError(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			in:          "String is not a struct",
			expectedErr: ErrObjectIsNotStruct,
		},
		{
			in: LenTestStruct{
				Value: "qwerty1",
			},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "Value",
					Err:   ErrInvalidStringLength,
				},
			},
		},
		{
			in: RegExpTestStruct{
				Value: "1234a",
			},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "Value",
					Err:   ErrInvalidStringLength,
				},
			},
		},
		{
			in: StringInTagTestStruct{
				Value: "notFoo",
			},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "Value",
					Err:   ErrValueIsNotInSet,
				},
			},
		},
		{
			in: MinIntTagTestStruct{
				Value: 0,
			},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "Value",
					Err:   ErrValueIsLess,
				},
			},
		},
		{
			in: MaxIntTagTestStruct{
				Value: 100,
			},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "Value",
					Err:   ErrValueIsGreater,
				},
			},
		},
		{
			in: IntInTagTestStruct{
				Value: 100,
			},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "Value",
					Err:   ErrValueIsNotInSet,
				},
			},
		},
		{
			in: MinMaxIntTagTestStruct{
				Value: 55,
			},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "Value",
					Err:   ErrValueIsGreater,
				},
			},
		},
		{
			in: MinMaxIntTagTestStruct{
				Value: 3,
			},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "Value",
					Err:   ErrValueIsLess,
				},
			},
		},
	}

	for i, testCase := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			t.Parallel()
			actualErrors := Validate(testCase.in)
			if testCase.expectedErr == nil {
				require.NoError(t, actualErrors)
			} else {
				var validationErrors ValidationErrors
				if errors.As(actualErrors, &validationErrors) {
					var expectedErrors ValidationErrors
					require.ErrorAs(t, testCase.expectedErr, &expectedErrors)
					for j, err := range validationErrors {
						require.ErrorIs(t, err, validationErrors[j])
					}
				} else {
					require.ErrorIs(t, actualErrors, testCase.expectedErr)
				}
			}
		})
	}
}

func TestValidateSeveralErrors(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			in: User{
				ID:     "0",
				Name:   "Name",
				Age:    10,
				Email:  "mail@email.com.1",
				Role:   UserRole("not admin"),
				Phones: []string{"1234", "0123"},
				meta:   nil,
			},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "ID",
					Err:   ErrInvalidStringLength,
				},
				ValidationError{
					Field: "Age",
					Err:   ErrValueIsLess,
				},
				ValidationError{
					Field: "Email",
					Err:   ErrStringNotMatchRegexp,
				},
				ValidationError{
					Field: "Role",
					Err:   ErrValueIsNotInSet,
				},
				ValidationError{
					Field: "Phones",
					Err:   ErrInvalidStringLength,
				},
				ValidationError{
					Field: "Phones",
					Err:   ErrInvalidStringLength,
				},
			},
		},
		{
			in: NotValidTagTestStruct{
				StringValue: "something",
				IntValue:    0,
			},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "StringValue",
					Err:   ErrBadValidateTag,
				},
				ValidationError{
					Field: "IntValue",
					Err:   ErrBadValidateTag,
				},
			},
		},
	}

	for i, testCase := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			t.Parallel()
			actualErrors := Validate(testCase.in)
			if testCase.expectedErr == nil {
				require.NoError(t, actualErrors)
			} else {
				var validationErrors ValidationErrors
				if errors.As(actualErrors, &validationErrors) {
					var expectedErrors ValidationErrors
					require.ErrorAs(t, testCase.expectedErr, &expectedErrors)
					for j, err := range validationErrors {
						require.ErrorIs(t, err, validationErrors[j])
					}
				} else {
					require.ErrorIs(t, actualErrors, testCase.expectedErr)
				}
			}
		})
	}
}
