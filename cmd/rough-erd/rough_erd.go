package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli"
)

func main() {

	app := cli.NewApp()
	app.Name = "rough-erd"
	app.Usage = "make rough ER diagram."
	app.Version = "v1.0.0"

	app.Commands = []cli.Command{
		{
			Name:    "make",
			Aliases: []string{"m"},
			Usage:   "make ER diagram.",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name: "database, d",
					Value: "mysql",
					Usage: "database type",
				},
				cli.StringFlag{
					Name: "user, u",
					Usage: "database user",
				},
				cli.StringFlag{
					Name: "password, p",
					Usage: "database password",
				},
				cli.StringFlag{
					Name: "port, P",
					Usage: "database port",
				},
			},
			Action: func(c *cli.Context) error {
				fmt.Println("run!!!")
				return nil
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
