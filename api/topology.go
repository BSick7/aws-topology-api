package api

import (
	"github.com/BSick7/aws-topology-api/services"
	"github.com/BSick7/aws-topology-api/types"
)

func GetVpcTopology(b *services.Broker, vpcId string) (types.Topology, error) {
	vpc, err := getVpcNode(b, vpcId)
	if err != nil {
		return types.Topology{}, err
	}
	if vpc == nil {
		return types.Topology{}, nil
	}

	resources := []*types.Resource{}

	pcxs, err := getVpcPeeringConnections(b)
	for _, pcx := range pcxs {
		if pcx.Metadata["accepter_vpc_id"] == vpc.Id || pcx.Metadata["requester_vpc_id"] == vpc.Id {
			resources = append(resources, pcx)
		}
	}

	return types.Topology{
		Vpc:       vpc,
		Resources: resources,
	}, nil
}
