package main

import (
	"fmt"
	"log"
	"os"

	"github.com/nesheep5/checo"
	"github.com/urfave/cli"
)

func main() {

	app := cli.NewApp()
	app.Name = "rough-er"
	app.Usage = "make rough ER diagram."
	app.Version = "v1.0.0"

	app.Commands = []cli.Command{
		{
			Name:    "list",
			Aliases: []string{"l"},
			Usage:   "Shows checking SNS list",
			Action: func(c *cli.Context) error {
				fmt.Println("Checking SNS list:")
				for _, ch := range checo.CheckerMap {
					fmt.Printf("  %v", ch.Name)
				}
				return nil
			},
		},
		{
			Name:    "check",
			Aliases: []string{"c"},
			Usage:   "Checking SNS account",
			Action: func(c *cli.Context) error {
				account := c.Args().Get(0)
				if account == "" {
					fmt.Errorf("account is required.")
				}

				fmt.Printf("Search Account: %v \n\n", account)

				for _, c := range checo.CheckerMap {
					exists, err := c.Exists(account)
					if err != nil {
						return err
					}

					var msg string
					if exists {
						msg = "Oops! Alrady Exists..."
					} else {
						msg = "OK! No Exists!"
					}

					fmt.Printf("%v : %v \n", c.Name, msg)
				}
				return nil
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
