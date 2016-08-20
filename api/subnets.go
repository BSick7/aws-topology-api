package api

/// Collects ec2 instances in a vpc
/// IAM Requirements:
///   - ec2:DescribeSubnets

import (
	"github.com/BSick7/aws-topology-api/services"
	"github.com/BSick7/aws-topology-api/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/hashicorp/go-multierror"
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

	var errs error
	for _, subnet := range out.Subnets {
		resource, err := mapSubnet(subnet)
		if err != nil {
			errs = multierror.Append(errs, err)
		}
		all = append(all, resource)
	}
	return all, errs
}

func mapSubnet(subnet *ec2.Subnet) (*types.Resource, error) {
	resource, err := types.NewResource(*subnet.SubnetId, "", types.ResourceTypeSubnet)
	if err != nil {
		return nil, err
	}

	m := resource.Metadata
	m["vpc_id"] = *subnet.VpcId
	m["cidr"] = *subnet.CidrBlock
	m["availability_zone"] = *subnet.AvailabilityZone

	return resource, nil
}
