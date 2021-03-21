package conf

import (
	"errors"
	"strconv"
	"strings"
)

var (
	ErrNotFound        = errors.New("key not found")
	ErrParsingBool     = errors.New("value can not be interpreted as bool")
	ErrCouldNotConvert = errors.New("could not cast one or more values to required type")
)

var valuesSeparator = ","

var boolMap = map[string]bool{
	// what evaluates to true
	"true": true,
	"yes":  true,
	"on":   true,
	"1":    true,
	// what evaluates to false
	"false": false,
	"no":    false,
	"off":   false,
	"0":     false,
}

// Setting represents key-value pair read from config file.
// It's Value field holds the value of key parsed from the configuration
type Setting struct {
	Value string
	found bool
}

// Int converts Setting's Value to int if possible
// If setting's key was not found in the config this method will
// return ErrNotFound
func (st Setting) Int() (n int, err error) {
	switch st.found {
	case true:
		n, err = strconv.Atoi(st.Value)
	default:
		err = ErrNotFound
	}
	return
}

// IntSlice splits Setting's Value (separator is ",") and adds
// each of resulting values to []int if possible.
// If value can not be converted to int it will be dropped. Check
// len to be sure that all required values were parsed
// If Setting's key was not found in the config this method will
// return ErrNotFound
func (st Setting) IntSlice() (slice []int, err error) {
	switch st.found {
	case true:
		digits := tidySplit(st.Value, valuesSeparator)
		for _, d := range digits {
			if n, e := strconv.Atoi(d); e == nil {
				slice = append(slice, n)
			} else {
				err = ErrCouldNotConvert
			}
		}
	default:
		err = ErrNotFound
	}
	return
}

// Float64 converts Setting's Value to float64 if possible
// If setting's key was not found in the config this method will
// return ErrNotFound
func (st Setting) Float64() (n float64, err error) {
	switch st.found {
	case true:
		n, err = strconv.ParseFloat(st.Value, 64)
	default:
		err = ErrNotFound
	}
	return
}

// IntSlice splits Setting's Value (separator is ",") and adds
// each of resulting values to []float64 if possible.
// If value can not be converted to float64 it will be dropped. Check
// len to be sure that all required values were parsed
// If Setting's key was not found in the config this method will
// return ErrNotFound
func (st Setting) Float64Slice() (slice []float64, err error) {
	switch st.found {
	case true:
		digits := tidySplit(st.Value, valuesSeparator)
		for _, d := range digits {
			if n, e := strconv.ParseFloat(d, 64); e == nil {
				slice = append(slice, n)
			} else {
				err = ErrCouldNotConvert
			}
		}
	default:
		err = ErrNotFound
	}
	return
}

// String returns option Value as string
func (st Setting) String() (s string, err error) {
	switch st.found {
	case true:
		s, err = st.Value, nil
	default:
		err = ErrNotFound
	}
	return
}

// StringSlice splits Setting's Value (separator is ",") and adds
// each of resulting values to []string trimming leading and trailing spaces
// from each string.
func (st Setting) StringSlice() (slice []string, err error) {
	switch st.found {
	case true:
		slice, err = tidySplit(st.Value, valuesSeparator), nil
	default:
		err = ErrNotFound
	}
	return
}

// Bool tries to interpret Setting's Value as bool
// "1", "true", "on", "yes" (case insensitive) yields true
// "0", "false", "off", "no" (case insensitive) yields false
func (st Setting) Bool() (value bool, err error) {
	switch st.found {
	case true:
		value, err = parseBool(st.Value)
	default:
		err = ErrNotFound
	}
	return
}

func parseBool(s string) (value bool, err error) {
	s = strings.ToLower(s)
	value, ok := boolMap[s]
	switch ok {
	case true:
		err = nil
	default:
		err = ErrParsingBool
	}
	return
}

func tidySplit(s, sep string) []string {
	splitted := strings.Split(s, sep)
	for i, str := range splitted {
		splitted[i] = strings.Trim(str, " ")
	}
	return splitted
}
