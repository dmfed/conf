// module conf
//
// import "github.com/dmfed/conf"
//
// Module conf implements a very simple config parser with two types of
// values: key value pairs and single word options. Each of these must be
// put on separate line in a file like this:
//	key1 = value
//	key2 = value1, value2, value3
//	option1
//	option2
// Values can also be read from any io.Reader or io.ReadCloser
//
// Typical use case would look like this:
//
//	config, err := conf.ParseFile("filename")
//	if err != nil {
//		// Means we failed to read from file
//		// config variable is now nil and unusable
//	}
//	value, err := config.GetSetting("mykey").Float64()
//	if err != nil {
//		// Means that value has not been found
//		// or can not be cast to desired type
//	}
//	// value now holds float64.
//
//	value2, _ := config.GetSetting("otherkey").String()
//	// value2 now holds string if "otherkey" was parsed, else an empty string.
//
// Trying to extract non existing value will always return default value for
// the type.
//
// See description of module's types and methods which are
// quite self-explanatory.
//
package conf

// Config holds parsed keys and values. Settings and Options can be
// accessed with Config.Settings and Config.Options maps directly.
type Config struct {
	// Settings store key value pairs ("key = value" in config file)
	// all key value pairs found when parsing input are accumulated in this map.
	Settings map[string]string
	// Options map stores single word options ("option" in config file)
	Options map[string]struct{}
}

// Get returns a Setting. If key was not found the returned Setting's Value
// will be empty string and Setting's methods to extract specific type of Value
// like Int(), String(), IntSlice() etc. will return errors.
func (c *Config) Get(key string) (s Setting) {
	s.Value, s.found = c.Settings[key]
	return
}

// Find returns a Setting. It is similar to Get method, but it return an error
// if key was not found. In this case Setting's Value field will be empty string
// and Setting's methods to extract value like Int(), String() etc. will return errors.
func (c *Config) Find(key string) (s Setting, err error) {
	s.Value, s.found = c.Settings[key]
	if !s.found {
		err = ErrNotFound
	}
	return
}

// GetDefault looks up for requested key and Setting with its value.
// If lookup fails it returns Setting with Value field set to def.
// Extraction methods
func (c *Config) GetDefault(key, def string) (s Setting) {
	s.Value, s.found = c.Settings[key]
	if !s.found {
		s.Value = def
		s.found = true
	}
	return
}

// HasSetting returns true if line:
// 	"key = somevalue"
// was found in the parsed data
func (c *Config) Has(key string) (exists bool) {
	_, exists = c.Settings[key]
	return
}

// HasOption returns true if line:
//	"key"
// was found in the parsed file
func (c *Config) HasOption(option string) (exists bool) {
	_, exists = c.Options[option]
	return
}
