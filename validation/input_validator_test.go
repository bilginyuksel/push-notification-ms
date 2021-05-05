package validation

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"testing"
	"time"
)

var fieldValidator map[string]Validator

type Tok struct {
	Username   string    `json:"username" optional:"true" empty:"true" default:"somethin"`
	Email      string    `json:"email"`
	Password   string    `json:"password"`
	Age        int       `json:"age"`
	CreateTime time.Time `json:"createTime"`
}

type stringValidator struct {
}

func (sv *stringValidator) Validate(name, value string, tag reflect.StructTag) bool {
	var (
		optional = false
		empty    = false
		def      = ""
	)

	val, ok := tag.Lookup("optional")
	if ok {
		optional = val == "true"
	}

	val, ok = tag.Lookup("empty")
	if ok {
		empty = val == "true"
	}

	val, ok = tag.Lookup("default")
	if ok {
		def = val
	}

	if !optional && !empty && value == "" {
		fmt.Println("validation failed, because optional and empty rules are not satisfied")
		return false
	}

	if value == "" {
		// initialize defalt value
		value = def
		fmt.Println("default value initialized")
	}

	fmt.Printf("name: %v, value: %v, tags: %v\n", name, value, tag)

	return false
}

func Test(t *testing.T) {
	// t1 := TestValidation{}
	// t2 := TestValidation{fsd: "notEmpty"}

	// type1 := reflect.TypeOf(t1)
	// type2 := reflect.TypeOf(t2)

	// t1Field := type1.Field(0)
	// t2Field := type2.Field(0)

	// fmt.Println(t1Field)
	// fmt.Println(t2Field)
	fieldValidator = make(map[string]Validator)
	fieldValidator["string"] = &stringValidator{}

	s := &Tok{
		Username:   "bilyuk",
		Email:      "bilyuk@gmail.com",
		Password:   "pass",
		Age:        20,
		CreateTime: time.Now(),
	}

	b, _ := json.Marshal(s)
	tValidate(b, &Tok{})
}

func tValidate(bytes []byte, inp interface{}) bool {
	err := json.Unmarshal(bytes, inp)

	if err != nil {
		log.Printf("json unmarshal failed for input, err: %v", err)
		return false
	}

	// check input tags ...
	structType := reflect.TypeOf(inp)
	r := reflect.ValueOf(inp).Elem()

	fmt.Printf("rf: %+v\n", r)
	fmt.Printf("st: %+v\n", structType.Elem().Name())
	numberOfFields := structType.Elem().NumField()

	for i := 0; i < numberOfFields; i++ {
		currentField := structType.Elem().Field(i)
		fieldType := currentField.Type.Name()
		validator, ok := fieldValidator[fieldType]
		if !ok {
			fmt.Printf("no validator found\n")
			continue
		}
		v := validator.Validate(currentField.Name, r.Field(i).String(), currentField.Tag)
		if v {
			fmt.Printf("input is valid\n")
		}
		fmt.Printf("currentField: %v\n", currentField)
	}

	return true
}
