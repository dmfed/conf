// module conf
//
// import "github.com/dmfed/conf"
//
// Module conf implements a very simple config parser with two types of
// values: key value pairs and single word options. Each of these must be
// put on separate line in a file like this:
//	key1 = value
//	key2 = value1, value2, value3
//	key3=v1,v2,v3
//	option1
//	option2
// Values can also be read from any io.Reader or io.ReadCloser. Note that
// values in key value pairs mey either be separated with spaces or not.
// Typical use case would look like this:
//	config, err := conf.ParseFile("filename")
//	if err != nil {
//		// Means we failed to read from file
//		// config variable is now nil and unusable
//	}
//	value, err := config.Find("mykey")
//	if err != nil {
//		// Means that key has not been found
//	}
//	// value now holds conf.Setting.
//	n, err := value.Float64() // tries to parse float from Setting.Value field.
//	//if err is nil then n holds float64.
// There is also a quicker Get() method which returns no errors
// ("i'm feeling lucky" way to lookup values). If it does not find
// requested key, the returned Setting has empty string in Value field.
//	value2 := config.Get("otherkey")
// 	mybool, err := value2.Bool() // tries to interpret Setting.Value field as bool
//	// mybool now holds boolean if "otherkey" was found and error returned
//	// by Bool() method is nil.
// Even shorter syntax would be:
// 	listnumbers, err := config.Get("numbers").IntSlice()
//	// Note that we'd still like to check for errors even if
//	// we're sure the key exists to make sure all values are converted.
//	// listnumbers holds slice of ints. If err is nil all of found values
//	// have been converted successfully.
// To check whether single-word options were found use:
//	if config.HasOption("wordoption") {
//		// do something
//	}
// See description of module's other methods which are
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

// Find looks up a Setting and returns it. If returned error is not nil
// the requested key was not found and returned Setting has empty string in Value
// field.
func (c *Config) Find(key string) (s Setting, err error) {
	v, ok := c.Settings[key]
	if !ok {
		err = ErrNotFound
	}
	s.Value = v
	return
}

// Get returns a Setting. If key was not found the returned Setting Value
// will be empty string.
func (c *Config) Get(key string) (s Setting) {
	s.Value = c.Settings[key]
	return
}

// GetDefault looks up a Setting with requested key.
// If lookup fails it returns Setting with Value field set to def.
func (c *Config) GetDefault(key, def string) (s Setting) {
	v, ok := c.Settings[key]
	switch ok {
	case true:
		s.Value = v
	default:
		s.Value = def
	}
	return
}

// HasOption returns true if line:
//	"key"
// was found in the parsed file
func (c *Config) HasOption(option string) (exists bool) {
	_, exists = c.Options[option]
	return
}
