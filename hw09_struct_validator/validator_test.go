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

	TestTag struct {
		CheckLenValue      string `validate:"len:10"`
		CheckRegExpValue   string `validate:"regexp:^\\d+$"`
		CheckEmailValue    string `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		CheckStringInValue string `validate:"in:foo,bar"`
		CheckMinValue      int    `validate:"min:18"`
		CheckMaxValue      int    `validate:"max:50"`
		CheckIntInValue    int    `validate:"in:100,500"`
	}
)

func TestValidate(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			in:          nil,
			expectedErr: nil,
		},
		{
			in:          "String is not a struct",
			expectedErr: ErrObjectIsNotStruct,
		},
		{
			in: SimpleTestStruct{
				Value: "qwerty",
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
			in: LenTestStruct{
				Value: "qwerty1",
			},
			expectedErr: ErrInvalidStringLength,
		},
		{
			in: RegExpTestStruct{
				Value: "1234a",
			},
			expectedErr: ErrStringNotMatchRegexp,
		},
		{
			in: StringInTagTestStruct{
				Value: "foo",
			},
			expectedErr: nil,
		},
		{
			in: TestTag{
				CheckLenValue:      "1234567890",
				CheckRegExpValue:   "1234567890",
				CheckEmailValue:    "mail@mail.com",
				CheckStringInValue: "foo",
				CheckMinValue:      20,
				CheckMaxValue:      20,
				CheckIntInValue:    100,
			},
			expectedErr: nil,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			t.Parallel()
			actualErrors := Validate(tt.in)

			if tt.expectedErr == nil {
				require.NoError(t, actualErrors)
			} else {
				require.NotEmpty(t, actualErrors)
				var validationErrors ValidationErrors

				if errors.As(actualErrors, &validationErrors) {
					var expectedErrors ValidationErrors
					require.ErrorAs(t, tt.expectedErr, &expectedErrors)

					for j, err := range validationErrors {
						require.ErrorIs(t, err, validationErrors[j])
					}
				}
				require.ErrorIs(t, actualErrors, tt.expectedErr)
			}
		})
	}
}
