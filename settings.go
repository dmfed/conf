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
	"t":    true,
	"1":    true,
	"yes":  true,
	"on":   true,
	// what evaluates to false
	"false": false,
	"f":     false,
	"0":     false,
	"no":    false,
	"off":   false,
}

// Setting represents key-value pair read from config file.
// It's Value field holds the value of key parsed from the configuration
type Setting struct {
	Value string
}

// Int converts Setting Value to int. Returned error
// will be non nil if convesion failed.
func (st Setting) Int() (int, error) {
	return parseInt(st.Value)
}

// IntSlice splits Setting Value (separator is ",") and adds
// each of resulting values to []int if possible.
// If one or more values can not be converted to float64 those will be dropped
// and method will return conf.ErrCouldNotConvert.
// Check error to be sure that all required values were parsed.
func (st Setting) IntSlice() ([]int, error) {
	return parseIntSlice(st.Value, valuesSeparator)
}

/* func (st Setting) split(sep string) Setting {
	st.sep = sep //Choose separator to split values ?
	return st
} */

// Float64 converts Setting Value to float64. Returned error
// will be non nil if convesion failed.
func (st Setting) Float64() (float64, error) {
	return parseFloat64(st.Value)
}

// Float64Slice splits Setting Value (separator is ",") and adds
// each of resulting values to []float64 if possible.
// If one or more values can not be converted to float64 those will be dropped
// and method will return conf.ErrCouldNotConvert.
// Check error to be sure that all required values were parsed.
func (st Setting) Float64Slice() ([]float64, error) {
	return parseFloat64Slice(st.Value, valuesSeparator)
}

// String returns option Value as string
// This method also implements Stringer interface from fmt module
func (st Setting) String() string {
	return st.Value
}

// StringSlice splits Setting's Value (separator is ",") and adds
// each of resulting values to []string trimming leading and trailing spaces
// from each string.
func (st Setting) StringSlice() []string {
	return tidySplit(st.Value, valuesSeparator)
}

// Bool tries to interpret Setting's Value as bool
// "1", "true", "on", "yes" (case insensitive) yields true
// "0", "false", "off", "no" (case insensitive) yields false
// If nothing matches will return false and conf.ErrParsingBool
func (st Setting) Bool() (bool, error) {
	return parseBool(st.Value)
}

func parseInt(s string) (n int, err error) {
	n, err = strconv.Atoi(s)
	return
}

func parseIntSlice(s, sep string) ([]int, error) {
	var (
		n      int
		err    error
		slice  []int
		digits []string
	)
	digits = tidySplit(s, sep)
	for _, d := range digits {
		n, err = strconv.Atoi(d)
		if err != nil {
			err = ErrCouldNotConvert
			break
		}
		slice = append(slice, n)
	}
	return slice, err
}

func parseFloat64(s string) (n float64, err error) {
	n, err = strconv.ParseFloat(s, 64)
	return
}

func parseFloat64Slice(s, sep string) ([]float64, error) {
	var (
		n      float64
		err    error
		slice  []float64
		digits []string
	)
	digits = tidySplit(s, sep)
	for _, d := range digits {
		n, err = strconv.ParseFloat(d, 64)
		if err != nil {
			err = ErrCouldNotConvert
			break
		}
		slice = append(slice, n)
	}
	return slice, err
}

func parseBool(s string) (value bool, err error) {
	s = strings.ToLower(s)
	value, ok := boolMap[s]
	if !ok {
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
