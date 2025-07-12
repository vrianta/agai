# Controller and View Guide

This guide explains how to create controllers and views in this framework, including available public variables and methods.

---

## 1. What is a Controller?

A **Controller** is a Go struct that handles HTTP requests for a specific route. It defines handler functions for HTTP methods (GET, POST, etc.), manages session data, and renders views (templates).

---

## 2. Controller Structure

### Public Fields

- **View**:  
  The name of the view (template) directory for this controller.  
  Example: `"home"` (looks for templates in `Views/home/`).

- **GET, POST, DELETE, PATCH, PUT, HEAD, OPTIONS**:  
  Handler functions for each HTTP method. Each receives the controller as `self` and returns a `Template.Response` (a map of data for the template).

### Public Methods

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

---

## 3. Example: Creating a Controller

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

---

## 4. Accessing Request Data

- `GetInput(key string)`: Returns a value from GET or POST parameters.
- `GetInputs()`: Returns all request parameters as a map.

---

## 5. Session Management

- `StoreData(key, value)`: Store data in the session (e.g., user ID).
- `GetStoredData(key)`: Retrieve data from the session.
- `Login()`, `Logout()`, `IsLoggedIn()`: Manage authentication state.

---

## 6. Redirects

- `Redirect(uri string)`: Redirects to another page.
- `WithCode(uri, code)`: Redirects with a custom HTTP status code.

---

## 7. Creating Views

### View Directory Structure

Each controller should have a corresponding directory under the `Views/` folder, named after the controller's `View` field.

**Example:**  
If `View: "home"`, templates should be in `Views/home/`.

### Supported Template Files

- `default.html`, `default.php`, or `index.html`/`index.php`: Default template for the controller.
- `get.html`, `post.html`, etc.: Templates for specific HTTP methods.

### Template Syntax

- PHP-style syntax is supported:
  - `<?= $var ?>` → `{{ .var }}`
  - `<?php if ($user): ?> ... <?php endif; ?>`
  - Loops: `<?php foreach ($items as $item): ?> ... <?php endforeach; ?>`

---

## 8. Registering Controllers

<<<<<<< HEAD
Register your controllers with the router using the `Router.Route` and `Router.RegisterRoutes` functions.

```go
import (
    "github.com/vrianta/agai/v1/router"
    "yourapp/Controller/Home"
)

func main() {
    router := Router.New("")
    router.RegisterRoutes(
        Router.Route("/home", Home.Home),
    )
}
```

---
=======
Register your controller in the router to make it available for handling requests.
>>>>>>> 1559b5cb4f55d903a1f6f41803258cd1614abb53

## 9. Summary of Public Variables and Methods

See above for a full list of fields and methods.

## References

- See the main project README for more details and advanced usage.