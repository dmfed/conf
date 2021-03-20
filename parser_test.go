package conf

import (
	"bytes"
	"fmt"
	"testing"
)

var testConf = []byte(`network
server=10.0.0.10
port=10000
token=test
editor=vim
color`)

func TestParser(t *testing.T) {
	r := bytes.NewReader(testConf)
	c := parseReader(r)
	if c.GetString("token") != "test" {
		fmt.Println("failed finding key value")
		t.Fail()
	}
	if c.GetInt("port") != 10000 {
		fmt.Println("failed finding int")
		t.Fail()
	}
	if c.GetBool("color") != true {
		fmt.Println("failed finding bool")
		t.Fail()
	}
}
