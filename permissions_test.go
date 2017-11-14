package go_pex

import (
	"testing"
	"strconv"
	"strings"
	"reflect"
)

// Simple struct test
type AStruct struct {
	Number int    `pex:"0123"`
	Text   string `pex:"0123" json:"Label"`
}

// Anonymous struct test
type BStruct struct {
	AStruct
	Boolean bool `pex:"0123"`
}

// Complex struct test
type CStruct struct {
	Struct    AStruct     `pex:"0123"`
	Pointer   *AStruct    `pex:"0123"`
	Interface interface{} `pex:"0123"`
}

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
		{permissionTag, 0, ActionWrite, false},
		{permissionTag, 0, ActionRead, false},
		{permissionTag, 1, ActionWrite, false},
		{permissionTag, 1, ActionRead, true},
		{permissionTag, 2, ActionWrite, true},
		{permissionTag, 2, ActionRead, false},
		{permissionTag, 3, ActionWrite, true},
		{permissionTag, 3, ActionRead, true},
	}

	for _, table := range tables {
		hasPermission := HasPermission(table.tag, table.userType, table.action)
		if hasPermission != table.result {
			t.Errorf("%s (tag = %s, userType = %d, action = %d) was incorrect, got: %t, want: %t.",
				t.Name(), table.tag, table.userType, table.action, hasPermission, table.result)
		}
	}
}

func TestGetJSONFieldName(t *testing.T) {
	tables := []struct {
		tag    string
		result string
	}{
		{"permissionTag", "permissionTag"},
		{"permissionTag,omitempty", "permissionTag"},
		{"UpperCase", "UpperCase"},
	}

	for _, table := range tables {
		fieldName := GetJSONFieldName(table.tag)
		if fieldName != table.result {
			t.Errorf("%s (tag = %s) was incorrect, got: %s, want: %s.",
				t.Name(), table.tag, fieldName, table.result)
		}
	}
}

func TestCleanSingleObject(t *testing.T) {
	t.Run("TestCleanSingleObjectSimple", testCleanSingleObjectSimple)
	t.Run("TestCleanSingleObjectAnonymousStruct", testCleanSingleObjectAnonymousStruct)
	t.Run("TestCleanSingleObjectStructField", testCleanSingleObjectStructField)
}

func testCleanSingleObjectSimple(t *testing.T) {
	t.Parallel()
	baseAStruct := AStruct{Number: 10, Text: "ABC"}

	tables := []struct {
		object   interface{}
		userType uint
		action   uint
		result   interface{}
	}{
		// Struct
		{baseAStruct, 0, ActionRead, map[string]interface{}{}},
		{baseAStruct, 1, ActionRead, map[string]interface{}{"Number": 10, "Label": "ABC"}},
		{baseAStruct, 2, ActionRead, map[string]interface{}{}},
		{baseAStruct, 3, ActionRead, map[string]interface{}{"Number": 10, "Label": "ABC"}},

		{baseAStruct, 0, ActionWrite, map[string]interface{}{}},
		{baseAStruct, 1, ActionWrite, map[string]interface{}{}},
		{baseAStruct, 2, ActionWrite, map[string]interface{}{"Number": 10, "Label": "ABC"}},
		{baseAStruct, 3, ActionWrite, map[string]interface{}{"Number": 10, "Label": "ABC"}},
		// Pointer
		{&baseAStruct, 0, ActionRead, map[string]interface{}{}},
		{&baseAStruct, 1, ActionRead, map[string]interface{}{"Number": 10, "Label": "ABC"}},
		{&baseAStruct, 2, ActionRead, map[string]interface{}{}},
		{&baseAStruct, 3, ActionRead, map[string]interface{}{"Number": 10, "Label": "ABC"}},

		{&baseAStruct, 0, ActionWrite, map[string]interface{}{}},
		{&baseAStruct, 1, ActionWrite, map[string]interface{}{}},
		{&baseAStruct, 2, ActionWrite, map[string]interface{}{"Number": 10, "Label": "ABC"}},
		{&baseAStruct, 3, ActionWrite, map[string]interface{}{"Number": 10, "Label": "ABC"}},
	}

	for _, table := range tables {
		cleanedObject := CleanSingleObject(table.object, table.userType, table.action)
		if !reflect.DeepEqual(cleanedObject, table.result) {
			t.Errorf("%s (object = %+v, userType = %d, action = %d) was incorrect, got: %+v, want: %+v.",
				t.Name(), table.object, table.userType, table.action, cleanedObject, table.result)
		}
	}
}

