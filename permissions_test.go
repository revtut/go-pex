package gopex

import (
	"reflect"
	"testing"
	"time"
	"database/sql"
)

// Simple struct
type AStruct struct {
	Number     int    `pex:"guest:,user:r,sys:w,admin:rw"`
	Text       string `pex:"guest:,user:r,sys:w,admin:rw" json:"Label"`
	otherField int
}

// Anonymous struct
type BStruct struct {
	AStruct
	Boolean bool `pex:"guest:,user:r,sys:w,admin:rw"`
}

// Complex struct
type CStruct struct {
	Struct    AStruct     `pex:"guest:,user:r,sys:w,admin:rw"`
	Pointer   *AStruct    `pex:"guest:,user:r,sys:w,admin:rw"`
	Interface interface{} `pex:"guest:,user:r,sys:w,admin:rw"`
}

// Empty struct
type DStruct struct {
}

// Special struct
type EStruct struct {
	Start  time.Time     `pex:"guest:,user:r,sys:w,admin:rw"`
	Stop   *time.Time    `pex:"guest:,user:r,sys:w,admin:rw"`
	Number sql.NullInt64 `pex:"guest:,user:r,sys:w,admin:rw"`
}

// Arrays and slices in struct fields
type FStruct struct {
	Name  string    `pex:"guest:,user:r,sys:w,admin:rw"`
	Array [2]int    `pex:"guest:,user:r,sys:w,admin:rw"`
	Slice []AStruct `pex:"guest:,user:r,sys:w,admin:rw"`
}

// Struct with shuffled permissions
type GStruct struct {
	Name    string `pex:"guest:,user:r,sys:w,admin:rw"`
	Version uint   `pex:"guest:rw,user:w,sys:r,admin:"`
}

// Struct with one special object inside
type HStruct struct {
	Boolean sql.NullBool `pex:"guest:,user:r,sys:w,admin:rw"`
	otherField int
}

func TestExtractSingleObjectFields(t *testing.T) {
	t.Run("TestExtractSingleObjectFieldsNonStruct", testExtractSingleObjectFieldsNonStruct)
	t.Run("TestExtractSingleObjectFieldsSimple", testExtractSingleObjectFieldsSimple)
	t.Run("TestExtractSingleObjectFieldsAnonymousStruct", testExtractSingleObjectFieldsAnonymousStruct)
	t.Run("TestExtractSingleObjectFieldsStructField", testExtractSingleObjectFieldsStructField)
	t.Run("TestExtractSingleObjectFieldsNil", testExtractSingleObjectFieldsNil)
	t.Run("TestExtractSingleObjectFieldsSpecial", testExtractSingleObjectFieldsSpecial)
}

func testExtractSingleObjectFieldsNonStruct(t *testing.T) {
	t.Parallel()

	baseValue := 10.0

	tables := []struct {
		object   interface{}
		userType string
		action   uint
		expected interface{}
	}{
		// Struct
		{baseValue, "guest", ActionRead, baseValue},
		{baseValue, "user", ActionRead, baseValue},
		{baseValue, "sys", ActionRead, baseValue},
		{baseValue, "admin", ActionRead, baseValue},

		{baseValue, "guest", ActionWrite, baseValue},
		{baseValue, "user", ActionWrite, baseValue},
		{baseValue, "sys", ActionWrite, baseValue},
		{baseValue, "admin", ActionWrite, baseValue},
		// Pointer
		{&baseValue, "guest", ActionRead, baseValue},
		{&baseValue, "user", ActionRead, baseValue},
		{&baseValue, "sys", ActionRead, baseValue},
		{&baseValue, "admin", ActionRead, baseValue},

		{&baseValue, "guest", ActionWrite, baseValue},
		{&baseValue, "user", ActionWrite, baseValue},
		{&baseValue, "sys", ActionWrite, baseValue},
		{&baseValue, "admin", ActionWrite, baseValue},
	}

	for _, table := range tables {
		actual := ExtractSingleObjectFields(table.object, table.userType, table.action)
		if !reflect.DeepEqual(actual, table.expected) {
			t.Errorf("%s (object = %+v, userType = %s, action = %d) was incorrect, got: %+v, want: %+v.",
				t.Name(), table.object, table.userType, table.action, actual, table.expected)
		}
	}
}

func testExtractSingleObjectFieldsSimple(t *testing.T) {
	t.Parallel()

	baseAStruct := AStruct{Number: 10, Text: "ABC"}

	tables := []struct {
		object   interface{}
		userType string
		action   uint
		expected interface{}
	}{
		// Struct
		{baseAStruct, "guest", ActionRead, map[string]interface{}{}},
		{baseAStruct, "user", ActionRead, map[string]interface{}{"Number": 10, "Label": "ABC"}},
		{baseAStruct, "sys", ActionRead, map[string]interface{}{}},
		{baseAStruct, "admin", ActionRead, map[string]interface{}{"Number": 10, "Label": "ABC"}},

		{baseAStruct, "guest", ActionWrite, map[string]interface{}{}},
		{baseAStruct, "user", ActionWrite, map[string]interface{}{}},
		{baseAStruct, "sys", ActionWrite, map[string]interface{}{"Number": 10, "Label": "ABC"}},
		{baseAStruct, "admin", ActionWrite, map[string]interface{}{"Number": 10, "Label": "ABC"}},
		// Pointer
		{&baseAStruct, "guest", ActionRead, map[string]interface{}{}},
		{&baseAStruct, "user", ActionRead, map[string]interface{}{"Number": 10, "Label": "ABC"}},
		{&baseAStruct, "sys", ActionRead, map[string]interface{}{}},
		{&baseAStruct, "admin", ActionRead, map[string]interface{}{"Number": 10, "Label": "ABC"}},

		{&baseAStruct, "guest", ActionWrite, map[string]interface{}{}},
		{&baseAStruct, "user", ActionWrite, map[string]interface{}{}},
		{&baseAStruct, "sys", ActionWrite, map[string]interface{}{"Number": 10, "Label": "ABC"}},
		{&baseAStruct, "admin", ActionWrite, map[string]interface{}{"Number": 10, "Label": "ABC"}},
	}

	for _, table := range tables {
		actual := ExtractSingleObjectFields(table.object, table.userType, table.action)
		if !reflect.DeepEqual(actual, table.expected) {
			t.Errorf("%s (object = %+v, userType = %s, action = %d) was incorrect, got: %+v, want: %+v.",
				t.Name(), table.object, table.userType, table.action, actual, table.expected)
		}
	}
}

