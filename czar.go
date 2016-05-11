package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	_ "github.com/aws/aws-sdk-go/aws/awsutil"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/codegangsta/cli"
	"github.com/codeskyblue/go-sh"
	"github.com/fatih/color"
	"github.com/feelobot/czar/cmds"
	"os"
	"strings"
)

func main() {
	configuration_file := fmt.Sprintf("%s/.czar.cfg.json", os.Getenv("HOME"))
	type Configuration struct {
		User     string
		Tag      string
		Metadata bool
	}

	cyan := color.New(color.FgCyan).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()

	svc := ec2.New(session.New(), &aws.Config{Region: aws.String("us-east-1")})
	app := cli.NewApp()

	app.Name = "Czar AWS EC2 CLI"
	app.Version = "0.0.3"
	app.Commands = []cli.Command{
		{
			Name:    "config",
			Aliases: []string{"c"},
			Usage:   "configure czar defaults",
			Action: func(c *cli.Context) error {
				configuration := Configuration{}

				reader := bufio.NewReader(os.Stdin)

				fmt.Print("Default User? (ubuntu/core/root): ")
				user, _ := reader.ReadString('\n')
				configuration.User = strings.TrimSpace(user)

				fmt.Print("Default Tag Search Filter? (Name): ")
				tag, _ := reader.ReadString('\n')
				configuration.Tag = strings.TrimSpace(tag)
				json_string, err := json.Marshal(configuration)
				if err != nil {
					fmt.Println(err)
				}
				fmt.Println(string(json_string))

				config_file, err := os.Create(configuration_file)
				if err != nil {
					panic(err)
				}
				defer config_file.Close()
				config_file.Write(json_string)
				config_file.Close()
				fmt.Printf("File created!")
				return nil
			},
		},
		{
			Name:  "ls",
			Usage: "execute commands accross ec2 instances",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "n, name",
					Usage: "Specify instance name",
				},
			},
			Action: func(c *cli.Context) {
				var name string
				if len(c.Args()[0]) > 0 {
					name = c.Args()[0]
				} else {
					name = c.String("n")
				}
				cmds.Ls(name)
			},
		},
		{
			Name:    "exec",
			Aliases: []string{"e"},
			Usage:   "execute commands accross ec2 instances",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "tag,t",
					Usage: "tag name to filter by",
				},
				cli.StringFlag{
					Name:  "value,v",
					Usage: "value of the tag key",
				},
				cli.StringFlag{
					Name:  "user,u",
					Usage: "user to log in with",
				},
				cli.BoolFlag{
					Name:  "metadata,m",
					Usage: "displays more data per instance if true",
				},
			}, // End of Flags
			// Execute Action
			Action: func(c *cli.Context) error {
				config := Configuration{}
				if len(c.String("t")) > 0 {
					config.Tag = c.String("t")
				}
				if len(c.String("u")) > 0 {
					config.User = c.String("u")
				}
				config.Metadata = c.Bool("m")
				if len(c.String("v")) > 0 && len(config.Tag) > 0 {
					params := &ec2.DescribeInstancesInput{
						Filters: []*ec2.Filter{
							{ // Required
								Name: aws.String(fmt.Sprintf("tag:%s", config.Tag)),
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
					session.ShowCMD = false

					fmt.Println(fmt.Sprintf("ReservationSets: %v", len(resp.Reservations)))
					for idx, _ := range resp.Reservations {
						for _, inst := range resp.Reservations[idx].Instances {
							for _, tag := range inst.Tags {
								if *tag.Key == config.Tag {
									fmt.Println(fmt.Sprintf("%s:", cyan(*tag.Value)))
									//fmt.Println(awsutil.Prettify(*inst))
									if config.Metadata {
										fmt.Println("Metadata")
										fmt.Println(yellow(fmt.Sprintf("%s %s %s", *inst.InstanceId, *inst.PublicDnsName, *inst.PrivateIpAddress)))
									}
								}
							}
							session.Command("ssh", "-t", "-t", "-o", "StrictHostKeyChecking=no", fmt.Sprintf("%s@%s", config.User, *inst.PrivateIpAddress), fmt.Sprintf("%s", c.Args()[0])).Run()

						}
					}
				}
				return nil
			},
		},
	}
	app.Run(os.Args)
}
