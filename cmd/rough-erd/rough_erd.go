package main

import (
	"log"
	"os"

	"github.com/nesheep5/rough-erd"

	"github.com/urfave/cli"
)

func main() {

	app := cli.NewApp()
	app.Name = "rough-erd"
	app.Usage = "This tool creates a rough ER diagram."
	app.Version = "v1.0.0"

	app.Commands = []cli.Command{
		{
			Name:    "make",
			Aliases: []string{"m"},
			Usage:   "make ER diagram.",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "databasetype",
					Value: "mysql",
					Usage: "database type",
				},
				cli.StringFlag{
					Name:  "user, u",
					Usage: "database user",
				},
				cli.StringFlag{
					Name:  "password, p",
					Usage: "database password",
				},
				cli.IntFlag{
					Name:  "port, P",
					Usage: "database port",
				},
				cli.StringFlag{
					Name:  "protocol",
					Usage: "database protocol",
				},
				cli.StringFlag{
					Name:  "name, n",
					Usage: "database name",
				},
			},
			Action: func(c *cli.Context) error {
				o := &rough_erd.Option{
					Database: c.String("databasetype"),
					User:     c.String("user"),
					Password: c.String("password"),
					Port:     c.Int("port"),
					Protocol: c.String("protocol"),
					Name:     c.String("name"),
				}
				return rough_erd.Run(o)
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
