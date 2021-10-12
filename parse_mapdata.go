package mapparser

import (
	"encoding/binary"
	"fmt"
)

const MapDataSize = 16384

func parseMapdata(rawdata []byte, mapblock *MapBlock) error {
	if len(rawdata) < MapDataSize {
		return fmt.Errorf("mapdata length invalid: %d", len(rawdata))
	}

	mapd := MapData{
		ContentId: make([]int, 4096),
		Param1:    make([]int, 4096),
		Param2:    make([]int, 4096),
	}
	mapblock.Mapdata = &mapd

	for i := 0; i < 4096; i++ {
		mapd.ContentId[i] = int(binary.BigEndian.Uint16(rawdata[i*2:]))
		mapd.Param1[i] = int(rawdata[(4096*2)+i])
		mapd.Param2[i] = int(rawdata[(4096*3)+i])
	}

	return nil
}
