package db

import (
	"bufio"
	"fmt"
	"os"

	"github.com/howeyc/gopass"
	"github.com/urfave/cli"
)

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

// CreateUserCommand creates a new user account.
var CreateUserCommand = cli.Command{
	Name:  "createuser",
	Usage: "create a user account",
	Flags: []cli.Flag{
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
	},
	Action: func(c *cli.Context) {
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