func testExtractSingleObjectFieldsAnonymousStruct(t *testing.T) {
	t.Parallel()
	baseBStruct := BStruct{AStruct: AStruct{Number: 10, Text: "ABC"}, Boolean: false}

	tables := []struct {
		object   interface{}
		userType string
		action   uint
		expected interface{}
	}{
		// Struct
		{baseBStruct, "guest", ActionRead, map[string]interface{}{}},
		{baseBStruct, "user", ActionRead, map[string]interface{}{"Number": 10, "Label": "ABC", "Boolean": false}},
		{baseBStruct, "sys", ActionRead, map[string]interface{}{}},
		{baseBStruct, "admin", ActionRead, map[string]interface{}{"Number": 10, "Label": "ABC", "Boolean": false}},

		{baseBStruct, "guest", ActionWrite, map[string]interface{}{}},
		{baseBStruct, "user", ActionWrite, map[string]interface{}{}},
		{baseBStruct, "sys", ActionWrite, map[string]interface{}{"Number": 10, "Label": "ABC", "Boolean": false}},
		{baseBStruct, "admin", ActionWrite, map[string]interface{}{"Number": 10, "Label": "ABC", "Boolean": false}},
		// Pointer
		{&baseBStruct, "guest", ActionRead, map[string]interface{}{}},
		{&baseBStruct, "user", ActionRead, map[string]interface{}{"Number": 10, "Label": "ABC", "Boolean": false}},
		{&baseBStruct, "sys", ActionRead, map[string]interface{}{}},
		{&baseBStruct, "admin", ActionRead, map[string]interface{}{"Number": 10, "Label": "ABC", "Boolean": false}},

		{&baseBStruct, "guest", ActionWrite, map[string]interface{}{}},
		{&baseBStruct, "user", ActionWrite, map[string]interface{}{}},
		{&baseBStruct, "sys", ActionWrite, map[string]interface{}{"Number": 10, "Label": "ABC", "Boolean": false}},
		{&baseBStruct, "admin", ActionWrite, map[string]interface{}{"Number": 10, "Label": "ABC", "Boolean": false}},
	}

	for _, table := range tables {
		actual := ExtractSingleObjectFields(table.object, table.userType, table.action)
		if !reflect.DeepEqual(actual, table.expected) {
			t.Errorf("%s (object = %+v, userType = %s, action = %d) was incorrect, got: %+v, want: %+v.",
				t.Name(), table.object, table.userType, table.action, actual, table.expected)
		}
	}
}

func testExtractSingleObjectFieldsStructField(t *testing.T) {
	t.Parallel()

	baseAStruct := AStruct{Number: 10, Text: "ABC"}
	baseCStruct := CStruct{Struct: baseAStruct, Pointer: &baseAStruct, Interface: baseAStruct}

	tables := []struct {
		object   interface{}
		userType string
		action   uint
		expected interface{}
	}{
		// Struct
		{baseCStruct, "guest", ActionRead, map[string]interface{}{}},
		{baseCStruct, "user", ActionRead, map[string]interface{}{
			"Struct":    map[string]interface{}{"Number": 10, "Label": "ABC"},
			"Pointer":   map[string]interface{}{"Number": 10, "Label": "ABC"},
			"Interface": map[string]interface{}{"Number": 10, "Label": "ABC"}}},
		{baseCStruct, "sys", ActionRead, map[string]interface{}{}},
		{baseCStruct, "admin", ActionRead, map[string]interface{}{
			"Struct":    map[string]interface{}{"Number": 10, "Label": "ABC"},
			"Pointer":   map[string]interface{}{"Number": 10, "Label": "ABC"},
			"Interface": map[string]interface{}{"Number": 10, "Label": "ABC"}}},

		{baseCStruct, "guest", ActionWrite, map[string]interface{}{}},
		{baseCStruct, "user", ActionWrite, map[string]interface{}{}},
		{baseCStruct, "sys", ActionWrite, map[string]interface{}{
			"Struct":    map[string]interface{}{"Number": 10, "Label": "ABC"},
			"Pointer":   map[string]interface{}{"Number": 10, "Label": "ABC"},
			"Interface": map[string]interface{}{"Number": 10, "Label": "ABC"}}},
		{baseCStruct, "admin", ActionWrite, map[string]interface{}{
			"Struct":    map[string]interface{}{"Number": 10, "Label": "ABC"},
			"Pointer":   map[string]interface{}{"Number": 10, "Label": "ABC"},
			"Interface": map[string]interface{}{"Number": 10, "Label": "ABC"}}},

		// Pointer
		{&baseCStruct, "guest", ActionRead, map[string]interface{}{}},
		{&baseCStruct, "user", ActionRead, map[string]interface{}{
			"Struct":    map[string]interface{}{"Number": 10, "Label": "ABC"},
			"Pointer":   map[string]interface{}{"Number": 10, "Label": "ABC"},
			"Interface": map[string]interface{}{"Number": 10, "Label": "ABC"}}},
		{&baseCStruct, "sys", ActionRead, map[string]interface{}{}},
		{&baseCStruct, "admin", ActionRead, map[string]interface{}{
			"Struct":    map[string]interface{}{"Number": 10, "Label": "ABC"},
			"Pointer":   map[string]interface{}{"Number": 10, "Label": "ABC"},
			"Interface": map[string]interface{}{"Number": 10, "Label": "ABC"}}},

		{&baseCStruct, "guest", ActionWrite, map[string]interface{}{}},
		{&baseCStruct, "user", ActionWrite, map[string]interface{}{}},
		{&baseCStruct, "sys", ActionWrite, map[string]interface{}{
			"Struct":    map[string]interface{}{"Number": 10, "Label": "ABC"},
			"Pointer":   map[string]interface{}{"Number": 10, "Label": "ABC"},
			"Interface": map[string]interface{}{"Number": 10, "Label": "ABC"}}},
		{&baseCStruct, "admin", ActionWrite, map[string]interface{}{
			"Struct":    map[string]interface{}{"Number": 10, "Label": "ABC"},
			"Pointer":   map[string]interface{}{"Number": 10, "Label": "ABC"},
			"Interface": map[string]interface{}{"Number": 10, "Label": "ABC"}}},
	}

	for _, table := range tables {
		actual := ExtractSingleObjectFields(table.object, table.userType, table.action)
		if !reflect.DeepEqual(actual, table.expected) {
			t.Errorf("%s (object = %+v, userType = %s, action = %d) was incorrect, got: %+v, want: %+v.",
				t.Name(), table.object, table.userType, table.action, actual, table.expected)
		}
	}
}

