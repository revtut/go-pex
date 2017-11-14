package go_pex

import (
	"testing"
	"strconv"
	"strings"
)

func TestHasPermission(t *testing.T) {
	permissionTag := strings.Join(
		[]string{
			strconv.Itoa(PermissionNone),
			strconv.Itoa(PermissionRead),
			strconv.Itoa(PermissionWrite),
			strconv.Itoa(PermissionReadWrite),
		}, "")

	tables := []struct {
		tag      string
		userType uint
		action   uint
		result   bool
	}{
		{permissionTag, 0, ActionRead, false},
		{permissionTag, 0, ActionWrite, false},
		{permissionTag, 1, ActionRead, true},
		{permissionTag, 1, ActionWrite, false},
		{permissionTag, 2, ActionRead, false},
		{permissionTag, 2, ActionWrite, true},
		{permissionTag, 3, ActionRead, true},
		{permissionTag, 3, ActionWrite, true},
	}

	for _, table := range tables {
		hasPermission := HasPermission(table.tag, table.userType, table.action)
		if hasPermission != table.result {
			t.Errorf("Has permission (tag = %s, userType = %d, action = %d) was incorrect, got: %t, want: %t.",
				table.tag, table.userType, table.action, hasPermission, table.result)
		}
	}
}

func TestGetJSONFieldName(t *testing.T) {
	tables := []struct {
		tag    string
		result string
	}{
		{"permissionTag", "permissionTag"},
		{"permissionTag,omitempty",  "permissionTag"},
		{"UpperCase", "UpperCase"},
	}

	for _, table := range tables {
		fieldName := GetJSONFieldName(table.tag)
		if fieldName != table.result {
			t.Errorf("Get JSON field name (tag = %s) was incorrect, got: %s, want: %s.",
				table.tag, fieldName, table.result)
		}
	}
}
