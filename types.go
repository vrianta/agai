package main

type flags struct {
	// App creation
	create_app bool
	app_name   string

	// Controller
	controllers_root           string
	controller_names_to_create []string

	// View
	create_view          bool
	view_root            string
	view_names_to_create []string

	// Model
	model_names_to_create []string

	// Component
	create_component          bool
	component_names_to_create []string

	// Start commands
	start_app     bool
	start_handler bool

	// Migrate commands
	migrate_model     bool
	migrate_component bool
}
