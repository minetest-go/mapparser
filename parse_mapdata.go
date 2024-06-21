package mapparser

import (
	"encoding/binary"
	"fmt"

	"github.com/minetest-go/types"
)

const MapDataSize = 16384

func parseMapdata(rawdata []byte, mapblock *types.MapBlock) error {
	if len(rawdata) < MapDataSize {
		return fmt.Errorf("mapdata length invalid: %d", len(rawdata))
	}

	mapd := types.MapData{
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