func testExtractSingleObjectFieldsNil(t *testing.T) {
	t.Parallel()
	baseAStruct := AStruct{Number: 0}
	var nilPointer *AStruct

	tables := []struct {
		object   interface{}
		userType string
		action   uint
		expected interface{}
	}{
		// Struct
		{baseAStruct, "guest", ActionRead, map[string]interface{}{}},
		{baseAStruct, "user", ActionRead, map[string]interface{}{"Number": 0, "Label": ""}},
		{baseAStruct, "sys", ActionRead, map[string]interface{}{}},
		{baseAStruct, "admin", ActionRead, map[string]interface{}{"Number": 0, "Label": ""}},

		{baseAStruct, "guest", ActionWrite, map[string]interface{}{}},
		{baseAStruct, "user", ActionWrite, map[string]interface{}{}},
		{baseAStruct, "sys", ActionWrite, map[string]interface{}{"Number": 0, "Label": ""}},
		{baseAStruct, "admin", ActionWrite, map[string]interface{}{"Number": 0, "Label": ""}},

		// Pointer
		{&baseAStruct, "guest", ActionRead, map[string]interface{}{}},
		{&baseAStruct, "user", ActionRead, map[string]interface{}{"Number": 0, "Label": ""}},
		{&baseAStruct, "sys", ActionRead, map[string]interface{}{}},
		{&baseAStruct, "admin", ActionRead, map[string]interface{}{"Number": 0, "Label": ""}},

		{&baseAStruct, "guest", ActionWrite, map[string]interface{}{}},
		{&baseAStruct, "user", ActionWrite, map[string]interface{}{}},
		{&baseAStruct, "sys", ActionWrite, map[string]interface{}{"Number": 0, "Label": ""}},
		{&baseAStruct, "admin", ActionWrite, map[string]interface{}{"Number": 0, "Label": ""}},

		// Nil pointer
		{nilPointer, "guest", ActionRead, nil},
		{nilPointer, "user", ActionRead, nil},
		{nilPointer, "sys", ActionRead, nil},
		{nilPointer, "admin", ActionRead, nil},

		{nilPointer, "guest", ActionWrite, nil},
		{nilPointer, "user", ActionWrite, nil},
		{nilPointer, "sys", ActionWrite, nil},
		{nilPointer, "admin", ActionWrite, nil},
	}

	for _, table := range tables {
		actual := ExtractSingleObjectFields(table.object, table.userType, table.action)
		if !reflect.DeepEqual(actual, table.expected) {
			t.Errorf("%s (object = %+v, userType = %s, action = %d) was incorrect, got: %+v, want: %+v.",
				t.Name(), table.object, table.userType, table.action, actual, table.expected)
		}
	}
}

func testExtractSingleObjectFieldsSpecial(t *testing.T) {
	t.Parallel()
	startTime := time.Now()
	stopTime := time.Now().Add(1000)
	baseEStruct := EStruct{Start: startTime, Stop: &stopTime, Number: sql.NullInt64{Int64: 10, Valid: true}}

	tables := []struct {
		object   interface{}
		userType string
		action   uint
		expected interface{}
	}{
		// Struct
		{baseEStruct, "guest", ActionRead, map[string]interface{}{}},
		{baseEStruct, "user", ActionRead, map[string]interface{}{"Start": startTime.String(), "Stop": stopTime.String(), "Number": int64(10)}},
		{baseEStruct, "sys", ActionRead, map[string]interface{}{}},
		{baseEStruct, "admin", ActionRead, map[string]interface{}{"Start": startTime.String(), "Stop": stopTime.String(), "Number": int64(10)}},

		{baseEStruct, "guest", ActionWrite, map[string]interface{}{}},
		{baseEStruct, "user", ActionWrite, map[string]interface{}{}},
		{baseEStruct, "sys", ActionWrite, map[string]interface{}{"Start": startTime.String(), "Stop": stopTime.String(), "Number": int64(10)}},
		{baseEStruct, "admin", ActionWrite, map[string]interface{}{"Start": startTime.String(), "Stop": stopTime.String(), "Number": int64(10)}},

		// Pointer
		{&baseEStruct, "guest", ActionRead, map[string]interface{}{}},
		{&baseEStruct, "user", ActionRead, map[string]interface{}{"Start": startTime.String(), "Stop": stopTime.String(), "Number": int64(10)}},
		{&baseEStruct, "sys", ActionRead, map[string]interface{}{}},
		{&baseEStruct, "admin", ActionRead, map[string]interface{}{"Start": startTime.String(), "Stop": stopTime.String(), "Number": int64(10)}},

		{&baseEStruct, "guest", ActionWrite, map[string]interface{}{}},
		{&baseEStruct, "user", ActionWrite, map[string]interface{}{}},
		{&baseEStruct, "sys", ActionWrite, map[string]interface{}{"Start": startTime.String(), "Stop": stopTime.String(), "Number": int64(10)}},
		{&baseEStruct, "admin", ActionWrite, map[string]interface{}{"Start": startTime.String(), "Stop": stopTime.String(), "Number": int64(10)}},
	}

	for _, table := range tables {
		actual := ExtractSingleObjectFields(table.object, table.userType, table.action)
		if !reflect.DeepEqual(actual, table.expected) {
			t.Errorf("%s (object = %+v, userType = %s, action = %d) was incorrect, got: %+v, want: %+v.",
				t.Name(), table.object, table.userType, table.action, actual, table.expected)
		}
	}
}

func TestExtractMultipleObjectsFields(t *testing.T) {
	t.Run("TestExtractMultipleObjectsFieldsNonSliceArray", testExtractMultipleObjectsFieldsNonSliceArray)
	t.Run("TestExtractMultipleObjectsFieldsBuiltin", testExtractMultipleObjectsFieldsBuiltin)
	t.Run("TestExtractMultipleObjectsFieldsStruct", testExtractMultipleObjectsFieldsStruct)
	t.Run("TestExtractMapObjectFields", testExtractMapObjectFields)
}

