package main

import (
  "os"
  "bufio"
  "github.com/codegangsta/cli"
  "encoding/json"
  "fmt"
  "strings"
  "github.com/aws/aws-sdk-go/aws"
  "github.com/aws/aws-sdk-go/service/ec2"
  "github.com/aws/aws-sdk-go/aws/awserr"
  "github.com/codeskyblue/go-sh"
)

func main() {

  type Configuration struct {
      Key    string
      User   string
      Tag    string
  }
  /*
  file, _ := os.Open("~/.czar.cfg.json")
  decoder := json.NewDecoder(file)
  configuration := Configuration{}
  err := decoder.Decode(&configuration)
  if err != nil {
    fmt.Println("error:", err)
  }
  //fmt.Println(configuration.Users) // output: [UserA, UserB]
  */
  // Note that you can also configure your region globally by
  // exporting the AWS_REGION environment variable
  svc := ec2.New(&aws.Config{Region: aws.String("us-east-1")})
  app := cli.NewApp()
  app.Name = "Czar AWS EC2 CLI"
  app.Version = "0.0.2"
  app.Commands = []cli.Command{
  {
    Name:      "config",
    Aliases:     []string{"c"},
    Usage:     "configure czar defaults",
    Action: func(c *cli.Context) {
      configuration := Configuration{}

      reader := bufio.NewReader(os.Stdin)
      fmt.Print("Enter Path to AWS Key (.pem): ")
      key_path, _ := reader.ReadString('\n')
      configuration.Key = strings.TrimSpace(key_path)

      fmt.Print("Default User? (ubuntu/core/root): ")
      user, _ := reader.ReadString('\n')
      configuration.User = strings.TrimSpace(user)

      fmt.Print("Default Tag Search Filter? (Name): ")
      tag, _ := reader.ReadString('\n')
      configuration.Tag = strings.TrimSpace(tag)
      json_string, err := json.Marshal(configuration)
      if err != nil {
          fmt.Println(err)
          return
      }
      fmt.Println(string(json_string))

      config_file, err := os.Create(fmt.Sprintf("%s/.czar.cfg.json",os.Getenv("HOME")))
      if err != nil {
         panic(err)
      }
      defer config_file.Close()
      config_file.Write(json_string)
      config_file.Close()
      fmt.Printf("File created!")
    },
  },
  {
    Name:      "exec",
    Aliases:     []string{"e"},
    Usage:     "execute commands accross ec2 instances",
    Flags: []cli.Flag{
      cli.StringFlag{
          Name:  "p, pem",
          Usage: "Specify pem file",
      },
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
    }, // End of Flags
    // Execute Action
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

        fmt.Println("> Number of reservation sets: ", len(resp.Reservations))
        for idx, res := range resp.Reservations {
            fmt.Println("  > Number of instances: ", len(res.Instances))
            for _, inst := range resp.Reservations[idx].Instances {
                fmt.Println("    - Instance ID: ", *inst.InstanceID)
                fmt.Println("    - DNS Name: ", *inst.PublicDNSName)
                for _, tag := range inst.Tags {
                  if *tag.Key == c.String("t") {
                    fmt.Println("    - Tag: ", *tag.Value)
                  }
                }
                session := sh.NewSession()
                session.ShowCMD = true
                session.Command("eval","`ssh-agent`").Run()
                session.Command("ssh-add",c.String("p")).Run()
                session.Command("ssh", "-o", "StrictHostKeyChecking=no", "-i" , c.String("p"), fmt.Sprintf("%s@%s",c.String("u"),*inst.PublicDNSName),fmt.Sprintf("%s",c.Args()[0])).Run()

            }
        }
      }
    },
  },
  }
  app.Run(os.Args)
}
