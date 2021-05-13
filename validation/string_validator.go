package validation

import (
	"errors"
	"fmt"
	"log"
	"reflect"
	"regexp"
	"strings"
)

const (
	notBlankErr    = "the string is blank"
	regExpMatchErr = "value: %v, does not match with the pattern"
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

	setDefault = func(field, value string, tags reflect.StructTag) error {
		return nil
	}
)
