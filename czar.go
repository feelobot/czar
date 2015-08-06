package main

import (
  "os"
  "github.com/codegangsta/cli"
  "fmt"
  "github.com/aws/aws-sdk-go/aws"
  "github.com/aws/aws-sdk-go/service/ec2"
)

func main() {
  app := cli.NewApp()
  app.Name = "Czar AWS EC2 CLI"
  app.Commands = []cli.Command{
  {
    Name:      "ssh",
    Aliases:     []string{"s"},
    Usage:     "options for task templates",
    Flags: []cli.Flag{
      cli.StringFlag{
        Name:  "f",
        Usage: "filter type",
        Action: func(c *cli.Context) {
          println("sshing into a box with a filter: ")
        },
      },
    },
    Action: func(c *cli.Context) {
      println("sshing into a box: ")
    },
  },
}
  app.Run(os.Args)
}
