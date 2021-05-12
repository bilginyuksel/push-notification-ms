package validation

import (
	"encoding/json"
	"reflect"
)

type Validator interface {
	Validate(field, value string, tags reflect.StructTag) bool
}

type KeyType string

type ValidationRule func(field, value string, tags reflect.StructTag) error

// optional:"false"
// empty:"false"
// between:"3,5" // float can be used also time can be used
// match:"AzI[0-9]" // match regular expression
// default: "default value" // default value
var (
	generalValidationRules = make(map[string][]ValidationRule)
)

const (
	KeyTypeString  KeyType = "string"
	KeyTypeTime    KeyType = "time"
	KeyTypeBoolean KeyType = "bool"
	KeyTypeInt     KeyType = "int"
	KeyTypeInt64   KeyType = "int64"
	KeyTypeFloat32 KeyType = "float32"
	KeyTypeFloat64 KeyType = "float64"
)

func NewValidationRule(key string, rule ValidationRule) {
	generalValidationRules[key] = append(generalValidationRules[key], rule)
}

func Validate(inp interface{}) error {
	structType := reflect.TypeOf(inp)
	r := reflect.ValueOf(inp).Elem()

	numberOfFields := structType.Elem().NumField()

	for i := 0; i < numberOfFields; i++ {
		field := structType.Elem().Field(i)
		fieldType := field.Type.Name()
		fieldName := field.Name
		fieldValue := r.Field(i).String()
		fieldTags := field.Tag

		if validationFuncList, ok := generalValidationRules[fieldType]; ok {
			err := executeValidationRules(validationFuncList, fieldName, fieldValue, fieldTags)

			if err != nil {
				return err
			}
		}
	}

	return nil
}

func ValidateWithBytes(bytes []byte, inp interface{}) error {
	err := json.Unmarshal(bytes, inp)

	if err != nil {
		return err
	}

	return Validate(inp)
}

func executeValidationRules(ruleFuncList []ValidationRule, fieldName, fieldValue string, tags reflect.StructTag) error {
	for _, ruleFunc := range ruleFuncList {
		err := ruleFunc(fieldName, fieldValue, tags)
		if err != nil {
			return err
		}
	}
	return nil
}
