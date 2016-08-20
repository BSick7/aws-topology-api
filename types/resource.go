package types

type Resource struct {
	Id           string                 `json:"id"`
	Arn          string                 `json:"arn"`
	Type         string                 `json:"type"`
	Metadata     map[string]interface{} `json:"metadata"`
	LinkUids     []string               `json:"linkUids"`
	ChildrenUids []string               `json:"childrenUids"`
}

func NewResource(awsId string, arn string, ntype string) *Resource {
	return &Resource{
		Arn:          arn,
		Id:           awsId,
		Type:         ntype,
		Metadata:     map[string]interface{}{},
		LinkUids:     []string{},
		ChildrenUids: []string{},
	}
}
