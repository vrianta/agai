package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/vrianta/agai/v1/log"
)

/*
This is a handler of my overall package
This will do following things

1. Create Default Application
2. Create Default Controller
3. Create Default Module
4. Controll the migration
5. Loop Over all the files to check for update and if any file updated then it will restart the application auto matically

----- Command line arguments
1. --new-application app_name / -na app_name (if will make sure the current directory is not a application folder by checking coltroller/model/views folders)
2. --new-controller controller_name / -nc (if controller for the same is not present then it will create it)
3. --new-model model_name / -nm (create new model if that is not already create - only the ID attribute will be there in the default mpdel)
4. --migrate-model / -mm (check the current model if the build is true in config then the it will confirm before doing any changes in the model)
5. --migrate-component / -mc (Migrate Components with the DB and ask for deletion or modification in Build mode true)
6. --start-server / -ss (start the server with the configaration mentioned in the config folder)

-------------------------------------------------------------------------------------------------------------------------------------------------

This file will be usefull for future expansion of extension and UI based Modifications
TODO:
IDK :P
*/

func main() {

	handle_args()

	create_controller() // creating all the controllers
	create_view()       // creating views
	create_models()     // create models
	create_components() // creating components

	create_configs() // create different configs
}

/*
handle_args parses command-line arguments manually without using the 'flag' package.
It sets internal f based on the arguments provided and extracts values for named parameters like app, controller, model, and component.
It also handles help-related f to print specific configuration help (web, database, session, SMTP).
This function exits the program after printing help or encountering invalid/missing arguments.
*/
func handle_args() {
	args := os.Args[1:]

	if len(args) == 0 {
		print_help()
		os.Exit(0)
		return
	}

	for i := 0; i < len(args); i++ {
		arg := args[i]

		switch arg {

		// --- Create View
		case "--create-view", "-cv":
			if i+1 < len(args) && !strings.HasPrefix(args[i+1], "-") {
				f.view_names_to_create = append(f.view_names_to_create, args[i+1])
				i++
			} else {
				f.create_view = true
			}

		// --- Create App
		case "--create-app", "-ca":
			f.create_app = true
			if i+1 < len(args) && !strings.HasPrefix(args[i+1], "-") {
				f.app_name = args[i+1]
				i++
			} else {
				fmt.Println("Missing app name for --create-app")
				os.Exit(1)
			}

		// --- Create Controller
		case "--create-controller", "-cc":
			if i+1 < len(args) && !strings.HasPrefix(args[i+1], "-") {
				f.controller_names_to_create = append(f.controller_names_to_create, args[i+1])
				i++
			} else {
				log.Error("Missing or invalid controller name for --create-controller")
				os.Exit(1)
			}

		// --- Create Model
		case "--create-model", "-cm":
			if i+1 < len(args) && !strings.HasPrefix(args[i+1], "-") {
				f.model_names_to_create = append(f.model_names_to_create, args[i+1])
				i++
			} else {
				fmt.Println("Missing or invalid model name for --create-model")
				os.Exit(1)
			}

		// --- Create Component
		case "--create-model-component", "-cmc":
			if i+1 < len(args) && !strings.HasPrefix(args[i+1], "-") {
				f.component_names_to_create = append(f.component_names_to_create, args[i+1])
				i++
			} else {
				f.create_component = true
			}

		// --- Start App
		case "--start-app", "-sa":
			f.start_app = true

		// --- Start Handler
		case "--start-handler", "-sh":
			f.start_handler = true

		// --- Migrate Model
		case "--migrate-model", "-mm":
			f.migrate_model = true

		// --- Migrate Component
		case "--migrate-component", "-mc":
			f.migrate_component = true

		// --- Help
		case "--help", "-h":
			print_help()
			os.Exit(0)

		// --- Help: Web Config
		case "--help-web-config", "-hwc":
			print_web_config_help()
			os.Exit(0)

		// --- Help: Database Config
		case "--help-database-config", "-hdc":
			print_database_config_help()
			os.Exit(0)

		// --- Help: Session Config
		case "--help-session-config", "-hsc":
			// print_session_config_help()
			os.Exit(0)

		// --- Help: SMTP Config
		case "--help-smtp-config", "-hsm":
			print_smtp_config_help()
			os.Exit(0)

		// --- Unknown Argument
		default:
			fmt.Printf("Unknown flag: %s\n", arg)
			fmt.Println("Use --help to see available options")
			os.Exit(1)
		}
	}
}
