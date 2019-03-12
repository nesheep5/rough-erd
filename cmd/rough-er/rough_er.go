package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli"
)

func main() {

	app := cli.NewApp()
	app.Name = "rough-er"
	app.Usage = "make rough ER diagram."
	app.Version = "v1.0.0"

	app.Commands = []cli.Command{
		{
			Name:    "run",
			Aliases: []string{"l"},
			Usage:   "run",
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
