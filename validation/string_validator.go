package validation

import (
	"errors"
	"fmt"
	"log"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

const (
	notBlankErr                  = "the string is blank"
	regExpMatchErr               = "value: %v, does not match with the pattern"
	betweenStrLengthIsNotInRange = "given string is not in range between the numbers in tag"
)

var (
	notBlank = func(value string, tags reflect.StructTag) error {
		blank, ok := tags.Lookup("blank")

		if !ok || blank == "true" || strings.TrimSpace(value) != "" {
			return nil
		}

		return errors.New(notBlankErr)
	}

	matchesRegExp = func(value string, tags reflect.StructTag) error {
		pattern, ok := tags.Lookup("pattern")

		if !ok {
			return nil
		}

		regExp, err := regexp.Compile(pattern)
		if err != nil {
			log.Printf("regular expression could not compiled, err: %v\n", err)
			return err
		}

		if !regExp.MatchString(value) {
			return errors.New(fmt.Sprintf(regExpMatchErr, value))
		}

		return nil
	}

	// Same with the int_validator between function
	// the only difference is when you use between with string it compares the length
	// but with int it compares the actual value
	lengthBetween = func(value string, tags reflect.StructTag) error {
		if lengthBetweenStr, ok := tags.Lookup("between"); ok {
			minMaxSplitted := strings.Split(lengthBetweenStr, "-")
			if len(minMaxSplitted) != 2 {
				return errors.New(betweenElemCountErr)
			}
			min, err := strconv.Atoi(minMaxSplitted[0])
			max, err := strconv.Atoi(minMaxSplitted[1])

			if err != nil {
				return errors.New(betweenElemNoIntErr)
			}

			if min > len(value) || max < len(value) {
				// log.Printf("min: %d, max: %d, value_length: %d", min, max, len(value))
				return errors.New(betweenStrLengthIsNotInRange)
			}
		}
		return nil
	}

	lengthEqual = func(value string, tags reflect.StructTag) error {
		if strSize, exists := tags.Lookup("size"); exists {
			size, err := strconv.Atoi(strSize)

			if err != nil {
				return errors.New("value of the size should be int")
			}

			if len(value) != size {
				return errors.New("value length is not equal to size")
			}
		}
		return nil
	}

	setDefault = func(reflectField reflect.StructField, reflectValue reflect.Value) error {
		tags := reflectField.Tag
		currValue := reflectValue.String()
		defaultValue, exists := tags.Lookup("default")

		// if tag does not exist or currValue is not empty do nothing
		if !exists || currValue != "" {
			return nil
		}

		if reflectValue.CanSet() {
			reflectValue.SetString(defaultValue)
			return nil
		}

		return errors.New("default value could not set")
	}
)
