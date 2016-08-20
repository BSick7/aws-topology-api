package api

/// Collects vpc peering connections for a vpc
/// IAM Requirements:
///   - ec2:DescribeVpcPeeringConnections

import (
	"github.com/BSick7/aws-topology-api/services"
	"github.com/BSick7/aws-topology-api/types"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/hashicorp/go-multierror"
)

func getVpcPeeringConnections(b *services.Broker) ([]*types.Resource, error) {
	resources := []*types.Resource{}
	res, err := b.EC2().DescribeVpcPeeringConnections(&ec2.DescribeVpcPeeringConnectionsInput{})
	if err != nil {
		return resources, err
	}

	var errs error
	for _, pcx := range res.VpcPeeringConnections {
		resource, err := mapPcx(pcx)
		if err != nil {
			errs = multierror.Append(errs, err)
		}
		resources = append(resources, resource)
	}

	return resources, errs
}

func mapPcx(pcx *ec2.VpcPeeringConnection) (*types.Resource, error) {
	resource, err := types.NewResource(*pcx.VpcPeeringConnectionId, "", types.ResourceTypeVpcPeeringConnection)
	if err != nil {
		return nil, err
	}

	avi := pcx.AccepterVpcInfo
	if avi != nil {
		resource.Metadata["accepter_vpc_id"] = *avi.VpcId
		resource.Metadata["accepter_vpc_cidr"] = *avi.CidrBlock
		resource.Metadata["accepter_owner_id"] = *avi.OwnerId
	}

	rvi := pcx.RequesterVpcInfo
	if rvi != nil {
		resource.Metadata["requester_vpc_id"] = *rvi.VpcId
		resource.Metadata["requester_vpc_cidr"] = *rvi.CidrBlock
		resource.Metadata["requester_owner_id"] = *rvi.OwnerId
	}

	if pcx.Status != nil {
		resource.Metadata["status_code"] = *pcx.Status.Code
		resource.Metadata["status_message"] = *pcx.Status.Message
	}

	return resource, nil
}
