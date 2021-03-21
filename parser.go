package conf

import (
	"bufio"
	"errors"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var (
	configKeyValueRe   = regexp.MustCompile(`(\w+) *= *(.+)\b[\t| |\n]*`)
	configOptionRe     = regexp.MustCompile(`(\w+)`)
	configCommentedOut = "#"
)

var ErrNotFound = errors.New("no such key")

// Config holds values read from file or any other
type Config struct {
	Values map[string]string
}

func (c *Config) Get(key string) (op Option) {
	op.Value, op.found = c.Values[key]
	return
}

func (c *Config) Has(key string) (present bool) {
	_, present = c.Values[key]
	return
}

type Option struct {
	Value string
	found bool
}

func (op Option) Int() (int, error) {
	switch op.found {
	case true:
		return strconv.Atoi(op.Value)
	default:
		return 0, ErrNotFound
	}
}

func (op Option) String() (string, error) {
	switch op.found {
	case true:
		return op.Value, nil
	default:
		return "", ErrNotFound
	}
}

func (op Option) Bool() bool {
	return op.found
}

func ParseFile(filename string) (*Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return parseReader(file), nil
}

func ParseReader(r io.Reader) *Config {
	return parseReader(r)
}

func parseReader(r io.Reader) *Config {
	var c Config
	c.Values = make(map[string]string)
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, configCommentedOut) {
			continue
		}
		switch {
		case configKeyValueRe.MatchString(line):
			kvpair := configKeyValueRe.FindStringSubmatch(line)
			c.Values[kvpair[1]] = kvpair[2]
		case configOptionRe.MatchString(line):
			opt := configOptionRe.FindString(line)
			c.Values[opt] = ""
		}
	}
	return &c
}
