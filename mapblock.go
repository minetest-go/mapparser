package mapparser

type MapBlock struct {
	Size         int            `json:"size"`
	Version      byte           `json:"version"`
	Underground  bool           `json:"underground"`
	Timestamp    uint32         `json:"timestamp"`
	Mapdata      *MapData       `json:"mapdata"`
	Metadata     *Metadata      `json:"metadata"`
	BlockMapping map[int]string `json:"blockmapping"`
}

func NewMapblock() *MapBlock {
	mb := MapBlock{}
	mb.Metadata = NewMetadata()
	mb.BlockMapping = make(map[int]string)
	return &mb
}

func (mb *MapBlock) IsEmpty() bool {
	if len(mb.BlockMapping) == 0 {
		// only air
		return true
	}

	if len(mb.BlockMapping) == 1 {
		for _, name := range mb.BlockMapping {
			if name == "vacuum:vacuum" {
				// only vacuum
				return true
			}
		}
	}

	// other stuff
	return false
}

func (mb *MapBlock) GetNodeId(x, y, z int) int {
	pos := GetNodePos(x, y, z)
	return mb.Mapdata.ContentId[pos]
}

func (mb *MapBlock) GetParam2(x, y, z int) int {
	pos := GetNodePos(x, y, z)
	return mb.Mapdata.Param2[pos]
}

func (mb *MapBlock) GetNodeName(x, y, z int) string {
	id := mb.GetNodeId(x, y, z)
	return mb.BlockMapping[id]
}
