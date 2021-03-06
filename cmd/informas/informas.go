package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/nathan-osman/informas/db"
	"github.com/nathan-osman/informas/server"
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
		cli.StringFlag{
			Name:  "http-addr",
			Value: ":8000",
			Usage: "address and port to listen on",
		},
		cli.StringFlag{
			Name:  "data-dir",
			Value: "data",
			Usage: "path to data directory",
		},
	}
	app.Action = func(c *cli.Context) {

		// Connect to the database
		if err := db.Connect(
			c.String("db-name"),
			c.String("db-user"),
			c.String("db-password"),
			c.String("db-host"),
			c.Int("db-port"),
		); err != nil {
			cli.HandleExitCoder(err)
		}

		// Perform all pending migrations
		if err := db.Migrate(); err != nil {
			cli.HandleExitCoder(err)
		}

		// Create the server
		s, err := server.New(
			c.String("http-addr"),
			c.String("data-dir"),
		)
		if err != nil {
			cli.HandleExitCoder(err)
		}

		// Start the server
		if err := s.Start(); err != nil {
			cli.HandleExitCoder(err)
		}

		// Wait for SIGINT
		q := make(chan os.Signal)
		signal.Notify(q, syscall.SIGINT)
		<-q

		// Shut everything down
		s.Stop()
	}
	app.Run(os.Args)
}