func testCleanSingleObjectAnonymousStruct(t *testing.T) {
	t.Parallel()
	baseBStruct := BStruct{AStruct: AStruct{Number: 10, Text: "ABC"}, Boolean: false}

	tables := []struct {
		object   interface{}
		userType uint
		action   uint
		result   interface{}
	}{
		// Struct
		{baseBStruct, 0, ActionRead, map[string]interface{}{}},
		{baseBStruct, 1, ActionRead, map[string]interface{}{"Number": 10, "Label": "ABC", "Boolean": false}},
		{baseBStruct, 2, ActionRead, map[string]interface{}{}},
		{baseBStruct, 3, ActionRead, map[string]interface{}{"Number": 10, "Label": "ABC", "Boolean": false}},

		{baseBStruct, 0, ActionWrite, map[string]interface{}{}},
		{baseBStruct, 1, ActionWrite, map[string]interface{}{}},
		{baseBStruct, 2, ActionWrite, map[string]interface{}{"Number": 10, "Label": "ABC", "Boolean": false}},
		{baseBStruct, 3, ActionWrite, map[string]interface{}{"Number": 10, "Label": "ABC", "Boolean": false}},
		// Pointer
		{&baseBStruct, 0, ActionRead, map[string]interface{}{}},
		{&baseBStruct, 1, ActionRead, map[string]interface{}{"Number": 10, "Label": "ABC", "Boolean": false}},
		{&baseBStruct, 2, ActionRead, map[string]interface{}{}},
		{&baseBStruct, 3, ActionRead, map[string]interface{}{"Number": 10, "Label": "ABC", "Boolean": false}},

		{&baseBStruct, 0, ActionWrite, map[string]interface{}{}},
		{&baseBStruct, 1, ActionWrite, map[string]interface{}{}},
		{&baseBStruct, 2, ActionWrite, map[string]interface{}{"Number": 10, "Label": "ABC", "Boolean": false}},
		{&baseBStruct, 3, ActionWrite, map[string]interface{}{"Number": 10, "Label": "ABC", "Boolean": false}},
	}

	for _, table := range tables {
		cleanedObject := CleanSingleObject(table.object, table.userType, table.action)
		if !reflect.DeepEqual(cleanedObject, table.result) {
			t.Errorf("%s (object = %+v, userType = %d, action = %d) was incorrect, got: %+v, want: %+v.",
				t.Name(), table.object, table.userType, table.action, cleanedObject, table.result)
		}
	}
}

func testCleanSingleObjectStructField(t *testing.T) {
	t.Parallel()

	baseAStruct := AStruct{Number: 10, Text: "ABC"}
	baseCStruct := CStruct{Struct: baseAStruct, Pointer: &baseAStruct, Interface: baseAStruct}

	tables := []struct {
		object   interface{}
		userType uint
		action   uint
		result   interface{}
	}{
		// Struct
		{baseCStruct, 0, ActionRead, map[string]interface{}{}},
		{baseCStruct, 1, ActionRead, map[string]interface{}{
			"Struct":    map[string]interface{}{"Number": 10, "Label": "ABC"},
			"Pointer":   map[string]interface{}{"Number": 10, "Label": "ABC"},
			"Interface": map[string]interface{}{"Number": 10, "Label": "ABC"},}},
		{baseCStruct, 2, ActionRead, map[string]interface{}{}},
		{baseCStruct, 3, ActionRead, map[string]interface{}{
			"Struct":    map[string]interface{}{"Number": 10, "Label": "ABC"},
			"Pointer":   map[string]interface{}{"Number": 10, "Label": "ABC"},
			"Interface": map[string]interface{}{"Number": 10, "Label": "ABC"},}},

		{baseCStruct, 0, ActionWrite, map[string]interface{}{}},
		{baseCStruct, 1, ActionWrite, map[string]interface{}{}},
		{baseCStruct, 2, ActionWrite, map[string]interface{}{
			"Struct":    map[string]interface{}{"Number": 10, "Label": "ABC"},
			"Pointer":   map[string]interface{}{"Number": 10, "Label": "ABC"},
			"Interface": map[string]interface{}{"Number": 10, "Label": "ABC"},}},
		{baseCStruct, 3, ActionWrite, map[string]interface{}{
			"Struct":    map[string]interface{}{"Number": 10, "Label": "ABC"},
			"Pointer":   map[string]interface{}{"Number": 10, "Label": "ABC"},
			"Interface": map[string]interface{}{"Number": 10, "Label": "ABC"},}},

		// Pointer
		{&baseCStruct, 0, ActionRead, map[string]interface{}{}},
		{&baseCStruct, 1, ActionRead, map[string]interface{}{
			"Struct":    map[string]interface{}{"Number": 10, "Label": "ABC"},
			"Pointer":   map[string]interface{}{"Number": 10, "Label": "ABC"},
			"Interface": map[string]interface{}{"Number": 10, "Label": "ABC"},}},
		{&baseCStruct, 2, ActionRead, map[string]interface{}{}},
		{&baseCStruct, 3, ActionRead, map[string]interface{}{
			"Struct":    map[string]interface{}{"Number": 10, "Label": "ABC"},
			"Pointer":   map[string]interface{}{"Number": 10, "Label": "ABC"},
			"Interface": map[string]interface{}{"Number": 10, "Label": "ABC"},}},

		{&baseCStruct, 0, ActionWrite, map[string]interface{}{}},
		{&baseCStruct, 1, ActionWrite, map[string]interface{}{}},
		{&baseCStruct, 2, ActionWrite, map[string]interface{}{
			"Struct":    map[string]interface{}{"Number": 10, "Label": "ABC"},
			"Pointer":   map[string]interface{}{"Number": 10, "Label": "ABC"},
			"Interface": map[string]interface{}{"Number": 10, "Label": "ABC"},}},
		{&baseCStruct, 3, ActionWrite, map[string]interface{}{
			"Struct":    map[string]interface{}{"Number": 10, "Label": "ABC"},
			"Pointer":   map[string]interface{}{"Number": 10, "Label": "ABC"},
			"Interface": map[string]interface{}{"Number": 10, "Label": "ABC"},}},
	}

	for _, table := range tables {
		cleanedObject := CleanSingleObject(table.object, table.userType, table.action)
		if !reflect.DeepEqual(cleanedObject, table.result) {
			t.Errorf("%s (object = %+v, userType = %d, action = %d) was incorrect, got: %+v, want: %+v.",
				t.Name(), table.object, table.userType, table.action, cleanedObject, table.result)
		}
	}
}

