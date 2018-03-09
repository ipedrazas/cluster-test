package cmd

import (
	"fmt"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/route53"
)

var svc *ec2.EC2
var svcR53 *route53.Route53

func initAWSSession() {
	sess := session.Must(session.NewSession())
	svc = ec2.New(sess, &aws.Config{Region: aws.String(awsRegion)})
	svcR53 = route53.New(sess, &aws.Config{Region: aws.String(awsRegion)})

}

// GetAllInstances returns EC2 instances
func GetAllInstances() []Instance {
	initAWSSession()
	filter := []*ec2.Filter{
		{
			Name: aws.String("tag:KubernetesCluster"),
			Values: []*string{
				aws.String(strings.Join([]string{"*", cluster, "*"}, "")),
			},
		},
	}
	return getInstances(filter)
}

// GetRunningMasters get all the masters and filters out the ones that are not running
func GetRunningInstances(ins []Instance) []Instance {
	var rIns []Instance
	for _, m := range ins {
		if m.State != "Running" {
			rIns = append(rIns, m)
		}
	}
	return rIns
}

// GetMasters uses clusterName to query
// AWS to get the instanceID.
func GetMasters() []Instance {
	initAWSSession()
	filter := []*ec2.Filter{
		{
			Name: aws.String("tag:KubernetesCluster"),
			Values: []*string{
				aws.String(strings.Join([]string{"*", cluster, "*"}, "")),
			},
		},
		{
			Name: aws.String("tag:k8s.io/role/master"),
			Values: []*string{
				aws.String(strings.Join([]string{"*", "1", "*"}, "")),
			},
		},
	}

	resources := getInstances(filter)
	return resources
}

// GetNodes receives a parameter (clusterName) that is used to query
// AWS to get the instanceID.
func GetNodes() []Instance {
	initAWSSession()
	filter := []*ec2.Filter{
		{
			Name: aws.String("tag:KubernetesCluster"),
			Values: []*string{
				aws.String(strings.Join([]string{"*", cluster, "*"}, "")),
			},
		},
		{
			Name: aws.String("tag:k8s.io/role/node"),
			Values: []*string{
				aws.String(strings.Join([]string{"*", "1", "*"}, "")),
			},
		},
	}
	resources := getInstances(filter)
	return resources
}

func getInstances(filters []*ec2.Filter) []Instance {
	if debug {
		fmt.Println(filters)
	}
	params := &ec2.DescribeInstancesInput{
		Filters: filters,
	}
	resp, err := svc.DescribeInstances(params)
	if err != nil {
		fmt.Println("there was an error listing instances in", err.Error())
		log.Fatal(err.Error())
	}
	var resources []Instance
	for _, reservation := range resp.Reservations {

		for _, instance := range reservation.Instances {
			i := &Instance{
				ID:   *instance.InstanceId,
				Name: *instance.KeyName,

				LaunchTime: *instance.LaunchTime,
				State:      *instance.State.Name,
			}
			// check if instance has IPs
			if instance.PublicIpAddress != nil {
				i.PublicIP = *instance.PublicIpAddress
			}
			if instance.PrivateIpAddress != nil {
				i.PrivateIP = *instance.PrivateIpAddress
			}

			for _, t := range instance.Tags {
				if *t.Key == "k8s.io/role/master" {
					i.IsMaster = true
				}
			}

			resources = append(resources, *i)
		}
	}
	return resources
}

func deleteInstance(instanceID string) (string, error) {
	res := fmt.Sprintf("Deleting EC2 instance %v - dry-run: %v\n", instanceID, dryrun)
	if !dryrun {
		request := &ec2.TerminateInstancesInput{
			InstanceIds: []*string{&instanceID},
		}
		_, err := svc.TerminateInstances(request)
		if err != nil {
			return fmt.Sprintf("error deleting instance %q", instanceID), err
		}

	}

	return res, nil
}

func listCNAMES() ([]RR, error) {

	if debug {
		fmt.Printf("ZoneID: %v", zoneId)
	}
	listParams := &route53.ListResourceRecordSetsInput{
		HostedZoneId: aws.String(zoneId), // Required
	}
	respList, err := svcR53.ListResourceRecordSets(listParams)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return nil, err
	}
	var result []RR
	if debug {
		fmt.Println("All records:")
		fmt.Println(cluster)
	}

	for _, r := range respList.ResourceRecordSets {
		if strings.HasSuffix(*r.Name, cluster+".") {
			var ips []string
			entry := &RR{
				Name: *r.Name,
			}
			for _, rr := range r.ResourceRecords {
				ips = append(ips, *rr.Value)
			}
			entry.Ips = ips
			result = append(result, *entry)
		}
	}

	return result, nil
}
