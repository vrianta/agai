package config

import (
	"fmt"
	"os"
)

// storing all the argument flags

// this will only be in consideration if the build mode is enabled
var SyncDatabaseEnabled = false
var RunServer = false
var SyncComponentsEnabled = false

func init() {
	// go through all the arugments and enable some flags
	if len(os.Args) < 2 {
		fmt.Println("[MASSSAGE] No Argument to Process")
		return
	}
	for _, arg := range os.Args[1:] {
		switch arg {
		case "--migrate-model", "-mm":
			SyncDatabaseEnabled = true
		case "--migrate-component", "-mc":
			SyncComponentsEnabled = true
		default:
			panic("Unknown Argument has been passed please check")
		}
	}
}
