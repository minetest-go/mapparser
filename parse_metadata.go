package mapparser

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"strconv"
	"strings"
)

/*
lua vm: https://github.com/yuin/gopher-lua
*/

const (
	INVENTORY_TERMINATOR = "EndInventory"
	INVENTORY_END        = "EndInventoryList"
	INVENTORY_START      = "List"
)

func parseMetadata(metadata []byte, mapblock *MapBlock) error {
	offset := 0
	version := metadata[offset]

	if version == 0 {
		//No data?
		return nil
	}

	offset++
	count := int(binary.BigEndian.Uint16(metadata[offset:]))

	offset += 2

	for i := 0; i < count; i++ {
		position := int(binary.BigEndian.Uint16(metadata[offset:]))
		pairsMap := mapblock.Metadata.GetPairsMap(position)

		offset += 2
		valuecount := int(binary.BigEndian.Uint32(metadata[offset:]))

		offset += 4
		for j := 0; j < valuecount; j++ {
			keyLength := int(binary.BigEndian.Uint16(metadata[offset:]))
			offset += 2

			key := string(metadata[offset : keyLength+offset])
			offset += keyLength

			valueLength := int(binary.BigEndian.Uint32(metadata[offset:]))
			offset += 4

			if len(metadata) <= valueLength+offset {
				return errors.New("metadata too short: " + strconv.Itoa(len(metadata)) +
					", valuelength: " + strconv.Itoa(int(valueLength)))
			}

			value := string(metadata[offset : valueLength+offset])
			offset += valueLength

			pairsMap[key] = value

			if version >= 2 { /* private tag doesn't exist in version=1 */
				offset++
			}
		}

		var currentInventoryName *string
		var currentInventory *Inventory

		scanner := bufio.NewScanner(bytes.NewReader(metadata[offset:]))
		for scanner.Scan() {
			txt := scanner.Text()
			offset += len(txt) + 1

			if strings.HasPrefix(txt, INVENTORY_START) {
				pairs := strings.Split(txt, " ")
				currentInventoryName = &pairs[1]
				currentInventory = mapblock.Metadata.GetInventory(position, *currentInventoryName)
				currentInventory.Size = 0

			} else if txt == INVENTORY_END {
				currentInventoryName = nil
				currentInventory = nil
			} else if currentInventory != nil {
				//content
				if strings.HasPrefix(txt, "Item") {
					item := Item{}
					parts := strings.Split(txt, " ")

					if len(parts) >= 2 {
						item.Name = parts[1]
					}

					if len(parts) >= 3 {
						val, err := strconv.ParseInt(parts[2], 10, 32)
						if err != nil {
							return err
						}
						item.Count = int(val)
					}

					if len(parts) >= 4 {
						val, err := strconv.ParseInt(parts[3], 10, 32)
						if err != nil {
							return err
						}
						item.Count = int(val)
					}

					currentInventory.Items = append(currentInventory.Items, &item)
					currentInventory.Size += 1

				}

			} else if txt == INVENTORY_TERMINATOR {
				break

			} else {
				return errors.New("Malformed inventory: " + txt)
			}
		}

		//TODO

	}

	return nil
}
