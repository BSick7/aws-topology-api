package api

/// Collects security groups in a vpc
/// IAM Requirements:
///   - ec2:DescribeSecurityGroups

import (
	"fmt"
	"github.com/BSick7/aws-topology-api/services"
	"github.com/BSick7/aws-topology-api/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func getSecurityGroups(b *services.Broker, vpcId string) ([]*types.Resource, error) {
	all := []*types.Resource{}
	out, err := b.EC2().DescribeSecurityGroups(&ec2.DescribeSecurityGroupsInput{
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

	for _, grp := range out.SecurityGroups {
		all = append(all, mapSecurityGroup(b, grp))
	}
	return all, nil
}

func genSecurityGroupArn(region string, accountId string, secGroupId string) string {
	return fmt.Sprintf("arn:aws:ec2:%s:%s:security-group/%s", region, accountId, secGroupId)
}

func mapSecurityGroup(b *services.Broker, secGroup *ec2.SecurityGroup) *types.Resource {
	arn := genSecurityGroupArn(b.Region(), b.AccountId(), *secGroup.GroupId)
	resource := types.NewResource(*secGroup.GroupId, arn, types.ResourceTypeSecurityGroup)

	m := resource.Metadata

	m["vpc_id"] = *secGroup.VpcId
	m["name"] = *secGroup.GroupName
	m["description"] = *secGroup.Description

	ingressRules := []map[string]string{}
	for _, ipPerm := range secGroup.IpPermissions {
		ingressRules = append(ingressRules, mapIpPermission(ipPerm)...)
	}
	m["ingress_rules"] = ingressRules

	egressRules := []map[string]string{}
	for _, ipPerm := range secGroup.IpPermissionsEgress {
		egressRules = append(egressRules, mapIpPermission(ipPerm)...)
	}
	m["egress_rules"] = egressRules

	return resource
}

func mapIpPermission(ipPerm *ec2.IpPermission) []map[string]string {
	all := []map[string]string{}

	for _, prefix := range ipPerm.PrefixListIds {
		rule := map[string]string{
			"protocol":       safeStringValue(ipPerm.IpProtocol),
			"prefix_list_id": safeStringValue(prefix.PrefixListId),
		}
		mapIpPort(ipPerm, rule)
		all = append(all, rule)
	}

	for _, ipRange := range ipPerm.IpRanges {
		rule := map[string]string{
			"protocol": safeStringValue(ipPerm.IpProtocol),
			"cidr":     safeStringValue(ipRange.CidrIp),
		}
		mapIpPort(ipPerm, rule)
		all = append(all, rule)
	}

	for _, otherGrp := range ipPerm.UserIdGroupPairs {
		rule := map[string]string{
			"protocol":       safeStringValue(ipPerm.IpProtocol),
			"group_id":       safeStringValue(otherGrp.GroupId),
			"group_name":     safeStringValue(otherGrp.GroupName),
			"user_id":        safeStringValue(otherGrp.UserId),
			"vpc_id":         safeStringValue(otherGrp.VpcId),
			"peering_id":     safeStringValue(otherGrp.VpcPeeringConnectionId),
			"peering_status": safeStringValue(otherGrp.PeeringStatus),
		}
		mapIpPort(ipPerm, rule)
		all = append(all, rule)
	}

	return all
}

func mapIpPort(ipPerm *ec2.IpPermission, m map[string]string) {
	if ipPerm.FromPort != nil {
		m["from_port"] = fmt.Sprintf("%d", *ipPerm.FromPort)
	}
	if ipPerm.ToPort != nil {
		m["to_port"] = fmt.Sprintf("%d", *ipPerm.ToPort)
	}
}

func safeStringValue(p *string) string {
	if p == nil {
		return ""
	}
	return *p
}
