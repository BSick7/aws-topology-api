package api

import (
	"github.com/BSick7/aws-topology-api/services"
	"github.com/BSick7/aws-topology-api/types"
)

func GetVpcTopology(b *services.Broker, vpcId string) (*types.Node, error) {
	vpc, err := getVpcNode(b, vpcId)
	if err != nil {
		return nil, err
	}

	return vpc, nil
}
