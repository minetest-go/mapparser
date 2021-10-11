package mapparser

import (
	"encoding/binary"
	"errors"
	"strconv"
)

var ErrNoData = errors.New("no data")
var ErrMapblockVersion = errors.New("mapblock version unsupported")

type MapdataCompressionType int

const (
	MapdataZlibCompression MapdataCompressionType = 1
	MapdataZstdCompression MapdataCompressionType = 2
)

func Parse(data []byte) (*MapBlock, error) {
	if len(data) == 0 {
		return nil, ErrNoData
	}

	mapblock := NewMapblock()
	mapblock.Size = len(data)

	// version
	mapblock.Version = data[0]

	if mapblock.Version < 25 || mapblock.Version > 29 {
		return mapblock, ErrMapblockVersion
	}

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

	if mapblock.Version <= 28 {
		// zlib compressed mapblock

		content_width := data[offset]
		params_width := data[offset+1]

		if content_width != 2 {
			return nil, errors.New("content_width = " + strconv.Itoa(int(content_width)))
		}

		if params_width != 2 {
			return nil, errors.New("params_width = " + strconv.Itoa(int(params_width)))
		}

		//mapdata (blocks)
		if mapblock.Version >= 27 {
			offset = 6
		} else {
			offset = 4
		}

		//metadata
		err := parseMapdata(data[offset:], &offset, mapblock, MapdataZlibCompression)
		if err != nil {
			return nil, err
		}

		// metadata
		err = parseMetadata(data[offset:], &offset, mapblock)
		if err != nil {
			return nil, err
		}

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
		mapblock.Timestamp = binary.BigEndian.Uint32(data[offset:])
		offset += 4

		return mapblock, nil
	}

	return nil, nil
}
