package cmds

import (
	"fmt"
	"github.com/codegangsta/cli"
)

func init() {
	// initialization code here
}

func Ls(c *cli.Context) {
	fmt.Println("using ls", c.String("v"))
}

func Exec(cmd string) {

}
