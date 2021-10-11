package mapparser

import "encoding/binary"

func parseBlockMapping(data []byte, offset *int, mapblock *MapBlock) {
	if len(data) > (*offset + 2) {
		// disk-data has per-block mapping, network-data has a global mapping
		numMappings := int(binary.BigEndian.Uint16(data[*offset:]))
		*offset += 2
		for i := 0; i < numMappings; i++ {
			nodeId := int(binary.BigEndian.Uint16(data[*offset:]))
			*offset += 2

			nameLen := int(binary.BigEndian.Uint16(data[*offset:]))
			*offset += 2

			blockName := string(data[*offset : *offset+nameLen])
			*offset += nameLen

			mapblock.BlockMapping[nodeId] = blockName
		}
	}
}
