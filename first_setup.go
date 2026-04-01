package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/vrianta/agai/v1/utils"
)

// this function will update the user.component.json.template with the username and password provided by the user
func first_setup() {
	var first_setup_username string
	var first_setup_password string
	// read the user.component.json.template file
	if data, err := templates.ReadFile("templates/components/user.component.json"); err == nil {
		content := string(data)
		// ask for the username and password
		fmt.Print("Enter the username for the admin user: ")
		fmt.Scanln(first_setup_username)
		fmt.Print("Enter the password for the admin user: ")
		fmt.Scanln(first_setup_password)
		// replace the {{username}} and {{password}} with the values provided by the user
		content = strings.ReplaceAll(content, "{{username}}", first_setup_username)
		hashedPassword, err := utils.HashPassword(first_setup_password)
		if err != nil {
			fmt.Println("Error hashing password:", err)
			os.Exit(1)
		}
		content = strings.ReplaceAll(content, "{{password}}", hashedPassword)

		// write the updated content to user.component.json
		if err := os.WriteFile("components/user.component.json", []byte(content), 0644); err != nil {
			fmt.Println("Error writing user.component.json:", err)
			os.Exit(1)
		}
	} else {
		fmt.Println("Error reading user.component.json.template:", err)
		os.Exit(1)
	}
}

func reConfigureAdminAccount() {
	if f.initiateAdmin {
		first_setup()
	}
}
