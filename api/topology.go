package api

import (
	"github.com/BSick7/aws-topology-api/services"
	"github.com/BSick7/aws-topology-api/types"
	"github.com/hashicorp/go-multierror"
)

func GetVpcTopology(b *services.Broker, vpcId string) (types.Topology, error) {
	if err := b.Init(); err != nil {
		return types.Topology{}, err
	}

	vpc, err := getVpcNode(b, vpcId)
	if err != nil {
		return types.Topology{}, err
	}
	if vpc == nil {
		return types.Topology{}, nil
	}

	resources := []*types.Resource{}
	var errs error

	// Peering Connections
	pcxs, err := getPeeringConnections(b)
	if err != nil {
		errs = multierror.Append(errs, err)
	}
	for _, pcx := range pcxs {
		if pcx.Metadata["accepter_vpc_id"] == vpc.Id || pcx.Metadata["requester_vpc_id"] == vpc.Id {
			resources = append(resources, pcx)
		}
	}

	// Instances
	instances, err := getInstances(b, vpc.Id)
	if err != nil {
		errs = multierror.Append(errs, err)
	}
	resources = append(resources, instances...)

	// Subnets
	subnets, err := getSubnets(b, vpc.Id)
	if err != nil {
		errs = multierror.Append(errs, err)
	}
	resources = append(resources, subnets...)

	// Security Groups
	secGroups, err := getSecurityGroups(b, vpc.Id)
	if err != nil {
		errs = multierror.Append(errs, err)
	}
	resources = append(resources, secGroups...)

	return types.Topology{
		Vpc:       vpc,
		Resources: resources,
	}, errs
}