func testExtractMultipleObjectsFieldsNonSliceArray(t *testing.T) {
	t.Parallel()

	baseAStruct := AStruct{Number: 10, Text: "ABC"}

	tables := []struct {
		object   interface{}
		userType string
		action   uint
		expected interface{}
	}{
		// Struct
		{baseAStruct, "guest", ActionRead, baseAStruct},
		{baseAStruct, "user", ActionRead, baseAStruct},
		{baseAStruct, "sys", ActionRead, baseAStruct},
		{baseAStruct, "admin", ActionRead, baseAStruct},

		{baseAStruct, "guest", ActionWrite, baseAStruct},
		{baseAStruct, "user", ActionWrite, baseAStruct},
		{baseAStruct, "sys", ActionWrite, baseAStruct},
		{baseAStruct, "admin", ActionWrite, baseAStruct},

		// Pointer
		{&baseAStruct, "guest", ActionRead, baseAStruct},
		{&baseAStruct, "user", ActionRead, baseAStruct},
		{&baseAStruct, "sys", ActionRead, baseAStruct},
		{&baseAStruct, "admin", ActionRead, baseAStruct},

		{&baseAStruct, "guest", ActionWrite, baseAStruct},
		{&baseAStruct, "user", ActionWrite, baseAStruct},
		{&baseAStruct, "sys", ActionWrite, baseAStruct},
		{&baseAStruct, "admin", ActionWrite, baseAStruct},
	}

	for _, table := range tables {
		actual := ExtractMultipleObjectsFields(table.object, table.userType, table.action)
		if !reflect.DeepEqual(actual, table.expected) {
			t.Errorf("%s (object = %+v, userType = %s, action = %d) was incorrect, got: %+v, want: %+v.",
				t.Name(), table.object, table.userType, table.action, actual, table.expected)
		}
	}
}

func testExtractMultipleObjectsFieldsBuiltin(t *testing.T) {
	t.Parallel()

	baseArray := [3]int{1, 2, 3}
	baseSlice := []float32{1, 2, 3}

	tables := []struct {
		object   interface{}
		userType string
		action   uint
		expected interface{}
	}{
		// Struct
		{baseArray, "guest", ActionRead, []interface{}{1, 2, 3}},
		{baseArray, "user", ActionRead, []interface{}{1, 2, 3}},
		{baseArray, "sys", ActionRead, []interface{}{1, 2, 3}},
		{baseArray, "admin", ActionRead, []interface{}{1, 2, 3}},

		{baseArray, "guest", ActionWrite, []interface{}{1, 2, 3}},
		{baseArray, "user", ActionWrite, []interface{}{1, 2, 3}},
		{baseArray, "sys", ActionWrite, []interface{}{1, 2, 3}},
		{baseArray, "admin", ActionWrite, []interface{}{1, 2, 3}},

		{baseSlice, "guest", ActionRead, []interface{}{float32(1), float32(2), float32(3)}},
		{baseSlice, "user", ActionRead, []interface{}{float32(1), float32(2), float32(3)}},
		{baseSlice, "sys", ActionRead, []interface{}{float32(1), float32(2), float32(3)}},
		{baseSlice, "admin", ActionRead, []interface{}{float32(1), float32(2), float32(3)}},

		{baseSlice, "guest", ActionWrite, []interface{}{float32(1), float32(2), float32(3)}},
		{baseSlice, "user", ActionWrite, []interface{}{float32(1), float32(2), float32(3)}},
		{baseSlice, "sys", ActionWrite, []interface{}{float32(1), float32(2), float32(3)}},
		{baseSlice, "admin", ActionWrite, []interface{}{float32(1), float32(2), float32(3)}},
		// Pointer
		{&baseArray, "guest", ActionRead, []interface{}{1, 2, 3}},
		{&baseArray, "user", ActionRead, []interface{}{1, 2, 3}},
		{&baseArray, "sys", ActionRead, []interface{}{1, 2, 3}},
		{&baseArray, "admin", ActionRead, []interface{}{1, 2, 3}},

		{&baseArray, "guest", ActionWrite, []interface{}{1, 2, 3}},
		{&baseArray, "user", ActionWrite, []interface{}{1, 2, 3}},
		{&baseArray, "sys", ActionWrite, []interface{}{1, 2, 3}},
		{&baseArray, "admin", ActionWrite, []interface{}{1, 2, 3}},

		{&baseSlice, "guest", ActionRead, []interface{}{float32(1), float32(2), float32(3)}},
		{&baseSlice, "user", ActionRead, []interface{}{float32(1), float32(2), float32(3)}},
		{&baseSlice, "sys", ActionRead, []interface{}{float32(1), float32(2), float32(3)}},
		{&baseSlice, "admin", ActionRead, []interface{}{float32(1), float32(2), float32(3)}},

		{&baseSlice, "guest", ActionWrite, []interface{}{float32(1), float32(2), float32(3)}},
		{&baseSlice, "user", ActionWrite, []interface{}{float32(1), float32(2), float32(3)}},
		{&baseSlice, "sys", ActionWrite, []interface{}{float32(1), float32(2), float32(3)}},
		{&baseSlice, "admin", ActionWrite, []interface{}{float32(1), float32(2), float32(3)}},
	}

	for _, table := range tables {
		actual := ExtractMultipleObjectsFields(table.object, table.userType, table.action)
		if !reflect.DeepEqual(actual, table.expected) {
			t.Errorf("%s (object = %+v, userType = %s, action = %d) was incorrect, got: %+v, want: %+v.",
				t.Name(), table.object, table.userType, table.action, actual, table.expected)
		}
	}
}

func testExtractMultipleObjectsFieldsStruct(t *testing.T) {
	t.Parallel()

	baseHStruct := HStruct{Boolean: sql.NullBool{Bool: true, Valid: true}}
	baseSlice := []HStruct{baseHStruct, baseHStruct}

	tables := []struct {
		object   interface{}
		userType string
		action   uint
		expected interface{}
	}{
		// Struct
		{baseSlice, "guest", ActionRead, []interface{}{
			map[string]interface{}{},
			map[string]interface{}{},
		}},
		{baseSlice, "user", ActionRead, []interface{}{
			map[string]interface{}{"Boolean": true},
			map[string]interface{}{"Boolean": true},
		}},
		{baseSlice, "sys", ActionRead, []interface{}{
			map[string]interface{}{},
			map[string]interface{}{},
		}},
		{baseSlice, "admin", ActionRead, []interface{}{
			map[string]interface{}{"Boolean": true},
			map[string]interface{}{"Boolean": true},
		}},

		{baseSlice, "guest", ActionWrite, []interface{}{
			map[string]interface{}{},
			map[string]interface{}{},
		}},
		{baseSlice, "user", ActionWrite, []interface{}{
			map[string]interface{}{},
			map[string]interface{}{},
		}},
		{baseSlice, "sys", ActionWrite, []interface{}{
			map[string]interface{}{"Boolean": true},
			map[string]interface{}{"Boolean": true},
		}},
		{baseSlice, "admin", ActionWrite, []interface{}{
			map[string]interface{}{"Boolean": true},
			map[string]interface{}{"Boolean": true},
		}},

		// Pointer
		{&baseSlice, "guest", ActionRead, []interface{}{
			map[string]interface{}{},
			map[string]interface{}{},
		}},
		{&baseSlice, "user", ActionRead, []interface{}{
			map[string]interface{}{"Boolean": true},
			map[string]interface{}{"Boolean": true},
		}},
		{&baseSlice, "sys", ActionRead, []interface{}{
			map[string]interface{}{},
			map[string]interface{}{},
		}},
		{&baseSlice, "admin", ActionRead, []interface{}{
			map[string]interface{}{"Boolean": true},
			map[string]interface{}{"Boolean": true},
		}},

		{&baseSlice, "guest", ActionWrite, []interface{}{
			map[string]interface{}{},
			map[string]interface{}{},
		}},
		{&baseSlice, "user", ActionWrite, []interface{}{
			map[string]interface{}{},
			map[string]interface{}{},
		}},
		{&baseSlice, "sys", ActionWrite, []interface{}{
			map[string]interface{}{"Boolean": true},
			map[string]interface{}{"Boolean": true},
		}},
		{&baseSlice, "admin", ActionWrite, []interface{}{
			map[string]interface{}{"Boolean": true},
			map[string]interface{}{"Boolean": true},
		}},
	}

	for _, table := range tables {
		actual := ExtractMultipleObjectsFields(table.object, table.userType, table.action)
		if !reflect.DeepEqual(actual, table.expected) {
			t.Errorf("%s (object = %+v, userType = %s, action = %d) was incorrect, got: %+v, want: %+v.",
				t.Name(), table.object, table.userType, table.action, actual, table.expected)
		}
	}
}

