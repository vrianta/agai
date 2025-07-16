# Agai Handler

Program to manage `Agai` Projects and it's files and functionality to make it smooth for the users

- [Agai Handler](#agai-handler)
    - [⚙️ Runtime \& Utility Flags](#️-runtime--utility-flags)
    - [🔧 Creation Flags](#-creation-flags)
    - [🆘 Help \& Configuration Flags](#-help--configuration-flags)
    - [🧪 Example](#-example)
      - [🏗️ Create a Full Application with Controller, Model \& View](#️-create-a-full-application-with-controller-model--view)
      - [🧱 Create Only a Controller and View (for existing app)](#-create-only-a-controller-and-view-for-existing-app)
      - [🔍 Generate Multiple Controllers in One Go](#-generate-multiple-controllers-in-one-go)
      - [🧬 Generate a Component (e.g., Navbar or Footer)](#-generate-a-component-eg-navbar-or-footer)
      - [🛠️ Run Dev Handler for Auto Reloading](#️-run-dev-handler-for-auto-reloading)
      - [🧳 Migrate Models and Components](#-migrate-models-and-components)
      - [📄 Show Config Help](#-show-config-help)
      - [🆘 Show General Help](#-show-general-help)
      - [🔁 Combine Any Flags](#-combine-any-flags)

### ⚙️ Runtime & Utility Flags

| Flag | Description |
| --- | --- |
| --start-app, -sa | Launch the application server. |
| --start-handler, -sh | Start development handler (for auto-reload, background tasks). |
| --migrate-model, -mm | Apply database migrations for models. |
| --migrate-component, -mc | Migrate components to keep in sync with the database. |

* * *

### 🔧 Creation Flags

| Flag | Description |
| --- | --- |
| --create-app, -ca | Create a new application.Example: --create-app blog |
| --create-controller, -cc | Generate a new controller.Example: --create-controller post |
| --create-model, -cm | Generate a new model.Example: --create-model user |
| --create-component | Create a component (like a reusable layout part).Example: --create-component nav |
| --create-view, -cv | Create views for a controller.Optionally specify controller: --create-view post |

* * *

### 🆘 Help & Configuration Flags

| Flag | Description |
| --- | --- |
| --help, -h | Display general help and available flags. |
| --help-web-config, -hwc | Show help for web server configuration (port, host, static folders, etc.). |
| --help-database-config, -hdc | Show help for database configuration and environment overrides. |
| --help-session-config, -hsc | Show help for session store settings. |
| --help-smtp-config, -hsm | Show help for configuring SMTP settings for email. |

* * *

### 🧪 Example

```bash 
go run . --create-app blog \
        --create-controller post \
        --create-model post \
        --create-view \
        --start-app
```

* * *

#### 🏗️ Create a Full Application with Controller, Model & View

```bash
go run . --create-app blog \
         --create-controller post \
         --create-model post \
         --create-view \
         --start-app
```

**What it does:**

*   Creates a new app called `blog`
*   Adds a `post` controller and model
*   Generates default views for the `post` controller
*   Immediately starts the **server**

#### 🧱 Create Only a Controller and View (for existing app)

```bash
go run . --create-controller comment \
         --create-view comment
```

**What it does:**

*   Adds a new `comment` controller
*   Generates views for `comment` controller only

#### 🔍 Generate Multiple Controllers in One Go

```bash
go run . --create-controller post \
         --create-controller comment \
         --create-controller user
```

**What it does:**

*   Adds three new controllers: `post`, `comment`, and `user`
    

You can also generate views in the same command:

```bash
go run . --create-controller user --create-view user
```

* * *

#### 🧬 Generate a Component (e.g., Navbar or Footer)

```bash
go run . --create-component navbar
```

**What it does:**

*   Creates a reusable component named `navbar`
    
* * *

#### 🛠️ Run Dev Handler for Auto Reloading

```bash
go run . --start-handler
```

**What it does:**

*   Starts the CLI’s built-in dev handler to automatically reload on changes or run background tasks
    
* * *

#### 🧳 Migrate Models and Components

```bash
go run . --migrate-model --migrate-component
```

**What it does:**

*   Runs database migrations for all models and syncs component data

* * *

#### 📄 Show Config Help

```bash
go run . --help-web-config
```

```bash
go run . --help-database-config
```

```bash
go run . --help-smtp-config
```

**What it does:**

*   Prints helpful explanations and examples for each config file (`web.json`, `database.json`, etc.)
    
* * *

#### 🆘 Show General Help

```bash
go run . --help
```

**What it does:**

*   Lists all available CLI flags with short descriptions
    
* * *

#### 🔁 Combine Any Flags

```bash
go run . --create-app shop \
        --create-controller product \
        --create-model product \
        --create-component cart \
        --create-view \
        --start-handler
```

This sets up a full eCommerce scaffold with product model, controller, component, and auto-handler.

Let me know if you'd like a **web-config snippet**, **flag groupings by use-case**, or badges (Go version, build status) added to the README!