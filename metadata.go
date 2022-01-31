package mapparser

type Metadata struct {
	Inventories map[int]map[string]*Inventory `json:"inventories"`
	Pairs       map[int]map[string]string     `json:"pairs"`
}

func NewMetadata() *Metadata {
	md := Metadata{}
	md.Inventories = make(map[int]map[string]*Inventory)
	md.Pairs = make(map[int]map[string]string)
	return &md
}

func (md *Metadata) GetMetadata(x, y, z int) map[string]string {
	return md.GetPairsMap(GetNodePos(x, y, z))
}

func (md *Metadata) GetPairsMap(pos int) map[string]string {
	pairsMap := md.Pairs[pos]
	if pairsMap == nil {
		pairsMap = make(map[string]string)
		md.Pairs[pos] = pairsMap
	}

	return pairsMap
}

func (md *Metadata) GetInventoryMap(pos int) map[string]*Inventory {
	invMap := md.Inventories[pos]
	if invMap == nil {
		invMap = make(map[string]*Inventory)
		md.Inventories[pos] = invMap
	}

	return invMap
}

func (md *Metadata) GetInventoryMapAtPos(x, y, z int) map[string]*Inventory {
	return md.GetInventoryMap(GetNodePos(x, y, z))
}

func (md *Metadata) GetInventory(pos int, name string) *Inventory {
	m := md.GetInventoryMap(pos)
	inv := m[name]
	if inv == nil {
		inv = &Inventory{}
		m[name] = inv
	}

	return inv
}