func testExtractMapObjectFields(t *testing.T) {
	baseAStruct := AStruct{Number: 10, Text: "ABC"}
	baseMap := map[string]interface{}{"foo": baseAStruct, "bar": &baseAStruct}

	tables := []struct {
		object   interface{}
		userType string
		action   uint
		expected interface{}
	}{
		// Struct
		{baseMap, "guest", ActionRead, map[interface{}]interface{}{
			"foo": map[string]interface{}{},
			"bar": map[string]interface{}{},
		}},
		{baseMap, "user", ActionRead, map[interface{}]interface{}{
			"foo": map[string]interface{}{"Number": 10, "Label": "ABC"},
			"bar": map[string]interface{}{"Number": 10, "Label": "ABC"},
		}},
		{baseMap, "sys", ActionRead, map[interface{}]interface{}{
			"foo": map[string]interface{}{},
			"bar": map[string]interface{}{},
		}},
		{baseMap, "admin", ActionRead, map[interface{}]interface{}{
			"foo": map[string]interface{}{"Number": 10, "Label": "ABC"},
			"bar": map[string]interface{}{"Number": 10, "Label": "ABC"},
		}},

		{baseMap, "guest", ActionWrite, map[interface{}]interface{}{
			"foo": map[string]interface{}{},
			"bar": map[string]interface{}{},
		}},
		{baseMap, "user", ActionWrite, map[interface{}]interface{}{
			"foo": map[string]interface{}{},
			"bar": map[string]interface{}{},
		}},
		{baseMap, "sys", ActionWrite, map[interface{}]interface{}{
			"foo": map[string]interface{}{"Number": 10, "Label": "ABC"},
			"bar": map[string]interface{}{"Number": 10, "Label": "ABC"},
		}},
		{baseMap, "admin", ActionWrite, map[interface{}]interface{}{
			"foo": map[string]interface{}{"Number": 10, "Label": "ABC"},
			"bar": map[string]interface{}{"Number": 10, "Label": "ABC"},
		}},

		// Pointer
		{&baseMap, "guest", ActionRead, map[interface{}]interface{}{
			"foo": map[string]interface{}{},
			"bar": map[string]interface{}{},
		}},
		{&baseMap, "user", ActionRead, map[interface{}]interface{}{
			"foo": map[string]interface{}{"Number": 10, "Label": "ABC"},
			"bar": map[string]interface{}{"Number": 10, "Label": "ABC"},
		}},
		{&baseMap, "sys", ActionRead, map[interface{}]interface{}{
			"foo": map[string]interface{}{},
			"bar": map[string]interface{}{},
		}},
		{&baseMap, "admin", ActionRead, map[interface{}]interface{}{
			"foo": map[string]interface{}{"Number": 10, "Label": "ABC"},
			"bar": map[string]interface{}{"Number": 10, "Label": "ABC"},
		}},

		{&baseMap, "guest", ActionWrite, map[interface{}]interface{}{
			"foo": map[string]interface{}{},
			"bar": map[string]interface{}{},
		}},
		{&baseMap, "user", ActionWrite, map[interface{}]interface{}{
			"foo": map[string]interface{}{},
			"bar": map[string]interface{}{},
		}},
		{&baseMap, "sys", ActionWrite, map[interface{}]interface{}{
			"foo": map[string]interface{}{"Number": 10, "Label": "ABC"},
			"bar": map[string]interface{}{"Number": 10, "Label": "ABC"},
		}},
		{&baseMap, "admin", ActionWrite, map[interface{}]interface{}{
			"foo": map[string]interface{}{"Number": 10, "Label": "ABC"},
			"bar": map[string]interface{}{"Number": 10, "Label": "ABC"},
		}},
	}

	for _, table := range tables {
		actual := ExtractMapObjectsFields(table.object, table.userType, table.action)
		if !reflect.DeepEqual(actual, table.expected) {
			t.Errorf("%s (object = %+v, userType = %s, action = %d) was incorrect, got: %+v, want: %+v.",
				t.Name(), table.object, table.userType, table.action, actual, table.expected)
		}
	}
}

func TestExtractObjectsFeatures(t *testing.T) {
	t.Run("TestExtractFieldsBuiltin", testExtractFieldsBuiltin)
	t.Run("TestExtractFieldsStruct", testExtractFieldsStruct)
	t.Run("TestExtractFieldsArraySlice", testExtractFieldsArraySlice)
	t.Run("TestExtractFieldsStructWithSliceArray", testExtractFieldsStructWithSliceArray)
}

