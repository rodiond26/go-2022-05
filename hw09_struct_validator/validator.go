package hw09structvalidator

import (
	"fmt"
	"reflect"
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

const (
	validateTag = "validate"
)

func (v ValidationErrors) Error() string {
	panic("implement me")
}

func Validate(v interface{}) error {
	s := reflect.ValueOf(v)
	if s.IsNil() {
		return nil // TODO fix
	}
	if s.Kind() != reflect.Struct {
		return nil // TODO fix
	}

	vType := s.Type()
	for i := 0; i < vType.NumField(); i++ {
		field := vType.Field(i)
		tag := field.Tag.Get(validateTag)
		if tag == "" {
			// TODO ?
			continue
		}
		err := validateField(field.Name, validateTag, vType.Field(i))
		fmt.Println(err) // TODO delete
	}

	return nil // TODO fix
}

func validateField(fieldName, validateTag string, field reflect.StructField) error {
	return nil // TODO fix

}

type App1 struct {
	Version string `validate:"len:5"`
}

func main() {
	app1 := App1{Version: "123"}
	err := Validate(app1)
	fmt.Println(err)
}
