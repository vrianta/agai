# Changelog – v0.2.1

This release introduces major architectural improvements, including model migration, disk-based session storage, component syncing, new CLI commands, and several core fixes. It lays the foundation for future growth with a modular internal structure and flexible configuration.

---

## Framework Core

### Added
- `handler.go`: new CLI entry point for commands like:
  - `--new-application`
  - `--new-controller`
  - `--new-model`
  - `--migrate-model`
  - `--migrate-component`
  - `--start-server`
- Command-line flag parser in `v1/config/args.go`
- CLI-based scaffolding for apps, models, controllers
- Model migration engine (`sync.go`, `schema.go`) with interactive prompts
- Component syncing system (DB ↔ JSON) with reload and dump
- Pluggable session storage (memory or disk)
- Environment variable overrides for all config values

---

## Code Structure

### Added
- Modular structure under `/v1/` for versioned evolution
- Separated config files: `config.go`, `args.go`, `database_config.go`, `var.go`, `type.go`
- Internal session logic moved to `v1/internal/session/`
- New `log/`, `controller/`, and `component/` folders
- Multiple `readme.md` files in major packages

### Removed
- Legacy global packages: `Log`, `Cookies`, `Template/registerTemplate.go`
- Deprecated `var.go` files and static config globals

---

## Models & ORM

- `ModelsHandler` now supports:
  - `SELECT`, `INSERT`, `UPDATE`, `DELETE`
  - Pagination (`Page`, `Offset`)
  - Sorting (`OrderBy`, `GroupBy`)
  - `.Clone()` for query reuse
- Schema migration tool detects field mismatches and applies diffs
- `Build: true` flag controls whether migrations prompt for confirmation

---

## Component System

- Loads components from `.components.json` or falls back to DB
- Supports key-value mapping and custom structs via generics
- File operations:
  - `ReloadComponents()`
  - `DumpComponentToJSON()`
- Example: `components/sample.components.json`

---

## Sessions

- LRU-based session heap with buffered updates
- Thread-safe access using `sync.RWMutex`
- Switch between in-memory or file-based session storage
- `Session.LoggedIn` exposed publicly

---

## Views & Templates

- Supports default and method-specific templates (e.g. `get.html`, `post.html`)
- Safe fallback when template missing
- PHP-style template syntax (`<?= $var ?>`) supported and documented
- Shared view support via `Views/shared/`

---

## Configuration

- New config structure:
  - `config.web.json`: web server and session settings
  - `config.database.json`: DB connection
- All values can be overridden via environment variables
- `Build` flag toggles dev/prod behaviors

---

## Documentation

- Expanded `readme.md` covering:
  - Project structure
  - Controllers, views, models, sessions
  - Config files and environment overrides
  - Components and migrations
- Package-level READMEs added for quick reference

---

## Bug Fixes

- Resolved high-traffic crash caused by unbuffered session LRU channel
- Missing templates now log fallback instead of panicking
- Safer session cleanup and expiration flow

---

## Coming Soon (Next Major Version)

- CLI: view generators, project scaffolds
- Middleware pipeline
- DB query layer upgrades
- Template includes/partials
- PostgreSQL driver support
- In-browser GUI editor for components (planned)
- More session storage supports - Raddis/similer tool, DB, JWT Token session management(where the session ID will be the token of the JWT)
- Support for different Return Types for Controller and difernt Json Return Support for functions 