func testExtractFieldsBuiltin(t *testing.T) {
	t.Parallel()

	baseValue := "example string"

	tables := []struct {
		object   interface{}
		userType string
		action   uint
		expected interface{}
	}{
		// Struct
		{baseValue, "guest", ActionRead, baseValue},
		{baseValue, "user", ActionRead, baseValue},
		{baseValue, "sys", ActionRead, baseValue},
		{baseValue, "admin", ActionRead, baseValue},

		{baseValue, "guest", ActionWrite, baseValue},
		{baseValue, "user", ActionWrite, baseValue},
		{baseValue, "sys", ActionWrite, baseValue},
		{baseValue, "admin", ActionWrite, baseValue},
		// Pointer
		{&baseValue, "guest", ActionRead, baseValue},
		{&baseValue, "user", ActionRead, baseValue},
		{&baseValue, "sys", ActionRead, baseValue},
		{&baseValue, "admin", ActionRead, baseValue},

		{&baseValue, "guest", ActionWrite, baseValue},
		{&baseValue, "user", ActionWrite, baseValue},
		{&baseValue, "sys", ActionWrite, baseValue},
		{&baseValue, "admin", ActionWrite, baseValue},
	}

	for _, table := range tables {
		actual := ExtractFields(table.object, table.userType, table.action)
		if !reflect.DeepEqual(actual, table.expected) {
			t.Errorf("%s (object = %+v, userType = %s, action = %d) was incorrect, got: %+v, want: %+v.",
				t.Name(), table.object, table.userType, table.action, actual, table.expected)
		}
	}
}

func testExtractFieldsStruct(t *testing.T) {
	t.Parallel()

	baseStruct := DStruct{}

	tables := []struct {
		object   interface{}
		userType string
		action   uint
		expected interface{}
	}{
		// Struct
		{baseStruct, "guest", ActionRead, map[string]interface{}{}},
		{baseStruct, "user", ActionRead, map[string]interface{}{}},
		{baseStruct, "sys", ActionRead, map[string]interface{}{}},
		{baseStruct, "admin", ActionRead, map[string]interface{}{}},

		{baseStruct, "guest", ActionWrite, map[string]interface{}{}},
		{baseStruct, "user", ActionWrite, map[string]interface{}{}},
		{baseStruct, "sys", ActionWrite, map[string]interface{}{}},
		{baseStruct, "admin", ActionWrite, map[string]interface{}{}},
		// Pointer
		{&baseStruct, "guest", ActionRead, map[string]interface{}{}},
		{&baseStruct, "user", ActionRead, map[string]interface{}{}},
		{&baseStruct, "sys", ActionRead, map[string]interface{}{}},
		{&baseStruct, "admin", ActionRead, map[string]interface{}{}},

		{&baseStruct, "guest", ActionWrite, map[string]interface{}{}},
		{&baseStruct, "user", ActionWrite, map[string]interface{}{}},
		{&baseStruct, "sys", ActionWrite, map[string]interface{}{}},
		{&baseStruct, "admin", ActionWrite, map[string]interface{}{}},
	}

	for _, table := range tables {
		actual := ExtractFields(table.object, table.userType, table.action)
		if !reflect.DeepEqual(actual, table.expected) {
			t.Errorf("%s (object = %+v, userType = %s, action = %d) was incorrect, got: %+v, want: %+v.",
				t.Name(), table.object, table.userType, table.action, actual, table.expected)
		}
	}
}

