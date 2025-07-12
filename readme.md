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
  - [Project Structure](#project-structure)
  - [ModelsHandler (ORM-like queryBuilder Builder)](#modelshandler-orm-like-querybuilder-builder)
    - [How to Create a Model (Quick Example)](#how-to-create-a-model-quick-example)
  - [Installation](#installation-1)
  - [Project Structure](#project-structure-1)
  - [Configuration (`web.config.json`)](#configuration-webconfigjson)
    - [Database Configuration (`Database.Config.json`)](#database-configuration-databaseconfigjson)
      - [Supported Environment Variables](#supported-environment-variables)
  - [Configuration and the Config Package](#configuration-and-the-config-package)
  - [Database Initialization](#database-initialization)
    - [How to Initialize the Database](#how-to-initialize-the-database)
  - [Server Creation \& Routing](#server-creation--routing)
    - [1. Define Route Handlers](#1-define-route-handlers)
  - [Creating Controllers and Views](#creating-controllers-and-views)
    - [What is a Controller?](#what-is-a-controller)
    - [Controller Structure](#controller-structure)
      - [Public FieldTypes](#public-fieldtypes)
      - [Public Methods](#public-methods)
    - [Example: Creating a Controller](#example-creating-a-controller)
    - [Accessing Request Data](#accessing-request-data)
    - [Session Management](#session-management)
    - [Redirects](#redirects)
    - [Creating Views](#creating-views)
      - [View Directory Structure](#view-directory-structure)
      - [Supported Template Files](#supported-template-files)
      - [Template Syntax](#template-syntax)
  - [Views Folder Structure](#views-folder-structure)
    - [Location](#location)
    - [Organization](#organization)
    - [Template Files](#template-files)
    - [Example Structure](#example-structure)
    - [Including Shared Templates](#including-shared-templates)
  - [Session Management](#session-management-1)
  - [Static, CSS, and JS File Serving](#static-css-and-js-file-serving)
  - [SMTP/Email Support](#smtpemail-support)
  - [Console Commands](#console-commands)
  - [Template Engine \& PHP Parsing Syntax](#template-engine--php-parsing-syntax)
    - [Write Templates in PHP-Style!](#write-templates-in-php-style)
      - [Supported Syntax](#supported-syntax)
      - [Example Template](#example-template)
      - [How It Works](#how-it-works)
  - [API Reference](#api-reference)
  - [License](#license)
  - [Frequently Asked Questions (FAQ)](#frequently-asked-questions-faq)

## Features

* **Custom HTTP Server**: Easily start and stop the server with CLI or embedded logic.
* **Routing System**: Register URL paths with HTTP method handlers.
* **Controller Architecture**: Clean modular logic with view bindings.
* **Session Management**: Now supports in-memory and disk-based storage with LRU logic.
* **Static File Serving**: Serve static, CSS, and JS files from configurable folders.
* **Advanced Template Engine**: PHP-style syntax rendered as Go templates.
* **Request Parsing**: Automatically parses GET and POST parameters.
* **Response Rendering**: Return a `*Template.Response` with typed data.
* **Component System**: Define reusable JSON or DB-synced structured data. mostly usefull for web components like navigation items settings etc
* **Model Migration Support**: Auto schema diffing and sync via `Build: false`.
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

---

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

---


## ModelsHandler (ORM-like queryBuilder Builder)

The framework includes a powerful, human-friendly queryBuilder builder called **ModelsHandler** for working with your database using Go structs. ModelsHandler provides a chainable API for building and executing SQL queries (SELECT, UPDATE, DELETE) in a style similar to popular ORMs.

- Define your model as Go structs and map them to database tables.
- Build queries using a fluent, readable API (e.g., `Users.Get().Where("age").GreaterThan(18).OrderBy("name").Fetch()`).
- Supports WHERE, AND, OR, IN, NOT IN, BETWEEN, LIKE, IS NULL, LIMIT, OFFSET, ORDER BY, GROUP BY, and more.
- Makes database access easy to read, write, and maintain.
- **Automatic Database Migration:** If the `Build` flag is set to `false` in your `web.config.json`, ModelsHandler will automatically migrate your database schema to match your model definitions. This means tables and columns are created or updated as needed, so you don't have to write migration scripts manually.

> **Note:** Migration only happens if the `Build` flag is `false`. In production, set it to `true` to prevent accidental schema changes.

### How to Create a Model (Quick Example)

To define a model, use the `models_handler.New` function, specifying the table name and a map of FieldTypes:

```go
import models_handler "github.com/vrianta/agai/v1/modelsHandler"

var Users = models_handler.New(
    "users", // Table name
    map[string]models_handler.Field{
        "userId": {
            Name:     "userId",
            Type:     models_handler.FieldTypesTypes.VarChar,
            Length:   20,
            Nullable: false,
        },
        "userName": {
            Name:     "userName",
            Type:     models_handler.FieldTypesTypes.VarChar,
            Length:   30,
            Nullable: false,
        },
    },
)
```
- The first argument is the table name in your database.
- The second argument is a map where each key is a column name and the value is a `Field` struct describing the column.
- You can add more FieldTypes and options (like indexes, types, etc.) as needed.

**For full documentation, advanced usage, and API reference, see [`modelsHandler/readme.md`](modelsHandler/readme.md).**

---

## Installation

1. Clone the repository or add it to your Go project:
   ```sh
   go get github.com/vrianta/agai
   ```
2. Import the package:
   ```go
   import "github.com/vrianta/agai"
   ```

---

## Project Structure

```
.
├── Config/           # Configuration loader (Config.go, type.go, var.go)
├── Controller/       # Route handler logic (Controller.go, type.go, var.go)
├── Cookies/          # Cookie utilities (Cookies.go, type.go, var.go)
├── Log/              # Logging utilities (Write.go, type.go, var.go)
├── Redirect/         # HTTP redirects (Redirect.go, type.go)
├── RenderEngine/     # Template engine (PHP-like syntax) (RenderEngine.go, type.go, var.go)
├── Response/         # Response codes/types (type.go, var.go)
├── Router/           # HTTP router (Router.go, type.go, var.go)
├── Session/          # Session management (Session.go, type.go)
├── smtp/             # SMTP client (client.go)
├── Template/         # Template helpers (template.go)
├── Utils/            # File and utility helpers (file.handler.go, util.go)
├── console.go        # Interactive console
├── server.go         # Server entry point
├── types.go          # Core types
├── vars.go           # Global variables
└── readme.md         # This guide
```

---

## Configuration (`web.config.json`)

> **For complete configuration details, environment variable reference, and advanced usage, see [`Config/readme.md`](Config/readme.md). The summary below covers the basics; the linked documentation provides authoritative and up-to-date information.**

Create a `Config.json` file in your project root. Example:
```json
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

- **SessionStoreType**: Type of session store to use (e.g., "memory", "redis", "database").

You can also override these values using environment variables:
- `SERVER_PORT`
- `SERVER_HOST`
- `SERVER_HTTPS`
- `BUILD`
- `MAX_SESSION_COUNT`
- `SESSION_STORE_TYPE`

Environment variables take precedence over values in `Config.json`.

### Database Configuration (`Database.Config.json`)

To enable database support, create a `Database.Config.json` file with your database settings, or set the appropriate environment variables.

Example:
```json
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

Environment variables take precedence over values in `database.config.json`.

---

## Configuration and the Config Package

All configuration for the server is managed by the **Config package**. This package loads settings from config files (such as `Config.json` and `Database.Config.json`) and supports overriding them with environment variables. The Config package ensures that your server is flexible and easy to configure for different environments (development, production, etc.).

- **Config file names, supported environment variables, and override order are fully documented in [`Config/readme.md`](Config/readme.md).**
- **Environment variables always take precedence over config file values.**
- For best practices, advanced usage, and troubleshooting, see the [Config package documentation](Config/readme.md).

**Quick Reference:**
- Main server config: `Config.json` (see example above)
- Database config: `Database.Config.json` (see example above)
- Supported environment variables: see [`Config/readme.md`](Config/readme.md)
- Override order: Environment variables > Config files > Defaults

For a complete guide to all configuration options, environment variable names, and advanced usage, please refer to [`Config/readme.md`](Config/readme.md).

---

## Database Initialization

To enable database support in your project, you need to initialize the database connection before starting the server. This is done by calling the `InitDatabase` method on your server instance.

### How to Initialize the Database

1. **Create a `Database.Config.json` file** in your project root with your database settings, or set the appropriate environment variables (see above for details).
2. **Call `InitDatabase()` before starting the server:**

   ```go
   package main

   import (
       "github.com/vrianta/agai"
   )

   func main() {
       srv := Server.New()
       srv.InitDatabase() // Initialize the database connection
       srv.Start()        // Start the server
   }
   ```

3. **If you do not want to use a database**, simply do not call `InitDatabase()`. The server will run without attempting to connect to any database.

> **Note:** If `InitDatabase()` is not called, the database will not be initialized. This allows you to run the server without any database connection if desired.

---

## Server Creation & Routing

### 1. Define Route Handlers

Each handler is a Go package (usually in `Controller/`) that exports a variable of type `Controller.Struct` with FieldTypes for the view and HTTP methods. Methods are functions that receive a pointer to the controller struct and return a `*Template.Response` (for GET) or handle logic for POST/DELETE.

Example:
```go
package Home

import (
	components "github.com/pritam-is-next/resume/Components"
	Controller "github.com/vrianta/agai/v1/Controller"
	"github.com/vrianta/agai/v1/internal/session"
	"github.com/vrianta/agai/v1/Template"
)

var Home = Controller.Struct{
	View: "home.php",
	GET:  func(self *Controller.Struct) *Template.Response {
	response := &Template.Response{
		"Title":          "Pritam Dutta",
		"Heading":        "Pritam Dutta",
		"NavItems":       components.NavItems,
		"Hero":           components.Hero,
		"AboutMe":        components.AboutMe,
		"Skills":         components.Skills,
		"Experiences":    components.Experiences,
		"Projects":       components.Projects,
		"ContactDetails": components.ContactDetails,
	}
	return response
},
}

```

- The `View` field specifies the template to render (e.g., `home.php`).
- The `GET`, `POST`, and `DELETE` FieldTypes are function handlers for each HTTP method.
- The `GET` handler returns a `*Template.Response` (a map of data for the template).
- You can import and use components or data as needed.

---

## Creating Controllers and Views

This section explains how to create controllers and views, including available public variables and methods.

### What is a Controller?
A **Controller** is a Go struct that handles HTTP requests for a specific route. It defines handler functions for HTTP methods (GET, POST, etc.), manages session data, and renders views (templates).

### Controller Structure
A controller is defined as a variable of type `Controller.Struct`. The main public FieldTypes and methods are:

#### Public FieldTypes
- **View**: The name of the view (template) directory for this controller. Example: `"home"` (looks for templates in `Views/home/`).
- **GET, POST, DELETE, PATCH, PUT, HEAD, OPTIONS**: Handler functions for each HTTP method. Each receives the controller as `self` and returns a `Template.Response` (a map of data for the template).

#### Public Methods
- `InitWR(w http.ResponseWriter, r *http.Request)`
- `InitSession(session *Session.Struct)`
- `RunRequest(session *Session.Struct)`
- `RegisterTemplate() error`
- `ExecuteTemplate(template *Template.Struct, response *Template.Response) error`
- `GetInput(key string) interface{}`
- `GetInputs() *map[string]interface{}`
- `StoreData(key string, value any)`
- `GetStoredData(key string) any`
- `Redirect(uri string)`
- `WithCode(uri string, code Response.Code)`
- `IsLoggedIn() bool`
- `Login() bool`
- `Logout()`

### Example: Creating a Controller
```go
package Home

import (
    "github.com/vrianta/agai/v1/Controller"
    "github.com/vrianta/agai/v1/Template"
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

> **For full details, advanced usage, and API reference for controllers, see [`Controller/readme.md`](Controller/readme.md).**

### Accessing Request Data
- `GetInput(key string)`: Returns a value from GET or POST parameters.
- `GetInputs()`: Returns all request parameters as a map.

### Session Management
- `StoreData(key, value)`: Store data in the session (e.g., user ID).
- `GetStoredData(key)`: Retrieve data from the session.
- `Login()`, `Logout()`, `IsLoggedIn()`: Manage authentication state.

### Redirects
- `Redirect(uri string)`: Redirects to another page.
- `WithCode(uri, code)`: Redirects with a custom HTTP status code.

### Creating Views
#### View Directory Structure
Each controller should have a corresponding directory under the `Views/` folder, named after the controller's `View` field.

**Example:** If `View: "home"`, templates should be in `Views/home/`.

#### Supported Template Files
- `default.html`, `default.php`, or `index.html`/`index.php`: Default template for the controller.
- `get.html`, `post.html`, etc.: Templates for specific HTTP methods.

#### Template Syntax
- PHP-style syntax is supported:
  - `<?= $var ?>` → `{{ .var }}`
  - `<?php if ($user): ?> ... <?php endif; ?>`
  - Loops: `<?php foreach ($items as $item): ?> ... <?php endforeach; ?>`

---

## Views Folder Structure

The `Views` folder contains all your HTML/PHP templates. Proper organization is important for the framework to locate and render the correct templates for each controller and HTTP method.

### Location
- The `Views` folder should be in your project root (or as configured in `Config.json` with the `Views_folder` key).

### Organization
- Each controller should have its own subfolder inside `Views`, named after the controller's `View` field (without file extension).
- For example, if your controller has `View: "home"`, create a folder `Views/home/`.

### Template Files
- Place your template files inside the corresponding controller subfolder.
- You can use the following naming conventions:
  - `default.html` or `default.php`: The default template for the controller (used if no method-specific template is found).
  - `index.html` or `index.php`: Also treated as the default template.
  - `get.html`, `post.html`, `delete.html`, etc.: Templates for specific HTTP methods (GET, POST, DELETE, etc.).

### Example Structure
```
Views/
├── home/
│   ├── default.php
│   ├── get.php
│   └── post.php
├── user/
│   ├── default.html
│   └── get.html
└── shared/
    └── header.php
```
- In this example, the `Home` controller with `View: "home"` will use templates from `Views/home/`.
- The framework will automatically select the correct template based on the HTTP method and file availability.

### Including Shared Templates
- You can create a `shared/` folder for partials like headers, footers, etc., and include them in your main templates using Go template syntax.

---

## Session Management

Session data and helpers are accessed via the `self.Session` variable inside your controller methods.

- Access POST/GET data:
  ```go
  uid := self.GetInput("uid")
  token := self.GetInput("token")
  ```
- Store/retrieve session variables:
  ```go
  self.Session.Store["uid"] = "user123"
  user := self.Session.Store["uid"]
  ```
- Login/logout:
  ```go
  self.Login()
  self.Logout()
  if self.IsLoggedIn() { /* ... */ }
  ```

---

## Static, CSS, and JS File Serving

- Place static files in folders listed in `Config.json`.
- Access them via `/Static/filename.ext`, `/Css/style.css`, etc.
- Files are cached for performance.

---

## SMTP/Email Support

Send emails easily:
```go
import "github.com/vrianta/agai/v1/smtp"
smtp.Client.InitSMTPClient("smtp.example.com", 587, "user", "pass")
err := smtp.Client.SendMail([]string{"to@example.com"}, "Subject", "Body")
```

---

## Console Commands

When the server is running, use these commands in the console:
- `start`    - Start the server
- `stop`     - Stop the server
- `restart`  - Restart the server
- `r`        - Shortcut for restart
- `exit`     - Stop and exit
- `-h`       - Help

---

## Template Engine & PHP Parsing Syntax

### Write Templates in PHP-Style!
- Place templates in the `Views` folder (or as configured).
- Use `.php` or `.html` extensions.

#### Supported Syntax
- **Echo:** `<?= $var ?>` → `{{ .var }}`
- **PHP Block:** `<?php ... ?>` for logic
- **Variables:**
  - `$$var` refers to variables passed from your Go controller as part of the `Template.Response` (the data map returned from your handler).
  - `$var` refers to variables created and used locally within the PHP template file itself (e.g., inside a foreach or assigned in the template logic).
  - You can also use `$obj->prop`, `$arr['key']`, `$arr[0]` for object and array access.
- **If/Else:**
  ```php
  <?php if ($$user): ?>
    Hello, <?= $$user ?>
  <?php elseif ($$guest): ?>
    Welcome, Guest!
  <?php else: ?>
    Please log in.
  <?php endif; ?>
  ```
- **Loops:**
  ```php
  <?php foreach ($$items as $item): ?>
    <?= $item ?>
  <?php endforeach; ?>
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
- **Comments:** `// comment` inside `<?php ... ?>` blocks

#### Example Template
```php
<!-- Views/home.php -->
<h1>Welcome</h1>
<?php if ($user): ?>
  <p>Hello, <?= $user ?>!</p>
<?php else: ?>
  <p>Please log in.</p>
<?php endif; ?>
<ul>
<?php foreach ($items as $item): ?>
  <li><?= $item ?></li>
<?php endforeach; ?>
</ul>
```

#### How It Works
- The server automatically parses PHP-style templates and converts them to Go's `html/template` syntax.
- You can use all Go template features in addition to the PHP-like syntax.

---

## API Reference

- `server.New(host, port, routes, config)` - Create a new server instance. `host` and `port` specify the address, `routes` is a map of URL paths to controller structs, and `config` is a pointer to your configuration (or nil for default).
- `server.Start()` - Start the server and launch the interactive console.
- `Controller.Struct` - The base struct for all controllers. FieldTypes include:
  - `View`: The template file to render (e.g., `home.php`).
  - `GET`, `POST`, `DELETE`: Function handlers for each HTTP method.
  - `Session`: The session object for the current request.
- `self.Session` - Access session data and helpers inside controller methods.
  - `.POST` / `.GET`: Maps for POST/GET parameters.
  - `.Store`: Map for session variables (persisted across requests).
  - `.Login(userID)`, `.Logout(redirectURL)`, `.IsLoggedIn()`: Session authentication helpers.
- `Session.RenderEngine` - For rendering responses and templates.
  - `.Render(str)`: Write a string to the response.
  - `.RenderTemplate(view, data)`: Render a template with data.
  - `.RenderError(msg, code)`: Render an error response.
  - `.StartRender()`: Flush the response buffer.
- `server.RemoveSession(sessionID)` - Remove a session by its ID.
- `smtp.Client` - Built-in SMTP client for sending emails.
  - `.InitSMTPClient(host, port, user, pass)`
  - `.SendMail(to, subject, body)`

For more details, see the source files in each package (Controller, Session, RenderEngine, etc.).

---

## License
See [LICENSE](LICENSE) for GPLv3 license details.

---

## Frequently Asked Questions (FAQ)

**Q: How should I name and organize my controller's View field and view folder?**
- The `View` field in your controller should match the subfolder name in the `Views` directory (case-sensitive). For example, `View: "home"` expects templates in `Views/home/`.
- Do not include a file extension in the `View` field; it should be just the folder name.

**Q: How does template resolution work?**
- For each HTTP method (GET, POST, etc.), the framework looks for a file named `get.html`, `post.html`, etc., in the controller's view folder.
- If no method-specific template is found, it falls back to `default.html`, `default.php`, `index.html`, or `index.php`.
- If none of these exist, the server will panic with an error indicating a missing default view.

**Q: Can I use multiple views per controller?**
- Each controller is associated with a single view folder. If you need multiple views, create multiple controllers or use conditional logic in your handler to select data/templates.

**Q: Are templates hot-reloaded in development?**
- In development mode (`Build: false` in `Config.json`), templates are reloaded on each request. In build/production mode, templates are cached for performance.

**Q: How do I include shared templates (partials)?**
- Place shared templates (e.g., header, footer) in a `Views/shared/` folder.
- Use Go template syntax to include them: `{{ template "shared/header.html" . }}`

**Q: What happens if a template file is missing or there is a rendering error?**
- If the default view is missing, the server will panic and log an error. Rendering errors will also panic and log details. Always check your logs for troubleshooting.

**Q: Can I pass complex data (arrays, maps, structs) to templates?**
- Yes, you can pass any Go data structure in the `Template.Response` map. Use dot notation in templates to access nested FieldTypes.

**Q: How is session data secured?**
- Sessions use secure, random IDs stored in HTTP-only cookies. For best security, run your server over HTTPS and set the `Secure` flag in your cookie configuration.
- There is no built-in CSRF protection; you should implement CSRF tokens in your forms if needed.

**Q: Can I serve static files from subfolders?**
- Yes, you can organize static files in subfolders within your static directories. There are no restrictions on file types or folder depth.

**Q: Is there middleware support?**
- The framework does not have a formal middleware system, but you can add logic in your controller handlers or wrap the main handler in `server.go` for global middleware.

**Q: How do I deploy to production?**
- Set `Build: true` in `Config.json` for template caching and performance.
- Use a reverse proxy (like Nginx) for HTTPS and static file serving if needed.
- Monitor logs and handle panics gracefully in production.

**Q: How can I extend the framework?**
- You can add new packages, extend controllers, or modify the template engine. For advanced features (like WebSockets), integrate with Go's standard libraries and register your handlers in `server.go`.

---

For more details on configuration options, environment variables, and advanced usage, see [`Config/readme.md`](Config/readme.md).

---
