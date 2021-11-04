package conf

import (
	"bytes"
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
		t.Log("Float64Slice accepting incorrect values")
		t.Fail()
	}
	if _, err := c.Get("floats").Float64Slice(); err != nil {
		t.Log("Float64Slice failing on correct values")
		t.Fail()
	}
	if v := c.Get("token").String(); v != "test" {
		t.Log("failed finding key value")
		t.Fail()
	}
	if v := c.Get("editor").String(); v != "vim" {
		t.Log("failed finding key value")
		t.Fail()
	}
	if v, _ := c.Get("port").Int(); v != 10000 {
		t.Log("failed finding int")
		t.Fail()
	}
	if v, _ := c.Get("distance").Float64(); v != 13.42 {
		t.Log("failed finding key value")
		t.Fail()
	}
	if c.HasOption("color") != true {
		t.Log("failed finding option")
		t.Fail()
	}
	if v, _ := c.Get("bool1").Bool(); v != true {
		t.Log("failed finding bool1 value")
		t.Fail()
	}
	if v, _ := c.Get("bool2").Bool(); v != true {
		t.Log("failed finding bool2 value")
		t.Fail()
	}
	if v := c.Get("two").String(); v != "two words" {
		t.Log("failed finding key value")
		t.Fail()
	}
	if v := c.Get("commas").String(); v != "abc, def, ghi" {
		t.Log("failed finding key value")
		t.Fail()
	}
	if v := c.Get("nonexistent").String(); v != "" {
		t.Log("returned non-empty string for nonexistent key")
		t.Fail()
	}
	if c.HasOption("removed") {
		t.Log("commented out line shows up in config")
		t.Fail()
	}
	splitted := c.Get("commas").StringSlice()
	if len(splitted) != 3 {
		t.Log("could not split string")
		t.Fail()
	}
	abc := splitted[0]
	ghi := splitted[2]
	if abc != "abc" || ghi != "ghi" {
		t.Log("Split() returned incorrect values")
		t.Fail()
	}
	if c.HasOption("no") != true {
		t.Log("should capture one option per line even if line holds two")
		t.Fail()
	}
	if c.HasOption("missedme") == true {
		t.Log("should only capture one option per line")
		t.Fail()
	}
	st := Setting{}
	if v, err := st.Float64(); v != 0.0 || err == nil {
		t.Log("empty string erroneously converts to float")
		t.Fail()
	}
	if def := c.GetDefault("non-existant-key", "myvalue"); def.Value != "myvalue" {
		t.Log("GetDefault fails to apply default value")
		t.Fail()
	}
}
