# AGAI Templating System Documentation

## Overview

The AGAI templating system is a flexible, chainable template rendering engine that supports multiple template formats (HTML, PHP, Go templates) with built-in theme support, hot-reload during development, and a comprehensive set of built-in template functions.

### Key Features

- **Multiple Template Formats**: Support for HTML, PHP, and Go templates
- **Theme System**: Organize templates by themes with automatic registration
- **Hot Reload**: Live reload during development without server restart
- **Built-in Functions**: Rich set of template helper functions
- **Component System**: Include and reuse template components
- **Type Safety**: Leverages Go's `html/template` package for XSS protection

---

## Architecture Overview

### Core Components

```
template.go                    # Core template context and execution
register_template.go           # Template and theme registration
default_template_functions.go  # Built-in template functions
├── Template Types
│   ├── GoTemplate (0)
│   ├── HtmlTemplate (1)
│   └── PhpTemplate (2)
└── Template Storage
    └── templateComponents map[string]*Context
```

### Data Structures

#### Context
Represents a single template file with metadata and rendering capabilities.

```go
type Context struct {
    uri          string                 // Full path to template file
    name         string                 // Template name (filename without extension)
    lastModified time.Time              // Track file changes for hot reload
    Html         *htmltemplate.Template // Compiled HTML template
    Php          *htmltemplate.Template // Compiled PHP-converted template
    ViewType     ViewType               // Template type (HTML/PHP/Go)
}
```

#### Contexts Structure
Holds method-specific templates (GET, POST, PUT, DELETE, etc.)

```go
type Contexts struct {
    index   *Context // Default template
    get     *Context // GET request template
    post    *Context // POST request template
    delete  *Context // DELETE request template
    patch   *Context // PATCH request template
    put     *Context // PUT request template
    head    *Context // HEAD request template
    options *Context // OPTIONS request template
}
```

---

## Directory Structure

```
views/                          # Main view folder (configurable)
├── index.html                  # Root-level templates
├── 404.html                    # 404 error page
├── dashboard.html
└── admin/                      # Theme folder
    ├── index.html
    ├── dashboard.html
    ├── users/
    │   ├── list.html
    │   └── edit.html
    └── settings/
        └── index.html

# Resulting template keys:
# - "index"
# - "404"
# - "dashboard"
# - "admin.index"
# - "admin.dashboard"
# - "admin.users.list"
# - "admin.users.edit"
# - "admin.settings.index"
```

---

## Template Registration System

### Automatic Registration Flow

1. **Initialization (`init()` function)**
   - Reads the view folder on application startup
   - Processes all files and directories recursively
   - Registers templates in the `templateComponents` map

2. **File Registration**
   - Files in the root view folder become top-level templates
   - Key format: `filename` (e.g., `index.html` → `"index"`)

3. **Theme Registration**
   - Subdirectories are treated as themes
   - Templates within themes are namespaced
   - Key format: `themename.filename` (e.g., `admin/dashboard.html` → `"admin.dashboard"`)

4. **Recursive Directory Processing**
   - Nested directories create nested namespaces
   - Key format: `theme.subfolder.filename` (e.g., `admin/users/list.html` → `"admin.users.list"`)

### Special Handling

#### 404 Pages

```go
// Root-level 404.html
http.HandleFunc("/404/", func(w http.ResponseWriter, r *http.Request) {
    // Renders /views/404.html
})

// Theme-specific 404 pages
http.HandleFunc("/" + theme + "/404/", func(w http.ResponseWriter, r *http.Request) {
    // Renders /views/{theme}/404.html
    // Falls back to default 404 if not found
})
```

---

## Built-in Template Functions

All functions are available in both HTML and PHP templates via the `ReponseFuncMaps` function map.

### String Functions

#### `upper`
Convert string to uppercase.
```html
{{ "hello" | upper }}
<!-- Output: HELLO -->
```

