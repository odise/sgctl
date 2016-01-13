package main

import (
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/spf13/cobra"
)

var region string

func main() {

	var instance string
	region = os.Getenv("AWS_REGION")

	var cmdAdd = &cobra.Command{
		Use:   "add <security group>",
		Short: "Add one or more security groups",
		Long: `Add one or more security groups to an EC2 instance. If the scurity
group is alraedy assigned it will be ignored.`,
		Run: func(cmd *cobra.Command, args []string) {
			if instance == "" {
				region = getRegion()
				instance = findInstanceID()
			}
			var new []*string
			for i := range args {
				new = append(new, aws.String(args[i]))
			}
			groups := getSg(instance)

			fmt.Printf("Adding %d new security groups to %d existing.\n", len(new), len(groups))
			modify(instance, generateSgSlice(groups, new))
		},
	}

	var cmdDelete = &cobra.Command{
		Use:   "del <security group>",
		Short: "Delete one or more security groups from an instance",
		Long: `Delete one or more security groups from an instance. In case the
security group is not assigned to the instance it will be ignored.`,
		Run: func(cmd *cobra.Command, args []string) {
			if instance == "" {
				region = getRegion()
				instance = findInstanceID()
			}
			var new []*string
			existing := getSg(instance)
			fmt.Printf("Removing %d security groups from %d existing.\n", len(args), len(existing))
			for i := range args {
				for j := range existing {
					if *existing[j].GroupId == args[i] {
						// remove it
						existing = append(existing[:j], existing[j+1:]...)
						break
					}
				}
			}
			modify(instance, generateSgSlice(existing, new))
		},
	}

	cmdAdd.Flags().StringVarP(&instance, "instance", "i", "", "EC2 instance identifier. If not defined this will be evaluated from the EC2 metadata.")
	cmdDelete.Flags().StringVarP(&instance, "instance", "i", "", "EC2 instance identifier. If not defined this will be evaluated from the EC2 metadata.")

	var rootCmd = &cobra.Command{Use: "app"}
	rootCmd.AddCommand(cmdAdd, cmdDelete)
	rootCmd.Execute()
}

// get all security groups from an instance
func getSg(instance string) []*ec2.GroupIdentifier {
	svc := ec2.New(session.New(), &aws.Config{Region: aws.String(region)})

	params := &ec2.DescribeInstanceAttributeInput{
		Attribute:  aws.String(ec2.InstanceAttributeNameGroupSet), // Required
		InstanceId: aws.String(instance),                          // Required
		DryRun:     aws.Bool(false),
	}
	resp, err := svc.DescribeInstanceAttribute(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return nil
	}
	return resp.Groups
}

// flatten the given array og ec2.GroupIdentifier
func generateSgSlice(existing []*ec2.GroupIdentifier, new []*string) []*string {
	var result []*string

	for i := range existing {
		result = append(result, existing[i].GroupId)
	}
	return append(new, result...)
}

// set the security groups of an EC2 instance
func modify(instance string, sg []*string) {
	svc := ec2.New(session.New(), &aws.Config{Region: aws.String(region)})

	params := &ec2.ModifyInstanceAttributeInput{
		InstanceId: aws.String(instance),
		Groups:     sg,
	}
	if _, err := svc.ModifyInstanceAttribute(params); err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		log.Fatal(err)
	}
}

func findInstanceID() string {
	metaClient := ec2metadata.New(session.New(&aws.Config{}))
	instanceID, err := metaClient.GetMetadata("instance-id")

	if err != nil {
		log.Fatal(err)
	}
	return instanceID
}

func getRegion() string {
	metaClient := ec2metadata.New(session.New(&aws.Config{}))
	region, err := metaClient.Region()
	if err != nil {
		log.Fatal(err)
	}
	return region
}
