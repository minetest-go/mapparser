package mapparser

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"strconv"
	"strings"

	"github.com/minetest-go/types"
)

/*
lua vm: https://github.com/yuin/gopher-lua
*/

const (
	INVENTORY_TERMINATOR = "EndInventory"
	INVENTORY_END        = "EndInventoryList"
	INVENTORY_START      = "List"
	ITEM_PREFIX          = "Item"
	ITEM_EMPTY           = "Empty"
)

func parseMetadata(metadata []byte, mapblock *types.MapBlock) error {
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
		fields := map[string]string{}
		mapblock.Fields[position] = fields

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

			fields[key] = value

			if version >= 2 { /* private tag doesn't exist in version=1 */
				offset++
			}
		}

		var currentInventoryName *string
		var currentInventory []string

		scanner := bufio.NewScanner(bytes.NewReader(metadata[offset:]))
		for scanner.Scan() {
			txt := scanner.Text()
			offset += len(txt) + 1

			if strings.HasPrefix(txt, INVENTORY_START) {
				pairs := strings.Split(txt, " ")
				currentInventoryName = &pairs[1]
				currentInventory = make([]string, 0)
				inv := mapblock.Inventory[position]
				if inv == nil {
					inv = map[string][]string{}
					mapblock.Inventory[position] = inv
				}
				inv[*currentInventoryName] = currentInventory

			} else if txt == INVENTORY_END {
				currentInventoryName = nil
				currentInventory = nil
			} else if currentInventory != nil {
				//content
				inv := mapblock.Inventory[position]
				if strings.HasPrefix(txt, ITEM_PREFIX) {
					item := txt[len(ITEM_PREFIX)+1:]
					currentInventory = append(currentInventory, item)
				} else if txt == ITEM_EMPTY {
					currentInventory = append(currentInventory, "")
				}
				inv[*currentInventoryName] = currentInventory

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
