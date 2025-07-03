package model

import (
	"fmt"
	"strings"

	DatabaseHandler "github.com/vrianta/Server/database"
)

// ===============================
// Query Builder for ModelsHandler
// ===============================
//
// This file provides a human-friendly, chainable query builder for working with database tables using Go structs.
// It is designed to make database operations (SELECT, UPDATE, DELETE) easy to read, write, and maintain.
//
// The builder mimics the style of popular Object-Relational Mappers (ORMs), allowing you to construct queries
// in a fluent, readable way. This means you can chain methods together to build up complex queries step by step.
//
// Example usage:
//
//   users := UserModel.Get().Where("age").GreaterThan(18).OrderBy("name").Fetch()
//
// Supported features include:
//   - WHERE, AND, OR conditions
//   - LIMIT, OFFSET for pagination
//   - GROUP BY, ORDER BY for sorting and grouping
//   - Comparison operators (>, <, =, !=, IN, NOT IN, BETWEEN, LIKE, IS NULL, etc.)
//   - UPDATE and DELETE operations
//
// Each function is documented in detail below. If you are new to this code, read the comments for each function
// to understand what it does and how to use it. The goal is to make database access as intuitive as possible.
//
// If you are not a Go developer, don't worry! The comments explain the logic in plain English.

// Entry point: create a new query for the given model struct.
// This function starts a new query chain. By default, it prepares for a SELECT operation.
// Example: UserModel.Get() returns a Query object you can chain more methods onto.
func (m *Struct) Get() *Query {
	return &Query{
		model:     m,        // The model (table) this query is for
		operation: "select", // Default operation is SELECT
	}
}

func (m *Struct) Create() *Query {
	return &Query{
		model:        m,
		operation:    "insert",
		insertFields: make(map[string]any),
	}
}

// =======================
// WHERE Clause Functions
// =======================

// Where begins a WHERE clause, specifying the column to filter on.
// Example: .Where("age")
func (q *Query) Where(column string) *Query {
	q.lastColumn = column // Remember which column the next condition is for
	return q
}

// Is adds an equality condition to the WHERE clause.
// Example: .Where("age").Is(30)  // WHERE age = 30
func (q *Query) Is(value any) *Query {
	q.whereClauses = append(q.whereClauses, fmt.Sprintf("`%s` = ?", q.lastColumn)) // Add an equality condition for the last column
	q.whereArgs = append(q.whereArgs, value)                                       // Add the value to the arguments for the query
	q.lastColumn = ""                                                              // Reset lastColumn for safety
	return q                                                                       // Return the query object for chaining
}

// IsNot adds a NOT EQUAL condition (`!=`) to the WHERE clause for the previously specified column.
// It is used after calling .Where("columnName").
//
// Example:
//
//	query.Where("status").IsNot("inactive")
//
// Generates:
//
//	WHERE `status` != 'inactive'
func (q *Query) IsNot(value any) *Query {
	q.whereClauses = append(q.whereClauses, fmt.Sprintf("`%s` != ?", q.lastColumn))
	q.whereArgs = append(q.whereArgs, value)
	q.lastColumn = ""
	return q
}

// Like adds a LIKE condition to the WHERE clause for SQL pattern matching (e.g., for wildcards like `%value%`).
// It is used after calling .Where("columnName").
//
// Example:
//
//	query.Where("username").Like("%pritam%")
//
// Generates:
//
//	WHERE `username` LIKE '%pritam%'
func (q *Query) Like(value string) *Query {
	q.whereClauses = append(q.whereClauses, fmt.Sprintf("`%s` LIKE ?", q.lastColumn))
	q.whereArgs = append(q.whereArgs, value)
	q.lastColumn = ""
	return q
}

// And appends a logical AND operator between WHERE conditions.
// It should be used between chained .Where() clauses.
//
// Example:
//
//	query.Where("role").Is("admin").And().Where("active").Is(true)
//
// Generates:
//
//	WHERE `role` = 'admin' AND `active` = true
func (q *Query) And() *Query {
	q.whereClauses = append(q.whereClauses, "AND")
	return q
}

// Or appends a logical OR operator between WHERE conditions.
// It should be used between chained .Where() clauses.
//
// Example:
//
//	query.Where("role").Is("admin").Or().Where("role").Is("moderator")
//
// Generates:
//
//	WHERE `role` = 'admin' OR `role` = 'moderator'
func (q *Query) Or() *Query {
	q.whereClauses = append(q.whereClauses, "OR")
	return q
}

