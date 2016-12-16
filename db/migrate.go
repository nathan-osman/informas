package db

import (
	"github.com/urfave/cli"
)

// Migrate performs all database migrations. This also includes the initial
// creation of the database tables.
var MigrateCommand = cli.Command{
	Name:  "migrate",
	Usage: "perform all database migrations",
	Action: func(c *cli.Context) {
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