#### `lower`
Convert string to lowercase.
```html
{{ "HELLO" | lower }}
<!-- Output: hello -->
```

#### `strlen`
Get string length.
```html
{{ "hello" | strlen }}
<!-- Output: 5 -->
```

### Collection Functions

#### `len` / `count`
Get length of any collection (array, slice, map, string, channel).
```html
{{ .items | len }}
{{ .users | count }}
```

### Output Functions

#### `print`
Format any data as a string.
```html
{{ .data | print }}
```

#### `date`
Get current date/time formatted.
```html
{{ "2006-01-02 15:04:05" | date }}
<!-- Output: 2026-02-22 14:30:00 -->
```

### Component Functions

#### `include`
Include another template component with optional data.
```html
{{ include "header" . }}
{{ include "admin.menu" .UserData }}
```

---

## Template Types & Formats

### 1. HTML Templates (`.html`, `.gohtml`)

Pure HTML files with Go template syntax. Automatically sanitized against XSS.

```html
<!-- views/dashboard.html -->
<!DOCTYPE html>
<html>
<head>
    <title>{{ .PageTitle }}</title>
</head>
<body>
    <h1>{{ "welcome" | upper }}</h1>
    <p>Total users: {{ .Users | count }}</p>
    
    <!-- Include header component -->
    {{ include "components.header" . }}
    
    {{range .Users}}
        <div>{{ .Name }} - {{ .Email }}</div>
    {{end}}
</body>
</html>
```

### 2. PHP Templates (`.php`, `.gophp`)

PHP-like templates automatically converted to Go template syntax.

```php
<!-- views/admin/dashboard.php -->
<!DOCTYPE html>
<html>
<head>
    <title><?php echo $PageTitle; ?></title>
</head>
<body>
    <h1><?php echo strtoupper("welcome"); ?></h1>
    <p>Total items: <?php echo count($Items); ?></p>
    
    <?php foreach($Items as $item): ?>
        <div><?php echo $item['name']; ?></div>
    <?php endforeach; ?>
</body>
</html>
```

#### PHP to Go Template Conversion

The system automatically converts PHP syntax:
- `<?php ... ?>` → Go template syntax
- PHP variables and functions → Go template equivalents
- Loop structures converted appropriately
- **NEW**: Complex array access with property expressions → Proper Go template index syntax

##### Complex Array Access Feature

The templating system now supports complex array accesses where array keys contain property expressions.

**Example**:
```php
<!-- PHP Template -->
<?php echo $event_themes[$enquiry->Aesthetic]->ThemeName; ?>

<!-- Converts to: -->
{{ (index .event_themes .enquiry.Aesthetic).ThemeName }}
```

**Supported Patterns**:
```php
<!-- Simple variable index (existing) -->
<?= $items[$index] ?>
<!-- Converts to: -->
{{ index .items $index }}

<!-- String key (existing) -->
<?= $themes['default']->Name ?>
<!-- Converts to: -->
{{ .themes.default.Name }}

<!-- Complex expression as index (NEW) -->
<?= $data[$config->Theme]->Settings ?>
<!-- Converts to: -->
{{ (index .data .config.Theme).Settings }}

<!-- Object property after array access (NEW) -->
<?= $event_themes[$enquiry->Aesthetic]->ThemeName ?>
<!-- Converts to: -->
{{ (index .event_themes .enquiry.Aesthetic).ThemeName }}
```

### 3. Go Templates

Standard Go html/template syntax with security baked-in.

---

## End-to-End Scenarios & Examples

### Scenario 1: Basic Page Rendering

**Requirement**: Render a simple user dashboard page.

**File Structure**:
```
views/
└── dashboard.html
```

**Template File** (`views/dashboard.html`):
```html
<!DOCTYPE html>
<html>
<head>
    <title>Dashboard</title>
</head>
<body>
    <h1>Welcome, {{ .UserName }}</h1>
    <p>Registration Date: {{ .RegisteredAt }}</p>
</body>
</html>
```

