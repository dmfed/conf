package main

import (
	"bytes"
	"fmt"

	"github.com/dmfed/conf"
)

var testConf = []byte(`
# commented
port=10000
servers = 10.0.0.1, 10.0.0.2, 10.0.0.3
bool0=0
booltrue = true
distance=13.42
color`)

func main() {
	r := bytes.NewReader(testConf) // creating io.Reader from []byte()
	config := conf.ParseReader(r)  // we could call conf.ParseFile("filename") here

	// First of all we can access parsed values directly:
	for key, value := range config.Settings {
		fmt.Println(key, value)
	}
	for opt := range config.Options {
		fmt.Println(opt)
	}
	fmt.Println()

	// GetSetting( key string) returns instance of Setting:
	// type Setting struct {
	//		Key		string	// requested key
	// 		Value 	string	// value if found
	// 		Found 	bool 	// true if requested key-value pair was found else false
	// }
	port := config.GetSetting("port")
	fmt.Printf("variable port has type: %T\nport.Key == %v, port.Value == %v, port.Found == %v\n", port, port.Key, port.Value, port.Found)

	// We can cast Setting.Value to a desired type including int, float64,
	// bool and string. If requested Setting was not found trying to extract
	// its value will return an error.
	// Method below will return error if "port = somevalue" was not in config
	// file or if conversion to int fails.
	n, _ := port.Int()
	fmt.Printf("n has value: %v, type: %T\n", n, n)

	// Split() method gets comma separated values from key-value pair:
	// 		key = value1, value2, value3
	// and retuns slice of []Setting. It will be empty if requested key
	// was not found and will have only one element if values were not actually
	// comma separated. Note that Split() never returns nil.
	servers := config.GetSetting("servers").Split()
	fmt.Println("Found servers:")
	for i, s := range servers {
		ip, _ := s.String()
		fmt.Println(i+1, "\t", ip)
	}

	if config.HasOption("color") {
		fmt.Println("Hooray, we've found option \"color\"!")
		// do something useful
	}

	distance, _ := config.GetSetting("distance").Float64()
	fmt.Printf("distance has value: %v, type: %T\n\n", distance, distance)

	var t, f bool
	t, _ = config.GetSetting("booltrue").Bool()

	// We can also check if key-value pair exists prior to
	// actually trying to get it.
	if config.HasSetting("bool0") {
		f, _ = config.GetSetting("bool0").Bool()
	}
	fmt.Printf("t's type is: %T, value: %v, f's type is: %T, value: %v\n\n", t, t, f, f)

	_, err := config.GetSetting("commented").String()
	fmt.Printf("Commented out string \"# commented\" does not appear in Config.\nTrying to extract value will return error: %v", err)
}
