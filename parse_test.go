package mapparser

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
	"testing"
)

func TestParse(t *testing.T) {

	data, err := ioutil.ReadFile("testdata/0.0.0")
	if err != nil {
		t.Error(err)
	}

	mapblock, err := Parse(data)

	if err != nil {
		t.Error(err)
	}

	if mapblock.IsEmpty() {
		t.Error("mapblock is empty")
	}

	if mapblock.Version != 28 {
		t.Error("wrong mapblock version: " + strconv.Itoa(int(mapblock.Version)))
	}

	if !mapblock.Underground {
		t.Error("Underground flag")
	}

	if len(mapblock.Mapdata.ContentId) != 4096 {
		t.Error("Mapdata length wrong")
	}

	if len(mapblock.Mapdata.Param2) != 4096 {
		t.Error("Mapdata length wrong")
	}

	if len(mapblock.Mapdata.Param1) != 4096 {
		t.Error("Mapdata length wrong")
	}

	pairs := mapblock.Metadata.GetPairsMap(0)
	if pairs["owner"] != "pipo" {
		t.Error(pairs["owner"])
	}

	if mapblock.GetNodeId(0, 0, 0) != 0 {
		t.Error("nodeid mismatch")
	}

	if mapblock.GetNodeName(0, 0, 0) != "travelnet:travelnet" {
		t.Error("nodename mismatch")
	}

	if mapblock.GetParam2(0, 0, 0) != 0 {
		t.Error("param2 mismatch")
	}

}

func TestParse2(t *testing.T) {

	data, err := ioutil.ReadFile("testdata/11.0.2")
	if err != nil {
		t.Error(err)
	}

	mapblock, err := Parse(data)

	if err != nil {
		t.Error(err)
	}

	if mapblock.IsEmpty() {
		t.Error("mapblock empty")
	}

	for k, v := range mapblock.BlockMapping {
		fmt.Println("Key", k, "Value", v)
	}
}

func TestParse3(t *testing.T) {

	data, err := ioutil.ReadFile("testdata/0.1.0")
	if err != nil {
		t.Error(err)
	}

	_, err = Parse(data)

	if err != nil {
		t.Error(err)
	}
}

func TestParseMetadata(t *testing.T) {

	data, err := ioutil.ReadFile("testdata/mb-with-metadata.bin")
	if err != nil {
		t.Error(err)
	}

	mb, err := Parse(data)

	if err != nil {
		t.Error(err)
	}

	str, err := json.MarshalIndent(mb, "", "	")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(string(str))
}
