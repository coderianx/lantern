package buckets

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/coderianx/lantern/internal/config"
	"github.com/coderianx/lantern/internal/store"
)

func DeleteBucket(name string) {
	if name == "" {
		fmt.Println("Please specify a name using --name")
		os.Exit(1)
	}

	var choise string
	fmt.Print("You want delete bucket [y/N]: ")
	fmt.Scanln(&choise)

	if strings.ToLower(choise) == "y" {
		_, err := store.DB.Exec(
			`
			DELETE FROM buckets
			WHERE name = ?
			`,
			name,
		)

		if err != nil {
			fmt.Println("Database error")
			fmt.Println("Error: ", err)
			os.Exit(1)
		}

		fmt.Println("Deleted successfully from database")

		lantern_path := config.GetLanternPath()

		path := filepath.Join(lantern_path, name)

		err = os.RemoveAll(path)

		if err != nil {
			if errors.Is(err, fs.ErrNotExist) {
				fmt.Println("Delete bucket folder error\nBucket folder not exists")
				os.Exit(1)
			}

			fmt.Println("Delete Bucket Folder error")
			fmt.Println("Error: ", err)
			os.Exit(1)
		}

		fmt.Println("Bucket deleted.")
	} else {
		os.Exit(1)
	}

}
