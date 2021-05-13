package validation

import (
	"encoding/json"
	"log"
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

var (
	// Store all of the known types to understand captured type.
	// If the captured type is a struct we need to validate it recursively.
	knownTypes = map[KeyType]bool{
		KeyTypeString: true,
		KeyTypeTime:   true,
		KeyTypeBool:   true,
	}
)

const (
	KeyTypeString  KeyType = "string"
	KeyTypeTime    KeyType = "time"
	KeyTypeBool    KeyType = "bool"
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
	NewBasicValidation(KeyTypeString, notBlank)
	NewBasicValidation(KeyTypeString, matchesRegExp)
	NewBasicValidation(KeyTypeString, lengthBetween)

	NewBasicValidation(KeyTypeInt, between)
}

func NewBasicValidation(key KeyType, fun BasicValidationFunc) {
	internalValidator := newInternalValidatorBasicFun(fun)
	generalValidationRules[string(key)] = append(generalValidationRules[string(key)], internalValidator)
}

func NewAdvValidation(key KeyType, fun AdvancedValidationFunc) {
	internalValidator := newInternalValidatorAdvFun(fun)
	generalValidationRules[string(key)] = append(generalValidationRules[string(key)], internalValidator)
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

func ValidateWithBytes(bytes []byte, inp interface{}) error {
	err := json.Unmarshal(bytes, inp)

	if err != nil {
		return err
	}

	return Validate(inp)
}

/* Validate function checks every field of the given interface recursively.
According to given tags it will execute the predefined rules for the type.
If any error occurs at the validation process it will return a meaningful error.

__Constraints__

1. Given interface and interfaces as fields should be pointers.

__Advanced Users__

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
		fieldType := field.Type.Name()

		// If current type is a struct recursively validate its fields
		if field.Type.Kind() == reflect.Struct {
			err := Validate(reflectValue.FieldByName(field.Name).Interface())

			if err != nil {
				return err
			}

		} else if validationFuncList, ok := generalValidationRules[fieldType]; ok {
			err := executeValidationRules(validationFuncList, field, reflectValue)

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
