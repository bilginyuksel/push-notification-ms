package validation

import (
	"encoding/json"
	"log"
	"reflect"
	"strings"
	"time"
)

type (
	Kind uint

	FuncValidateInt       func(value int, tags reflect.StructTag) error
	FuncValidateInt64     func(value int64, tags reflect.StructTag) error
	FuncValidateString    func(value string, tags reflect.StructTag) error
	FuncValidateFloat32   func(value float32, tags reflect.StructTag) error
	FuncValidateFloat64   func(value float64, tags reflect.StructTag) error
	FuncValidateBool      func(value bool, tags reflect.StructTag) error
	FuncValidateInterface func(value interface{}, tags reflect.StructTag) error
	FuncValidateTime      func(value time.Time, tags reflect.StructTag) error
	FuncValidateGlob      func(field reflect.StructField, value reflect.Value) error
)

type Validator interface {
	Validate(field reflect.StructField, value reflect.Value) error
}

type (
	intValidator struct {
		v FuncValidateInt
	}

	int64Validator struct {
		v FuncValidateInt64
	}

	float32Validator struct {
		v FuncValidateFloat32
	}

	float64Validator struct {
		v FuncValidateFloat64
	}

	boolValidator struct {
		v FuncValidateBool
	}

	stringValidator struct {
		v FuncValidateString
	}

	interfaceValidator struct {
		v FuncValidateInterface
	}

	timeValidator struct {
		v FuncValidateTime
	}

	globalValidator struct {
		v FuncValidateGlob
	}
)

const (
	Invalid Kind = iota
	Bool
	Int
	Int8
	Int16
	Int32
	Int64
	Uint
	Uint8
	Uint16
	Uint32
	Uint64
	Uintptr
	Float32
	Float64
	Complex64
	Complex128
	Array
	Chan
	Func
	Interface
	Map
	Ptr
	Slice
	String
	Struct
	UnsafePointer
	Time
)

var kindMap = map[string]Kind{
	"bool":      Bool,
	"int":       Int,
	"int64":     Int64,
	"float32":   Float32,
	"float64":   Float64,
	"array":     Array,
	"interface": Interface,
	"map":       Map,
	"slice":     Slice,
	"string":    String,
	"struct":    Struct,
	"time.Time": Time,
}

var (
	validatorMap = make(map[Kind][]Validator)
)

// NewKind when specifying a new kind name you need to check the correct name for your struct.
// For example if you work on package name validation in your project. The structs you have created
// in that package should named like validation.StructName...
func NewKind(fullKindName string, customKind Kind) {
	if _, ok := kindMap[fullKindName]; ok {
		log.Printf("kind is already exists")
	} else {
		kindMap[fullKindName] = customKind
	}
}

func init() {
	NewStringValidator(matchesRegExp)
	NewStringValidator(lengthBetween)
	NewStringValidator(notBlank)
	NewStringValidator(lengthEqual)
	NewValidator(String, setDefault)

	NewIntValidator(between)
}

func NewIntValidator(fun FuncValidateInt) {
	v := &intValidator{v: fun}
	validatorMap[Int] = append(validatorMap[Int], v)
}

func NewInt64Validator(fun FuncValidateInt64) {
	v := &int64Validator{v: fun}
	validatorMap[Int64] = append(validatorMap[Int64], v)
}

func NewFloat32Validator(fun FuncValidateFloat32) {
	v := &float32Validator{v: fun}
	validatorMap[Float32] = append(validatorMap[Float32], v)
}

func NewFloat64Validator(fun FuncValidateFloat64) {
	v := &float64Validator{v: fun}
	validatorMap[Float64] = append(validatorMap[Float64], v)
}

func NewStringValidator(fun FuncValidateString) {
	v := &stringValidator{v: fun}
	validatorMap[String] = append(validatorMap[String], v)
}

func NewBoolValidator(fun FuncValidateBool) {
	v := &boolValidator{v: fun}
	validatorMap[Bool] = append(validatorMap[Bool], v)
}

func NewTimeValidator(fun FuncValidateTime) {
	v := &timeValidator{v: fun}
	validatorMap[Time] = append(validatorMap[Time], v)
}

