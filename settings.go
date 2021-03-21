package conf

import (
	"errors"
	"strconv"
	"strings"
)

var (
	// ErrNotFound is returned when trying to get empty value from Setting.
	ErrNotFound = errors.New("key was not found")
	// ErrParsingBool is returned when Setting.Bool() method is called and Setting.Value
	// can not be interpreted as boolean.
	ErrParsingBool = errors.New("value can not be interpreted as bool")
)

// Option represents single valuse read from config
type Setting struct {
	// Key holds the name of key parsed from the configuration
	Key string
	// Value holds the value of key parsed from the configuration
	Value string
	// Found is set to true if Key was found in the parsed config, false otherwise
	Found bool
}

// Int converts Option's Value to int if possible
// If option was not found in the config this method will
// return ErrNotFound
func (st Setting) Int() (int, error) {
	switch st.Found {
	case true:
		return strconv.Atoi(st.Value)
	default:
		return 0, ErrNotFound
	}
}

// Float64 converts Option's Value to float64 if possible
// If setting was not found in the config this method will
// return ErrNotFound
func (st Setting) Float64() (float64, error) {
	switch st.Found {
	case true:
		return strconv.ParseFloat(st.Value, 64)
	default:
		return 0, ErrNotFound
	}
}

// String returns option Value as string
func (st Setting) String() (string, error) {
	switch st.Found {
	case true:
		return st.Value, nil
	default:
		return "", ErrNotFound
	}
}

// Bool tries to interpret option Value as bool
// "1", "true", "yes" (case insensitive) yields true
// "0", "false", "no" (case insensitive) yields false
func (st Setting) Bool() (bool, error) {
	switch st.Found {
	case true:
		return parseBool(st.Value)
	default:
		return false, ErrNotFound
	}
}

// Split splits Setting's value with separator sep and returns
// []Setting. If separator was not found the method returns slice
// with only one Setting. This method is intended for use when
// config file has comma (or any other delimiter) separated values
// like:
//    myoption = first,second,third
// Split(",") will return slice with 3 separate Setting each holding
// one of "first, second, third"
func (st Setting) Split(sep string) []Setting {
	options := []Setting{}
	if !st.Found {
		return options
	}
	for _, val := range strings.Split(st.Value, sep) {
		val = strings.Trim(val, " ")
		options = append(options, Setting{Value: val, Found: true})
	}
	return options
}

func parseBool(s string) (value bool, err error) {
	switch strings.ToLower(s) {
	case "1":
		fallthrough
	case "yes":
		fallthrough
	case "true":
		value = true
	case "0":
		fallthrough
	case "no":
		fallthrough
	case "false":
		value = false
	default:
		err = ErrParsingBool
	}
	return
}
