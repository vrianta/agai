# Go Server Package

This Go server package provides a flexible way to create and start a web server with configurable routes, session management, static file serving, and more.

## Features

- **Configurable Routes**: Map URLs to handler structs with GET, POST, DELETE methods.
- **Session Handling**: Built-in session management with secure session IDs and cookie handling.
- **Static, CSS, and JS File Serving**: Serve static assets from configurable folders.
- **Template Rendering**: Render HTML templates with caching and auto-reload on file changes.
- **Customizable**: Extend the server with your own handlers and logic.
- **Console Control**: Start, stop, and restart the server interactively.

## Installation

Import the package into your Go project:

```go
import "github.com/vrianta/Server"
```

## Usage

### Step 1: Define Route Handlers

Each route handler should be a struct with `GET`, `POST`, and/or `DELETE` methods that accept a `*server.Session`:

```go
package controllers

import (
    server "github.com/vrianta/Server"
)

type home struct{}

var Home home

func (h *home) GET(Session *server.Session) {
    Session.RenderEngine.Render("Welcome to Home Page!")
}

func (h *home) POST(Session *server.Session) {
    Session.RenderEngine.Render("POST request received!")
}

func (h *home) DELETE(Session *server.Session) {
    Session.RenderEngine.Render("DELETE request received!")
}
```

### Step 2: Create a New Server Instance

Create a new server by calling [`server.New`](server.go):

```go
srv := server.New("", "8080", server.Routes{
    "/": &controllers.Home{},
    // Add more routes here
}, nil) // Pass nil for default config or provide a *server.Config
```

### Step 3: Start the Server

Start the server with:

```go
srv.Start()
```

This will also launch the interactive console for start/stop/restart commands.

## Session Management

Sessions are managed via the [`Session`](types.go) struct. You can access POST/GET data, store session variables, and check login status:

```go
func (h *home) GET(Session *server.Session) {
    if Session.IsLoggedIn() {
        Session.RenderEngine.Render("Welcome, " + Session.Store["uid"].(string))
    } else {
        Session.RenderEngine.Render("Please log in.")
    }
}
```

### Accessing POST and GET Data

```go
uid, uidExists := Session.POST["uid"]
token, tokenExists := Session.POST["token"]
```

### Setting and Checking Login

```go
Session.Login("user123")
if Session.IsLoggedIn() {
    // User is logged in
}
Session.Logout("/login")
```

## Rendering Responses

Use the [`RenderEngine`](types.go) to send responses:

```go
Session.RenderEngine.Render("Hello World")
Session.RenderEngine.StartRender()
```

To render HTML templates (from the `Views` folder):

```go
err := Session.RenderEngine.RenderTemplate("home.html", map[string]interface{}{
    "Title": "Home",
    "User":  Session.Store["uid"],
})
if err != nil {
    Session.RenderEngine.RenderError("Template error", server.ResponseCodes.InternalServerError)
}
Session.RenderEngine.StartRender()
```

## Static, CSS, and JS Files

Configure static, CSS, and JS folders in the [`Config`](types.go):

```go
cfg := &server.Config{
    Http: true,
    Static_folders: []string{"Static"},
    CSS_Folders:    []string{"Css"},
    JS_Folders:     []string{"Js"},
    Views_folder:   "Views",
}
srv := server.New("", "8080", routes, cfg)
```

Files in these folders are served automatically.

## SMTP Client

Send emails using the SMTP client in [`smtp/client.go`](smtp/client.go):

```go
import "github.com/vrianta/Server/smtp"

smtp.Client.InitSMTPClient("smtp.example.com", 587, "user", "pass")
err := smtp.Client.SendMail([]string{"to@example.com"}, "Subject", "Body")
if err != nil {
    // handle error
}
```

## Console Commands

When the server is running, you can use the following commands in the console:

- `start`    - Start the server
- `stop`     - Stop the server
- `restart`  - Restart the server
- `r`        - Shortcut for restart
- `exit`     - Stop the server and exit the program
- `-h`       - Display available commands

## Example Project Structure

```
.
├── config.go
├── console.go
├── cookies.go
├── file.handler.go
├── go.mod
├── LICENSE
├── log.go
├── readme.md
├── redirect.go
├── render.go
├── routing.go
├── server.go
├── session_manager.go
├── types.go
├── util.go
├── vars.go
└── smtp/
    └── client.go
```

## API Reference

- [`server.New`](server.go): Create a new server instance.
- [`server.Start`](server.go): Start the server.
- [`Session`](types.go): Session object for each request.
- [`Session.RenderEngine`](types.go): For rendering responses.
- [`Session.Login`](session_manager.go): Log in a user.
- [`Session.Logout`](session_manager.go): Log out a user.
- [`Session.IsLoggedIn`](session_manager.go): Check login status.
- [`Session.ParseRequest`](session_manager.go): Parse GET/POST data.
- [`Session.POST`](types.go): POST parameters.
- [`Session.GET`](types.go): GET parameters.
- [`Session.Store`](types.go): Session variables.
- [`server.RemoveSession`](server.go): Remove a session by ID.

For more details, see the source files:
- [server.go](server.go)
- [session_manager.go](session_manager.go)
- [types.go](types.go)
- [render.go](render.go)
- [smtp/client.go](smtp/client.go)

---

**License:**  
See [LICENSE](LICENSE) for GPLv3 license details.