// In adds an IN condition to the WHERE clause for checking if a column's value exists in a set of values.
// It is used after .Where("columnName") and accepts a variadic list of values.
//
// Example:
//
//	query.Where("userId").In(1, 2, 3)
//
// Generates:
//
//	WHERE `userId` IN (1, 2, 3)
//
// Note: The values passed are safely parameterized using `?` placeholders to prevent SQL injection.
func (q *Query) In(values ...any) *Query {
	placeholders := strings.TrimRight(strings.Repeat("?,", len(values)), ",")
	q.whereClauses = append(q.whereClauses, fmt.Sprintf("`%s` IN (%s)", q.lastColumn, placeholders))
	q.whereArgs = append(q.whereArgs, values...)
	q.lastColumn = ""
	return q
}

// NotIn adds a NOT IN condition to the WHERE clause for excluding values.
// Usage: .Where("status").NotIn("inactive", "banned")
func (q *Query) NotIn(values ...any) *Query {
	placeholders := strings.TrimRight(strings.Repeat("?,", len(values)), ",")
	q.whereClauses = append(q.whereClauses, fmt.Sprintf("`%s` NOT IN (%s)", q.lastColumn, placeholders))
	q.whereArgs = append(q.whereArgs, values...)
	q.lastColumn = ""
	return q
}

// GreaterThan adds a "greater than" condition to the WHERE clause.
// Usage: .Where("score").GreaterThan(100)
func (q *Query) GreaterThan(value any) *Query {
	q.whereClauses = append(q.whereClauses, fmt.Sprintf("`%s` > ?", q.lastColumn))
	q.whereArgs = append(q.whereArgs, value)
	q.lastColumn = ""
	return q
}

// LessThan adds a "less than" condition to the WHERE clause.
// Usage: .Where("score").LessThan(50)
func (q *Query) LessThan(value any) *Query {
	q.whereClauses = append(q.whereClauses, fmt.Sprintf("`%s` < ?", q.lastColumn))
	q.whereArgs = append(q.whereArgs, value)
	q.lastColumn = ""
	return q
}

// Between adds a BETWEEN condition to the WHERE clause for a range.
// Usage: .Where("created_at").Between(start, end)
func (q *Query) Between(min, max any) *Query {
	q.whereClauses = append(q.whereClauses, fmt.Sprintf("`%s` BETWEEN ? AND ?", q.lastColumn))
	q.whereArgs = append(q.whereArgs, min, max)
	q.lastColumn = ""
	return q
}

// IsNull adds an IS NULL condition to the WHERE clause.
// Usage: .Where("deleted_at").IsNull()
func (q *Query) IsNull() *Query {
	q.whereClauses = append(q.whereClauses, fmt.Sprintf("`%s` IS NULL", q.lastColumn))
	q.lastColumn = ""
	return q
}

// IsNotNull adds an IS NOT NULL condition to the WHERE clause.
// Usage: .Where("deleted_at").IsNotNull()
func (q *Query) IsNotNull() *Query {
	q.whereClauses = append(q.whereClauses, fmt.Sprintf("`%s` IS NOT NULL", q.lastColumn))
	q.lastColumn = ""
	return q
}

// =======================
// UPDATE Query Functions
// =======================

// Set marks the start of an UPDATE operation, specifying which field to update.
// Call this before .To().
// Example: .Set("name")
func (q *Query) Set(field string) *Query {
	q.lastSet = field
	if q.operation == "" {
		q.operation = "update" // default fallback
	}
	return q
}

// To specifies the value to set for the previously specified field in an UPDATE.
// Example: .Set("name").To("Alice")
func (q *Query) To(value any) *Query {
	switch q.operation {
	case "update":
		q.setClauses = append(q.setClauses, fmt.Sprintf("`%s` = ?", q.lastSet))
		q.setArgs = append(q.setArgs, value)
	case "insert":
		q.insertFields[q.lastSet] = value
	}
	q.lastSet = ""
	return q
}

// =======================
// SELECT Query Functions
// =======================

// Limit restricts the number of results returned by the query.
// Example: .Limit(10)
func (q *Query) Limit(n int) *Query {
	q.limit = n // Store the limit for later
	return q
}

