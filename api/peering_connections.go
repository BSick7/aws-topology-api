package api

/// Collects vpc peering connections for a vpc
/// IAM Requirements:
///   - ec2:DescribeVpcPeeringConnections

import (
	"fmt"
	"github.com/BSick7/aws-topology-api/services"
	"github.com/BSick7/aws-topology-api/types"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func getVpcPeeringConnections(b *services.Broker) ([]*types.Resource, error) {
	resources := []*types.Resource{}
	res, err := b.EC2().DescribeVpcPeeringConnections(&ec2.DescribeVpcPeeringConnectionsInput{})
	if err != nil {
		return resources, err
	}

	var errs error
	for _, pcx := range res.VpcPeeringConnections {
		resources = append(resources, mapPcx(b, pcx))
	}

	return resources, errs
}

func mapPcx(b *services.Broker, pcx *ec2.VpcPeeringConnection) *types.Resource {
	arn := fmt.Sprintf("arn:aws:ec2:%s:%s:vpc-peering-connection/%s", b.Region(), b.AccountId(), *pcx.VpcPeeringConnectionId)
	resource := types.NewResource(*pcx.VpcPeeringConnectionId, arn, types.ResourceTypeVpcPeeringConnection)

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

	return resource
}
