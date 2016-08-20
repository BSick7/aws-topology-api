package api

/// Collects ec2 instances in a vpc
/// IAM Requirements:
///   - ec2:DescribeInstancesPages

import (
	"github.com/BSick7/aws-topology-api/services"
	"github.com/BSick7/aws-topology-api/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/hashicorp/go-multierror"
)

func getInstances(b *services.Broker, vpcId string) ([]*types.Resource, error) {
	all := []*types.Resource{}
	input := &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			&ec2.Filter{
				Name:   aws.String("vpc-id"),
				Values: []*string{aws.String(vpcId)},
			},
		},
	}
	var errs error
	mapInstanceFn := func(p *ec2.DescribeInstancesOutput, lastPage bool) bool {
		for _, reservation := range p.Reservations {
			for _, instance := range reservation.Instances {
				resource, err := mapInstance(instance)
				if err != nil {
					errs = multierror.Append(errs, err)
				}
				all = append(all, resource)
			}
		}
		return true
	}
	err := b.EC2().DescribeInstancesPages(input, mapInstanceFn)
	if err != nil {
		errs = multierror.Append(errs, err)
	}
	return all, errs
}

func mapInstance(instance *ec2.Instance) (*types.Resource, error) {
	resource, err := types.NewResource(*instance.InstanceId, "", types.ResourceTypeInstance)
	if err != nil {
		return nil, err
	}

	m := resource.Metadata
	m["vpc_id"] = *instance.VpcId
	m["subnet_id"] = *instance.SubnetId
	m["image_id"] = *instance.ImageId
	m["instance_type"] = *instance.InstanceType
	m["public_ip"] = *instance.PublicIpAddress
	m["public_dns"] = *instance.PublicDnsName
	m["private_ip"] = *instance.PrivateIpAddress
	m["private_dns"] = *instance.PrivateDnsName
	m["virtualization_type"] = *instance.VirtualizationType

	if instance.Placement != nil {
		m["availability_zone"] = *instance.Placement.AvailabilityZone
	}
	if instance.State != nil {
		m["state_code"] = *instance.State.Code
		m["state_name"] = *instance.State.Name
	}
	if instance.IamInstanceProfile != nil {
		m["instance_profile_id"] = *instance.IamInstanceProfile.Id
		m["instance_profile_arn"] = *instance.IamInstanceProfile.Arn
	}

	securityGroups := []string{}
	if instance.SecurityGroups != nil && len(instance.SecurityGroups) > 0 {
		for _, sg := range instance.SecurityGroups {
			securityGroups = append(securityGroups, *sg.GroupId)
		}
	}
	m["security_groups"] = securityGroups

	return resource, nil
}
