package mapparser

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
	"testing"
)

func TestReadU16(t *testing.T) {
	v := readU16([]byte{0x00, 0x00}, 0)
	if v != 0 {
		t.Error(v)
	}

	v = readU16([]byte{0x00, 0x01}, 0)
	if v != 1 {
		t.Error(v)
	}

	v = readU16([]byte{0x01, 0x00}, 0)
	if v != 256 {
		t.Error(v)
	}

}
func TestReadU32(t *testing.T) {
	v := readU32([]byte{0x00, 0x00, 0x00, 0x00}, 0)
	if v != 0 {
		t.Error(v)
	}
}

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
	fmt.Println(string(str))
}
