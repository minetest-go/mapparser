package mapparser

import (
	"encoding/binary"
	"errors"
	"fmt"

	"github.com/minetest-go/types"
)

var ErrNoData = errors.New("no data")
var ErrMapblockVersion = errors.New("mapblock version unsupported")

func Parse(data []byte) (*types.MapBlock, error) {
	if len(data) == 0 {
		return nil, ErrNoData
	}

	mapblock := types.NewMapblock()
	mapblock.Size = len(data)

	// version
	mapblock.Version = data[0]

	if mapblock.Version < 25 || mapblock.Version > 29 {
		return mapblock, ErrMapblockVersion
	}

	if mapblock.Version <= 28 {
		// zlib compressed mapblock

		//flags
		flags := data[1]
		mapblock.Underground = (flags & 0x01) == 0x01

		var offset int

		if mapblock.Version >= 27 {
			offset = 4
		} else {
			//u16 lighting_complete not present
			offset = 2
		}

		content_width := data[offset]
		params_width := data[offset+1]

		if content_width != 2 {
			return nil, fmt.Errorf("content_width unexpected: %d", content_width)
		}

		if params_width != 2 {
			return nil, fmt.Errorf("params_width unexpected: %d", params_width)
		}

		// mapdata offset
		if mapblock.Version >= 27 {
			offset = 6
		} else {
			offset = 4
		}

		// mapdata
		mapdata, skip_bytes, err := decompress_zlib(data[offset:])
		if err != nil {
			return nil, err
		}
		err = parseMapdata(mapdata, mapblock)
		if err != nil {
			return nil, err
		}
		offset += skip_bytes

		// metadata
		metadata, skip_bytes, err := decompress_zlib(data[offset:])
		if err != nil {
			return nil, err
		}
		err = parseMetadata(metadata, mapblock)
		if err != nil {
			return nil, err
		}
		offset += skip_bytes

		//static objects
		parseStaticObjects(data, &offset)

		//timestamp
		offset += 4

		//mapping version
		offset++
		parseBlockMapping(data, &offset, mapblock)

		return mapblock, nil
	}

	if mapblock.Version >= 29 {
		uncompressed_data, _, err := decompress_zstd(data[1:])
		if err != nil {
			return nil, err
		}

		offset := 0

		//flags
		flags := uncompressed_data[1]
		mapblock.Underground = (flags & 0x01) == 0x01
		offset++

		//lighting complete
		offset += 2

		mapblock.Timestamp = binary.BigEndian.Uint32(uncompressed_data[offset:])
		offset += 4

		//mapping version
		offset++
		parseBlockMapping(uncompressed_data, &offset, mapblock)

		// content width * 2
		offset += 2

		// mapdata
		err = parseMapdata(uncompressed_data[offset:], mapblock)
		if err != nil {
			return nil, err
		}
		offset += MapDataSize

		// metadata
		err = parseMetadata(uncompressed_data[offset:], mapblock)
		if err != nil {
			return nil, err
		}
		return mapblock, nil
	}

	return nil, nil
}
