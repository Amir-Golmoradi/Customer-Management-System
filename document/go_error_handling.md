# Error Handling in Go — A Practical Lecture (with Layered Guidance)

## Core principles

Go error handling is intentionally explicit and lightweight. The goal is **clarity and control**, not cleverness.

Key principles:

- **Errors are values**: check them immediately, return them early.
- **Add context, don’t destroy it**: wrap errors with `%w` so callers can inspect the root cause.
- **Separate concerns**: logging is usually done at the edge (API/CLI), not in every layer.
- **Use typed/sentinel errors** for control flow only when necessary.

## The three tools you mentioned

### `errors`

Use it to create and test errors.

- `errors.New("...")` to create a simple error.
- `errors.Is(err, target)` to check wrapped errors.
- `errors.As(err, &typed)` to extract typed errors.

### `fmt`

Use it to format and wrap errors with context.

- `fmt.Errorf("...: %w", err)` keeps the original error accessible.

### `log`

Use it to emit a message. Logging is **not** the same as returning errors.

- Use `log` at the boundary of your system (API/CLI/worker) where you decide what to show or store.
- Avoid logging in lower layers, or you’ll get duplicate logs and lose context.

## Baseline: the default pattern

```go
if err != nil {
    return nil, fmt.Errorf("create customer: %w", err)
}
```

Why this is good:

- Adds **context** for the caller.
- Preserves the original error for `errors.Is` / `errors.As`.

## When to use `errors.New`

Use `errors.New` for **sentinel errors** that represent a stable category, such as “not found.”

```go
var ErrCustomerNotFound = errors.New("customer not found")
```

Then in your code:

```go
if errors.Is(err, ErrCustomerNotFound) {
    // handle not found
}
```

## When to log

Log only at the boundaries:

- HTTP handlers / API layer
- CLI entrypoints
- background job runners

Lower layers (service/repository) should **return** errors, not log them.

Reason: logging everywhere creates duplicates and removes control over how errors are reported.

## Layered error handling (service, repository, API)

This section shows reliable patterns for each layer.

### Repository layer (data access)

Responsibilities:

- Translate driver-specific errors into stable, domain-level errors.
- Wrap unexpected failures with context.

Example:

```go
var ErrCustomerNotFound = errors.New("customer not found")

func (r *Repository) FindCustomerByID(ctx context.Context, id int32) (*Customer, error) {
    c, err := r.queries.GetCustomerByID(ctx, id)
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, ErrCustomerNotFound
        }
        return nil, fmt.Errorf("get customer by id: %w", err)
    }
    return &c, nil
}
```

Why this works:

- Normalizes database errors into a domain signal (`ErrCustomerNotFound`).
- Keeps unexpected errors wrapped for debugging.

### Service layer (business logic)

Responsibilities:

- Orchestrate multiple repositories.
- Apply business rules.
- Avoid logging (leave that to the edge).

Pattern:

- If you add context, wrap the error.
- If not, return it directly.

Example:

```go
func (s *Service) GetCustomerByID(ctx context.Context, id int32) (*Customer, error) {
    c, err := s.repository.FindCustomerByID(ctx, id)
    if err != nil {
        return nil, fmt.Errorf("service get customer by id: %w", err)
    }
    return c, nil
}
```

If you need to map errors:

```go
if errors.Is(err, ErrCustomerNotFound) {
    return nil, ErrCustomerNotFound
}
return nil, fmt.Errorf("service get customer by id: %w", err)
```

### API layer (HTTP/GRPC)

Responsibilities:

- Convert domain errors into user-facing responses.
- Decide what to log and how.

Example (HTTP):

```go
customer, err := service.GetCustomerByID(ctx, id)
if err != nil {
    switch {
    case errors.Is(err, ErrCustomerNotFound):
        http.Error(w, "not found", http.StatusNotFound)
    default:
        log.Printf("get customer: %v", err)
        http.Error(w, "internal server error", http.StatusInternalServerError)
    }
    return
}
```

This preserves correctness while keeping logs and responses clean.

## Good rules of thumb

- **Return early** on errors.
- **Wrap with `%w`** when adding context.
- **Use `errors.Is/As`** to check wrapped errors.
- **Avoid logging** in non-edge layers.
- **Prefer explicit errors** over string matching.

## Anti-patterns to avoid

- **Swallowing errors** (returning `nil` when error happened).
- **Replacing errors** with `errors.New("...")` in every layer.
- **Logging at every layer** (causes duplicates).
- **String comparisons** to detect error types.

## Final mental model

- Repository: normalize low-level errors into domain errors.
- Service: orchestrate and wrap for context, don’t log.
- API: translate errors into responses and log as needed.

This yields reliable, debuggable, and maintainable error handling.
