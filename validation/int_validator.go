package validation

import (
	"errors"
	"log"
	"reflect"
	"strconv"
	"strings"
)

const (
	betweenElemCountErr      = "between tag should contains floor and ceil values seperated by '-'"
	betweenElemNoIntErr      = "between elements should be int"
	betweenValueIsNotInRange = "given value is not between the range"
)

var (
	between = func(field, value string, tags reflect.StructTag) error {
		if betweenStr, ok := tags.Lookup("between"); ok {
			minMaxSplitted := strings.Split(betweenStr, "-")
			if len(minMaxSplitted) != 2 {
				return errors.New(betweenElemCountErr)
			}
			min, err := strconv.Atoi(minMaxSplitted[0])
			max, err := strconv.Atoi(minMaxSplitted[1])

			if err != nil {
				return errors.New(betweenElemNoIntErr)
			}

			// value should be converted to int automatically
			// maybe we can send data as byte array too. Because it is simplistic
			// and also for struct types it is more efficient to transfer data in byte format
			intVal, _ := strconv.Atoi(value)
			if min > intVal || max < intVal {
				log.Printf("min: %d, max: %d, value: %d\n", min, max, intVal)
				return errors.New(betweenValueIsNotInRange)
			}
		}
		return nil
	}
)
