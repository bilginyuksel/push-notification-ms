package validation

import (
	"encoding/json"
	"log"
	"reflect"
)

type (
	ValidateInt            func(field string, value int, tags reflect.StructTag) error
	ValidateInt64          func(field string, value int64, tags reflect.StructTag) error
	ValidateString         func(field string, value string, tags reflect.StructTag) error
	ValidateFloat32        func(field string, value float32, tags reflect.StructTag) error
	ValidateFloat64        func(field string, value float64, tags reflect.StructTag) error
	AdvancedValidationFunc func(field reflect.StructField, value reflect.Value) error
)

type internalValidator interface {
	basicValidate(field string, value int, tags reflect.StructTag) error
	advValidate(field reflect.StructField, value reflect.Value) error
	isBasic() bool
}

type internalValidatorImpl struct {
	basicValidator    ValidateInt
	advancedValidator AdvancedValidationFunc
	basic             bool
}

// const (
// 	String  Type = Type(reflect.String)
// 	Time    Type = "time"
// 	Bool    Type = "bool"
// 	Int     Type = "int"
// 	Int64   Type = "int64"
// 	Float32 Type = "float32"
// 	Float64 Type = "float64"
// )

const (
	Time   reflect.Kind = 99
	String reflect.Kind = reflect.String
)

var (
	validatorMap = make(map[reflect.Kind][]internalValidator)
)

func init() {
	NewAdvValidation(reflect.String, notBlank)
	NewAdvValidation(reflect.String, matchesRegExp)
	NewAdvValidation(reflect.String, lengthBetween)

	NewBasicValidation(reflect.Int, between)
}

func NewBasicValidation(key reflect.Kind, fun ValidateInt) {
	internalValidator := newInternalValidatorBasicFun(fun)
	validatorMap[key] = append(validatorMap[key], internalValidator)
}
func NewValidateInt(fun ValidateInt) {
	internalValidator := newInternalValidatorBasicFun(fun)
	validatorMap[reflect.Int] = append(validatorMap[reflect.Int], internalValidator)
}

func NewAdvValidation(key reflect.Kind, fun AdvancedValidationFunc) {
	internalValidator := newInternalValidatorAdvFun(fun)
	validatorMap[key] = append(validatorMap[key], internalValidator)
}

func newInternalValidatorBasicFun(fun ValidateInt) internalValidator {
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

func ValidateWithBytes(bytes []byte, inp interface{}) error {
	err := json.Unmarshal(bytes, inp)

	if err != nil {
		return err
	}

	return Validate(inp)
}

/*
Validate function checks every field of the given interface recursively.
According to given tags it will execute the predefined rules for the type.
If any error occurs at the validation process it will return a meaningful error.

	Constraints

1. Given interface and interfaces as fields should be pointers.

	Advanced Users

Users can define new tags and rules for existing types easily.
There are multiple default tags already defined and ready to use.
But if the developer wants to define a new custom tag, he/she can easily do that.


	Example:


	import (
		"log"
		"regexp"
		"reflect"
		"github.com/bilginyuksel/validation"
	)

	// TestString is a struct to test simple custom validation rule
	// we define this to simply restrict input as email format
	type TestString struct {
		Email string `json:"email" pattern:"^[\\w-\\.]+@([\\w-]+\\.)+[\\w-]{2,4}$"`
	}

	var (
		matchesRegExp = func(field, value string, tags reflect.StructTag) error {
			pattern, ok := tags.Lookup("pattern")

			if !ok {
				return nil
			}

			regExp, err := regexp.Compile(pattern)
			if err != nil {
				return err
			}

			if !regExp.MatchString(value) {
				log.Printf("value is not a match with the pattern, value: %v", value)
				return errors.New("value is not a match with the pattern")
			}

			return nil
		}
	)

	func main() {
		// Create the custom validation tag
		validation.NewBasicValidation(validation.KeyTypeString, matchesRegExp)

		testInp := &TestString{Email: "test@gmail.com"}

		err := validation.Validate(testInp) // the output will be nil
		log.Println(err)
	}


*/
func Validate(inp interface{}) error {
	log.Printf("inp: %+v", inp)
	// Structs can be nil so check if it is nil or not...
	// And also structs can have some special validations also so maybe
	// we can close this nil check,.... but it is important because if the
	// struct is nil, that means we can't access its fields, so I need to check
	// it immediately if I have to.
	if inp == nil {
		return nil
	}

	structType := reflect.TypeOf(inp)
	reflectValue := reflect.ValueOf(inp).Elem()

	numFields := structType.Elem().NumField()

	for i := 0; i < numFields; i++ {
		field := structType.Elem().Field(i)
		kind := field.Type.Kind()

		if validatorFunList, ok := validatorMap[kind]; ok {
			err := executeValidationRules(validatorFunList, field, reflectValue)

			if err != nil {
				return err
			}
		}

		// If current type is a struct recursively validate its fields
		if kind == reflect.Struct {
			err := Validate(reflectValue.FieldByName(field.Name).Interface())

			if err != nil {
				return err
			}
		}
	}

	return nil
}

func executeValidationRules(ruleFuncList []internalValidator, field reflect.StructField, value reflect.Value) error {
	for _, validator := range ruleFuncList {
		var err error
		if validator.isBasic() {
			err = validator.basicValidate(field.Name, int(value.FieldByName(field.Name).Int()), field.Tag)
		} else {
			err = validator.advValidate(field, value.FieldByName(field.Name))
		}

		if err != nil {
			return err
		}
	}
	return nil
}

func (iv *internalValidatorImpl) basicValidate(field string, value int, tags reflect.StructTag) error {
	return iv.basicValidator(field, value, tags)
}

func (iv *internalValidatorImpl) advValidate(field reflect.StructField, value reflect.Value) error {
	return iv.advancedValidator(field, value)
}

func (iv *internalValidatorImpl) isBasic() bool {
	return iv.basic
}
