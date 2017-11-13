package go_pex

const PERMISSION_TAG = "pex"

const (
	// ActionWrite is used when the action is writing
	ActionWrite = 0
	// ActionRead is used when the action is reading
	ActionRead = 1
)

const (
	// PermissionNone means it hasn't any permission
	PermissionNone = 0
	// PermissionRead means it has reading permissions
	PermissionRead = 1
	// PermissionWrite means it has write permissions
	PermissionWrite = 2
	// PermissionReadWrite means it has read and write permissions
	PermissionReadWrite = 3
)
