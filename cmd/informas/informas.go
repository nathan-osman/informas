package main

import (
	"os"

	"github.com/nathan-osman/informas/db"
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
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "db-name",
			Value: "postgres",
			Usage: "PostgreSQL database name",
		},
		cli.StringFlag{
			Name:  "db-user",
			Value: "postgres",
			Usage: "PostgreSQL database user",
		},
		cli.StringFlag{
			Name:  "db-password",
			Value: os.Getenv("POSTGRES_ENV_POSTGRES_PASSWORD"),
			Usage: "PostgreSQL database password",
		},
		cli.StringFlag{
			Name:  "db-host",
			Value: "postgres",
			Usage: "PostgreSQL database host",
		},
		cli.IntFlag{
			Name:  "db-port",
			Value: 5432,
			Usage: "PostgreSQL database port",
		},
	}
	app.Before = func(c *cli.Context) error {
		db.Connect(
			c.String("db-name"),
			c.String("db-user"),
			c.String("db-password"),
			c.String("db-host"),
			c.Int("db-port"),
		)
		return nil
	}
	app.Commands = []cli.Command{
		db.CreateUserCommand,
		db.MigrateCommand,
	}
	app.Run(os.Args)
}
