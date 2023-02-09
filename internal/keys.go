package internal

// ContextKey represents context key.
type ContextKey string

// A collection of context keys.
const (
	KeyFactory       ContextKey = "factory"
	KeyApp           ContextKey = "app"
	KeyActiveProfile ContextKey = "active-profile"
	KeyActiveRegion  ContextKey = "active-region"
	KeySession       ContextKey = "session"
	BucketName       ContextKey = "bucket_name"
	ObjectName       ContextKey = "object_name"
	FolderName       ContextKey = "folder_name"
	KeyAliases       ContextKey = "aliases"
	UserName         ContextKey = "user_name"
	GroupName        ContextKey = "group_name"
)
