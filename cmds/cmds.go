package cmds

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/awsutil"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/codegangsta/cli"
	"github.com/codeskyblue/go-sh"
	"github.com/fatih/color"
)

func init() {
}

func Ls(c *cli.Context) {
	svc := ec2.New(session.New(), &aws.Config{Region: aws.String("us-east-1")})
	cyan := color.New(color.FgCyan).SprintFunc()
	orange := color.New(color.FgMagenta).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	params := &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			{ // Required
				Name: aws.String(fmt.Sprintf("tag:%s", c.String("t"))),
				Values: []*string{
					aws.String(c.String("v")), // Required
					// More values...
				},
			},
			// More values...
		},
	}
	resp, err := svc.DescribeInstances(params)

	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			// Generic AWS error with Code, Message, and original error (if any)
			fmt.Println(awsErr.Code(), awsErr.Message(), awsErr.OrigErr())
			if reqErr, ok := err.(awserr.RequestFailure); ok {
				// A service error occurred
				fmt.Println(reqErr.Code(), reqErr.Message(), reqErr.StatusCode(), reqErr.RequestID())
			}
		} else {
			// This case should never be hit, the SDK should always return an
			// error which satisfies the awserr.Error interface.
			fmt.Println(err.Error())
		}
	}

	for idx, _ := range resp.Reservations {
		for _, inst := range resp.Reservations[idx].Instances {
			for _, tag := range inst.Tags {
				if *tag.Key == c.String("t") {
					fmt.Printf(fmt.Sprintf("%s: ", cyan(*tag.Value)))
					if c.Bool("d") {
						fmt.Println(awsutil.Prettify(*inst))
					}
					fmt.Printf(*inst.InstanceId)
					fmt.Printf(yellow(fmt.Sprintf(" %s ", *inst.PublicDnsName)))
					fmt.Println(orange(fmt.Sprintf("%s", *inst.PrivateIpAddress)))
				}
			}
		}
	}
}

func Ssh(c *cli.Context) {
	if c.Bool("d") {
		fmt.Println("tag: ", c.String("t"))
		fmt.Println("value: ", c.String("v"))
		fmt.Println("args: ", c.Args()[0])
		fmt.Println("user: ", c.String("u"))
	}
	svc := ec2.New(session.New(), &aws.Config{Region: aws.String("us-east-1")})
	cyan := color.New(color.FgCyan).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	if len(c.String("v")) > 0 && len(c.String("t")) > 0 {
		params := &ec2.DescribeInstancesInput{
			Filters: []*ec2.Filter{
				{ // Required
					Name: aws.String(fmt.Sprintf("tag:%s", c.String("t"))),
					Values: []*string{
						aws.String(c.String("v")), // Required
						// More values...
					},
				},
				// More values...
			},
		}
		resp, err := svc.DescribeInstances(params)

		if err != nil {
			if awsErr, ok := err.(awserr.Error); ok {
				// Generic AWS error with Code, Message, and original error (if any)
				fmt.Println(awsErr.Code(), awsErr.Message(), awsErr.OrigErr())
				if reqErr, ok := err.(awserr.RequestFailure); ok {
					// A service error occurred
					fmt.Println(reqErr.Code(), reqErr.Message(), reqErr.StatusCode(), reqErr.RequestID())
				}
			} else {
				// This case should never be hit, the SDK should always return an
				// error which satisfies the awserr.Error interface.
				fmt.Println(err.Error())
			}
		}

		//Start a Shell Session
		session := sh.NewSession()
		session.ShowCMD = c.Bool("d")

		fmt.Println(fmt.Sprintf("ReservationSets: %v", len(resp.Reservations)))
		for idx, _ := range resp.Reservations {
			for _, inst := range resp.Reservations[idx].Instances {
				for _, tag := range inst.Tags {
					if *tag.Key == c.String("t") {
						fmt.Println(fmt.Sprintf("%s:", cyan(*tag.Value)))
						if c.Bool("d") {
							fmt.Println(awsutil.Prettify(*inst))
						}
						if c.Bool("m") {
							fmt.Println("Metadata")
							fmt.Println(yellow(fmt.Sprintf("%s %s %s", *inst.InstanceId, *inst.PublicDnsName, *inst.PrivateIpAddress)))
						}
					}
				}
				session.Command("ssh", "-q", "-t", "-t", "-o", "StrictHostKeyChecking=no", fmt.Sprintf("%s@%s", c.String("u"), *inst.PrivateIpAddress), fmt.Sprintf("%s", c.Args()[0])).Run()

			}
		}
	}
}

func exec() {

}
