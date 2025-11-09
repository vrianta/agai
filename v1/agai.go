package agai

import (
	"fmt"
	"os"

	"github.com/vrianta/agai/v1/internal/flags"
)

/*
 * File to handle Arguments from the user
 * storing all the argument flags
 */

func init() {
	// go through all the arugments and enable some flags
	if len(os.Args) < 2 {
		print_help()
		os.Exit(0)
		return
	}
	for _, arg := range os.Args[1:] {
		switch arg {
		case "--migrate-model", "-mm":
			flags.SyncDatabaseEnabled = true
		case "--migrate-component", "-mc":
			flags.SyncComponentsEnabled = true
		case "--start-server", "-ss":
			flags.StartServer = true
		case "--show-dsn", "-sdn":
			flags.ShowDsn = true
		case "--help", "-h":
			print_help()
			os.Exit(1)

		default:
			println("Wrong Argument Passed plesae use go run . --help/-h to get the list of arguments")
			os.Exit(1)
		}
	}
}

func New() app {
	return app{}
}

func print_help() {
	fmt.Println("Flags:")
	fmt.Println("  --migrate-model,     -mm   Run model database migrations")
	fmt.Println("  --migrate-component, -mc   Sync components with the database")
	fmt.Println("  --start-server,      -ss   Start the HTTP server")
	fmt.Println("  --show-dns,          -sdn  Show Dsn if the database connnection failed")
	fmt.Println("  --help,              -h    Show this help message")
	fmt.Println()
	fmt.Println("Example:")
	fmt.Println("  go run . --migrate-model --start-server")
}
