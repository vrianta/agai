# ModelsHandler: Go ORM-like Query Builder

ModelsHandler is a human-friendly, chainable query builder for working with database tables using Go structs. It allows you to easily create models, build queries, and interact with your database in a readable, maintainable way—similar to popular Object-Relational Mappers (ORMs).

## Features
- Chainable, fluent API for building queries
- Supports SELECT, CREATE, UPDATE, DELETE operations
- WHERE, AND, OR, IN, NOT IN, BETWEEN, LIKE, IS NULL, etc.
- LIMIT, OFFSET, ORDER BY, GROUP BY
- Easy to extend and understand
- **Automatic Database Migration**: When the `Build` flag is set to `false` in your configuration, ModelsHandler will automatically migrate your database schema to match your model definitions. This means tables and columns are created or updated as needed, so you don't have to write migration scripts manually.

> **Note:** Migration only happens if the `Build` flag is `false`. In development, set `Build: false` in your `web.config.json` to enable auto-migration. In production, set it to `true` to prevent accidental schema changes.

---

## 1. Defining a Model

To use ModelsHandler, define your model struct and register it. For example, to create a `User` model:

```go
package models

import models_handler "github.com/vrianta/Server/modelsHandler"

var Users = models_handler.New(
    "users", // Table name
    map[string]models_handler.Field{
        "userId": {
            Name:     "userId",
            Type:     models_handler.FieldsTypes.VarChar,
            Length:   20,
            Nullable: false,
            Index: models_handler.Index{
                PrimaryKey: true,
                Unique:     false,
                Index:      true,
            },
        },
        "userName": {
            Name:     "userName",
            Type:     models_handler.FieldsTypes.VarChar,
            Length:   30,
            Nullable: false,
            Index: models_handler.Index{
                Unique: true,
                Index:  true,
            },
        },
        "password": {
            Name:     "password",
            Type:     models_handler.FieldsTypes.Text,
            Nullable: false,
        },
        "firstName": {
            Name:     "firstName",
            Type:     models_handler.FieldsTypes.Text,
            Nullable: false,
        },
    },
)
```

---

## 2. Building and Executing Queries

### Creating Records (INSERT)

```go
// Create a new user
err := Users.Create().
    Set("userId").To("u123").
    Set("userName").To("alice").
    Set("password").To("securepass").
    Set("firstName").To("Alice").
    Exec()
```

### Fetching Data (SELECT)

```go
// Get all users older than 18, ordered by name
users, err := Users.Get().Where("age").GreaterThan(18).OrderBy("userName").Fetch()
```

### Fetching a Single Row

```go
// Get the first user with the name 'Alice'
user, err := Users.Get().Where("userName").Is("Alice").First()
```

### Updating Data (UPDATE)

```go
// Update the password of a user with userId = 'u123'
err := Users.Get().Set("password").To("newpass").Where("userId").Is("u123").Exec()
```

### Deleting Data (DELETE)

```go
// Delete users with a specific first name
err := Users.Get().Where("firstName").Is("John").Delete()
```

---

## 3. Query Builder API Reference

- `.Create()` — Start a new INSERT query
- `.Get()` — Start a new query (default SELECT)
- `.Where(column)` — Add a WHERE condition
- `.Is(value)` — WHERE column = value
- `.IsNot(value)` — WHERE column != value
- `.Like(pattern)` — WHERE column LIKE pattern
- `.And()`, `.Or()` — Combine conditions
- `.In(values...)`, `.NotIn(values...)` — WHERE column IN/NOT IN (...)
- `.GreaterThan(value)`, `.LessThan(value)` — WHERE column > / < value
- `.Between(min, max)` — WHERE column BETWEEN min AND max
- `.IsNull()`, `.IsNotNull()` — WHERE column IS (NOT) NULL
- `.Set(field).To(value)` — For UPDATE queries
- `.Limit(n)`, `.Offset(n)`, `.Page(page, pageSize)` — Pagination
- `.OrderBy(clause)`, `.GroupBy(clause)` — Sorting and grouping
- `.Fetch()` — Execute SELECT and get all results
- `.First()` — Execute SELECT and get the first result
- `.Exec()` — Execute UPDATE
- `.Delete()` — Execute DELETE

---

## 4. Example: Full Query Chain

```go
// Create a new user with multiple fields
err := Users.Create().
    Set("userId").To("u123").
    Set("userName").To("bob").
    Set("age").To(25).
    Set("status").To("active").
    Exec()

// Update all users named 'Bob' to be age 30
err := Users.Get().Set("age").To(30).Where("userName").Is("Bob").Exec()

// Get the first 10 active users, ordered by creation date
users, err := Users.Get().Where("status").Is("active").OrderBy("created_at DESC").Limit(10).Fetch()
```

---

## 5. Notes
- Always check for errors after executing queries.
- You can chain as many query builder methods as you need.
- The API is designed to be self-explanatory and easy to read.

---

## 6. Contributing
Pull requests and suggestions are welcome! Please document your code and keep the API intuitive.

---

## 7. License
MIT
