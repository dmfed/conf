package conf

import (
	"bufio"
	"io"
	"os"
	"regexp"
	"strings"
)

var (
	configKeyValueRe   = regexp.MustCompile(`(\w+) *= *(.+)\b[\t| |\n]*`)
	configOptionRe     = regexp.MustCompile(`(\w+)`)
	configCommentedOut = "#"
)

// ParseFile reads values from file. It returns nil and error if os.Open(filename) fails.
// It would be wise to always check returned error.
// ParseFile captures two types of values: "key = value" and "option". Either key value pair or
// option must be put in its own line in the file.
// Key or option must be a single word. In example line:
// 	"option1 option2"
// only option1 will be captured. option2 needs to be in a separate line in the file to take effect.
// In line "key = value1,value2,value3" all of value1, value2, and value3 will be
// captured. They can be later accesed separately with Setting's Split() method.
func ParseFile(filename string) (*Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return parseReader(file), nil
}

// ParseReader reads from r and returns Config. See also ParseFile.
func ParseReader(r io.Reader) *Config {
	return parseReader(r)
}

// ParseReadCloser reads from r, returns Config and calls r.Close().
// See also ParseFile.
func ParseReadCloser(r io.ReadCloser) *Config {
	defer r.Close()
	return parseReader(r)
}

func parseReader(r io.Reader) *Config {
	var c Config
	c.Settings = make(map[string]string)
	c.Options = make(map[string]struct{})
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, configCommentedOut) {
			continue
		}
		switch {
		case configKeyValueRe.MatchString(line):
			kvpair := configKeyValueRe.FindStringSubmatch(line)
			c.Settings[kvpair[1]] = kvpair[2]
		case configOptionRe.MatchString(line):
			opt := configOptionRe.FindString(line)
			c.Options[opt] = struct{}{}
		}
	}
	return &c
}
