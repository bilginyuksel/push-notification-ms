package validation

import (
	"encoding/json"
	"reflect"
)

type Validator interface {
	Validate(field, value string, tags reflect.StructTag) bool
}

type KeyType string

type BasicValidationFunc func(field, value string, tags reflect.StructTag) error
type AdvancedValidationFunc func(field reflect.StructField, value reflect.Value) error

type internalValidator interface {
	basicValidate(field, value string, tags reflect.StructTag) error
	advValidate(field reflect.StructField, value reflect.Value) error
	isBasic() bool
}

type internalValidatorImpl struct {
	basicValidator    BasicValidationFunc
	advancedValidator AdvancedValidationFunc
	basic             bool
}

const (
	KeyTypeString  KeyType = "string"
	KeyTypeTime    KeyType = "time"
	KeyTypeBoolean KeyType = "bool"
	KeyTypeInt     KeyType = "int"
	KeyTypeInt64   KeyType = "int64"
	KeyTypeFloat32 KeyType = "float32"
	KeyTypeFloat64 KeyType = "float64"
)

// optional:"false"
// empty:"false"
// between:"3,5" // float can be used also time can be used
// match:"AzI[0-9]" // match regular expression
// default: "default value" // default value
var (
	generalValidationRules = make(map[string][]internalValidator)
)

func init() {
	NewBasicValidation(string(KeyTypeString), notBlank)
	NewBasicValidation(string(KeyTypeString), matchesRegExp)
	// NewBasicValidation(string(KeyTypeString), setDefault)
}

func NewBasicValidation(key string, fun BasicValidationFunc) {
	internalValidator := newInternalValidatorBasicFun(fun)
	generalValidationRules[key] = append(generalValidationRules[key], internalValidator)
}

func NewAdvValidation(key string, fun AdvancedValidationFunc) {
	internalValidator := newInternalValidatorAdvFun(fun)
	generalValidationRules[key] = append(generalValidationRules[key], internalValidator)
}

func newInternalValidatorBasicFun(fun BasicValidationFunc) internalValidator {
	return &internalValidatorImpl{
		basicValidator: fun,
		basic:          true,
	}
}

func newInternalValidatorAdvFun(fun AdvancedValidationFunc) internalValidator {
	return &internalValidatorImpl{
		advancedValidator: fun,
		basic:             false,
	}
}

func Validate(inp interface{}) error {
	structType := reflect.TypeOf(inp)
	reflectValue := reflect.ValueOf(inp).Elem()

	numFields := structType.Elem().NumField()

	for i := 0; i < numFields; i++ {
		field := structType.Elem().Field(i)
		fieldType := field.Type.Name()

		if validationFuncList, ok := generalValidationRules[fieldType]; ok {
			err := executeValidationRules(validationFuncList, field, reflectValue)

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

func executeValidationRules(ruleFuncList []internalValidator, field reflect.StructField, value reflect.Value) error {
	for _, validator := range ruleFuncList {
		var err error
		if validator.isBasic() {
			err = validator.basicValidate(field.Name, value.FieldByName(field.Name).String(), field.Tag)
		} else {
			err = validator.advValidate(field, value)
		}

		if err != nil {
			return err
		}
	}
	return nil
}

func (iv *internalValidatorImpl) basicValidate(field, value string, tags reflect.StructTag) error {
	return iv.basicValidator(field, value, tags)
}

func (iv *internalValidatorImpl) advValidate(field reflect.StructField, value reflect.Value) error {
	return iv.advancedValidator(field, value)
}

func (iv *internalValidatorImpl) isBasic() bool {
	return iv.basic
}
