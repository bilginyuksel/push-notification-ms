package validation

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
)

type TestValidation struct {
	fsd string `empty:"false"`
}

func main() {
	t1 := TestValidation{}
	t2 := TestValidation{fsd: "notEmpty"}

	type1 := reflect.TypeOf(t1)
	type2 := reflect.TypeOf(t2)

	t1Field := type1.Field(0)
	t2Field := type2.Field(0)

	fmt.Println(t1Field)
	fmt.Println(t2Field)
}

func Validate(bytes []byte, inp interface{}) bool {
	err := json.Unmarshal(bytes, inp)

	if err != nil {
		log.Printf("json unmarshal failed for input, err: %v", err)
		return false
	}

	// check input tags ...

	return true
}
