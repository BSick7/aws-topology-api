package data

/// Collects ec2 instances in a vpc
/// IAM Requirements:
///   - ec2:DescribeInstancesPages

import (
	"fmt"
	"github.com/BSick7/aws-topology-api/services"
	"github.com/BSick7/aws-topology-api/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
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
	mapInstanceFn := func(p *ec2.DescribeInstancesOutput, lastPage bool) bool {
		for _, reservation := range p.Reservations {
			for _, instance := range reservation.Instances {
				all = append(all, mapInstance(b, instance))
			}
		}
		return true
	}
	err := b.EC2().DescribeInstancesPages(input, mapInstanceFn)
	return all, err
}

func genInstanceArn(region string, accountId string, instanceId string) string {
	return fmt.Sprintf("arn:aws:ec2:%s:%s:instance/%s", region, accountId, instanceId)
}

func mapInstance(b *services.Broker, instance *ec2.Instance) *types.Resource {
	arn := genInstanceArn(b.Region(), b.AccountId(), *instance.InstanceId)
	resource := types.NewResource(*instance.InstanceId, arn, types.ResourceTypeInstance)

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

	return resource
}
