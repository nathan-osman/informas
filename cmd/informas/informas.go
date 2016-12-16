package main

import (
	"os"

	"github.com/nathan-osman/informas/models"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "informas"
	app.Usage = "centralized Twitter account manager"
	app.Version = "0.1.0"
	app.Authors = []cli.Author{
		cli.Author{
			Name:  "Nathan Osman",
			Email: "nathan@quickmediasolutions.com",
		},
	}
	app.Commands = []cli.Command{
		models.MigrateCommand,
	}
	app.Run(os.Args)
}
