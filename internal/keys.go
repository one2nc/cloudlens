package internal

// ContextKey represents context key.
type ContextKey string

// A collection of context keys.
const (
	KeyFactory ContextKey = "factory"
	KeyApp     ContextKey = "app"
	KeySession ContextKey = "session"
	BucketName ContextKey = "bucket_name"
	ObjectName ContextKey = "object_name"
	FolderName ContextKey = "folder_name"
	KeyAliases ContextKey = "aliases"
)