**Handler Code**:
```go
// In your router/request handler
func Dashboard(w http.ResponseWriter, r *http.Request) {
    // Get the template
    template, ok := template.GetTemplate("dashboard")
    if !ok {
        http.NotFound(w, r)
        return
    }
    
    // Prepare data
    data := map[string]interface{}{
        "UserName":     "John Doe",
        "RegisteredAt": "2026-02-22",
    }
    
    // Render and write response
    buf, err := template.Execute(data)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    w.Write(buf)
}
```

---

### Scenario 2: Theme-based Multi-tenant Application

**Requirement**: Support multiple themes for different organizations.

**File Structure**:
```
views/
├── index.html (default/shared)
├── admin/
│   ├── index.html
│   ├── dashboard.html
│   └── users/
│       ├── list.html
│       └── edit.html
├── client/
│   ├── index.html
│   └── dashboard.html
└── partner/
    ├── index.html
    └── reports.html
```

**Usage**:
```go
// Admin theme dashboard
func AdminDashboard(w http.ResponseWriter, r *http.Request) {
    template, ok := template.GetTemplate("admin.dashboard")
    if !ok {
        // Fallback to root dashboard
        template, _ = template.GetTemplate("dashboard")
    }
    
    data := map[string]interface{}{
        "Theme": "admin",
        "Stats": getAdminStats(),
    }
    
    buf, _ := template.Execute(data)
    w.Write(buf)
}

// Client theme dashboard
func ClientDashboard(w http.ResponseWriter, r *http.Request) {
    template, ok := template.GetTemplate("client.dashboard")
    if !ok {
        template, _ = template.GetTemplate("dashboard")
    }
    
    data := map[string]interface{}{
        "Theme": "client",
        "Stats": getClientStats(),
    }
    
    buf, _ := template.Execute(data)
    w.Write(buf)
}
```

---

### Scenario 3: Component Reuse with Include

**Requirement**: Create reusable header, navigation, and footer components.

**File Structure**:
```
views/
├── components/
│   ├── header.html
│   ├── navbar.html
│   └── footer.html
└── pages/
    ├── home.html
    ├── about.html
    └── contact.html
```

**Component Templates**:

```html
<!-- views/components/header.html -->
<header>
    <h1>{{ .SiteTitle }}</h1>
    <subtitle>{{ .SiteSubtitle }}</subtitle>
</header>
```

```html
<!-- views/components/navbar.html -->
<nav class="navbar">
    {{range .NavItems}}
    <a href="{{ .Href }}" class="nav-link">{{ .Label }}</a>
    {{end}}
</nav>
```

```html
<!-- views/components/footer.html -->
<footer>
    <p>&copy; {{ "2006" | date }} {{ .CompanyName }}</p>
</footer>
```

**Using Components**:

```html
<!-- views/pages/home.html -->
<!DOCTYPE html>
<html>
<head>
    <title>Home</title>
</head>
<body>
    {{ include "components.header" . }}
    {{ include "components.navbar" .Navigation }}
    
    <main>
        <h2>{{ .PageTitle }}</h2>
        <p>{{ .Content }}</p>
    </main>
    
    {{ include "components.footer" . }}
</body>
</html>
```

**Handler**:
```go
func HomePage(w http.ResponseWriter, r *http.Request) {
    template, _ := template.GetTemplate("pages.home")
    
    data := map[string]interface{}{
        "SiteTitle":    "My Website",
        "SiteSubtitle": "Welcome!",
        "PageTitle":    "Home",
        "Content":      "This is the home page content.",
        "Navigation": map[string]interface{}{
            "NavItems": []map[string]string{
                {"Label": "Home", "Href": "/"},
                {"Label": "About", "Href": "/about"},
                {"Label": "Contact", "Href": "/contact"},
            },
        },
        "CompanyName": "Acme Corp",
    }
    
    buf, _ := template.Execute(data)
    w.Write(buf)
}
```

