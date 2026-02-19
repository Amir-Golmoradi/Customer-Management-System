## Issue: Go Idioms, Error Handling, and `sqlc` Consistency in `sql_customer_repository.go`

### What you did
You've implemented the `CustomerRepository` interface in `sql_customer_repository.go`, but there are several areas that can be improved for Go idiomatic style, robust error handling, and correct interaction with your `sqlc` generated code.

```go
package repository

import (
	"context"
	"database/sql"
	"errors"
	"log"

	database "github.com/Amir-Golmoradi/Customer-Management-System/internal/database/generated"
)

type sqlCustomerRepository struct {
	queries *database.Queries
}

// NewSqlCustomerRepository Constructor
func NewSqlCustomerRepository(q *database.Queries) CustomerRepository {
	return &sqlCustomerRepository{queries: q}
}

func (s *sqlCustomerRepository) FindAllCustomers(ctx context.Context) ([]database.Customer, error) {
	return s.queries.ListCustomers(ctx)
}

func (s *sqlCustomerRepository) FindCustomerByID(ctx context.Context, id int32) (*database.Customer, error) {
	customer, err := s.queries.GetCustomerByID(ctx, id) // <-- Issue 1: Method name mismatch
	if err != nil {
		log.Fatalf("No customer found with id %v.\n", id) // <-- Issue 2: Using log.Fatalf
	}
	return &customer, nil
}

func (s *sqlCustomerRepository) FindCustomerByEmail(ctx context.Context, email string) (*database.Customer, error) {
	customerEmail, err := s.queries.GetCustomerByEmail(ctx, email)
	if err != nil {
		log.Fatalf("No customer found with email %s.", email) // <-- Issue 2: Using log.Fatalf
	}
	return &customerEmail, nil
}

func (s *sqlCustomerRepository) CreateNewCustomer(ctx context.Context, name, email, password string) (*database.Customer, error) {
	params := database.CreateCustomerParams{
		Name:     name,
		Email:    email,
		Password: password,
	}
	customer, err := s.queries.CreateCustomer(ctx, params)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) { // <-- Issue 3: Incorrect error check for Create
		}
		return err // <-- Issue 4: Return type mismatch
	}
	return nil // <-- Issue 4: Return type mismatch
}

func (s *sqlCustomerRepository) UpdateExistingCustomer(ctx context.Context, id int32, name, email, password string) error {
	params := database.UpdateCustomerParams{ // <-- Issue 5: UpdateCustomerParams not generated
		ID:       id,
		Name:     name,
		Email:    email,
		Password: password,
	}
	_, err := s.queries.UpdateCustomer(ctx, params) // <-- Issue 6: UpdateCustomer method not generated
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}
		return err
	}
	return nil
}

func (s *sqlCustomerRepository) DeleteCustomerByEmail(ctx context.Context, email string) error {
	_, err := s.queries.DeleteCustomerByEmail(ctx, email) // <-- Issue 7: DeleteCustomerByEmail method not generated
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}
		return err
	}
	return nil
}
```

### Why it's a problem

1.  **Constructor Naming and Comment**:
    *   `NewSqlCustomerRepository` should be `NewSQLCustomerRepository` (acronym capitalization).
    *   The comment `// NewSqlCustomerRepository Constructor` is redundant and non-idiomatic. The function name itself implies it's a constructor.

2.  **`FindCustomerByID` and `FindCustomerByEmail` - `log.Fatalf`**:
    *   Using `log.Fatalf` in a repository method is a critical error. `Fatalf` terminates the program immediately. A repository should *return* errors to its caller (e.g., a service layer), allowing the application to handle different error scenarios gracefully (e.g., return a 404 Not Found, log the error, etc.).
    *   The `sqlc` generated `GetCustomer` (by ID) and `GetCustomerByEmail` methods return `sql.ErrNoRows` if no customer is found. This error should be propagated.

3.  **`FindCustomerByID` - Method Name Mismatch**:
    *   Your `sqlc` generated `query.sql.go` contains `GetCustomer(ctx context.Context, id int32)` for fetching by ID, not `GetCustomerByID`. You're calling a non-existent method.

4.  **`CreateNewCustomer` - Incorrect Error Handling and Return Type**:
    *   `errors.Is(err, sql.ErrNoRows)` is generally not relevant for an `INSERT` operation. `sql.ErrNoRows` is returned when a `SELECT` query finds no rows. For `INSERT`, you might get errors like unique constraint violations.
    *   The interface expects `(*database.Customer, error)`, but your implementation attempts to `return err` (an `error` type) or `return nil` (a `*database.Customer` type, but `nil` for the error). This is a type mismatch. You must return both `*database.Customer` and `error`. If an error occurs, return `nil` for the customer and the error. If successful, return the created customer and `nil` for the error.

