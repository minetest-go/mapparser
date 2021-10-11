package mapparser

import "encoding/binary"

func parseStaticObjects(data []byte, offset *int) {
	*offset++ //static objects version
	staticObjectsCount := int(binary.BigEndian.Uint16(data[*offset:]))
	*offset += 2
	for i := 0; i < staticObjectsCount; i++ {
		*offset += 13
		dataSize := int(binary.BigEndian.Uint16(data[*offset:]))
		*offset += dataSize + 2
	}
}
