package conf

// Config holds keys and values parsed from file, io.Reader, or io.ReadCloser
// You can access Config.Settings and Config.Options directly.
type Config struct {
	// Settings store key value pairs ("key = value" in config file)
	Settings map[string]string
	// Options store single valued options ("option" in config file)
	Options map[string]struct{}
}

// Get returns Setting. If key was not found Setting's Value will be empty string
func (c *Config) GetSetting(key string) (op Setting) {
	op.Value, op.Found = c.Settings[key]
	op.Key = key
	return
}

// HasSetting returns true if line:
// 	"key = somevalue"
// was present in the parsed file
func (c *Config) HasSetting(key string) (exists bool) {
	_, exists = c.Settings[key]
	return
}

// HasOption returns true if line:
// 	"key"
// was present in the parsed file
func (c *Config) HasOption(option string) (exists bool) {
	_, exists = c.Options[option]
	return
}
