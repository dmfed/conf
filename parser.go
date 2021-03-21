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

// ParseFile reads Config from file. It returns nil and error if os.Open(filename) fails.
// It would be wise to always check returned error.
// ParseFile captures two types of values: "key = value" and "option". Either key value pair or
// option must be put in its own line in the file.
// key or option must be a single word. In example line "option1 option2" only option1
// will be captured. Thus option2 needs to be put in a separate line in the file.
// In line "key = value1 value2 value3" all of value1, value2, and value3 will be
// captured.
func ParseFile(filename string) (*Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return parseReader(file), nil
}

// ParseReader reads Config from r. See also ParseFile.
func ParseReader(r io.Reader) *Config {
	return parseReader(r)
}

// ParseReadCloser reads Config from r and calls r.Close(). See also ParseFile.
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
