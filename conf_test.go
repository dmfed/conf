package conf

import (
	"bytes"
	"fmt"
	"testing"
)

var testConf = []byte(`port=10000
# removed
two = two words
commas = abc, def, ghi
token=test
bool1=1
bool2=true

editor = vim
distance=13.42
floats=0.5,2.37,6
floatswithstring = 0.5, hello, 0.9
no missedme
color`)

func TestPackage(t *testing.T) {
	r := bytes.NewReader(testConf)
	c := parseReader(r)
	if _, err := c.Get("floatswithstring").Float64Slice(); err == nil {
		t.Fail()
	}
	if _, err := c.Get("floats").Float64Slice(); err != nil {
		t.Fail()
	}
	if v := c.Get("token").String(); v != "test" {
		fmt.Println("failed finding key value")
		t.Fail()
	}
	if v := c.Get("editor").String(); v != "vim" {
		fmt.Println("failed finding key value")
		t.Fail()
	}
	if v, _ := c.Get("port").Int(); v != 10000 {
		fmt.Println("failed finding int")
		t.Fail()
	}
	if v, _ := c.Get("distance").Float64(); v != 13.42 {
		fmt.Println("failed finding key value")
		t.Fail()
	}
	if c.HasOption("color") != true {
		fmt.Println("failed finding option")
		t.Fail()
	}
	if v, _ := c.Get("bool1").Bool(); v != true {
		fmt.Println("failed finding bool1 value")
		t.Fail()
	}
	if v, _ := c.Get("bool2").Bool(); v != true {
		fmt.Println("failed finding bool2 value")
		t.Fail()
	}
	if v := c.Get("two").String(); v != "two words" {
		fmt.Println("failed finding key value")
		t.Fail()
	}
	if v := c.Get("commas").String(); v != "abc, def, ghi" {
		fmt.Println("failed finding key value")
		t.Fail()
	}
	if v := c.Get("nonexistent").String(); v != "" {
		fmt.Println("returned non-empty string for nonexistent key")
		t.Fail()
	}
	if c.HasOption("removed") {
		fmt.Println("commented out line shows up in config")
		t.Fail()
	}
	splitted := c.Get("commas").StringSlice()
	if len(splitted) != 3 {
		fmt.Println("could not split string")
		t.Fail()
	}
	abc := splitted[0]
	ghi := splitted[2]
	if abc != "abc" || ghi != "ghi" {
		fmt.Println("Split() returned incorrect values")
		t.Fail()
	}
	if c.HasOption("no") != true {
		fmt.Println("should capture one option per line even if line holds two")
		t.Fail()
	}
	if c.HasOption("missedme") == true {
		fmt.Println("should only capture one option per line")
		t.Fail()
	}
	st := Setting{}
	if v, err := st.Float64(); v != 0.0 || err == nil {
		fmt.Println("empty string erroneously converts to float")
		t.Fail()
	}
	if def := c.GetDefault("non-existant-key", "myvalue"); def.Value != "myvalue" {
		fmt.Println("GetDefault fails to apply default value")
		t.Fail()
	}
}
