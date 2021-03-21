// module conf
//
// import "github.com/dmfed/conf"
//
// Module conf implements a very simple config parser with two types of
// values: key value pairs and single word options. Each of these must be
// put on separate line in a file like this:
//   key1 = value
//   key2 = value1, value2, value3
//   option1
//	 option2
//
// Values can also be read from any io.Reader or io.ReadCloser
//
// Typical use case would look like this:
//
// config, err := conf.ParseFile("filename")
// if err != nil {
//		// Means we failed to read from file
// 		// config variable is now nil and unusable
// }

// value, err := config.GetSetting("mykey").Float64()
// if err != nil {
// 		// Means that value has not been found
//		// or can not be cast to desired type
// }
// value now holds float64.
//
// value2, _ := config.GetSetting("otherkey").String()
//
// value2 now holds string if "otherkey" was parsed, else an empty string.
// Trying to extract non existing value will always return default value for
// the type.
//
// See description of module's types and methods which are
// quite self-explanatory.
//
// See also github.com/dmfed/conf/example/main.go for a working example
// of program parsing config.
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
// will be empty string and Setting's Found field will be set to false
func (c *Config) GetSetting(key string) (s Setting) {
	s.Value, s.Found = c.Settings[key]
	s.Key = key
	return
}

// HasSetting returns true if line:
// 	"key = somevalue"
// was found in the parsed data
func (c *Config) HasSetting(key string) (exists bool) {
	_, exists = c.Settings[key]
	return
}

// HasOption returns true if line:
// 	"key"
// was found in the parsed file
func (c *Config) HasOption(option string) (exists bool) {
	_, exists = c.Options[option]
	return
}
