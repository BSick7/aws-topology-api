package data

/// Collects ec2 instances in a vpc
/// IAM Requirements:
///   - ec2:DescribeSubnets

import (
	"fmt"
	"github.com/BSick7/aws-topology-api/services"
	"github.com/BSick7/aws-topology-api/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func getSubnets(b *services.Broker, vpcId string) ([]*types.Resource, error) {
	all := []*types.Resource{}
	out, err := b.EC2().DescribeSubnets(&ec2.DescribeSubnetsInput{
		Filters: []*ec2.Filter{
			&ec2.Filter{
				Name:   aws.String("vpc-id"),
				Values: []*string{aws.String(vpcId)},
			},
		},
	})
	if err != nil {
		return all, err
	}

	for _, subnet := range out.Subnets {
		all = append(all, mapSubnet(b, subnet))
	}
	return all, nil
}

func genSubnetArn(region string, accountId string, subnetId string) string {
	return fmt.Sprintf("arn:aws:ec2:%s:%s:subnet/%s", region, accountId, subnetId)
}

func mapSubnet(b *services.Broker, subnet *ec2.Subnet) *types.Resource {
	arn := genSubnetArn(b.Region(), b.AccountId(), *subnet.SubnetId)
	resource := types.NewResource(*subnet.SubnetId, arn, types.ResourceTypeSubnet)

	m := resource.Metadata
	m["vpc_id"] = *subnet.VpcId
	m["cidr"] = *subnet.CidrBlock
	m["availability_zone"] = *subnet.AvailabilityZone

	return resource
}
