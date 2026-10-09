package users

import (
	"fmt"
	"os"

	"github.com/coderianx/lantern/internal/store"
	"golang.org/x/crypto/bcrypt"
)

func CreateUser(username, password, role string) {
	if username == "" {
		fmt.Println("Please specify a name using --username.")
		os.Exit(1)
	}

	if password == "" {
		fmt.Println("Please specify a name using --password.")
	}

	passwordHash, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		12,
	)

	if err != nil {
		fmt.Println("Hash Password error")
		fmt.Println("Error: ", err)
		os.Exit(1)
	}

	_, err = store.DB.Exec(
		`
		INSERT INTO users (username, password role)
		VALUES (?,?,?)
		`,
		username,
		string(passwordHash),
		role,
	)

	if err != nil {
		fmt.Println("Database insert error")
		fmt.Println("Error: ", err)
		os.Exit(1)
	}
}
