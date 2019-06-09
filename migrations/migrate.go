// Wrap the `goose.Run` command to allow us to use configuration specific
// to the organization repo, which involves using the database settings from
// the config.
package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/cedric-parisi/fizzbuzz-api/config"
	"github.com/cedric-parisi/fizzbuzz-api/storage"

	_ "github.com/lib/pq"
	"github.com/pressly/goose"
)

const (
	driver = "postgres"
)

var (
	dir        = flag.String("dir", "migrations/sql", "directory with migration files")
	configPath = flag.String("config", ".", "config file")
)

func main() {
	flag.Usage = usage

	flag.Parse()
	args := flag.Args()

	if len(args) == 0 {
		log.Fatal("expected at least one arg")
	}

	command := args[0]

	if err := goose.SetDialect(driver); err != nil {
		log.Fatal(err)
	}

	// Load the configuration
	config, err := config.NewConfig()
	if err != nil {
		log.Fatalf("error loading configuration: %v", err)
	}

	// Create a new storage layer
	pq, err := storage.NewPostgres(config)
	if err != nil {
		log.Fatalf("unable to create storage layer: %s", err)
	}

	if err := goose.Run(command, pq.DB, *dir, args[1:]...); err != nil {
		log.Fatalf("goose run: %v", err)
	}
}

func usage() {
	fmt.Println(usageRun)
	flag.PrintDefaults()
	fmt.Println(usageCommands)
}

var (
	usageRun      = `goose [OPTIONS] COMMAND`
	usageCommands = `
Commands:
    up                   Migrate the DB to the most recent version available
    up-to VERSION        Migrate the DB to a specific VERSION
    down                 Roll back the version by 1
    down-to VERSION      Roll back to a specific VERSION
    redo                 Re-run the latest migration
    status               Dump the migration status for the current DB
    version              Print the current version of the database
    create NAME [sql|go] Creates new migration file with next version
`
)
