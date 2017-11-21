package gopex

import (
	"database/sql"
	"encoding/json"
	"reflect"
	"strings"
	"time"
)

// CleanObject is a function that receives an object, cleans it by removing the values that the user has not
// access for that action and returns a pointer to the cleaned object
func CleanObject(object interface{}, userType uint, action uint) interface{} {
	extractedFields := ExtractFields(object, userType, action)

	// Get the reflect value
	reflectValue := getReflectValue(object)
	if reflectValue == nil {
		return nil
	}

	// Create pointer to new object
	reflectType := reflect.TypeOf(reflectValue.Interface())
	result := reflect.New(reflectType).Interface()

	// Marshal
	data, err := json.Marshal(extractedFields)
	if err != nil {
		return nil
	}

	// Unmarshal
	err = json.Unmarshal(data, &result)
	if err != nil {
		return nil
	}

	return result
}

// ExtractFields extracts all the fields that a given user have access and
// returns a JSON interface of that object either its an array, slice or struct.
// It uses the json tag to get the field name, it it is not defined uses the field
// name of the struct.
func ExtractFields(object interface{}, userType uint, action uint) interface{} {
	reflectValue := getReflectValue(object)
	if reflectValue == nil {
		return nil
	}

	switch reflectValue.Kind() {
	case reflect.Struct:
		return ExtractSingleObjectFields(object, userType, action)
	case reflect.Slice, reflect.Array:
		return ExtractMultipleObjectsFields(object, userType, action)
	case reflect.Map:
		return ExtractMapObjectsFields(object, userType, action)
	default:
		return reflectValue.Interface()
	}
}

// ExtractSingleObjectFields extracts all the fields that a given user have access and
// returns a JSON interface of that object.
// It uses the json tag to get the field name, it it is not defined uses the field
// name of the struct.
func ExtractSingleObjectFields(object interface{}, userType uint, action uint) interface{} {
	reflectValue := getReflectValue(object)
	if reflectValue == nil {
		return nil
	}

	// If not struct or a special object just return the object
	if reflectValue.Kind() != reflect.Struct ||
		isSpecialObject(reflectValue.Interface()) {
		return reflectValue.Interface()
	}

	// Iterate through all the fields
	reflectType := reflect.TypeOf(reflectValue.Interface())
	resultObject := map[string]interface{}{}
	for i := 0; i < reflectValue.NumField(); i++ {
		resultField := extractField(reflectType.Field(i), reflectValue.Field(i), userType, action)
		for key, value := range resultField {
			resultObject[key] = value
		}
	}

	return resultObject
}

// ExtractMultipleObjectsFields extracts all the fields that a given user have access and
// returns a JSON interface of an array of objects.
// It uses the json tag to get the field name of each of the objects,
// it it is not defined uses the field name of the struct.
func ExtractMultipleObjectsFields(object interface{}, userType uint, action uint) interface{} {
	// Get the reflect value
	reflectValue := getReflectValue(object)
	if reflectValue == nil {
		return nil
	}

	// If not slice or array just return the object
	if reflectValue.Kind() != reflect.Slice &&
		reflectValue.Kind() != reflect.Array {
		return reflectValue.Interface()
	}
	// Multiple objects of builtin types, then no need to iterate
	if reflect.TypeOf(reflectValue.Interface()).Elem().Kind() != reflect.Struct {
		return reflectValue.Interface()
	}

	// Iterate through each single object in the slice
	resultObjects := make([]interface{}, reflectValue.Len())
	for i := 0; i < reflectValue.Len(); i++ {
		resultObjects[i] = ExtractFields(reflectValue.Index(i).Interface(), userType, action)
	}
	return resultObjects
}

// ExtractMapObjectsFields extracts all the fields that a given user have access and
// returns a JSON interface of an array of objects.
// It uses the json tag to get the field name of each of the objects,
// it it is not defined uses the field name of the struct.
func ExtractMapObjectsFields(object interface{}, userType uint, action uint) interface{} {
	// Get the reflect value
	reflectValue := getReflectValue(object)
	if reflectValue == nil {
		return nil
	}

	// If not slice or array just return the object
	if reflectValue.Kind() != reflect.Map {
		return reflectValue.Interface()
	}

	// Iterate through each single object in the slice
	resultObjects := make(map[interface{}]interface{}, reflectValue.Len())

	for _, key := range reflectValue.MapKeys() {
		resultObjects[key.Interface()] = ExtractFields(reflectValue.MapIndex(key).Interface(), userType, action)
	}

	return resultObjects
}

// GetJSONFieldName returns the field name given a JSON tag
func GetJSONFieldName(jsonTag string) string {
	if jsonTag == "" {
		return ""
	}

	return strings.Split(jsonTag, ",")[0]
}

// HasPermission returns true if the user has permission for that action on that field
// or false otherwise
func HasPermission(permissionTag string, userType uint, action uint) bool {
	// Get permissions tag
	if permissionTag == "" {
		return true
	}

	permission := int(permissionTag[userType] - '0')

	// Check permissions
	if action == ActionRead {
		return permission == PermissionRead || permission == PermissionReadWrite
	} else if action == ActionWrite {
		return permission == PermissionWrite || permission == PermissionReadWrite
	} else {
		return true
	}
}

// extractField extracts the field and returns a map from string to interface
func extractField(field reflect.StructField, value reflect.Value, userType uint, action uint) map[string]interface{} {
	resultField := map[string]interface{}{}

	if field.PkgPath != "" { // Field is exported or not
		return resultField
	}

	if !HasPermission(field.Tag.Get(PermissionTag), userType, action) {
		return resultField
	}

	// Get the field name
	fieldName := GetJSONFieldName(field.Tag.Get("json"))
	if fieldName == "" {
		fieldName = field.Name
	}

	cleanedField := ExtractFields(value.Interface(), userType, action)

	if field.Anonymous { // Anonymous fields
		subObjectMap, ok := cleanedField.(map[string]interface{})
		if ok {
			for key, value := range subObjectMap {
				resultField[key] = value
			}

			return resultField
		}
	}

	resultField[fieldName] = cleanedField
	return resultField
}

// isSpecialObject returns true if the given object is from a certain type
func isSpecialObject(object interface{}) bool {
	switch object.(type) {
	case time.Time:
		return true
	case sql.NullBool, sql.NullFloat64, sql.NullInt64, sql.NullString:
		return true
	default:
		return false
	}
}

// getReflectValue returns the reflect value of an interface it is exists
// and its valid
func getReflectValue(object interface{}) *reflect.Value {
	// Get the reflect value of the object
	reflectValue := reflect.ValueOf(object)
	for reflectValue.Kind() == reflect.Ptr || reflectValue.Kind() == reflect.Interface {
		reflectValue = reflectValue.Elem()
	}
	if !reflectValue.IsValid() {
		return nil
	}

	return &reflectValue
}
