package buckets

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/coderianx/lantern/internal/store"
)

func CreateBucket(name string) {
	if name == "" {
		fmt.Println("Please specify a name using --name.")
		os.Exit(1)
	}

	home, err := os.UserHomeDir()

	if err != nil {
		fmt.Println("Opss... Error")
		fmt.Println(err)
		os.Exit(1)
	}

	path := filepath.Join(filepath.Join(home, ".lantern"), name)

	err = os.Mkdir(path, 0777)

	if err != nil {
		fmt.Println("Create bucket folder error")
		fmt.Println(err)
		os.Exit(1)
	}

	_, err = store.DB.Exec(
		`
		INSERT INTO buckets (name)
		VALUES (?)
		`,
		name,
	)

	if err != nil {
		fmt.Println("Database table create error")
		fmt.Println("Error: ", err)

		err = os.RemoveAll(path)

		if err != nil {
			fmt.Println("Delete folder error")
			fmt.Println("Error: ", err)
		}

		os.Exit(1)
	}

	fmt.Printf("Bucket created\nBucket name: %s\n", name)
}
