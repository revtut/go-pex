package gopex

// PermissionTag is the tag to use in structs to specify the permissions of each field
const PermissionTag = "pex"

// Actions
const (
	// ActionRead is used when the action is writing
	ActionRead = 0
	// ActionWrite is used when the action is reading
	ActionWrite = 1
)

// Permissions
const (
	// PermissionRead means it has reading permissions
	PermissionRead = "r"
	// PermissionWrite means it has write permissions
	PermissionWrite = "w"
)