---

### Scenario 4: Dynamic Content with Loops and Conditionals

**Requirement**: Display a paginated user list with filtering.

**File Structure**:
```
views/
└── admin/
    └── users/
        └── list.html
```

**Template** (`views/admin/users/list.html`):
```html
<!DOCTYPE html>
<html>
<head>
    <title>User Management</title>
</head>
<body>
    <h1>User List</h1>
    
    {{if .Users}}
        <table>
            <thead>
                <tr>
                    <th>ID</th>
                    <th>Name</th>
                    <th>Email</th>
                    <th>Role</th>
                    <th>Status</th>
                    <th>Actions</th>
                </tr>
            </thead>
            <tbody>
                {{range .Users}}
                <tr>
                    <td>{{ .ID }}</td>
                    <td>{{ .Name }}</td>
                    <td>{{ .Email }}</td>
                    <td>{{ .Role | upper }}</td>
                    <td>
                        {{if eq .Status "active"}}
                            <span class="badge-success">Active</span>
                        {{else if eq .Status "inactive"}}
                            <span class="badge-warning">Inactive</span>
                        {{else}}
                            <span class="badge-danger">Suspended</span>
                        {{end}}
                    </td>
                    <td>
                        <a href="/admin/users/{{ .ID }}/edit">Edit</a>
                        <a href="/admin/users/{{ .ID }}/delete">Delete</a>
                    </td>
                </tr>
                {{end}}
            </tbody>
        </table>
        
        <!-- Pagination -->
        {{if gt .CurrentPage 1}}
            <a href="?page={{ sub .CurrentPage 1 }}">Previous</a>
        {{end}}
        
        <span>Page {{ .CurrentPage }} of {{ .TotalPages }}</span>
        
        {{if lt .CurrentPage .TotalPages}}
            <a href="?page={{ add .CurrentPage 1 }}">Next</a>
        {{end}}
    {{else}}
        <p>No users found.</p>
    {{end}}
</body>
</html>
```

**Handler**:
```go
func ListUsers(w http.ResponseWriter, r *http.Request) {
    template, _ := template.GetTemplate("admin.users.list")
    
    page := 1
    if p := r.URL.Query().Get("page"); p != "" {
        page, _ = strconv.Atoi(p)
    }
    
    // Fetch from database
    users := getUsersForPage(page)
    totalPages := getTotalPages()
    
    data := map[string]interface{}{
        "Users":       users,
        "CurrentPage": page,
        "TotalPages":  totalPages,
    }
    
    buf, _ := template.Execute(data)
    w.Write(buf)
}
```

---

### Scenario 5: Hot Reload During Development

**Requirement**: Automatically reload templates when files change.

**Configuration** (`config/web.go`):
```go
// Must be set to false during development
config.GetWebConfig().Build = false
```

**How It Works**:

1. Templates are NOT cached in development mode
2. `Update()` method checks file modification time
3. If file is newer, template is recompiled
4. Client-side EventSource script monitors for reload events
5. JavaScript auto-reloads page when changes detected

**Automatic Hot Reload Script** (injected into all templates during development):

```javascript
const source = new EventSource("http://localhost:8888/hot-reload");

source.onmessage = function(event) {
    if (event.data === "reload") {
        window.location.reload();
    }
};

source.onerror = function(err) {
    console.warn("[LiveReload] Disconnected from server", err);
};
```

**Handler Setup**:
```go
func UpdateTemplate(w http.ResponseWriter, r *http.Request) {
    if template, ok := template.GetTemplate("my-template"); ok {
        // Check for modifications
        if err := template.Update(); err != nil {
            log.Error("Failed to update template: %v", err)
        }
        
        buf, _ := template.Execute(data)
        w.Write(buf)
    }
}
```

---

### Scenario 6: Error Pages with Custom Themes

**Requirement**: Display branded 404 pages for different themes.

