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

func getVpcPeeringConnections(b *services.Broker) ([]*types.Node, error) {
	nodes := []*types.Node{}
	res, err := b.EC2().DescribeVpcPeeringConnections(&ec2.DescribeVpcPeeringConnectionsInput{})
	if err != nil {
		return nodes, err
	}

	var errs error
	for _, pcx := range res.VpcPeeringConnections {
		node, err := mapPcx(pcx)
		if err != nil {
			errs = multierror.Append(errs, err)
		}
		nodes = append(nodes, node)
	}

	return nodes, errs
}

func mapPcx(pcx *ec2.VpcPeeringConnection) (*types.Node, error) {
	node, err := types.NewNode(*pcx.VpcPeeringConnectionId, "", types.NodeTypeVpcPeeringConnection)
	if err != nil {
		return nil, err
	}

	avi := pcx.AccepterVpcInfo
	if avi != nil {
		node.Metadata["accepter_vpc_id"] = *avi.VpcId
		node.Metadata["accepter_vpc_cidr"] = *avi.CidrBlock
		node.Metadata["accepter_owner_id"] = *avi.OwnerId
	}

	rvi := pcx.RequesterVpcInfo
	if rvi != nil {
		node.Metadata["requester_vpc_id"] = *rvi.VpcId
		node.Metadata["requester_vpc_cidr"] = *rvi.CidrBlock
		node.Metadata["requester_owner_id"] = *rvi.OwnerId
	}

	if pcx.Status != nil {
		node.Metadata["status_code"] = pcx.Status.Code
		node.Metadata["status_message"] = pcx.Status.Message
	}

	return node, nil
}
