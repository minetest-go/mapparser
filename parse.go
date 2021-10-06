package mapparser

import (
	"encoding/binary"
	"errors"
	"strconv"
)

func Parse(data []byte) (*MapBlock, error) {
	if len(data) == 0 {
		return nil, errors.New("no data")
	}

	mapblock := NewMapblock()
	mapblock.Size = len(data)

	// version
	mapblock.Version = data[0]

	if mapblock.Version < 25 || mapblock.Version > 28 {
		return nil, errors.New("mapblock-version not supported: " + strconv.Itoa(int(mapblock.Version)))
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
	count, err := parseMapdata(mapblock, data[offset:])
	if err != nil {
		return nil, err
	}

	offset += count

	count, err = parseMetadata(mapblock, data[offset:])
	if err != nil {
		return nil, err
	}

	offset += count

	//static objects

	offset++ //static objects version
	staticObjectsCount := int(binary.BigEndian.Uint16(data[offset:]))
	offset += 2
	for i := 0; i < staticObjectsCount; i++ {
		offset += 13
		dataSize := int(binary.BigEndian.Uint16(data[offset:]))
		offset += dataSize + 2
	}

	//timestamp
	offset += 4

	//mapping version
	offset++

	numMappings := int(binary.BigEndian.Uint16(data[offset:]))
	offset += 2
	for i := 0; i < numMappings; i++ {
		nodeId := int(binary.BigEndian.Uint16(data[offset:]))
		offset += 2

		nameLen := int(binary.BigEndian.Uint16(data[offset:]))
		offset += 2

		blockName := string(data[offset : offset+nameLen])
		offset += nameLen

		mapblock.BlockMapping[nodeId] = blockName
	}

	return mapblock, nil
}