func NewInterfaceValidator(kind Kind, fun FuncValidateInterface) {
	v := &interfaceValidator{v: fun}
	validatorMap[kind] = append(validatorMap[kind], v)
}

// NewCustomValidator users can define custom validators if they want to. It is not necessary to do that,
// you can easily use interface returned methods etc. But as you wish you can use this api to define customized validators.
func NewCustomValidator(kind Kind, v Validator) {
	isExist := false
	for _, val := range kindMap {
		if val == kind {
			isExist = true
		}
	}

	if !isExist {
		log.Println("kind is not exist, please add the custom kind first")
		return
	}

	// if there is no kind exist yet, we need to create this kind for every map
	// we should get its string representation also.
	validatorMap[kind] = append(validatorMap[kind], v)
}

func NewValidator(key Kind, fun FuncValidateGlob) {
	v := &globalValidator{v: fun}
	validatorMap[key] = append(validatorMap[key], v)
}

func (iv *intValidator) Validate(field reflect.StructField, value reflect.Value) error {
	return iv.v(int(value.Int()), field.Tag)
}

func (iv *int64Validator) Validate(field reflect.StructField, value reflect.Value) error {
	return iv.v(value.Int(), field.Tag)
}

func (iv *float32Validator) Validate(field reflect.StructField, value reflect.Value) error {
	return iv.v(float32(value.Float()), field.Tag)
}

func (iv *float64Validator) Validate(field reflect.StructField, value reflect.Value) error {
	return iv.v(value.Float(), field.Tag)
}

func (iv *stringValidator) Validate(field reflect.StructField, value reflect.Value) error {
	return iv.v(value.String(), field.Tag)
}

func (iv *globalValidator) Validate(field reflect.StructField, value reflect.Value) error {
	return iv.v(field, value)
}

func (iv *boolValidator) Validate(field reflect.StructField, value reflect.Value) error {
	return iv.v(value.Bool(), field.Tag)
}

func (iv *timeValidator) Validate(field reflect.StructField, value reflect.Value) error {
	if field.Type.String() == "*time.Time" { // pointer time conversion
		timeVal := value.Interface().(*time.Time)
		return iv.v(*timeVal, field.Tag)
	}

	timeVal := value.Interface().(time.Time)
	return iv.v(timeVal, field.Tag)
}

func (iv *interfaceValidator) Validate(field reflect.StructField, value reflect.Value) error {
	return iv.v(value.Interface(), field.Tag)
}

// ValidateWithBytes this function gives you the chance to keep your code a little bit cleaner.
// It unmarshals the json inside so you should give the interface address and the bytes it will do the validation
// after filling the object.
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
	if inp == nil {
		return nil
	}

	structType := reflect.TypeOf(inp)
	reflectValue := reflect.ValueOf(inp).Elem()

	numFields := structType.Elem().NumField()

	for i := 0; i < numFields; i++ {
		field := structType.Elem().Field(i)
		kind := field.Type.Kind()
		// we use kind name without pointer prefix so remove it
		kindName := strings.TrimPrefix(field.Type.String(), "*") // time.Time, *time.Time

		if validatorFunList, ok := validatorMap[Kind(kind)]; ok {
			err := executeValidationRules(validatorFunList, field, reflectValue)

			if err != nil {
				return err
			}
		}

		// if kindName is a struct, then we should check for the kindName
		// because if there are any special validators used for this struct type
		// we should go for it otherwise we will recursively check for its fields...
		if k, ok := kindMap[kindName]; ok && kind == reflect.Struct {
			// There is a special kind and we need to find its validators
			// After we execute struct process if also it is exist
			if specialValidators, ok := validatorMap[k]; ok {
				err := executeValidationRules(specialValidators, field, reflectValue)

				if err != nil {
					return err
				}
			}
		} else if kind == reflect.Struct {
			// If current type is a struct then recursively Validate its fields
			err := Validate(reflectValue.FieldByName(field.Name).Interface())

			if err != nil {
				return err
			}
		}
	}

	return nil
}

func executeValidationRules(ruleFuncList []Validator, field reflect.StructField, value reflect.Value) error {
	for _, validator := range ruleFuncList {
		if err := validator.Validate(field, value.FieldByName(field.Name)); err != nil {
			return err
		}
	}

	return nil
}
