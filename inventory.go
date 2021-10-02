package mapparser

type Inventory struct {
	Size  int     `json:"size"`
	Items []*Item `json:"items"`
}

func (inv *Inventory) IsEmpty() bool {
	if len(inv.Items) == 0 {
		return true
	}

	for _, item := range inv.Items {
		if item.Name != "" && item.Count > 0 {
			return false
		}
	}

	return true
}