func testExtractFieldsArraySlice(t *testing.T) {
	t.Parallel()

	baseSlice := []int{1}
	baseArray := [1]bool{false}
	baseHStruct := HStruct{Boolean: sql.NullBool{Bool: true, Valid: true}}
	baseStructSlice := []HStruct{baseHStruct, baseHStruct}
	basePointerStructSlice := []*HStruct{&baseHStruct, &baseHStruct}

	tables := []struct {
		object   interface{}
		userType string
		action   uint
		expected interface{}
	}{
		// Struct
		{baseSlice, "guest", ActionRead, []interface{}{1}},
		{baseSlice, "user", ActionRead, []interface{}{1}},
		{baseSlice, "sys", ActionRead, []interface{}{1}},
		{baseSlice, "admin", ActionRead, []interface{}{1}},

		{baseSlice, "guest", ActionWrite, []interface{}{1}},
		{baseSlice, "user", ActionWrite, []interface{}{1}},
		{baseSlice, "sys", ActionWrite, []interface{}{1}},
		{baseSlice, "admin", ActionWrite, []interface{}{1}},

		{baseArray, "guest", ActionRead, []interface{}{false}},
		{baseArray, "user", ActionRead, []interface{}{false}},
		{baseArray, "sys", ActionRead, []interface{}{false}},
		{baseArray, "admin", ActionRead, []interface{}{false}},

		{baseArray, "guest", ActionWrite, []interface{}{false}},
		{baseArray, "user", ActionWrite, []interface{}{false}},
		{baseArray, "sys", ActionWrite, []interface{}{false}},
		{baseArray, "admin", ActionWrite, []interface{}{false}},

		{baseStructSlice, "guest", ActionRead, []interface{}{
			map[string]interface{}{},
			map[string]interface{}{},
		}},
		{baseStructSlice, "user", ActionRead, []interface{}{
			map[string]interface{}{"Boolean": true},
			map[string]interface{}{"Boolean": true},
		}},
		{baseStructSlice, "sys", ActionRead, []interface{}{
			map[string]interface{}{},
			map[string]interface{}{},
		}},
		{baseStructSlice, "admin", ActionRead, []interface{}{
			map[string]interface{}{"Boolean": true},
			map[string]interface{}{"Boolean": true},
		}},

		{baseStructSlice, "guest", ActionWrite, []interface{}{
			map[string]interface{}{},
			map[string]interface{}{},
		}},
		{baseStructSlice, "user", ActionWrite, []interface{}{
			map[string]interface{}{},
			map[string]interface{}{},
		}},
		{baseStructSlice, "sys", ActionWrite, []interface{}{
			map[string]interface{}{"Boolean": true},
			map[string]interface{}{"Boolean": true},
		}},
		{baseStructSlice, "admin", ActionWrite, []interface{}{
			map[string]interface{}{"Boolean": true},
			map[string]interface{}{"Boolean": true},
		}},

		{basePointerStructSlice, "guest", ActionRead, []interface{}{
			map[string]interface{}{},
			map[string]interface{}{},
		}},
		{basePointerStructSlice, "user", ActionRead, []interface{}{
			map[string]interface{}{"Boolean": true},
			map[string]interface{}{"Boolean": true},
		}},
		{basePointerStructSlice, "sys", ActionRead, []interface{}{
			map[string]interface{}{},
			map[string]interface{}{},
		}},
		{basePointerStructSlice, "admin", ActionRead, []interface{}{
			map[string]interface{}{"Boolean": true},
			map[string]interface{}{"Boolean": true},
		}},

		{basePointerStructSlice, "guest", ActionWrite, []interface{}{
			map[string]interface{}{},
			map[string]interface{}{},
		}},
		{basePointerStructSlice, "user", ActionWrite, []interface{}{
			map[string]interface{}{},
			map[string]interface{}{},
		}},
		{basePointerStructSlice, "sys", ActionWrite, []interface{}{
			map[string]interface{}{"Boolean": true},
			map[string]interface{}{"Boolean": true},
		}},
		{basePointerStructSlice, "admin", ActionWrite, []interface{}{
			map[string]interface{}{"Boolean": true},
			map[string]interface{}{"Boolean": true},
		}},
		// Pointer
		{&baseSlice, "guest", ActionRead, []interface{}{1}},
		{&baseSlice, "user", ActionRead, []interface{}{1}},
		{&baseSlice, "sys", ActionRead, []interface{}{1}},
		{&baseSlice, "admin", ActionRead, []interface{}{1}},

		{&baseSlice, "guest", ActionWrite, []interface{}{1}},
		{&baseSlice, "user", ActionWrite, []interface{}{1}},
		{&baseSlice, "sys", ActionWrite, []interface{}{1}},
		{&baseSlice, "admin", ActionWrite, []interface{}{1}},

		{&baseArray, "guest", ActionRead, []interface{}{false}},
		{&baseArray, "user", ActionRead, []interface{}{false}},
		{&baseArray, "sys", ActionRead, []interface{}{false}},
		{&baseArray, "admin", ActionRead, []interface{}{false}},

		{&baseArray, "guest", ActionWrite, []interface{}{false}},
		{&baseArray, "user", ActionWrite, []interface{}{false}},
		{&baseArray, "sys", ActionWrite, []interface{}{false}},
		{&baseArray, "admin", ActionWrite, []interface{}{false}},

		{&baseStructSlice, "guest", ActionRead, []interface{}{
			map[string]interface{}{},
			map[string]interface{}{},
		}},
		{&baseStructSlice, "user", ActionRead, []interface{}{
			map[string]interface{}{"Boolean": true},
			map[string]interface{}{"Boolean": true},
		}},
		{&baseStructSlice, "sys", ActionRead, []interface{}{
			map[string]interface{}{},
			map[string]interface{}{},
		}},
		{&baseStructSlice, "admin", ActionRead, []interface{}{
			map[string]interface{}{"Boolean": true},
			map[string]interface{}{"Boolean": true},
		}},

		{&baseStructSlice, "guest", ActionWrite, []interface{}{
			map[string]interface{}{},
			map[string]interface{}{},
		}},
		{&baseStructSlice, "user", ActionWrite, []interface{}{
			map[string]interface{}{},
			map[string]interface{}{},
		}},
		{&baseStructSlice, "sys", ActionWrite, []interface{}{
			map[string]interface{}{"Boolean": true},
			map[string]interface{}{"Boolean": true},
		}},
		{&baseStructSlice, "admin", ActionWrite, []interface{}{
			map[string]interface{}{"Boolean": true},
			map[string]interface{}{"Boolean": true},
		}},

		{&basePointerStructSlice, "guest", ActionRead, []interface{}{
			map[string]interface{}{},
			map[string]interface{}{},
		}},
		{&basePointerStructSlice, "user", ActionRead, []interface{}{
			map[string]interface{}{"Boolean": true},
			map[string]interface{}{"Boolean": true},
		}},
		{&basePointerStructSlice, "sys", ActionRead, []interface{}{
			map[string]interface{}{},
			map[string]interface{}{},
		}},
		{&basePointerStructSlice, "admin", ActionRead, []interface{}{
			map[string]interface{}{"Boolean": true},
			map[string]interface{}{"Boolean": true},
		}},

		{&basePointerStructSlice, "guest", ActionWrite, []interface{}{
			map[string]interface{}{},
			map[string]interface{}{},
		}},
		{&basePointerStructSlice, "user", ActionWrite, []interface{}{
			map[string]interface{}{},
			map[string]interface{}{},
		}},
		{&basePointerStructSlice, "sys", ActionWrite, []interface{}{
			map[string]interface{}{"Boolean": true},
			map[string]interface{}{"Boolean": true},
		}},
		{&basePointerStructSlice, "admin", ActionWrite, []interface{}{
			map[string]interface{}{"Boolean": true},
			map[string]interface{}{"Boolean": true},
		}},
	}

	for _, table := range tables {
		actual := ExtractFields(table.object, table.userType, table.action)
		if !reflect.DeepEqual(actual, table.expected) {
			t.Errorf("%s (object = %+v, userType = %s, action = %d) was incorrect, got: %+v, want: %+v.",
				t.Name(), table.object, table.userType, table.action, actual, table.expected)
		}
	}
}

func testExtractFieldsStructWithSliceArray(t *testing.T) {
	t.Parallel()

	baseAStruct := AStruct{Number: 10, Text: "DEF"}
	baseStruct := FStruct{Name: "ABC", Array: [2]int{1, 2}, Slice: []AStruct{baseAStruct, baseAStruct}}

	tables := []struct {
		object   interface{}
		userType string
		action   uint
		expected interface{}
	}{
		// Struct
		{baseStruct, "guest", ActionRead, map[string]interface{}{}},
		{baseStruct, "user", ActionRead, map[string]interface{}{
			"Name":  "ABC",
			"Array": []interface{}{1, 2},
			"Slice": []interface{}{
				map[string]interface{}{"Number": 10, "Label": "DEF"},
				map[string]interface{}{"Number": 10, "Label": "DEF"},
			}}},
		{baseStruct, "sys", ActionRead, map[string]interface{}{}},
		{baseStruct, "admin", ActionRead, map[string]interface{}{
			"Name":  "ABC",
			"Array": []interface{}{1, 2},
			"Slice": []interface{}{
				map[string]interface{}{"Number": 10, "Label": "DEF"},
				map[string]interface{}{"Number": 10, "Label": "DEF"},
			}}},
		{baseStruct, "guest", ActionWrite, map[string]interface{}{}},
		{baseStruct, "user", ActionWrite, map[string]interface{}{}},
		{baseStruct, "sys", ActionWrite, map[string]interface{}{
			"Name":  "ABC",
			"Array": []interface{}{1, 2},
			"Slice": []interface{}{
				map[string]interface{}{"Number": 10, "Label": "DEF"},
				map[string]interface{}{"Number": 10, "Label": "DEF"},
			}}},
		{baseStruct, "admin", ActionWrite, map[string]interface{}{
			"Name":  "ABC",
			"Array": []interface{}{1, 2},
			"Slice": []interface{}{
				map[string]interface{}{"Number": 10, "Label": "DEF"},
				map[string]interface{}{"Number": 10, "Label": "DEF"},
			}}},

		// Pointer
		{&baseStruct, "guest", ActionRead, map[string]interface{}{}},
		{&baseStruct, "user", ActionRead, map[string]interface{}{
			"Name":  "ABC",
			"Array": []interface{}{1, 2},
			"Slice": []interface{}{
				map[string]interface{}{"Number": 10, "Label": "DEF"},
				map[string]interface{}{"Number": 10, "Label": "DEF"},
			}}},
		{&baseStruct, "sys", ActionRead, map[string]interface{}{}},
		{&baseStruct, "admin", ActionRead, map[string]interface{}{
			"Name":  "ABC",
			"Array": []interface{}{1, 2},
			"Slice": []interface{}{
				map[string]interface{}{"Number": 10, "Label": "DEF"},
				map[string]interface{}{"Number": 10, "Label": "DEF"},
			}}},

		{&baseStruct, "guest", ActionWrite, map[string]interface{}{}},
		{&baseStruct, "user", ActionWrite, map[string]interface{}{}},
		{&baseStruct, "sys", ActionWrite, map[string]interface{}{
			"Name":  "ABC",
			"Array": []interface{}{1, 2},
			"Slice": []interface{}{
				map[string]interface{}{"Number": 10, "Label": "DEF"},
				map[string]interface{}{"Number": 10, "Label": "DEF"},
			}}},
		{&baseStruct, "admin", ActionWrite, map[string]interface{}{
			"Name":  "ABC",
			"Array": []interface{}{1, 2},
			"Slice": []interface{}{
				map[string]interface{}{"Number": 10, "Label": "DEF"},
				map[string]interface{}{"Number": 10, "Label": "DEF"},
			}}},
	}

	for _, table := range tables {
		actual := ExtractFields(table.object, table.userType, table.action)
		if !reflect.DeepEqual(actual, table.expected) {
			t.Errorf("%s (object = %+v, userType = %s, action = %d) was incorrect, got: %+v, want: %+v.",
				t.Name(), table.object, table.userType, table.action, actual, table.expected)
		}
	}
}

