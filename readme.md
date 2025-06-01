# Go Server Framework - Complete User Guide

Welcome to the Go Server Framework! This guide will help you set up, configure, and use every feature of your server, including advanced PHP-style template parsing, session management, static file serving, and more.

---

## Table of Contents
1. [Features](#features)
2. [Installation](#installation)
3. [Project Structure](#project-structure)
4. [Configuration (`Config.json`)](#configuration-configjson)
5. [Server Creation & Routing](#server-creation--routing)
6. [Session Management](#session-management)
7. [Static, CSS, and JS File Serving](#static-css-and-js-file-serving)
8. [SMTP/Email Support](#smtpemail-support)
9. [Console Commands](#console-commands)
10. [Template Engine & PHP Parsing Syntax](#template-engine--php-parsing-syntax)
11. [API Reference](#api-reference)
12. [License](#license)

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

## Session Management

Session data and helpers are accessed via the `self.Session` variable inside your controller methods.

- Access POST/GET data:
  ```go
  uid := self.Session.POST["uid"]
  token := self.Session.GET["token"]
  ```
- Store/retrieve session variables:
  ```go
  self.Session.Store["uid"] = "user123"
  user := self.Session.Store["uid"]
  ```
- Login/logout:
  ```go
  self.Session.Login("user123")
  self.Session.Logout("/login")
  if self.Session.IsLoggedIn() { /* ... */ }
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
  <?php if ($user): ?>
    Hello, <?= $user ?>
  <?php elseif ($guest): ?>
    Welcome, Guest!
  <?php else: ?>
    Please log in.
  <?php endif; ?>
  ```
- **Loops:**
  ```php
  <?php foreach ($items as $item): ?>
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
