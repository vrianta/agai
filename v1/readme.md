# Go Server Framework – v0.2.1 User Guide

The **Go Server Framework** is a lightweight, modular, and extendable web framework built from the ground up in Go. It’s designed for developers who prefer simplicity, structure, and full control over how their web server runs — without sacrificing flexibility.

This framework supports everything you need to build production-ready web applications:

* Component-based architecture
* Model-driven development with auto schema migration
* Clean routing system
* Customizable views with PHP-style templating
* Session storage (in-memory or disk)
* A growing CLI toolkit to bootstrap, migrate, and manage your application lifecycle

Version `v0.2.1` introduces major improvements including a new model migration engine, disk-based session storage, a per-model component system, crash-safe concurrency fixes, and cleaner project structure under the new `v1/` module path.

This documentation will walk you through every feature — from setup to advanced usage — while keeping the philosophy simple: **write less, control more, stay fast.**

## Table of Contents

- [Go Server Framework – v0.2.1 User Guide](#go-server-framework--v021-user-guide)
  - [Table of Contents](#table-of-contents)
  - [Features](#features)
  - [Installation](#installation)
  - [Run Application](#run-application)
    - [Flags](#flags)
  - [Project Structure](#project-structure)
  - [Configuration (`web.config.JSON`)](#configuration-webconfigjson)
    - [Database Configuration (`database.Config.JSON`)](#database-configuration-databaseconfigjson)
      - [Supported Environment Variables](#supported-environment-variables)
  - [Routing](#routing)
    - [Define Route Handlers](#define-route-handlers)
    - [Use Custom Router](#use-custom-router)
  - [Controllers](#controllers)
    - [What is a Controller?](#what-is-a-controller)
    - [Controller Structure](#controller-structure)
      - [Public Fields](#public-fields)
      - [Public Methods](#public-methods)
    - [Example: Creating a Controller](#example-creating-a-controller)
    - [Accessing Request Data](#accessing-request-data)
    - [Session Management](#session-management)
    - [Redirects](#redirects)
    - [Creating Views](#creating-views)
      - [View Directory Structure](#view-directory-structure)
      - [Supported Template Files](#supported-template-files)
    - [Location](#location)
    - [Template Files](#template-files)
    - [Example Structure](#example-structure)
  - [Models](#models)
    - [How to Initialize the Database](#how-to-initialize-the-database)
    - [How to Create a Model (Quick Example)](#how-to-create-a-model-quick-example)
    - [Model Query Builder System](#model-query-builder-system)
    - [Model Query Result](#model-query-result)
  - [Components](#components)
    - [Structure](#structure)
    - [Example](#example)
      - [File Name](#file-name)
      - [Component Model](#component-model)
      - [Component](#component)
    - [Get Controller](#get-controller)
  - [Static, CSS, and JS File Serving](#static-css-and-js-file-serving)
  - [SMTP/Email Support](#smtpemail-support)
  - [Template Engine \& PHP Parsing Syntax](#template-engine--php-parsing-syntax)
    - [Write Templates in PHP-Style!](#write-templates-in-php-style)
      - [Supported Syntax](#supported-syntax)
      - [Example Template](#example-template)
      - [How It Works](#how-it-works)
  - [Frequently Asked Questions (FAQ)](#frequently-asked-questions-faq)
  - [License](#license)

## Features

* **Custom HTTP Server**: Easily start and stop the server with CLI or embedded logic.
* **Routing System**: Register URL paths with HTTP method handlers.
* **Controller Architecture**: Clean modular logic with view bindings.
* **Session Management**: Now supports in-memory and disk-based storage with LRU logic.
* **Static File Serving**: Serve static, CSS, and JS files from configurable folders.
* **Advanced Template Engine**: PHP-style syntax rendered as Go templates.
* **Request Parsing**: Automatically parses GET and POST parameters.
* **Response Rendering**: Return a `*Template.Response` with typed data.
* **JSON Response**: Response JSON Data on Empty Views of controller, useful for api building, where it converts  `*Template.Response` to JSON
* **Component System**: Define reusable JSON or DB-synced structured data. mostly useful for web components like navigation items settings etc
* **Model migration support**: Auto schema diffing and sync via `Build: false`.
* **Logging**: Error-aware logging with helper functions.
* **Console Commands**: Launch server, generate app structure, run migrations.
* **Environment Overrides**: Override all config settings with environment variables. Supports containers

## Installation

```sh
go get github.com/vrianta/agai
```

Import as needed:

```go
import "github.com/vrianta/agai"
```

## Run Application

After Creating Folder structure

```go
go run . --flags
```

### Flags

You can control various behaviors of the framework using command-line flags when running the application. These flags are especially useful during development and deployment.

| Flag                  | Shortcut | Description                                      |
| --------------------- | -------- | ------------------------------------------------ |
| `--migrate-model`     | `-mm`    | Run model-level database migrations              |
| `--migrate-component` | `-mc`    | Sync components with the database                |
| `--start-server`      | `-ss`    | Start the HTTP server                            |
| `--show-dsn`          | `-sdn`   | Show the Dsn if the database connection failed   |
| `--help`              | `-h`     | Display the list of available command-line flags |

Example Usage
To run model migrations and start the server:

```go
go run . --migrate-model --start-server
```

To sync only the components with the database:

```go
go run . -mc
```

To display help:

```go
go run . -h
```

> **Note:** 
> If no flags are provided, the help message will be shown by default.
> Invalid flags will cause the program to exit with an error.

## Project Structure

```
├───.vscode
├───changelog
└───v1
    ├───component          # Component JSON <-> DB layer
    ├───controller         # Application logic handlers
    ├───cookies            # Cookie utilities
    ├───database           # DB-specific connection and metadata
    ├───internal
    │   └───session        # Session heap, LRU management
    |   └───config         # Configuration and environment overrides
    ├───log                # Logging utility
    ├───model              # Model definitions and schema logic
    ├───render_engine      # PHP-style parsing and template renderer
    ├───response           # HTTP response codes/types
    ├───router             # Internal request router
    ├───server             # Server and entry handler
    ├───smtp               # Email client
    ├───template           # View rendering helpers
    └───utils              # File I/O and helpers
```

## Configuration (`web.config.JSON`)

> **For complete configuration details, environment variable reference, and advanced usage, see [`Config/readme.md`](/v1/internal/config/readme.md). The summary below covers the basics; the linked documentation provides authoritative and up-to-date information.**

Create a `web.config.JSON` file in your project root. Example:
```JSON
{
  "Port": "8080",
  "Host": "localhost",
  "Https": false,
  "Build": false,
  "StaticFolders": ["Static"],
  "CssFolders": ["Css"],
  "JsFolders": ["Js"],
  "ViewFolder": "Views",
  "MaxSessionCount": 1000,
  "SessionStoreType": "memory"
}
```
- **Port**: The port number the server will listen on (e.g., "8080").
- **Host**: The hostname or IP address to bind the server to (e.g., "localhost").
- **Https**: Enable HTTPS server (true/false).
- **Build**: Enable/disable template caching and build mode (true/false).
- **StaticFolders**: List of folders for static files (e.g., ["Static"]).

- **CssFolders**: List of folders for CSS files (e.g., ["Css"]).

- **JsFolders**: List of folders for JS files (e.g., ["Js"]).

- **ViewFolder**: Folder for HTML/PHP templates (e.g., "Views").

- **MaxSessionCount**: Maximum number of concurrent sessions allowed (integer).

- **SessionStoreType**: Type of session store to use (e.g., "memory", "reddis", "database").

> **Note:** redis is not supported yet

You can also override these values using environment variables:
- `SERVER_PORT`
- `SERVER_HOST`
- `SERVER_HTTPS`
- `BUILD`
- `MAX_SESSION_COUNT`
- `SESSION_STORE_TYPE`

Environment variables take precedence over values in `web.config.JSON`.

### Database Configuration (`database.Config.JSON`)

To enable database support, create a `database.config.JSON` file with your database settings, or set the appropriate environment variables.

Example:
```JSON
{
  "Host": "localhost",
  "Port": "3306",
  "User": "root",
  "Password": "",
  "Database": "mydatabase",
  "Protocol": "tcp",
  "Driver": "mysql",
  "SSLMode": "disable"
}
```

#### Supported Environment Variables
- `DB_HOST`
- `DB_PORT`
- `DB_USER`
- `DB_PASSWORD`
- `DB_DATABASE`
- `DB_PROTOCOL`
- `DB_DRIVER`
- `DB_SSLMODE`

Environment variables take precedence over values in `database.config.JSON`.

> **Note:** SSL Mode is not supported yet

## Routing

Routing in this framework is minimal, declarative, and type-safe. You define routes by associating URL paths with controller functions using the `router` module.

### Define Route Handlers

To define a route, use the `Router.New(basePath).RegisterRoutes(...)` chain:

```go 
Router.New("/").RegisterRoutes(  
Router.Route("", Controllers.Home),  
Router.Route("home", Controllers.Home),  
Router.Route("admin", Controllers.Admin),  
Router.Route("login", Controllers.Login),  
Router.Route("logout", Controllers.Logout),  
)  
```

Each `Router.Route(path, handler)` defines a relative path from the base path.

> **Example:** if you want to create api router
```go
Router.New("/api/").RegisterRoutes() // You can add API routes here later
```

### Use Custom Router

To use a custom Routing Handler you need to use function `RegisterCustomRoutingHandler(func(w http.ResponseWriter, r *http.Request))`
> **Note:** To make your custom Handler function compatible with Our Controller Structure make sure it uses out controller.Struct to run the structs
> and to run the controller it should do 
```go
tempController = _controller.Copy() //  copy is important 
tempController.Init(w, r, sess) //  tempController.Run(w, r, sess) run is also valid
```

## Controllers

This section explains how to create controllers and views, including available public variables and methods.

### What is a Controller?
A **Controller** is a Go struct that handles HTTP requests for a specific route. It defines handler functions for HTTP methods (GET, POST, etc.), manages session data, and renders views (templates).

### Controller Structure
A controller is defined as a variable of type `Controller.Struct`. The main public field types and methods are:

#### Public Fields
- **View**: The name of the view (template) directory for this controller. Example: `"home"` (looks for templates in `Views/home/`).
- **GET, POST, DELETE, PATCH, PUT, HEAD, OPTIONS**: Handler functions for each HTTP method. Each receives the controller as `self` and returns a `Template.Response` (a map of data for the template).

#### Public Methods
- `GetInput(key string) interface{}`
- `GetInputs() *map[string]interface{}`
- `StoreData(key string, value any)`
- `GetStoredData(key string) any`
- `Redirect(uri string)`
- `RedirectWithCode(uri string, code Response.Code)`
- `IsLoggedIn() bool`
- `Login() bool`
- `Logout()`

### Example: Creating a Controller
```go
package Home

import (
    "github.com/vrianta/agai/v1/controller"
    "github.com/vrianta/agai/v1/template"
)

var Home = Controller.Struct{
    View: "home",
    GET: func(self *Controller.Struct) *Template.Response {
        return &Template.Response{
            "Title": "Welcome Home",
            "User":  self.GetStoredData("uid"),
        }
    },
    POST: func(self *Controller.Struct) *Template.Response {
        username := self.GetInput("username")
        self.StoreData("uid", username)
        return &Template.Response{
            "Message": "Logged in as " + username.(string),
        }
    },
}
```

> **For full details, advanced usage, and API reference for controllers, see [`controller/readme.md`](/v1/controller/readme.md).**

### Accessing Request Data
- `GetInput(key string)`: Returns a value from GET or POST parameters.
- `GetInputs()`: Returns all request parameters as a map.

### Session Management
- `StoreData(key, value)`: Store data in the session (e.g., user ID).
- `GetStoredData(key)`: Retrieve data from the session.
- `Login()`, `Logout()`, `IsLoggedIn()`: Manage authentication state.

### Redirects
- `Redirect(uri string)`: Redirects to another page.
- `RedirectWithCode(uri, code)`: Redirects with a custom HTTP status code.

### Creating Views

#### View Directory Structure
Each controller where view is mentioned should have a corresponding directory under the `views/` folder, named after the controller's `View` field.

**Example:** If `View: "home"`, templates should be in `Views/home/`.

#### Supported Template Files
- `default.html`, `default.PHP`, or `index.html`/`index.PHP`: Default template for the controller.
- `get.html`, `post.html`, etc.: Templates for specific HTTP methods.

> **Note:** All PHP Features are not supported. Only basic PHP syntaxes are supported to know more check [Write Templates in PHP-Style!](#write-templates-in-PHP-style)

### Location
- The `Views` folder should be in your project root (or as configured in `web.config.JSON` with the `Views_folder` key).

### Template Files
- Place your template files inside the corresponding controller subfolder.
- You can use the following naming conventions:
  - `default.html` or `default.PHP`: The default template for the controller (used if no method-specific template is found).
  - `index.html` or `index.PHP`: Also treated as the default template.
  - `get.html`, `post.html`, `delete.html`, etc.: Templates for specific HTTP methods (GET, POST, DELETE, etc.).

### Example Structure
```
Views/
├── home/
│   ├── default.PHP
│   ├── get.PHP
│   └── post.PHP
├── user/
│   ├── default.html
│   └── get.html
└── shared/
    └── header.PHP
```
- In this example, the `Home` controller with `View: "home"` will use templates from `Views/home/`.
- The framework will automatically select the correct template based on the HTTP method and file availability.

<!-- ### Including Shared Templates -->
<!-- - You can create a `shared/` folder for partials like headers, footers, etc., and include them in your main templates using Go template syntax. -->

## Models

### How to Initialize the Database

To enable database support in your project, you need to Mention host name in the `database.config.JSON`.

> **Note:** If `Host` is not Mentioned in he config, the database will not be initialized. This allows you to run the server without any database connection if desired.


### How to Create a Model (Quick Example)

To define a model, use the `models_handler.New` function, specifying the table name and a map of field types:

```go
import model "github.com/vrianta/agai/v1/model"

var Users = model.New("users", struct {
	UserId    model.Field
	UserName  model.Field
	Password  model.Field
	FirstName model.Field
}{
	UserId: model.Field{
		Type:     model.field types.VarChar,
		Length:   20,
		Nullable: false,
		Index: model.Index{
			PrimaryKey: true,
			Unique:     false,
			Index:      true,
		},
	},
	UserName: model.Field{
		Type:     model.field types.VarChar,
		Length:   30,
		Nullable: false,
		Index: model.Index{
			Unique: true,
			Index:  true,
		},
	},
	Password: model.Field{
		Type:     model.field types.Text,
		Nullable: false,
	},
	FirstName: model.Field{
		Type:     model.field types.Text,
		Nullable: false,
	},
})
```
- The first argument is the table name in your database.
- The second argument Struct where you define the table

**For full documentation, advanced usage, and API reference, see [`modelsHandler/readme.md`](/v1/model/readme.md).**



### Model Query Builder System

The framework includes a powerful, human-friendly queryBuilder builder called **ModelsHandler** for working with your database using Go structs. ModelsHandler provides a chainable API for building and executing SQL queries (SELECT, UPDATE, DELETE) in a style similar to popular ORMs.

- Define your model as Go structs and map them to database tables.
- Build queries using a fluent, readable API (e.g., `Users.Get().Where("age").GreaterThan(18).OrderBy("name").Fetch()`).
- Supports WHERE, AND, OR, IN, NOT IN, BETWEEN, LIKE, IS NULL, LIMIT, OFFSET, ORDER BY, GROUP BY, and more.
- Makes database access easy to read, write, and maintain.

> **Note:** Migration only happens if you pass the required flags while running the application

### Model Query Result

Result is going to be a `map[string]any` 

* **Get list of Data**: After Building the query if you use `Fetch()` -> it return a array of map[string]any which resembles the Fields
* **Get Single Element:** In some cases use might need only first elemet of the result then he/she should use `First()` after the query build
* **Delete Element:** `.Delete()` — Execute DELETE
* **Only Execute the Query:** `.Exec()` — Execute UPDATE

> **For More Information Check** [Model-Doc](/v1/model/readme.md)

## Components

`Components` are data holder for a specific table, where it can hold a set of the data which needs to be updated in the database
`Use Cases` are Navigation items, settings template

> **Note:** If the table is not created in the DB then the DB is populated with the template file, it will check content of the template has been changed then it will update the local
> If the Config file has a new element which is not present in the data base then it will update the data base

### Structure

Component structures have to follow strict rule to comply

* Every Component Model should have a `Primary Key`
* It is recommended to have the `Primary Key` as `Integer` so that it can follow same orientation as you mentioned in the config file
* Every `JSON` element should have a key which should be the value of `Primary Key`
* File Name should be like `table_name`.component.JSON

### Example

#### File Name 

```
navigation_items.component.JSON
```

#### Component Model

```go
package models

import (
	"github.com/vrianta/agai/v1/model"
)

var Nav_items = model.New("navigation_items", struct {
	Id       model.Field
	Name     model.Field
	Href     model.Field
	Disabled model.Field
	Dropdown model.Field
}{
	Id: model.Field{
		Nullable: false,
		Type:     model.field types.Int,
		Length:   10,
		Index: model.Index{
			PrimaryKey: true,
			Index:      true,
		},
	},
	Name: model.Field{
		Nullable: false,
		Type:     model.field types.VarChar,
		Length:   10,
	},
	Href: model.Field{
		Nullable: false,
		Type:     model.field types.Text,
	},
	Disabled: model.Field{
		Nullable:     false,
		DefaultValue: "0",
		Type:         model.field types.Bool,
	},
	Dropdown: model.Field{
		Nullable:     true,
		DefaultValue: "",
		Type:         model.field types.JSON,
	},
})

```

#### Component

```JSON
{
  "0": {
    "Disabled": 0,
    "Dropdown": null,
    "Href": "#home",
    "Id": 0,
    "Name": "Home"
  },
  "1": {
    "Disabled": 0,
    "Dropdown": null,
    "Href": "#about-me",
    "Id": 1,
    "Name": "About Me"
  },
  "2": {
    "Disabled": 0,
    "Dropdown": null,
    "Href": "#skills",
    "Id": 2,
    "Name": "Skills"
  },
  "3": {
    "Disabled": 0,
    "Dropdown": null,
    "Href": "#experience",
    "Id": 3,
    "Name": "Experience"
  },
  "4": {
    "Disabled": 0,
    "Dropdown": null,
    "Href": "#projects",
    "Id": 4,
    "Name": "Projects"
  },
  "5": {
    "Disabled": 0,
    "Dropdown": null,
    "Href": "#contact-me",
    "Id": 5,
    "Name": "Contact"
  }
}
```


### Get Controller 

To Get the values of controller do `model.GetComponents()` - it will return map[component_index]componentvalue

```go
models.Nav_items.GetComponents()
```


## Static, CSS, and JS File Serving

- Place static files in folders listed in `web.config.JSON`.
- Access them via `/static/filename.ext`, `/css/style.css`, etc.
- Files are cached for performance.

## SMTP/Email Support

Send emails easily:
```go
import "github.com/vrianta/agai/v1/smtp"
smtp.Client.InitSMTPClient("smtp.example.com", 587, "user", "pass")
err := smtp.Client.SendMail([]string{"to@example.com"}, "Subject", "Body")
```

## Template Engine & PHP Parsing Syntax

### Write Templates in PHP-Style!
- Place templates in the `Views` folder (or as configured).
- Use `.PHP` or `.html` extensions.

#### Supported Syntax
- **Echo:** `<?= $var ?>` → `{{ .var }}`
- **PHP Block:** `<?PHP ... ?>` for logic
- **Variables:**
  - `$$var` refers to variables passed from your Go controller as part of the `Template.Response` (the data map returned from your handler).
  - `$var` refers to variables created and used locally within the PHP template file itself (e.g., inside a foreach or assigned in the template logic).
  - You can also use `$obj->prop`, `$arr['key']`, `$arr[0]` for object and array access.
- **If/Else:**
  ```PHP
  <?PHP if ($$user): ?>
    Hello, <?= $$user ?>
  <?PHP elseif ($$guest): ?>
    Welcome, Guest!
  <?PHP else: ?>
    Please log in.
  <?PHP endif; ?>
  ```
- **Loops:**
  ```PHP
  <?PHP foreach ($$items as $item): ?>
    <?= $item ?>
  <?PHP endforeach; ?>
  ```
- **Function Calls:**
  - `strtoupper($var)` → `upper .var`
  - `strtolower($var)` → `lower .var`
  - `strlen($var)`/`count($arr)` → `len .arr`
  - `isset($var)` → `ne .var nil`
  - `empty($var)` → `eq .var ""`
- **Operators:**
  - `==` → `eq`, `!=` → `ne`, `<` → `lt`, `>` → `gt`, `<=` → `le`, `>=` → `ge`
  - `&&` → `and`, `||` → `or`, `!` → `not`
- **Comments:** `// comment` inside `<?PHP ... ?>` blocks

#### Example Template
```PHP
<!-- Views/home.PHP -->
<h1>Welcome</h1>
<?PHP if ($user): ?>
  <p>Hello, <?= $user ?>!</p>
<?PHP else: ?>
  <p>Please log in.</p>
<?PHP endif; ?>
<ul>
<?PHP foreach ($items as $item): ?>
  <li><?= $item ?></li>
<?PHP endforeach; ?>
</ul>
```

#### How It Works
- The server automatically parses PHP-style templates and converts them to Go's `html/template` syntax.
- You can use all Go template features in addition to the PHP-like syntax.

## Frequently Asked Questions (FAQ)

**Q: How should I name and organize my controller's View field and view folder?**
- The `View` field in your controller should match the subfolder name in the `Views` directory (case-sensitive). For example, `View: "home"` expects templates in `Views/home/`.
- Do not include a file extension in the `View` field; it should be just the folder name.

**Q: How does template resolution work?**
- For each HTTP method (GET, POST, etc.), the framework looks for a file named `get.html`, `post.html`, etc., in the controller's view folder.
- If no method-specific template is found, it falls back to `default.html`, `default.PHP`, `index.html`, or `index.PHP`.
- If none of these exist, the server will panic with an error indicating a missing default view.

**Q: Can I use multiple views per controller?**
- Each controller is associated with a single view folder. If you need multiple views, create multiple controllers or use conditional logic in your handler to select data/templates.

**Q: Are templates hot-reloaded in development?**
- In development mode (`Build: false` in `Config.JSON`), templates are reloaded on each request. In build/production mode, templates are cached for performance.

**Q: How do I include shared templates (partials)?**
- Place shared templates (e.g., header, footer) in a `Views/shared/` folder.
- Use Go template syntax to include them: `{{ template "shared/header.html" . }}`

**Q: What happens if a template file is missing or there is a rendering error?**
- If the default view is missing, the server will panic and log an error. Rendering errors will also panic and log details. Always check your logs for troubleshooting.

**Q: Can I pass complex data (arrays, maps, structs) to templates?**
- Yes, you can pass any Go data structure in the `Template.Response` map. Use dot notation in templates to access nested field types.

**Q: How is session data secured?**
- Sessions use secure, random IDs stored in HTTP-only cookies. For best security, run your server over HTTPS and set the `Secure` flag in your cookie configuration.
- There is no built-in CSRF protection; you should implement CSRF tokens in your forms if needed.

**Q: Can I serve static files from subfolders?**
- Yes, you can organize static files in subfolders within your static directories. There are no restrictions on file types or folder depth.

**Q: Is there middleware support?**
- The framework does not have a formal middleware system, but you can add logic in your controller handlers or wrap the main handler in `server.go` for global middleware.

**Q: How do I deploy to production?**
- Set `Build: true` in `Config.JSON` for template caching and performance.
- Use a reverse proxy (like Nginx) for HTTPS and static file serving if needed.
- Monitor logs and handle panics gracefully in production.

**Q: How can I extend the framework?**
- You can add new packages, extend controllers, or modify the template engine. For advanced features (like WebSockets), integrate with Go's standard libraries and register your handlers in `server.go`.

## License
See [LICENSE](LICENSE) for GPLv3 license details.