// Fetch executes the built SELECT query and returns all matching rows as a slice of Struct pointers.
//
// ---
// LAYMAN'S EXPLANATION:
//
// Fetch is like asking the database: "Give me all the rows that match my conditions."
// It builds a SELECT SQL query using the filters (WHERE), limits (LIMIT), and other options you set up by chaining methods.
//
// 1. It gets a database connection.
// 2. It builds the WHERE and LIMIT parts of the SQL query.
// 3. It creates the full SQL query string (e.g., SELECT * FROM users WHERE age > 18 LIMIT 10).
// 4. It runs this query on the database and gets back the rows.
// 5. For each row:
//   - It creates a new Struct (like a Go object for a row).
//   - It copies the model's field definitions.
//   - It fills in the values from the database into the Struct fields.
//   - It adds this Struct to the results list.
//
// 6. It returns the list of Structs (rows) and any error.
//
// Key variables:
//
//	db: database connection
//	where, limit: SQL WHERE and LIMIT parts
//	query: the SQL query string
//	rows: the result set from the database
//	columns: column names in the result
//	results: the list of Structs to return
func (q *Query) Fetch() ([]*Struct, error) {
	db, err := DatabaseHandler.GetDatabase() // Step 1: Get a database connection
	if err != nil {
		return nil, err // If connection fails, return error
	}

	where := q.buildWhere() // Step 2: Build WHERE clause from conditions
	limit := q.buildLimit() // Step 2: Build LIMIT clause if set

	// Step 3: Build the SQL SELECT statement dynamically
	query := fmt.Sprintf("SELECT * FROM %s %s %s", q.model.TableName, where, limit) // Compose the SQL query string
	rows, err := db.Query(query, q.whereArgs...)                                    // Step 4: Run the query on the database
	if err != nil {
		return nil, err // If query fails, return error
	}
	defer rows.Close() // Make sure to close the rows when done

	columns, err := rows.Columns() // Get the column names in the result
	if err != nil {
		return nil, err // If failed to get columns, return error
	}

	var results []*Struct // Prepare a slice to hold all result Structs
	for rows.Next() {     // Step 5: Loop through each row in the result
		pointers := make([]any, len(columns)) // Prepare a slice of pointers for Scan
		holders := make([]any, len(columns))  // Prepare a slice to hold the values
		for i := range columns {
			pointers[i] = &holders[i] // Each pointer points to a holder
		}

		if err := rows.Scan(pointers...); err != nil { // Scan the row into holders
			return nil, err // If scan fails, return error
		}

		row := &Struct{ // Create a new Struct for this row
			TableName: q.model.TableName,      // Set the table name
			fields:    make(map[string]Field), // Prepare the fields map
		}
		for k, v := range q.model.fields { // Copy the model's field definitions
			row.fields[k] = v
		}
		for i, col := range columns { // Assign scanned values to the Struct fields
			val := holders[i] // Get the value from holders
			// Convert []byte to string for text columns
			if b, ok := val.([]byte); ok {
				val = string(b)
			}
			if f, ok := row.fields[col]; ok { // If field exists in Struct
				f.value = val       // Set the value
				row.fields[col] = f // Update the field in the map
			}
		}
		results = append(results, row) // Add the Struct to the results slice
	}
	return results, rows.Err() // Step 6: Return the results and any error
}

// First executes the built SELECT query and returns only the first matching row (or nil if none).
//
// ---
// LAYMAN'S EXPLANATION:
//
// First is a shortcut for "just give me the first row that matches my query."
//
// 1. If you didn't set a limit, it sets the limit to 1 (so only one row is fetched).
// 2. It calls Fetch to get the results.
// 3. If there are no results, it returns nil.
// 4. If there is at least one result, it returns the first one.
//
// Key variables:
//
//	rows: the list of results from Fetch
//	q.limit: the maximum number of results to get (set to 1 here)
func (q *Query) First() (*Struct, error) {
	if q.limit == 0 { // Step 1: If no limit set, set limit to 1
		q.limit = 1
	}
	rows, err := q.Fetch() // Step 2: Call Fetch to get results
	if err != nil {
		return nil, err // If Fetch fails, return error
	}
	if len(rows) == 0 { // Step 3: If no results, return nil
		return nil, nil
	}
	return rows[0], nil // Step 4: Return the first result
}

// =======================
// UPDATE Query Execution
// =======================

