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
iamsure = hopefully you're not mistaken
float = 13.1984
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

	// Find( key string) returns instance of Setting and an
	// error if key was not found
	//
	// type Setting struct {
	// 		Value 	string	// value if found
	// }
	//
	// You can access Setting's Value field (type string) directly.
	port, err := config.Find("port")
	fmt.Printf("variable port has type: %T, port.Value == %v, type of port.Value is: %T, error returned: %v\n", port, port.Value, port.Value, err)

	// We can cast Setting.Value to a desired type including int, float64,
	// bool and string. Method will return an error if Setting Value field
	// can not be interpreted as desired type.
	n, err := port.Int()
	fmt.Printf("n has value: %v, type: %T, err: %v\n", n, n, err)

	// Another syntax for getting Setting instance is to use Get() method
	// which never returns errors. Get will return Setting with empty string
	// in Value filed if requested key was now found.
	d := config.Get("distance")
	distance, err := d.Float64()
	fmt.Printf("var distance has value: %v, type: %T, error value: %v\n", distance, distance, err)

	// Get() syntax can be is slightly shorter if we're "sure" that key exists in
	// the config.
	sure := config.Get("iamsure").String() // String method never returns errors
	// or simply sure := config.Get("iamsure").Value:
	fmt.Println(sure)

	// or like this:
	fl, _ := config.Get("float").Float64() // we're dropping error here. Bad if value fails to convert.
	fmt.Println("fl, _ := config.Get(\"float\").Float64() Value of var fl is:", fl)

	// We can access comma-separated values of key like this:
	servers := config.Get("servers").StringSlice()
	fmt.Println("Found servers:")
	for i, s := range servers {
		fmt.Println(i+1, "\t", s)
	}
	// There is also a GetDefault() method
	def := config.GetDefault("non-existant-key", "myvalue")
	fmt.Println(def.Value) // "myvalue"

	// You can use HasOption method to find whether single-word options were
	// present in the the config
	if config.HasOption("color") {
		fmt.Println("Hooray, we've found option \"color\"!")
		// do something useful
	}

	// Below code finds two keys with bool values in the Config
	// and outputs those.
	var t, f bool
	t, _ = config.Get("booltrue").Bool()
	f, _ = config.Get("bool0").Bool()
	fmt.Printf("t's type is: %T, value: %v, f's type is: %T, value: %v\n\n", t, t, f, f)
}
