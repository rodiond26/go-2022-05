package hw09structvalidator

import (
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
	if v == nil {
		return nil
	}
	vType := reflect.TypeOf(v)          // отражение типа
	if vType.Kind() != reflect.Struct { // тип данных
		return nil // TODO fix
	}

	for i := 0; i < vType.NumField(); i++ {
		field := vType.Field(i) // поле структуры по индексу
		if field.Tag == "" {
			// TODO ?
			continue
		}
		tag := field.Tag.Get(validateTag)
		if tag == "" {
			// TODO ?
			continue
		}

		value := reflect.ValueOf(v).Field(i)
		if field.Type.Kind() == reflect.String {
			err := validateString(field, value, tag)
			_ = err // TODO fix
		}

	}

	return nil // TODO fix
}

func validateString(field reflect.StructField, value reflect.Value, tag string) error {
	f := field.Name
	v := value.String()
	_ = f      // TODO fix
	_ = v      // TODO fix
	return nil // TODO fix
}