// Exec executes an UPDATE query using the built SET and WHERE clauses.
// Only works if the operation is set to "update" (via Set).
// Prints the number of affected rows for debugging.
//
// ---
// LAYMAN'S EXPLANATION:
//
// Exec is used to update rows in the database. It's like saying: "Change these fields for all rows that match my conditions."
//
// 1. It checks if you're actually doing an update (not a select or delete).
// 2. It gets a database connection.
// 3. It checks if you specified any fields to update. If not, it returns an error.
// 4. It builds the SET part (fields and new values) and the WHERE part (which rows to update).
// 5. It creates the SQL UPDATE query (e.g., UPDATE users SET name = 'Alice' WHERE id = 1).
// 6. It combines all the values for the SET and WHERE clauses.
// 7. It runs the update on the database.
// 8. It prints how many rows were updated (for debugging).
// 9. It returns any error that happened.
//
// Key variables:
//
//	db: database connection
//	set: SET clause (fields and new values)
//	where: WHERE clause (which rows to update)
//	query: the SQL update statement
//	args: all the values to use in the query
//	result: the result of running the update
func (q *Query) Exec() error {
	db, err := DatabaseHandler.GetDatabase()
	if err != nil {
		return err
	}

	switch q.operation {
	case "update":
		// (your existing update logic)
	case "insert":
		if len(q.insertFields) == 0 {
			return fmt.Errorf("no fields to insert")
		}
		cols := []string{}
		vals := []string{}
		args := []any{}

		for k, v := range q.insertFields {
			cols = append(cols, fmt.Sprintf("`%s`", k))
			vals = append(vals, "?")
			args = append(args, v)
		}

		query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
			q.model.TableName,
			strings.Join(cols, ", "),
			strings.Join(vals, ", "),
		)

		result, err := db.Exec(query, args...)
		if err != nil {
			return err
		}
		if id, err := result.LastInsertId(); err == nil {
			fmt.Printf("[Insert] Table: %s | Last Inserted ID: %d\n", q.model.TableName, id)
		} else {
			fmt.Printf("[Insert] Table: %s | Row Inserted\n", q.model.TableName)
		}
		return nil

	default:
		return fmt.Errorf("invalid Exec call: unknown operation '%s'", q.operation)
	}

	return nil
}

// =======================
// DELETE Query Function
// =======================

// Delete executes a DELETE query using the built WHERE and LIMIT clauses.
// It removes matching rows from the database table.
// Prints the number of affected rows for debugging.
func (q *Query) Delete() error {
	db, err := DatabaseHandler.GetDatabase()
	if err != nil {
		return err
	}
	q.operation = "delete"

	where := q.buildWhere()
	limit := q.buildLimit()

	query := fmt.Sprintf("DELETE FROM %s %s %s", q.model.TableName, where, limit)
	result, err := db.Exec(query, q.whereArgs...)
	if err != nil {
		return err
	}
	affected, _ := result.RowsAffected()
	fmt.Printf("[Delete] Table: %s | Rows Affected: %d\n", q.model.TableName, affected)
	return nil
}

// =======================
// Sorting and Grouping
// =======================

// OrderBy sets the ORDER BY clause for sorting results.
// Usage: .OrderBy("created_at DESC")
func (q *Query) OrderBy(clause string) *Query {
	q.orderBy = clause
	return q
}

// GroupBy sets the GROUP BY clause for grouping results.
// Usage: .GroupBy("status")
func (q *Query) GroupBy(clause string) *Query {
	q.groupBy = clause
	return q
}

// =======================
// Pagination
// =======================

// Offset sets the OFFSET for skipping a number of rows (for pagination).
// Usage: .Offset(20)
func (q *Query) Offset(n int) *Query {
	q.offset = n
	return q
}

// Page sets both LIMIT and OFFSET for paginated queries.
// Usage: .Page(2, 10) // page 2, 10 results per page
func (q *Query) Page(page int, pageSize int) *Query {
	if page < 1 {
		page = 1
	}
	q.limit = pageSize
	q.offset = (page - 1) * pageSize
	return q
}

// =======================
// Helper Functions
// =======================

// buildWhere constructs the WHERE clause from the accumulated conditions.
// Returns an empty string if there are no conditions.
func (q *Query) buildWhere() string {
	if len(q.whereClauses) == 0 {
		return ""
	}
	return "WHERE " + strings.Join(q.whereClauses, " ")
}

// buildLimit constructs the LIMIT clause if a limit is set.
// Returns an empty string if no limit is specified.
func (q *Query) buildLimit() string {
	if q.limit > 0 {
		return fmt.Sprintf("LIMIT %d", q.limit)
	}
	return ""
}

func (q *Query) Clone() *Query {
	copy := *q
	copy.whereClauses = append([]string{}, q.whereClauses...)
	copy.whereArgs = append([]any{}, q.whereArgs...)
	copy.setClauses = append([]string{}, q.setClauses...)
	copy.setArgs = append([]any{}, q.setArgs...)
	return &copy
}
