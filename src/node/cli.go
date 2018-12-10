package main

import (
	"fmt"
	//"log"
	"github.com/urfave/cli"
	"os"
)

func Run() {

	app := cli.NewApp()
	app.Name = "divvy"
	app.Usage = "fight the loneliness!"
	app.Commands = []cli.Command{
		{
			Name:    "show",
			Aliases: []string{"ls"},
			Usage:   "add a task to the list",
			Action: func(c *cli.Context) error {
				fmt.Println("added task: ", c.Args().First())
				return nil
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "\r\n%s", err.Error())
	}
	fmt.Fprintf(os.Stderr, "\r\n")
}
