package types

type Topology struct {
	Vpc       *Resource   `json:"vpc"`
	Resources []*Resource `json:"resources"`
}
