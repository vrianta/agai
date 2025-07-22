package session

import (
	"encoding/gob"
	"fmt"
	"os"
)

/*
Handles session Storage an retrival functions and the mechanism
*/

func loadAllSessionsFromDisk() error {
	// Check if sessions.data exists
	if _, err := os.Stat("sessions.data"); os.IsNotExist(err) {
		fmt.Println("[Sessions] sessions.data not found — skipping load")
		return nil
	} else if err != nil {
		// Some other filesystem error
		return fmt.Errorf("error checking session file: %w", err)
	}

	// Open file and decode
	f, err := os.Open("sessions.data")
	if err != nil {
		return fmt.Errorf("failed to open session store file: %w", err)
	}
	defer f.Close()

	var loaded map[string]*Instance
	dec := gob.NewDecoder(f)
	if err := dec.Decode(&loaded); err != nil {
		return fmt.Errorf("failed to decode session map: %w", err)
	}

	sessionStoreMutex.Lock()
	instances = loaded
	sessionStoreMutex.Unlock()

	fmt.Printf("[Sessions] Loaded %d sessions from sessions.data\n", len(instances))
	return nil
}
