package conf

import (
	"bytes"
	"fmt"
	"testing"
)

var testConf = []byte(`network
server=10.0.0.10
port=10000
# removed
token=test
editor = vim
color`)

func TestParser(t *testing.T) {
	r := bytes.NewReader(testConf)
	c := parseReader(r)
	if v, _ := c.Get("token").String(); v != "test" {
		fmt.Println("failed finding key value")
		t.Fail()
	}
	if v, _ := c.Get("editor").String(); v != "vim" {
		fmt.Println("failed finding key value")
		t.Fail()
	}
	if v, _ := c.Get("port").Int(); v != 10000 {
		fmt.Println("failed finding int")
		t.Fail()
	}
	if v := c.Get("color").Bool(); v != true {
		fmt.Println("failed finding bool")
		t.Fail()
	}
	if c.Has("removed") {
		fmt.Println("commented out line shows up in config")
		t.Fail()
	}
}
