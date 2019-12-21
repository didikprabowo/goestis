package main

import (
	"github.com/urfave/cli"
	"github.com/k/amqp"
	"github.com/k/app"
	"os"
)

func main() {

	clientApp := cli.NewApp()
	clientApp.Name = ""
	clientApp.Version = "0.0.1"
	clientApp.Commands = []cli.Command{
		{
			Name:        "producer",
			Description: "start producer",
			Action: func(c *cli.Context) {
				app.RunningApp()
			},
		},
		{
			Name:        "consumer",
			Description: "start consumer",
			Action: func(c *cli.Context) {
				amqp.PubArticle()
			},
		},
	}

	clientApp.Run(os.Args)

}
