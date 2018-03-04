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

func getMasters(clusterName string) []string {
	resources := getInstances(clusterName, "masters")
	return resources
}

func getNodes(clusterName string) []string {
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

func deleteInstance(instanceID string) error {
	fmt.Printf("Deleting EC2 instance %v - dry-run: %v\n", instanceID, dryrun)
	if !dryrun {
		request := &ec2.TerminateInstancesInput{
			InstanceIds: []*string{&instanceID},
		}
		_, err := svc.TerminateInstances(request)
		if err != nil {
			return fmt.Errorf("error deleting instance %q: %v", instanceID, err)
		}
	}
	return nil
}
