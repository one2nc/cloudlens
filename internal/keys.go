package internal

// ContextKey represents context key.
type ContextKey string

// A collection of context keys.
const (
	KeyApp     ContextKey = "app"
	KeySession ContextKey = "session"
)
