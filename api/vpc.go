package api

/// Collects vpc
/// IAM Requirements:
///   - ec2:DescribeVpcs

import (
	"fmt"
	"github.com/BSick7/aws-topology-api/services"
	"github.com/BSick7/aws-topology-api/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func getVpcNode(b *services.Broker, vpcId string) (*types.Resource, error) {
	res, err := b.EC2().DescribeVpcs(&ec2.DescribeVpcsInput{
		VpcIds: []*string{aws.String(vpcId)},
	})
	if err != nil {
		if isVpcMissing(err) {
			return nil, nil
		}
		return nil, err
	}

	if len(res.Vpcs) < 1 {
		return nil, nil
	}

	vpc := res.Vpcs[0]
	arn := fmt.Sprintf("arn:aws:ec2:%s:%s:vpc/%s", b.Region(), b.AccountId(), *vpc.VpcId)
	node := types.NewResource(*vpc.VpcId, arn, types.ResourceTypeVpc)

	node.Metadata["cidr"] = *vpc.CidrBlock
	node.Metadata["tenancy"] = *vpc.InstanceTenancy

	return node, nil
}

func isVpcMissing(err error) bool {
	if aerr, ok := err.(awserr.Error); ok {
		return aerr.Code() == "InvalidVpcID.NotFound"
	}
	return false
}
