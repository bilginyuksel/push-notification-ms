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
	notBlank = func(field, value string, tags reflect.StructTag) error {
		if blank, ok := tags.Lookup("blank"); ok && blank == "false" && strings.TrimSpace(value) != "" {
			return nil
		}
		return errors.New(notBlankErr)
	}

	matchesRegExp = func(field, value string, tags reflect.StructTag) error {
		pattern, ok := tags.Lookup("pattern")

		if !ok {
			return nil
		}

		regExp, err := regexp.Compile(pattern)
		if err != nil {
			log.Printf("[regExpMatch] regular expression could not compiled, err: %v\n", err)
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
	lengthBetween = func(field, value string, tags reflect.StructTag) error {
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

	setDefault = func(field, value string, tags reflect.StructTag) error {
		return nil
	}
)
