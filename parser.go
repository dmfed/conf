package conf

import (
	"bufio"
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

// Config holds values read from file or any other
type Config struct {
	Values map[string]string
}

func (c *Config) GetString(key string) (s string) {
	if val, ok := c.Values[key]; ok {
		s = val
	}
	return
}

func (c *Config) GetBool(key string) (b bool) {
	_, b = c.Values[key]
	return
}

func (c *Config) GetInt(key string) (n int) {
	if val, ok := c.Values[key]; ok {
		num, err := strconv.Atoi(val)
		if err == nil {
			n = num
		}
	}
	return
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
