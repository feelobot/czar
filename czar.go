package main

import (
  "os"
  "github.com/codegangsta/cli"
  "fmt"
  "github.com/aws/aws-sdk-go/aws"
  "github.com/aws/aws-sdk-go/service/ec2"
  "github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/awsutil"
)

func main() {
  // Note that you can also configure your region globally by
  // exporting the AWS_REGION environment variable
  svc := ec2.New(&aws.Config{Region: aws.String("us-east-1")})
  app := cli.NewApp()
  app.Name = "Czar AWS EC2 CLI"
  app.Commands = []cli.Command{
  {
    Name:      "ssh",
    Aliases:     []string{"s"},
    Usage:     "options for task templates",
    Flags: []cli.Flag{
      cli.StringFlag{
        Name:  "tag,t",
        Usage: "tag name",
      },
      cli.StringFlag{
        Name:  "value,v",
        Usage: "tag value",
      },
    },
    Action: func(c *cli.Context) {
      if len(c.String("v")) > 0 && len(c.String("t")) > 0 {
        params := &ec2.DescribeInstancesInput{
      		Filters: []*ec2.Filter{
      			{ // Required
      				Name: aws.String(fmt.Sprintf("tag:%s",c.String("t"))),
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

      	// Pretty-print the response data.
      	fmt.Println(awsutil.Prettify(resp))
      }
    },
  },
  }
  app.Run(os.Args)
}
