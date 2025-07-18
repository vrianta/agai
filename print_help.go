package main

import "fmt"

/*
--create-app/-ca -> To create application (--create-app app_name)
--create-controller/-cc -> to create controller (--create-controller controller_name)
--create-model/-cm -> create model (--create-model model_name)
--create-component/ -> crete component (--create-component component_name)
--start-app/-sa -> To Start the Application (--start-app)
--start-handler/-sh -> To controll apps various functions. it is good to have it running to auto update and handle multiple functions (--start-handler)
--migrate-model/-mm -> to migrate models (--migrate-model)
--migrate-component/-mc -> to migrate components (--migrate-component)
--help/-h -> to print help
--help-web-config/-hwc -> to print details on web config
--help-database-config/-hdc -> to print config details for database
--help-session-config/-hsc -> to print config details of session
--help-smtp-config/-hsc -> to print config details for smtp config
*/
func print_help() {
	fmt.Println("Flags:")
	fmt.Printf("  %-32s %s\n", "--create-app, -ca", "Create an application (e.g. --create-app blog)")
	fmt.Printf("  %-32s %s\n", "--create-controller, -cc", "Create a controller (e.g. --create-controller post)")
	fmt.Printf("  %-32s %s\n", "--create-model, -cm", "Create a model (e.g. --create-model user)")
	fmt.Printf("  %-32s %s\n", "--create-model-component, -cmc", "Create a component (e.g. --create-component nav)")
	fmt.Printf("  %-32s %s\n", "--create-view, -cv", "Create a view (optionally pass controller/view name)")
	fmt.Printf("  %-32s %s\n", "--start-app, -sa", "Start the application server")
	fmt.Printf("  %-32s %s\n", "--start-handler, -sh", "Start the dev handler (auto-reloader / multi-handler)")
	fmt.Printf("  %-32s %s\n", "--migrate-model, -mm", "Run model migrations")
	fmt.Printf("  %-32s %s\n", "--migrate-component, -mc", "Sync components with the database")
	fmt.Printf("  %-32s %s\n", "--help, -h", "Show this help message")
	fmt.Printf("  %-32s %s\n", "--help-web-config, -hwc", "Show help for web config structure")
	fmt.Printf("  %-32s %s\n", "--help-database-config, -hdc", "Show help for database config structure")
	fmt.Printf("  %-32s %s\n", "--help-session-config, -hsc", "Show help for session config structure")
	fmt.Printf("  %-32s %s\n", "--help-smtp-config, -hsm", "Show help for SMTP config structure")
	fmt.Println()
	fmt.Println("Example:")
	fmt.Println("  go run . --create-app blog --create-model post --start-app")
}

// print_web_config_help prints an explanation of the web server configuration structure,
// including available options, their sources (config file or environment), and expected values.
func print_web_config_help() {
	fmt.Println("Web Configuration Help")
	fmt.Println("-----------------------")
	fmt.Println("Configuration can be defined in config/web.json or overridden via environment variables:")
	fmt.Println()

	fmt.Println("Environment Variables:")
	fmt.Println("  SERVER_PORT           - Port the server runs on (e.g. 8080)")
	fmt.Println("  SERVER_HOST           - Host address (e.g. 127.0.0.1 or 0.0.0.0)")
	fmt.Println("  SERVER_HTTPS          - Enable HTTPS (true/false)")
	fmt.Println("  BUILD                 - Enable build mode for development (true/false)")
	fmt.Println("  MAX_SESSION_COUNT     - Max number of active sessions (integer)")
	fmt.Println("  SESSION_STORE_TYPE    - Type of session store (e.g. disk, memory)")
	fmt.Println()

	fmt.Println("File-based Config (config/web.json):")
	fmt.Println(`  {
    "port": "8080",
    "host": "127.0.0.1",
    "https": false,
    "build": true,
    "maxSessionCount": 100,
    "sessionStoreType": "disk",
    "staticFolders": ["public", "assets"],
    "cssFolders": ["css"],
    "jsFolders": ["js"],
    "viewFolder": "views"
  }`)
	fmt.Println()
	fmt.Println("Notes:")
	fmt.Println("- Environment variables override config file values.")
	fmt.Println("- You can customize which folders are used to serve static files, CSS, JS, and views.")
}

// print_database_config_help prints the supported database configuration fields,
// how they can be set via environment variables, and provides a sample config file structure.
func print_database_config_help() {
	fmt.Println("Database Configuration Help")
	fmt.Println("---------------------------")
	fmt.Println("Configuration can be defined in config/database.json or overridden via environment variables.")
	fmt.Println()

	fmt.Println("Environment Variables:")
	fmt.Println("  DB_HOST        - Database server host (e.g. localhost)")
	fmt.Println("  DB_PORT        - Database server port (e.g. 5432 for PostgreSQL)")
	fmt.Println("  DB_USER        - Database username")
	fmt.Println("  DB_PASSWORD    - Database password")
	fmt.Println("  DB_DATABASE    - Target database name")
	fmt.Println("  DB_PROTOCOL    - Network protocol (usually 'tcp')")
	fmt.Println("  DB_DRIVER      - Database driver (e.g. postgres, mysql, sqlite3)")
	fmt.Println("  DB_SSLMODE     - SSL mode (e.g. disable, require, verify-full)")
	fmt.Println()

	fmt.Println("File-based Config (config/database.json):")
	fmt.Println(`  {
    "host": "localhost",
    "port": "5432",
    "user": "your_user",
    "password": "your_password",
    "database": "your_db",
    "protocol": "tcp",
    "driver": "postgres",
    "sslmode": "disable"
  }`)
	fmt.Println()
	fmt.Println("Notes:")
	fmt.Println("- Environment variables override values in the JSON config.")
	fmt.Println("- Use 'driver' to set the SQL dialect: postgres, mysql, sqlite3, etc.")
	fmt.Println("- 'sslmode' is required for PostgreSQL and can be left empty for others.")
}

// print_smtp_config_help prints the structure and usage of the SMTP configuration,
// including environment variables and JSON-based configuration format.
func print_smtp_config_help() {
	fmt.Println("SMTP Configuration Help")
	fmt.Println("------------------------")
	fmt.Println("SMTP settings are read from config/smtp.json, but can be overridden via environment variables.")
	fmt.Println()

	fmt.Println("Environment Variables:")
	fmt.Println("  SMTP_HOST        - SMTP server (e.g. smtp.gmail.com)")
	fmt.Println("  SMTP_PORT        - SMTP port (e.g. 587 or 465)")
	fmt.Println("  SMTP_USERNAME    - SMTP login username (usually an email)")
	fmt.Println("  SMTP_PASSWORD    - Password or app-specific SMTP password")
	fmt.Println("  SMTP_USE_TLS     - Whether to use TLS (true/false)")
	fmt.Println()

	fmt.Println("File-based Config (config/smtp.json):")
	fmt.Println(`  {
    "Host": "smtp.example.com",
    "Port": 587,
    "Username": "your@email.com",
    "Password": "yourpassword",
    "UseTLS": true
  }`)
	fmt.Println()
	fmt.Println("Notes:")
	fmt.Println("- Environment variables override JSON config values.")
	fmt.Println("- UseTLS should be true for most modern SMTP services.")
	fmt.Println("- Use port 587 for STARTTLS, or 465 for implicit TLS.")
}