func TestCleanObject(t *testing.T) {
	t.Run("TestCleanObjectStruct", testCleanObjectStruct)
	t.Run("TestCleanObjectSlice", testCleanObjectSlice)
}

func testCleanObjectStruct(t *testing.T) {
	t.Parallel()
	baseStruct := GStruct{Name: "ABC", Version: 1}

	tables := []struct {
		object   interface{}
		userType string
		action   uint
		expected interface{}
	}{
		// Struct
		{baseStruct, "guest", ActionRead, &GStruct{Version: 1}},
		{baseStruct, "user", ActionRead, &GStruct{Name: "ABC"}},
		{baseStruct, "sys", ActionRead, &GStruct{Version: 1}},
		{baseStruct, "admin", ActionRead, &GStruct{Name: "ABC"}},

		{baseStruct, "guest", ActionWrite, &GStruct{Version: 1}},
		{baseStruct, "user", ActionWrite, &GStruct{Version: 1}},
		{baseStruct, "sys", ActionWrite, &GStruct{Name: "ABC"}},
		{baseStruct, "admin", ActionWrite, &GStruct{Name: "ABC"}},

		// Pointer
		{&baseStruct, "guest", ActionRead, &GStruct{Version: 1}},
		{&baseStruct, "user", ActionRead, &GStruct{Name: "ABC"}},
		{&baseStruct, "sys", ActionRead, &GStruct{Version: 1}},
		{&baseStruct, "admin", ActionRead, &GStruct{Name: "ABC"}},

		{&baseStruct, "guest", ActionWrite, &GStruct{Version: 1}},
		{&baseStruct, "user", ActionWrite, &GStruct{Version: 1}},
		{&baseStruct, "sys", ActionWrite, &GStruct{Name: "ABC"}},
		{&baseStruct, "admin", ActionWrite, &GStruct{Name: "ABC"}},
	}

	for _, table := range tables {
		actual := CleanObject(table.object, table.userType, table.action)
		if !reflect.DeepEqual(actual, table.expected) {
			t.Errorf("%s (object = %+v, userType = %s, action = %d) was incorrect, got: %+v, want: %+v.",
				t.Name(), table.object, table.userType, table.action, actual, table.expected)
		}
	}
}

func testCleanObjectSlice(t *testing.T) {
	t.Parallel()
	baseStruct := GStruct{Name: "ABC", Version: 1}
	baseSlice := []GStruct{baseStruct, baseStruct}

	tables := []struct {
		object   interface{}
		userType string
		action   uint
		expected interface{}
	}{
		// Struct
		{baseSlice, "guest", ActionRead, &[]GStruct{{Version: 1}, {Version: 1}}},
		{baseSlice, "user", ActionRead, &[]GStruct{{Name: "ABC"}, {Name: "ABC"}}},
		{baseSlice, "sys", ActionRead, &[]GStruct{{Version: 1}, {Version: 1}}},
		{baseSlice, "admin", ActionRead, &[]GStruct{{Name: "ABC"}, {Name: "ABC"}}},

		{baseSlice, "guest", ActionWrite, &[]GStruct{{Version: 1}, {Version: 1}}},
		{baseSlice, "user", ActionWrite, &[]GStruct{{Version: 1}, {Version: 1}}},
		{baseSlice, "sys", ActionWrite, &[]GStruct{{Name: "ABC"}, {Name: "ABC"}}},
		{baseSlice, "admin", ActionWrite, &[]GStruct{{Name: "ABC"}, {Name: "ABC"}}},

		// Pointer
		{&baseSlice, "guest", ActionRead, &[]GStruct{{Version: 1}, {Version: 1}}},
		{&baseSlice, "user", ActionRead, &[]GStruct{{Name: "ABC"}, {Name: "ABC"}}},
		{&baseSlice, "sys", ActionRead, &[]GStruct{{Version: 1}, {Version: 1}}},
		{&baseSlice, "admin", ActionRead, &[]GStruct{{Name: "ABC"}, {Name: "ABC"}}},

		{&baseSlice, "guest", ActionWrite, &[]GStruct{{Version: 1}, {Version: 1}}},
		{&baseSlice, "user", ActionWrite, &[]GStruct{{Version: 1}, {Version: 1}}},
		{&baseSlice, "sys", ActionWrite, &[]GStruct{{Name: "ABC"}, {Name: "ABC"}}},
		{&baseSlice, "admin", ActionWrite, &[]GStruct{{Name: "ABC"}, {Name: "ABC"}}},
	}

	for _, table := range tables {
		actual := CleanObject(table.object, table.userType, table.action)
		if !reflect.DeepEqual(actual, table.expected) {
			t.Errorf("%s (object = %+v, userType = %s, action = %d) was incorrect, got: %+v, want: %+v.",
				t.Name(), table.object, table.userType, table.action, actual, table.expected)
		}
	}
}
