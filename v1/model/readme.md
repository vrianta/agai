# ModelsHandler: Go ORM-like queryBuilder Builder

ModelsHandler is a human-friendly, chainable queryBuilder builder for working with database tables using Go structs. It allows you to easily create models, build queries, and interact with your database in a readable, maintainable way—similar to popular Object-Relational Mappers (ORMs).

---

# Table of Contents
- [ModelsHandler: Go ORM-like queryBuilder Builder](#modelshandler-go-orm-like-querybuilder-builder)
- [Table of Contents](#table-of-contents)
  - [Features](#features)
  - [1. Defining a Model](#1-defining-a-model)
  - [2. Building and Executing Queries](#2-building-and-executing-queries)
    - [Creating Records (InsertRow)](#creating-records-insertrow)
    - [Fetching Data (SELECT)](#fetching-data-select)
    - [Fetching a Single Row](#fetching-a-single-row)
    - [Updating Data (UPDATE)](#updating-data-update)
    - [Deleting Data (DELETE)](#deleting-data-delete)
  - [3. queryBuilder Builder API Reference](#3-querybuilder-builder-api-reference)
  - [4. Example: Full queryBuilder Chain](#4-example-full-querybuilder-chain)
  - [5. Notes](#5-notes)
  - [6. Contributing](#6-contributing)
  - [7. License](#7-license)

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
package model

import models_handler "github.com/vrianta/agai/v1/modelsHandler"

var Users = models_handler.New(
    "users", // Table name
    map[string]models_handler.Field{
        "userId": {
            Name:     "userId",
            Type:     models_handler.FieldTypesTypes.VarChar,
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
            Type:     models_handler.FieldTypesTypes.VarChar,
            Length:   30,
            Nullable: false,
            Index: models_handler.Index{
                Unique: true,
                Index:  true,
            },
        },
        "password": {
            Name:     "password",
            Type:     models_handler.FieldTypesTypes.Text,
            Nullable: false,
        },
        "firstName": {
            Name:     "firstName",
            Type:     models_handler.FieldTypesTypes.Text,
            Nullable: false,
        },
    },
)
```

---

## 2. Building and Executing Queries

### Creating Records (InsertRow)

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

## 3. queryBuilder Builder API Reference

- `.Create()` — Start a new InsertRow queryBuilder
- `.Get()` — Start a new queryBuilder (default SELECT)
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

## 4. Example: Full queryBuilder Chain

```go
// Create a new user with multiple FieldTypes
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
- You can chain as many queryBuilder builder methods as you need.
- The API is designed to be self-explanatory and easy to read.

---

## 6. Contributing
Pull requests and suggestions are welcome! Please document your code and keep the API intuitive.

---

## 7. License
MIT
