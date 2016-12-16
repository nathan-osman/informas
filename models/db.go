package models

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/Sirupsen/logrus"
	_ "github.com/lib/pq"
	"github.com/urfave/cli"
)

var (
	db          *sql.DB
	log         = logrus.WithField("context", "db")
	commonFlags = []cli.Flag{
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
			Name:  "db-name",
			Value: "postgres",
			Usage: "PostgreSQL database name",
		},
	}
)

// ConnectToDatabase establishes a connection to the PostgreSQL database used
// for all SQL queries. This function should be called before using any other
// types or functions in the package.
func ConnectToDatabase(c *cli.Context) error {
	d, err := sql.Open(
		"postgres",
		fmt.Sprintf(
			"dbname=%s user=%s password=%s host=%s port=%d",
			c.String("db-name"),
			c.String("db-user"),
			c.String("db-password"),
			c.String("db-host"),
			c.Int("db-port"),
		),
	)
	if err != nil {
		return err
	}
	db = d
	return nil
}

// Migrate performs all database migrations. This also includes the initial
// creation of the database tables.
var MigrateCommand = cli.Command{
	Name:  "migrate",
	Usage: "perform all database migrations",
	Flags: commonFlags,
	Action: func(c *cli.Context) {
		if err := ConnectToDatabase(c); err != nil {
			log.Fatal(err.Error())
		}
		tableMigrations := map[string]func() error{
			"Users": migrateUserTable,
		}
		for n, f := range tableMigrations {
			if err := f(); err != nil {
				log.Fatal(err.Error())
			} else {
				log.Infof("migrated %s table", n)
			}
		}
	},
}
