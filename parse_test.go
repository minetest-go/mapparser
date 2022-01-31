package mapparser

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
	"testing"
)

func validateMapblock(t *testing.T, mapblock *MapBlock) {
	if mapblock == nil {
		t.Error("no data")
		return
	}

	if mapblock.Mapdata == nil {
		t.Error("mapdata is nil")
		return
	}

	if mapblock.Mapdata.ContentId == nil {
		t.Error("contentid is nil")
		return
	}

	if mapblock.Mapdata.Param1 == nil {
		t.Error("param1 is nil")
		return
	}

	if mapblock.Mapdata.Param2 == nil {
		t.Error("param2 is nil")
		return
	}

	if len(mapblock.Mapdata.ContentId) != 4096 {
		t.Error("invalid contentid size")
	}

	if len(mapblock.Mapdata.Param1) != 4096 {
		t.Error("invalid param1 size")
	}

	if len(mapblock.Mapdata.Param2) != 4096 {
		t.Error("invalid param2 size")
	}

	for _, nodeid := range mapblock.Mapdata.ContentId {
		nodename := mapblock.BlockMapping[nodeid]
		if nodename == "" {
			t.Error(fmt.Sprintf("Nodename not found for id: %d", nodeid))
		}
	}

	/*
		for i, param1 := range mapblock.Mapdata.Param1 {
			if param1 > 15 {
				t.Error(fmt.Sprintf("Invalid param1: %d @ %d", param1, i))
			}
		}
	*/

}

func TestParse(t *testing.T) {

	data, err := ioutil.ReadFile("testdata/0.0.0")
	if err != nil {
		t.Error(err)
	}

	mapblock, err := Parse(data)

	validateMapblock(t, mapblock)

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

	json, err := json.Marshal(mapblock)
	if err != nil {
		t.Error(err)
	}

	fmt.Println(string(json))
}

func TestParseError(t *testing.T) {
	data, err := Parse([]byte{})

	if data != nil {
		t.Error("data is set")
	}

	if err == nil {
		t.Error("error expected")
	}

	if err != ErrNoData {
		t.Error("wrong error")
	}
}

func TestParseError2(t *testing.T) {
	data, err := Parse([]byte{24})

	if data == nil {
		t.Error("data is not set")
		return
	}

	if data.Version != 24 {
		t.Error("mapblock version wrong")
	}

	if err == nil {
		t.Error("error expected")
	}

	if err != ErrMapblockVersion {
		t.Error("wrong error")
	}
}

func TestParseZstd(t *testing.T) {

	data, err := ioutil.ReadFile("testdata/zstd-block.bin")
	if err != nil {
		t.Error(err)
	}

	mapblock, err := Parse(data)

	if err != nil {
		t.Error(err)
	}

	fmt.Println(mapblock)

	validateMapblock(t, mapblock)
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

	validateMapblock(t, mapblock)

	if mapblock.IsEmpty() {
		t.Error("mapblock empty")
	}

	for id, name := range mapblock.BlockMapping {
		fmt.Printf("%d = %s\n", id, name)
	}
}

func TestParse3(t *testing.T) {

	data, err := ioutil.ReadFile("testdata/0.1.0")
	if err != nil {
		t.Error(err)
	}

	mapblock, err := Parse(data)

	if err != nil {
		t.Error(err)
	}

	validateMapblock(t, mapblock)
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

func TestParseNetworkBlock(t *testing.T) {
	data, err := ioutil.ReadFile("testdata/network-blockdata.bin")
	if err != nil {
		t.Error(err)
	}

	offsetData := make([]byte, len(data)+1)
	offsetData[0] = 28
	copy(offsetData[1:], data)

	mb, err := Parse(offsetData)

	if err != nil {
		t.Error(err)
	}

	str, err := json.MarshalIndent(mb, "", "	")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(string(str))
}
