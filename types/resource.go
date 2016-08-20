package types

import "github.com/hashicorp/go-uuid"

type Node struct {
	Uid          string            `json:"uid"`
	Id           string            `json:"id"`
	Arn          string            `json:"arn"`
	Type         string            `json:"type"`
	Metadata     map[string]string `json:"metadata"`
	LinkUids     []string          `json:"linkUids"`
	ChildrenUids []string          `json:"childrenUids"`
}

func NewNode(awsId string, arn string, ntype string) (*Node, error) {
	uid, err := uuid.GenerateUUID()
	return &Node{
		Uid:          uid,
		Id:           awsId,
		Arn:          arn,
		Type:         ntype,
		Metadata:     map[string]interface{}{},
		LinkUids:     []string{},
		ChildrenUids: []string{},
	}, err
}
