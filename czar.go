package main

import (
	"github.com/codegangsta/cli"
	"github.com/feelobot/czar/cmds"
	"os"
)

func main() {

	app := cli.NewApp()

	app.Name = "Czar AWS EC2 CLI"
	app.Version = "0.0.4"
	app.Commands = []cli.Command{
		{
			Name:  "ls",
			Usage: "execute commands accross ec2 instances",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "tag,t",
					Usage: "tag name to filter by",
				},
				cli.StringFlag{
					Name:  "value,v",
					Usage: "filter value",
				},
				cli.BoolFlag{
					Name:  "debug,d",
					Usage: "show debugging info",
				},
			},
			Action: func(c *cli.Context) error {
				cmds.Ls(c)
				return nil
			},
		},
		{
			Name:  "ssh",
			Usage: "execute commands accross ec2 instances",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "tag,t",
					Usage: "tag name to filter by",
				},
				cli.StringFlag{
					Name:  "value,v",
					Usage: "filter value",
				},
				cli.BoolFlag{
					Name:  "debug,d",
					Usage: "show debugging info",
				},
				cli.StringFlag{
					Name:  "user,u",
					Usage: "user to log in with",
				},
				cli.StringFlag{
					Name:  "concurrency,c",
					Usage: "how many servers to execute on concurrently",
				},
				cli.StringFlag{
					Name:  "limit,l",
					Usage: "limit of servers to run on",
				},
				cli.BoolFlag{
					Name:  "metadata,m",
					Usage: "displays more data per instance if true",
				},
			}, // End of Flags
			// Execute Action
			Action: func(c *cli.Context) error {
				cmds.Ssh(c)
				return nil
			},
		},
	}
	app.Run(os.Args)
}
