package config

import (
	"fmt"
	"os"
)

// storing all the argument flags

// this will only be in consideration if the build mode is enabled
var SyncDatabase = false
var RunServer = false

func init() {
	// go through all the arugments and enable some flags
	if len(os.Args) < 2 {
		fmt.Println("[MASSSAGE] No Argument to Process")
		return
	}
	for _, arg := range os.Args[1:] {
		switch arg {
		case "--migrate-model", "-m":
			SyncDatabase = true
		default:
			panic("Unknown Argument has been passed please check")
		}
	}
}
