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
	RoleName         ContextKey = "role_name"
)

// TODO rename the file name to constants.go
// TODO move the constants in the package where they used. e.g KeyFactory is used in view package so create constants.go in view package and move view related consts in that file.