5.  **`UpdateExistingCustomer` and `DeleteCustomerByEmail` - Missing `sqlc` Generated Methods**:
    *   Your `sqlc` generated `query.sql.go` (which I reviewed earlier) does *not* contain `UpdateCustomer` or `DeleteCustomerByEmail` methods, nor does it define `database.UpdateCustomerParams`.
    *   This means you need to add `UPDATE` and `DELETE` queries to your `internal/database/sql/query.sql` file and then re-run `sqlc generate` to create these methods and structs. Without them, your repository implementation will not compile or function correctly.

### Correct Go-idiomatic way

1.  **Constructor**: Rename to `NewSQLCustomerRepository` and remove the redundant comment.
2.  **Error Handling**: Always return errors from repository methods. Let the caller decide how to handle them. Check for `sql.ErrNoRows` specifically for "not found" scenarios and return a custom error or `nil` for the entity.
3.  **`sqlc` Consistency**: Ensure the methods you call on `s.queries` actually exist in the `database.Queries` interface (generated by `sqlc`). If they don't, you need to define the corresponding SQL queries and regenerate.
4.  **Return Types**: Match the return types of your implementation methods exactly to the interface definition.

### Fixed example (assuming `UpdateCustomer` and `DeleteCustomerByEmail` queries are added to `query.sql` and `sqlc` is regenerated)

First, you would need to add these to `internal/database/sql/query.sql`:

```sql
-- name: UpdateCustomer :one
UPDATE customers
SET name = $2, email = $3, password = $4
WHERE id = $1
RETURNING id, name, email, password, created_at;

-- name: DeleteCustomerByEmail :exec
DELETE FROM customers
WHERE email = $1;
```

Then, after running `sqlc generate`, your `query.sql.go` would have `UpdateCustomer` and `DeleteCustomerByEmail` methods, and `UpdateCustomerParams` struct.

Now, the corrected `internal/repository/sql_customer_repository.go`:

```go
package repository

import (
	"context"
	"database/sql" // Import sql for sql.ErrNoRows
	"errors"

	database "github.com/Amir-Golmoradi/Customer-Management-System/internal/database/generated"
)

// Define a custom error for when a customer is not found
var ErrCustomerNotFound = errors.New("customer not found")

type sqlCustomerRepository struct {
	queries *database.Queries
}

// NewSQLCustomerRepository creates and returns a new CustomerRepository interface
// backed by a SQL database implementation.
func NewSQLCustomerRepository(q *database.Queries) CustomerRepository {
	return &sqlCustomerRepository{queries: q}
}

func (s *sqlCustomerRepository) FindAllCustomers(ctx context.Context) ([]database.Customer, error) {
	customers, err := s.queries.ListCustomers(ctx)
	if err != nil {
		return nil, err
	}
	return customers, nil
}

func (s *sqlCustomerRepository) FindCustomerByID(ctx context.Context, id int32) (*database.Customer, error) {
	// Correctly call the generated GetCustomer method
	customer, err := s.queries.GetCustomer(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrCustomerNotFound // Return a specific error for not found
		}
		return nil, err // Propagate other errors
	}
	return &customer, nil
}

func (s *sqlCustomerRepository) FindCustomerByEmail(ctx context.Context, email string) (*database.Customer, error) {
	customer, err := s.queries.GetCustomerByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrCustomerNotFound // Return a specific error for not found
		}
		return nil, err // Propagate other errors
	}
	return &customer, nil
}

func (s *sqlCustomerRepository) CreateNewCustomer(ctx context.Context, name, email, password string) (*database.Customer, error) {
	params := database.CreateCustomerParams{
		Name:     name,
		Email:    email,
		Password: password,
	}
	customer, err := s.queries.CreateCustomer(ctx, params)
	if err != nil {
		// Handle specific errors like unique constraint violation if needed
		// For now, just return the error
		return nil, err
	}
	return &customer, nil // Return the created customer and nil error
}

func (s *sqlCustomerRepository) UpdateExistingCustomer(ctx context.Context, id int32, name, email, password string) error {
	// Assuming UpdateCustomerParams and UpdateCustomer method are generated by sqlc
	params := database.UpdateCustomerParams{
		ID:       id,
		Name:     name,
		Email:    email,
		Password: password,
	}
	_, err := s.queries.UpdateCustomer(ctx, params) // Assuming this returns a Customer or similar
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrCustomerNotFound // Or a more specific "no rows updated" error
		}
		return err
	}
	return nil
}

func (s *sqlCustomerRepository) DeleteCustomerByEmail(ctx context.Context, email string) error {
	// Assuming DeleteCustomerByEmail method is generated by sqlc
	err := s.queries.DeleteCustomerByEmail(ctx, email) // :exec queries return only error
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// If the delete query returns ErrNoRows, it means no row was found to delete.
			// Depending on requirements, you might return ErrCustomerNotFound or just nil.
			return ErrCustomerNotFound
		}
		return err
	}
	return nil
}
```