**File Structure**:
```
views/
├── 404.html (default/fallback)
└── admin/
    └── 404.html (admin-specific)
```

**Default 404** (`views/404.html`):
```html
<!DOCTYPE html>
<html>
<head>
    <title>404 - Page Not Found</title>
</head>
<body>
    <h1>404 - Page Not Found</h1>
    <p>The page you're looking for doesn't exist.</p>
    <a href="/">Back to Home</a>
</body>
</html>
```

**Admin 404** (`views/admin/404.html`):
```html
<!DOCTYPE html>
<html>
<head>
    <title>404 - Page Not Found | Admin Panel</title>
</head>
<body>
    <h1>404 - Admin Page Not Found</h1>
    <p>This admin page doesn't exist.</p>
    <a href="/admin/">Back to Admin Home</a>
</body>
</html>
```

**Automatic HTTP Handlers**:

```go
// Registered automatically during init()
http.HandleFunc("/404/", func(w http.ResponseWriter, r *http.Request) {
    t, _ := templateComponents["404"]
    buf, _ := t.Execute("")
    w.Write(buf)
})

http.HandleFunc("/admin/404/", func(w http.ResponseWriter, r *http.Request) {
    t, _ := templateComponents["admin.404"]
    if t == nil {
        t, _ = templateComponents["404"] // fallback
    }
    buf, _ := t.Execute("")
    w.Write(buf)
})
```

---

### Scenario 7: Complex Array Access with Property Expressions

**Requirement**: Handle theme selection where theme ID is based on a property value from another object.

**File Structure**:
```
views/
└── product/
    └── details.html
```

**Handler Code**:
```go
func ProductDetails(w http.ResponseWriter, r *http.Request) {
    product := getProductFromDB()
    enquiry := getEnquiryPreferences()
    
    // Themes keyed by aesthetic level
    themes := map[string]ThemeConfig{
        "minimalist": ThemeConfig{
            ThemeName: "Clean & Modern",
            Colors: []string{"#fff", "#333"},
        },
        "elaborate": ThemeConfig{
            ThemeName: "Rich & Detailed",
            Colors:    []string{"#ffd700", "#8b0000"},
        },
    }
    
    data := map[string]interface{}{
        "Product":      product,
        "Enquiry":      enquiry,
        "EventThemes":  themes,
    }
    
    template, _ := template.GetTemplate("product.details")
    buf, _ := template.Execute(data)
    w.Write(buf)
}
```

**PHP Template** (`views/product/details.html`):
```php
<!DOCTYPE html>
<html>
<head>
    <title><?php echo $Product->Name; ?></title>
</head>
<body>
    <!-- Example of complex array access with property expression -->
    <div class="theme-header" style="background-color: <?php echo $EventThemes[$Enquiry->Aesthetic]->Colors[0]; ?>;">
        <h1><?php echo $Product->Name; ?></h1>
        <p>Theme: <strong><?php echo $EventThemes[$Enquiry->Aesthetic]->ThemeName; ?></strong></p>
    </div>
    
    <div class="product-content">
        <p><?php echo $Product->Description; ?></p>
        <price><?php echo $Product->Price; ?></price>
    </div>
</body>
</html>
```

**Converted to Go Template**:
```html
<!DOCTYPE html>
<html>
<head>
    <title>{{ .Product.Name }}</title>
</head>
<body>
    <!-- Array with complex property-based index -->
    <div class="theme-header" style="background-color: {{ index (index .EventThemes .Enquiry.Aesthetic).Colors 0 }};">
        <h1>{{ .Product.Name }}</h1>
        <p>Theme: <strong>{{ (index .EventThemes .Enquiry.Aesthetic).ThemeName }}</strong></p>
    </div>
    
    <div class="product-content">
        <p>{{ .Product.Description }}</p>
        <price>{{ .Product.Price }}</price>
    </div>
</body>
</html>
```

