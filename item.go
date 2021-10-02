package mapparser

type Item struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
	Wear  int    `json:"wear"`
	//TODO: metadata
}

func (i *Item) IsEmpty() bool {
	return i.Name == "" && i.Count == 0
}
