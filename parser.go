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
	configKeyValueRe = regexp.MustCompile(`(\w+) *= *(.+)\b[\t| |\n]*`)
	configOptionRe   = regexp.MustCompile(`(\w+)`)
)

type Config struct {
	Values map[string]string
}

func (c *Config) Get(key string) (s string) {
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

func Parse(filename string) (*Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	c := parseReader(file)
	return c, nil
}

func parseReader(r io.Reader) *Config {
	var c Config
	c.Values = make(map[string]string)
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#") {
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