func TestCleanMultipleObjects(t *testing.T) {
	t.Run("TestCleanMultipleObjectsArray", testCleanMultipleObjectsArray)
}

func testCleanMultipleObjectsArray(t *testing.T) {
	t.Parallel()
	baseArray := [3]int{1, 2, 3}

	tables := []struct {
		object   interface{}
		userType uint
		action   uint
		result   interface{}
	}{
		// Struct
		{baseArray, 0, ActionRead, baseArray},
		{baseArray, 1, ActionRead, baseArray},
		{baseArray, 2, ActionRead, baseArray},
		{baseArray, 3, ActionRead, baseArray},

		{baseArray, 0, ActionWrite, baseArray},
		{baseArray, 1, ActionWrite, baseArray},
		{baseArray, 2, ActionWrite, baseArray},
		{baseArray, 3, ActionWrite, baseArray},
		// Pointer
		{&baseArray, 0, ActionRead, baseArray},
		{&baseArray, 1, ActionRead, baseArray},
		{&baseArray, 2, ActionRead, baseArray},
		{&baseArray, 3, ActionRead, baseArray},

		{&baseArray, 0, ActionWrite, baseArray},
		{&baseArray, 1, ActionWrite, baseArray},
		{&baseArray, 2, ActionWrite, baseArray},
		{&baseArray, 3, ActionWrite, baseArray},
	}

	for _, table := range tables {
		cleanedObject := CleanMultipleObjects(table.object, table.userType, table.action)
		if !reflect.DeepEqual(cleanedObject, table.result) {
			t.Errorf("%s (object = %+v, userType = %d, action = %d) was incorrect, got: %+v, want: %+v.",
				t.Name(), table.object, table.userType, table.action, cleanedObject, table.result)
		}
	}
}

/*baseAStruct := AStruct{Number: 10, Text: "ABC"}
	baseSlice := []AStruct{baseAStruct, baseAStruct}

	tables := []struct {
		object   interface{}
		userType uint
		action   uint
		result   interface{}
	}{
		// Struct
		{baseSlice, 3, ActionRead, []map[string]interface{}{
			{"Number": 10, "Label": "ABC"},
			{"Number": 10, "Label": "ABC"},
		}},
	}

	for _, table := range tables {
		cleanedObjects := CleanMultipleObjects(table.object, table.userType, table.action)
		if !reflect.DeepEqual(cleanedObjects, table.result) {
			t.Errorf("%s (object = %+v, userType = %d, action = %d) was incorrect, got: %+v, want: %+v.",
				t.Name(), table.object, table.userType, table.action, cleanedObjects, table.result)
		}
	}*/