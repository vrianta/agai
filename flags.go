package main

var (
	// Create commands
	create_app bool
	app_name   string

	create_controller bool
	controller_name   string

	create_model bool
	model_name   string

	create_component bool
	component_name   string

	// Start commands
	start_app     bool
	start_handler bool

	// Migrate commands
	migrate_model     bool
	migrate_component bool
)
