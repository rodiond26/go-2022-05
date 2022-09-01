package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	panic("")
}

func (v ValidationError) Error() string {
	return fmt.Sprintf("Field [%s] error [%s]", v.Field, v.Error())
}

func Validate(v interface{}) error {
	if v == nil {
		return nil
	}
	typeOfV := reflect.TypeOf(v)
	valueOfV := reflect.ValueOf(v)
	validationErrors := make(ValidationErrors, 0)

	if typeOfV.Kind() != reflect.Struct {
		return ValidationError{Field: "", Err: errors.New("the object is not a struct")}
	}

	for i := 0; i < typeOfV.NumField(); i++ {
		field := typeOfV.Field(i)
		if field.Tag == "" {
			continue
		}
		tag := field.Tag.Get(validateTag)
		if tag == "" || tag == "-" {
			continue
		}

		tags := toStringSlice(tag, tagSeparator)
		if field.Type.Kind() == reflect.String {
			err := validateString(field.Name, valueOfV.Field(i).String(), tags)
			if err != nil {
				validationErrors = append(validationErrors, err)
			}
			continue
		}
		if field.Type.Kind() == reflect.Int {
			err := validateInt64(field.Name, valueOfV.Field(i).Int(), tags)
			if err != nil {
				validationErrors = append(validationErrors, err)
			}
			continue
		}
		if field.Type.Kind() == reflect.Slice {
			err := validateSlice(field, valueOfV.Field(i), tags)
			if err != nil {
				validationErrors = append(validationErrors, err)
			}
		}
	}
	return nil
}

type TagType int

const (
	NotValid TagType = iota
	Length
	In
	RegExp
	Min
	Max
)

const (
	validateTag  = "validate"
	tagBorder    = ":"
	tagSeparator = "|"
	inSeparator  = ","
	lengthTag    = "len"
	inTag        = "in"
	regexpTag    = "regexp"
	minTag       = "min"
	maxTag       = "max"
	empty        = ""
)

func getTypeAndValueOfTag(tag string) (TagType, string) {
	prefix := before(tag, tagBorder)
	switch prefix {
	case lengthTag:
		return Length, after(tag, tagBorder)
	}
}

func validateString(fieldName, fieldValue string, tags []string) ValidationError {
	for _, tag := range tags {
		tagType, tagValue := getTypeAndValueOfTag(tag)
		switch tagType {
		case NotValid:
			return ValidationError{Field: fieldName, Err: fmt.Errorf("cannot validate value [%s] using tag [%s]", fieldName, tag)}
		case Length:
			validLength, err := strconv.Atoi(tagValue)
			if err != nil {
				return ValidationError{Field: fieldName, Err: fmt.Errorf("cannot validate value [%s] using tag [%s]", fieldName, tag)} // TODO fix
			}
			if len(fieldValue) != validLength {
				return ValidationError{Field: fieldName, Err: fmt.Errorf("cannot validate value [%s] using tag [%s]", fieldName, tag)} // TODO fix
			}
		case RegExp:
			isMatched, _ := regexp.MatchString(tagValue, fieldValue)
			if !isMatched {
				return ValidationError{Field: fieldName, Err: fmt.Errorf("cannot validate value [%s] using tag [%s]", fieldName, tag)} // TODO fix

			}
		case In:
			tagValues := toStringSlice(tagValue, inSeparator)
			if !containsString(tagValues, fieldValue) {
				return ValidationError{Field: fieldName, Err: fmt.Errorf("cannot validate value [%s] using tag [%s]", fieldName, tag)} // TODO fix
			}
		default:
			return ValidationError{Field: fieldName, Err: fmt.Errorf("cannot validate value [%s] using tag [%s]", fieldName, tag)} // TODO fix
		}
	}
	return ValidationError{}
}

func validateInt64(fieldName string, fieldValue int64, tags []string) ValidationError {
	for _, tag := range tags {
		tagType, tagValue := getTypeAndValueOfTag(tag)
		switch tagType {
		case NotValid:
			return ValidationError{Field: fieldName, Err: fmt.Errorf("cannot validate value [%s] using tag [%s]", fieldName, tag)}
		case Min:
			min, err := strconv.Atoi(tagValue)
			if err != nil {
				return ValidationError{Field: fieldName, Err: fmt.Errorf("cannot validate value [%s] using tag [%s]", fieldName, tag)} // TODO fix
			}
			if fieldValue < int64(min) {
				return ValidationError{Field: fieldName, Err: fmt.Errorf("cannot validate value [%s] using tag [%s]", fieldName, tag)} // TODO fix
			}
		case Max:
			max, err := strconv.Atoi(tagValue)
			if err != nil {
				return ValidationError{Field: fieldName, Err: fmt.Errorf("cannot validate value [%s] using tag [%s]", fieldName, tag)} // TODO fix
			}
			if fieldValue > int64(max) {
				return ValidationError{Field: fieldName, Err: fmt.Errorf("cannot validate value [%s] using tag [%s]", fieldName, tag)} // TODO fix
			}
		case In:
			tagValues := toInt64Slice(tagValue, inSeparator)
			if !containsInt64(tagValues, fieldValue) {
				return ValidationError{Field: fieldName, Err: fmt.Errorf("cannot validate value [%s] using tag [%s]", fieldName, tag)} // TODO fix
			}
		default:
			return ValidationError{Field: fieldName, Err: fmt.Errorf("cannot validate value [%s] using tag [%s]", fieldName, tag)} // TODO fix
		}
	}
	return ValidationError{}
}

func validateSlice(field reflect.StructField, sliceValue reflect.Value, tags []string) ValidationError {
	for i := 0; i < sliceValue.Len(); i++ {
		value := sliceValue.Index(i)
		if value.Kind() == reflect.String {
			err := validateString(field.Name, value.String(), tags)
			return err
		}
		if value.Kind() == reflect.Int {
			err := validateInt64(field.Name, value.Int(), tags)
			return err
		}
	}
	return ValidationError{}
}

// before returns substring before a string.
func before(str, substr string) string {
	pos := strings.Index(str, substr)
	if pos == -1 {
		return empty
	}
	return strings.ToLower(str[0:pos])
}

// after returns substring after a string.
func after(str, substr string) string {
	pos := strings.LastIndex(str, substr)
	if pos == -1 {
		return empty
	}
	adjustedPos := pos + len(substr)
	if adjustedPos >= len(str) {
		return empty
	}
	return str[adjustedPos:]
}

func containsString(elements []string, element string) bool {
	for _, e := range elements {
		if e == element {
			return true
		}
	}
	return false
}

func containsInt64(elements []int64, element int64) bool {
	for _, e := range elements {
		if e == element {
			return true
		}
	}
	return false
}

func toStringSlice(str, sep string) []string {
	rawValues := strings.Split(str, sep)
	result := make([]string, 0)
	for _, s := range rawValues {
		result = append(result, strings.TrimSpace(s))
	}
	return result
}

func toInt64Slice(str, sep string) []int64 {
	rawValues := strings.Split(str, sep)
	result := make([]int64, 0)
	for _, s := range rawValues {
		st := strings.TrimSpace(s)
		v, err := strconv.Atoi(st)
		if err != nil {
			continue
		}
		result = append(result, int64(v))
	}
	return result
}
