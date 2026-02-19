# Go `context.Context` — A Practical Lecture

## What `context.Context` is

`context.Context` is a standard library interface used to carry **request-scoped signals** and **metadata** across API boundaries and goroutines. It lets you:

- **Cancel work** when the caller no longer needs the result.
- **Enforce deadlines/timeouts** for operations.
- **Propagate request-scoped values** (like request IDs) across layers.

In short: it is the control plane for *when* work should stop and *which* metadata should follow the work.

## The interface (minimal but powerful)

At its core, `Context` exposes four methods:

- `Deadline() (time.Time, bool)`
- `Done() <-chan struct{}`
- `Err() error`
- `Value(key any) any`

These enable cancellation and value propagation without coupling your code to a specific implementation.

## Cancellation and deadlines

Cancellation is central to `context.Context`. It prevents wasted work and resource leaks.

### How cancellation flows

- A parent context is created (often from `context.Background()` or `context.TODO()`).
- Child contexts are derived with `context.WithCancel`, `context.WithTimeout`, or `context.WithDeadline`.
- Cancelling a parent cancels all its children.

### Typical flow

1. Handler creates a context with a timeout.
2. It passes the context to the repository or external API.
3. If the timeout expires or the client disconnects, the context is canceled.
4. Lower layers stop work and return `context.Canceled` or `context.DeadlineExceeded`.

## Propagating request-scoped values

`Context` can carry small pieces of metadata like:

- Request IDs / correlation IDs
- User identity / auth claims (sparingly)
- Feature flags (sparingly)

**Rules of thumb:**

- Only store data that is truly request-scoped.
- Avoid passing large objects or business data.
- Prefer explicit parameters for core data — `Value` is for cross-cutting concerns.

## Where to use it

`Context` should be the **first parameter** in functions that do I/O or long-running work:

- Database queries
- HTTP calls
- RPCs
- Background jobs

In Go, the convention is:

```go
func (s *Service) DoSomething(ctx context.Context, ...) error
```

## Common patterns

### Creating a base context

```go
ctx := context.Background()
```

### With timeout

```go
ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
defer cancel()
```

### Cancel on shutdown

```go
ctx, cancel := context.WithCancel(context.Background())
// call cancel() on shutdown
```

### Listening for cancellation

```go
select {
case <-ctx.Done():
    return ctx.Err()
default:
    // continue work
}
```

## Method-by-method usage and examples

This section shows **when** to use each `Context` method and **how** it typically looks in real code.

### `Deadline()` — when you need to inspect time limits

**When to use:**

- You want to log or adapt behavior based on remaining time.
- You need to set downstream timeouts to be smaller than the parent deadline.

**Example:**

```go
if deadline, ok := ctx.Deadline(); ok {
    remaining := time.Until(deadline)
    if remaining < 200*time.Millisecond {
        // skip expensive work, return a fast path
        return fastPath(ctx)
    }
}
```

### `Done()` — when you need to stop work promptly

**When to use:**

- You are in a loop or long-running process and must exit on cancellation.
- You are waiting on multiple events and want cancellation to win.

**Example (loop):**

```go
for {
    select {
    case <-ctx.Done():
        return ctx.Err()
    default:
        // do chunk of work
    }
}
```

**Example (select over work + cancellation):**

```go
select {
case result := <-workCh:
    return result, nil
case <-ctx.Done():
    return nil, ctx.Err()
}
```

### `Err()` — when you need to know why it ended

**When to use:**

- You want to return the canonical reason for cancellation.
- You need to distinguish timeout vs manual cancel.

**Example:**

```go
if err := ctx.Err(); err != nil {
    // err is context.Canceled or context.DeadlineExceeded
    return nil, err
}
```

### `Value()` — when you need request-scoped metadata

**When to use:**

- You need cross-cutting metadata (request IDs, auth context) without changing every function signature.

**Example (defining and reading a key):**

```go
type ctxKey string

const requestIDKey ctxKey = "request_id"

func WithRequestID(ctx context.Context, id string) context.Context {
    return context.WithValue(ctx, requestIDKey, id)
}

func RequestIDFromContext(ctx context.Context) (string, bool) {
    v := ctx.Value(requestIDKey)
    id, ok := v.(string)
    return id, ok
}
```

**Avoid:** storing large objects, config, or business entities in context values.

## Errors you will see

- `context.Canceled` — the context was explicitly canceled.
- `context.DeadlineExceeded` — the context hit its timeout or deadline.

Your code should propagate these errors upward so the caller can handle them.

## Best practices

- Always pass `context.Context` down the call chain.
- Do not store contexts in structs; pass them as parameters.
- Do not use `context.Background()` inside libraries or lower layers; accept a `ctx` from callers.
- Always call the `cancel()` function when you create a derived context.
- Keep values small and scoped to infrastructure concerns.

## Why it matters

Without `context.Context`, cancellation and timeouts become ad-hoc and unreliable. In real systems, this causes:

- Leaked goroutines
- Stuck DB connections
- Slow shutdowns
- Poor observability

`Context` gives you a uniform way to coordinate lifetime and metadata across concurrent, layered code.

## Final mental model

Think of `context.Context` as:

- **A signal**: “Stop if you see this is done.”
- **A deadline**: “Don’t go past this time.”
- **A metadata carrier**: “This request has these tags.”

It is not a general-purpose data store or business model carrier. It is the control and traceability channel for work in Go.
