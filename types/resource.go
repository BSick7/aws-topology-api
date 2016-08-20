package types

import "github.com/hashicorp/go-uuid"

type Resource struct {
	Uid          string                 `json:"uid"`
	Id           string                 `json:"id"`
	Arn          string                 `json:"arn"`
	Type         string                 `json:"type"`
	Metadata     map[string]interface{} `json:"metadata"`
	LinkUids     []string               `json:"linkUids"`
	ChildrenUids []string               `json:"childrenUids"`
}

func NewResource(awsId string, arn string, ntype string) (*Resource, error) {
	uid, err := uuid.GenerateUUID()
	return &Resource{
		Uid:          uid,
		Id:           awsId,
		Arn:          arn,
		Type:         ntype,
		Metadata:     map[string]interface{}{},
		LinkUids:     []string{},
		ChildrenUids: []string{},
	}, err
}
