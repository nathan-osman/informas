package models

import (
	"bufio"
	"database/sql"
	"fmt"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/howeyc/gopass"
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

// readLine reads a line from STDIN, optionally masking input. Any error will
// result in the application being terminated.
func readLine(prompt string, r *bufio.Reader, mask bool) string {
	fmt.Print(prompt)
	if mask {
		b, err := gopass.GetPasswd()
		if err != nil {
			log.Fatal(err.Error())
		}
		return string(b)
	} else {
		s, err := r.ReadString('\n')
		if err != nil {
			log.Fatal(err.Error())
		}
		return s
	}
}

var (
	// Migrate performs all database migrations. This also includes the initial
	// creation of the database tables.
	MigrateCommand = cli.Command{
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

	// CreateUserCommand creates a new user account.
	CreateUserCommand = cli.Command{
		Name:  "createuser",
		Usage: "create a user account",
		Flags: append(commonFlags,
			cli.StringFlag{
				Name:  "username",
				Value: "",
				Usage: "username for the new account",
			},
			cli.StringFlag{
				Name:  "password",
				Value: "",
				Usage: "password for the new account",
			},
			cli.StringFlag{
				Name:  "email",
				Value: "",
				Usage: "email for the new account",
			},
			cli.BoolFlag{
				Name:  "admin",
				Usage: "create an admin account",
			},
			cli.BoolFlag{
				Name:  "disabled",
				Usage: "create a disabled account",
			},
		),
		Action: func(c *cli.Context) {
			if err := ConnectToDatabase(c); err != nil {
				log.Fatal(err.Error())
			}
			u := &User{
				Username:   c.String("username"),
				Email:      c.String("email"),
				IsAdmin:    c.Bool("admin"),
				IsDisabled: c.Bool("disabled"),
			}
			p := c.String("password")
			r := bufio.NewReader(os.Stdin)
			if len(u.Username) == 0 {
				u.Username = readLine("Username: ", r, false)
			}
			if len(p) == 0 {
				p = readLine("Password: ", r, true)
			}
			if len(u.Email) == 0 {
				u.Email = readLine("Email: ", r, false)
			}
			if err := u.SetPassword(p); err != nil {
				log.Fatal(err.Error())
			}
			if err := u.Save(); err != nil {
				log.Fatal(err.Error())
			}
			log.Infof("user ID:%d created", u.ID)
		},
	}
)
