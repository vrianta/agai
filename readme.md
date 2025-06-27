# Go Server Framework - Complete User Guide

Welcome to the Go Server Framework! This guide will help you set up, configure, and use every feature of your server, including advanced PHP-style template parsing, session management, static file serving, and more.

---

## Table of Contents
1. [Features](#features)
2. [Installation](#installation)
3. [Project Structure](#project-structure)
4. [Configuration (`Config.json`)](#configuration-configjson)
5. [Server Creation & Routing](#server-creation--routing)
6. [Creating Controllers and Views](#creating-controllers-and-views)
7. [Session Management](#session-management)
8. [Static, CSS, and JS File Serving](#static-css-and-js-file-serving)
9. [SMTP/Email Support](#smtpemail-support)
10. [Console Commands](#console-commands)
11. [Template Engine & PHP Parsing Syntax](#template-engine--php-parsing-syntax)
12. [API Reference](#api-reference)
13. [License](#license)

---

## Features
- **Custom HTTP Server**: Easily start and stop the server with interactive console commands. Configurable via `Config.json` (HTTP/HTTPS, static/view folders, build mode).
- **Routing System**: Map URL paths to controller structs. Supports GET, POST, DELETE HTTP methods. Dynamic handler invocation based on request method.
- **Controller Architecture**: Modular controller packages. Each controller defines its own view and HTTP method handlers. Handlers return data for templates or perform logic.
- **Session Management**: Secure, cookie-based session tracking. Session creation, retrieval, update, and destruction. Session variables (`Store` map) for user data. Login/logout helpers and authentication checks. Session expiry and cleanup mechanism.
- **Static File Serving**: Serve static, CSS, and JS files from configurable folders. Static file caching with last-modified checks. Efficient file read and cache update logic.
- **Advanced Template Engine**: Write templates in PHP-style syntax (`<?= $var ?>`, `<?php ... ?>`). Automatic conversion to Go’s `html/template` syntax. Supports variables, loops, conditionals, and custom operators. Template caching and reloading on file change.
- **Request Parsing**: Automatic parsing of GET and POST parameters. Easy access to request data in controllers.
- **Response Rendering**: Render templates with data from controllers. Render plain strings or error responses. Disable client-side caching for sensitive pages.
- **SMTP/Email Support**: Built-in SMTP client for sending emails. Configurable SMTP server, port, user, and password.
- **Logging**: Centralized logging for errors and server events.
- **Console Commands**: Start, stop, restart, and exit server from the console. Help command for available options.
- **Utilities**: Cookie management helpers. File utilities for reading and caching. Type definitions for routes, sessions, and templates.
- **Extensible & Secure**: Easily add new controllers, routes, and templates. Secure session IDs, cache control, and best practices. Template and static file caching, mutex-protected maps, and efficient request handling.

---

## Installation

1. Clone the repository or add it to your Go project:
   ```sh
   go get github.com/vrianta/Server
   ```
2. Import the package:
   ```go
   import "github.com/vrianta/Server"
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

## Configuration (`Config.json`)

Create a `Config.json` file in your project root. Example:
```json
{
  "Http": true,
  "Static_folders": ["Static"],
  "CSS_Folders": ["Css"],
  "JS_Folders": ["Js"],
  "Views_folder": "Views",
  "Build": false
}
```
- **Http**: Enable HTTP server
- **Build**: Ensure if it is a Build Version or not
- **Static_folders**: List of folders for static files
- **CSS_Folders**: List of folders for CSS
- **JS_Folders**: List of folders for JS
- **Views_folder**: Folder for HTML/PHP templates
- **Build**: Enable/disable template caching

---

## Server Creation & Routing

### 1. Define Route Handlers

Each handler is a Go package (usually in `Controller/`) that exports a variable of type `Controller.Struct` with fields for the view and HTTP methods. Methods are functions that receive a pointer to the controller struct and return a `*Template.Response` (for GET) or handle logic for POST/DELETE.

Example:
```go
package Home

import (
	components "github.com/pritam-is-next/resume/Components"
	Controller "github.com/vrianta/Server/Controller"
	"github.com/vrianta/Server/Session"
	"github.com/vrianta/Server/Template"
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
- The `GET`, `POST`, and `DELETE` fields are function handlers for each HTTP method.
- The `GET` handler returns a `*Template.Response` (a map of data for the template).
- You can import and use components or data as needed.

---

## Creating Controllers and Views

This section explains how to create controllers and views, including available public variables and methods.

### What is a Controller?
A **Controller** is a Go struct that handles HTTP requests for a specific route. It defines handler functions for HTTP methods (GET, POST, etc.), manages session data, and renders views (templates).

### Controller Structure
A controller is defined as a variable of type `Controller.Struct`. The main public fields and methods are:

#### Public Fields
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
    "github.com/vrianta/Server/Controller"
    "github.com/vrianta/Server/Template"
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
import "github.com/vrianta/Server/smtp"
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
- `Controller.Struct` - The base struct for all controllers. Fields include:
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
- Yes, you can pass any Go data structure in the `Template.Response` map. Use dot notation in templates to access nested fields.

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