**How It Works**:
1. PHP syntax: `$$EventThemes[$Enquiry->Aesthetic]` = access EventThemes map using a property value
2. Converted to: `(index .EventThemes .Enquiry.Aesthetic)` = use Go's index function with property access
3. Property access: `->ThemeName` converts to `.ThemeName` and works on the result
4. Final output: `(index .EventThemes .Enquiry.Aesthetic).ThemeName` = safe Go template syntax

---

## Template Execution Flow

### Step-by-Step Execution Process

```
1. Request Received
   ↓
2. Route Handler Triggered
   ↓
3. GetTemplate(name) called
   ↓
4. Lookup in templateComponents map
   ↓
5. If development mode:
   - Call Update() to check file modifications
   - Recompile if file changed
   ↓
6. Prepare data (map[string]interface{})
   ↓
7. Call Execute(data)
   ↓
8. Go's html/template renders template with data
   - All output is HTML-escaped (XSS protection)
   - Built-in functions applied
   - Components included recursively
   ↓
9. Buffer returned with rendered HTML
   ↓
10. Write to HTTP response
```

---

## Best Practices

### Do's ✅

- Use the `include` function to create reusable components
- Organize complex templates in theme directories
- Always prepare data as `map[string]interface{}` for flexibility
- Use template functions for string manipulation (upper, lower)
- Leverage Go's built-in template security (auto-escaping)
- Keep template logic simple; move complex logic to handlers

### Don'ts ❌

- Don't store sensitive data in templates
- Don't use `html/template.HTML` unless absolutely necessary (bypasses escaping)
- Don't put complex business logic in templates
- Don't forget to handle errors from `Execute()`
- Don't hardcode URLs; pass them through data
- Don't mix multiple template formats randomly

---

## Configuration

### View Folder Configuration

```go
// In config/application.config.go
config.GetViewFolder() // Returns configured views directory
                       // Default: "views"
```

### Build Mode

```go
// Development mode (hot reload enabled)
config.GetWebConfig().Build = false

// Production mode (templates cached)
config.GetWebConfig().Build = true
```

---

## Troubleshooting

### Template Not Found

```
Error: GetTemplate returns false for template name
```

**Solution**: Verify the template file exists and check the template key format:
- Files in root: `filename` (no extension)
- Files in subdirs: `dirname.filename`
- Nested dirs: `dir1.dir2.filename`

### Hot Reload Not Working

```
Changes to template files not reflected
```

**Solution**: 
1. Verify `config.GetWebConfig().Build = false` in development
2. Check that EventSource endpoint (`/hot-reload`) is responding
3. Verify file system has file change detection enabled

### XSS Vulnerabilities

```
HTML output showing escaped characters
```

**Solution**: Only use `html/template` for auto-escaping. Never use `text/template`.

### Component Include Failing

```
Error: No Template Found: componentname
```

**Solution**: Ensure component template name matches the key:
- Include uses exact template name lookup
- Use correct namespace: `include "theme.component"`

---

## Advanced Features

### Buffer Pool Optimization

The system uses a `sync.Pool` for efficient buffer reuse:

```go
template_bufPool = sync.Pool{
    New: func() any { return new(bytes.Buffer) },
}
```

This reduces garbage collection pressure on high-traffic servers.

### Thread Safety

Template registration uses `sync.RWMutex` for safe concurrent access:

```go
templateRecordsMutex = &sync.RWMutex{}
```

Safe for concurrent reads from multiple goroutines.

---

## Summary

The AGAI templating system provides a modern, secure, and flexible way to render web content. With automatic theme registration, hot-reload capabilities, and a rich set of built-in functions, it enables developers to build dynamic, maintainable web applications with minimal boilerplate.

Key takeaways:
- **Easy to use**: Simple, intuitive API
- **Secure**: XSS protection built-in
- **Scalable**: Theme-based organization
- **Developer-friendly**: Hot reload and clear error messages
- **Performant**: Buffer pooling and efficient caching

For more examples and use cases, refer to the scenario sections above.
