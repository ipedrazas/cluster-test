package cmd

import (
	"fmt"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

var svc *ec2.EC2

func initAWSSession() {
	sess := session.Must(session.NewSession())
	svc = ec2.New(sess, &aws.Config{Region: aws.String(awsRegion)})

}

// GetMasters receives a parameter (clusterName) that is used to query
// AWS to get the instanceID.
func GetMasters(clusterName string) []string {
	initAWSSession()
	resources := getInstances(clusterName, "masters")
	return resources
}

// GetMasters receives a parameter (clusterName) that is used to query
// AWS to get the instanceID.
func GetNodes(clusterName string) []string {
	initAWSSession()
	resources := getInstances(clusterName, "nodes")
	return resources
}

func getInstances(clusterName string, role string) []string {
	fmt.Printf("listing instances with tag %v in %v: \n", role, clusterName)
	params := &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			{
				Name: aws.String("tag:Name"),
				Values: []*string{
					aws.String(strings.Join([]string{"*", clusterName, "*"}, "")),
				},
			},
			{
				Name: aws.String("tag:Name"),
				Values: []*string{
					aws.String(strings.Join([]string{"*", role, "*"}, "")),
				},
			},
		},
	}
	resp, err := svc.DescribeInstances(params)
	if err != nil {
		fmt.Println("there was an error listing instances in", err.Error())
		log.Fatal(err.Error())
	}
	var resources []string
	for _, reservation := range resp.Reservations {
		for _, instance := range reservation.Instances {
			// fmt.Println(*instance.InstanceId)
			resources = append(resources, *instance.InstanceId)
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
