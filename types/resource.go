package types

type Resource struct {
	Id           string                 `json:"id"`
	Arn          string                 `json:"arn"`
	Type         string                 `json:"type"`
	Metadata     map[string]interface{} `json:"metadata"`
	LinkArns     []string               `json:"linkArns"`
	ChildrenArns []string               `json:"childrenArns"`
}

func NewResource(awsId string, arn string, ntype string) *Resource {
	return &Resource{
		Arn:          arn,
		Id:           awsId,
		Type:         ntype,
		Metadata:     map[string]interface{}{},
		LinkArns:     []string{},
		ChildrenArns: []string{},
	}
}
