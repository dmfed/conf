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
no missedme
color`)

func TestPackage(t *testing.T) {
	r := bytes.NewReader(testConf)
	c := parseReader(r)
	if v, _ := c.GetSetting("token").String(); v != "test" {
		fmt.Println("failed finding key value")
		t.Fail()
	}
	if v, _ := c.GetSetting("editor").String(); v != "vim" {
		fmt.Println("failed finding key value")
		t.Fail()
	}
	if v, _ := c.GetSetting("port").Int(); v != 10000 {
		fmt.Println("failed finding int")
		t.Fail()
	}
	if v, _ := c.GetSetting("distance").Float64(); v != 13.42 {
		fmt.Println("failed finding key value")
		t.Fail()
	}
	if c.HasOption("color") != true {
		fmt.Println("failed finding option")
		t.Fail()
	}
	if v, _ := c.GetSetting("bool1").Bool(); v != true {
		fmt.Println("failed finding bool1 value")
		t.Fail()
	}
	if v, _ := c.GetSetting("bool2").Bool(); v != true {
		fmt.Println("failed finding bool2 value")
		t.Fail()
	}
	if v, _ := c.GetSetting("two").String(); v != "two words" {
		fmt.Println("failed finding key value")
		t.Fail()
	}
	if v, _ := c.GetSetting("commas").String(); v != "abc, def, ghi" {
		fmt.Println("failed finding key value")
		t.Fail()
	}
	if v, _ := c.GetSetting("nonexistent").String(); v != "" {
		fmt.Println("returned non-empty string for nonexistent key")
		t.Fail()
	}
	if c.HasOption("removed") || c.HasSetting("removed") {
		fmt.Println("commented out line shows up in config")
		t.Fail()
	}
	v := c.GetSetting("commas")
	splitted := v.Split()
	if len(splitted) != 3 {
		fmt.Println("could not split Option")
		t.Fail()
	}
	abc, _ := splitted[0].String()
	ghi, _ := splitted[2].String()
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
}
