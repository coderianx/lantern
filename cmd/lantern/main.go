package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/coderianx/lantern/internal/buckets"
	"github.com/coderianx/lantern/internal/store"
	"github.com/coderianx/lantern/internal/users"
)

var (
	LanternPath string
)

func main() {
	bucket := flag.String("bucket", "", "Bucket works")
	name := flag.String("name", "", "bucket name")

	user := flag.String("user", "", "users works")
	username := flag.String("username", "", "user username")
	password := flag.String("password", "", "user password")
	role := flag.String("role", "user", "user role")

	flag.Parse()

	home, err := os.UserHomeDir()

	if err != nil {
		fmt.Println("Opss... Error;\n\n", err)
		os.Exit(1)
	}

	switch runtime.GOOS {
	case "linux", "darwin":
		LanternPath = filepath.Join(home, ".lantern")
	case "windows":
		LanternPath = filepath.Join(home, ".lantern")
	}

	// Create folder
	os.Mkdir(LanternPath, 0777)

	// Connect database
	err = store.ConnectDB(
		filepath.Join(LanternPath, "lantern.db"),
	)

	if err != nil {
		fmt.Println("Database connect error")
		fmt.Println("Error: ", err)
		os.Exit(1)
	}

	defer store.CloseDB()

	// Create Tables
	if err := store.CreateTables(); err != nil {
		fmt.Println("Create tables error")
		fmt.Println("Error: ", err)
		os.Exit(1)
	}

	switch *bucket {
	case "create":
		buckets.CreateBucket(*name)
	case "delete":
		buckets.DeleteBucket(*name)
	}

	switch *user {
	case "create":
		users.CreateUser(*username, *password, *role)
	}
}

func GetLanternPath() string {
	return LanternPath
}
