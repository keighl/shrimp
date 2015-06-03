package main

import (
	"github.com/codegangsta/cli"
	c "github.com/keighl/shrimp/config"
	"log"
	"os"
)

var (
	Config *c.Configuration
	Logger = log.New(os.Stdout, "toolz: ", log.Lshortfile)
)

func main() {
	Config = c.Config(os.Getenv("SHRIMP_ENV"))
	newCLI().Run(os.Args)
}

func newCLI() *cli.App {
	app := cli.NewApp()

	app.Commands = []cli.Command{
		{
			Name: "something",
			Action: func(c *cli.Context) {
				ID := c.Args()[0]
				Logger.Println("---> ID", ID)
			},
		},
	}
	app.Name = "shrimp-toolz"
	app.Version = "1.0"
	return app
}
