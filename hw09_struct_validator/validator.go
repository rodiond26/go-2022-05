package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

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

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	if len(v) == 0 {
		return "no validation errors"
	}

	var sb strings.Builder
	for i := 0; i < len(v); i++ {
		sb.WriteString(fmt.Sprintf("Field: %s; Error(s): %s", v[i].Field, v[i].Err))
		sb.WriteString("\n")
	}
	return sb.String()
}

func (v ValidationError) Error() string {
	return fmt.Sprintf("Field [%s] error [%s]", v.Field, v.Err)
}

func Validate(v interface{}) error {
	if v == nil {
		return nil
	}
	typeOfV := reflect.TypeOf(v)
	valueOfV := reflect.ValueOf(v)
	errs := make(ValidationErrors, 0)

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
				for j := range err {
					errs = append(errs, toValidationError(field.Name, err[j]))
				}
			}
			continue
		}
		if field.Type.Kind() == reflect.Int {
			err := validateInt64(field.Name, valueOfV.Field(i).Int(), tags)
			if err != nil {
				for j := range err {
					errs = append(errs, toValidationError(field.Name, err[j]))
				}
			}
			continue
		}
		if field.Type.Kind() == reflect.Slice {
			err := validateSlice(field, valueOfV.Field(i), tags)
			if err != nil {
				for j := range err {
					errs = append(errs, toValidationError(field.Name, err[j]))
				}
			}
		}
	}
	return errs
}

type TagType int

func getTypeAndValueOfTag(tag string) (TagType, string) {
	prefix := before(tag, tagBorder)
	switch prefix {
	case lengthTag:
		return Length, after(tag, tagBorder)
	case regexpTag:
		return RegExp, after(tag, tagBorder)
	case inTag:
		return In, after(tag, tagBorder)
	case minTag:
		return Min, after(tag, tagBorder)
	case maxTag:
		return Max, after(tag, tagBorder)
	default:
		return NotValid, tag
	}
}

func validateString(fieldName, fieldValue string, tags []string) []error {
	errs := make([]error, 0)
	for _, tag := range tags {
		tagType, tagValue := getTypeAndValueOfTag(tag)
		switch tagType {
		case NotValid:
			errs = append(errs, validationError(fieldName, fieldValue, tagValue))
		case Length:
			validLength, err := strconv.Atoi(tagValue)
			if err != nil {
				errs = append(errs, validationError(fieldName, string(fieldValue), tagValue))
				continue
			}
			if len(fieldValue) != validLength {
				errs = append(errs, validationError(fieldName, string(fieldValue), tagValue))
			}
		case RegExp:
			isMatched, _ := regexp.MatchString(tagValue, fieldValue)
			if !isMatched {
				errs = append(errs, validationError(fieldName, string(fieldValue), tagValue))
			}
		case In:
			tagValues := toStringSlice(tagValue, inSeparator)
			if !containsString(tagValues, fieldValue) {
				errs = append(errs, validationError(fieldName, string(fieldValue), tagValue))
			}
		default:
			continue
		}
	}
	return errs
}

func validateInt64(fieldName string, fieldValue int64, tags []string) []error {
	errs := make([]error, 0)
	for _, tag := range tags {
		tagType, tagValue := getTypeAndValueOfTag(tag)
		switch tagType {
		case NotValid:
			errs = append(errs, validationError(fieldName, string(fieldValue), tagValue))
		case Min:
			min, err := strconv.Atoi(tagValue)
			if err != nil {
				errs = append(errs, validationError(fieldName, string(fieldValue), tagValue))
				continue
			}
			if fieldValue < int64(min) {
				errs = append(errs, validationError(fieldName, string(fieldValue), tagValue))
			}
		case Max:
			max, err := strconv.Atoi(tagValue)
			if err != nil {
				errs = append(errs, validationError(fieldName, string(fieldValue), tagValue))
				continue
			}
			if fieldValue > int64(max) {
				errs = append(errs, validationError(fieldName, string(fieldValue), tagValue))
			}
		case In:
			tagValues := toInt64Slice(tagValue, inSeparator)
			if !containsInt64(tagValues, fieldValue) {
				errs = append(errs, validationError(fieldName, string(fieldValue), tagValue))
			}
		default:
			continue
		}
	}
	return errs
}

func validateSlice(field reflect.StructField, sliceValue reflect.Value, tags []string) ValidationErrors {
	errs := make([]ValidationError, 0)
	for i := 0; i < sliceValue.Len(); i++ {
		value := sliceValue.Index(i)
		if value.Kind() == reflect.String {
			err := validateString(field.Name, value.String(), tags)
			if err != nil {
				for j := range err {
					errs = append(errs, toValidationError(field.Name, err[j]))
				}
			}
		}
		if value.Kind() == reflect.Int {
			err := validateInt64(field.Name, value.Int(), tags)
			if err != nil {
				for j := range err {
					errs = append(errs, toValidationError(field.Name, err[j]))
				}
			}
		}
	}
	return errs
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

func validationError(fieldName, fieldValue, tag string) error {
	return fmt.Errorf("cannot validate field [%s] value [%s] using tag [%s]", fieldName, fieldValue, tag)
}

func toValidationError(fieldName string, errorMsg error) ValidationError {
	return ValidationError{Field: fieldName, Err: errorMsg}
}
