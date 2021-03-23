# module conf

**go get github.com/dmfed/conf** to download.

**import "github.com/dmfed/conf"** to use in your code.

Module conf implements a very simple config parser with two types of
values: key value pairs and single word options. Each of these must be
put on separate line in a file like this:
```
key1 = value
key2 = value1, value2, value3
option1
option2
```
Values can also be read from any io.Reader or io.ReadCloser. Note that
values in key value pairs mey either be separated with spaces or not.

Typical use case would look like this:
```go
config, err := conf.ParseFile("filename")
if err != nil {
	// Means we failed to read from file
	// config variable is now nil and unusable
}
value, err := config.Find("mykey")
if err != nil {
	// Means that key has not been found
}
// value now holds conf.Setting type which can be accessed directly.
fmt.Println(value.Value)
n, err := value.Float64() // tries to parse float from Setting.Value field.
// if err is nil then n holds float64.
```
Shorter syntax is also available with Get() method which drops errors if 
key was not found and simply returns Setting with empty Value field:
```go
listnumbers, err := config.Get("numbers").IntSlice()
// listnumbers holds slice of ints. If err is nil all of found values
// have converted successfully.
```
Note that we'd still like to check for errors in above example even if
we're sure the key exists.  This way we'll configrm that all values have been converted.

There is also a GetDefault() method:
```go
def := config.GetDefault("non-existant-key", "myvalue")
fmt.Println(def.Value) // "myvalue"
```
To check whether single-word options were found use:
```go
if config.HasOption("wordoption") {
	// do something
}
```
See description of module's types and methods which are
quite self-explanatory.

See also **[https://pkg.go.dev/github.com/dmfed/conf](https://pkg.go.dev/github.com/dmfed/conf)** for a complete description of module's 
functions.

Below is listing of a working example of a program parsing config (also found **example/main.go** in the repository).
```go
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
	// There is also a GetDefault() method which sets output Setting value
	// to requested default if key was not found.
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
```
Here is the full output of the above code:
```
distance 13.42
iamsure hopefully you're not mistaken
float 13.1984
port 10000
servers 10.0.0.1, 10.0.0.2, 10.0.0.3
bool0 0
booltrue true
color

variable port has type: conf.Setting, port.Value == 10000, type of port.Value is: string, error returned: <nil>
n has value: 10000, type: int, err: <nil>
var distance has value: 13.42, type: float64, error value: <nil>
hopefully you're not mistaken
fl, _ := config.Get("float").Float64() Value of var fl is: 13.1984
Found servers:
1 	 10.0.0.1
2 	 10.0.0.2
3 	 10.0.0.3
Hooray, we've found option "color"!
t's type is: bool, value: true, f's type is: bool, value: false

```
