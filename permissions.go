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
func CleanObject(object interface{}, userType string, action uint) interface{} {
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
func ExtractFields(object interface{}, userType string, action uint) interface{} {
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
func ExtractSingleObjectFields(object interface{}, userType string, action uint) interface{} {
	reflectValue := getReflectValue(object)
	if reflectValue == nil {
		return nil
	}

	// If not struct return the object
	if reflectValue.Kind() != reflect.Struct {
		return reflectValue.Interface()
	}

	// If special object, extract value
	if isSpecialObject(reflectValue.Interface()) {
		return getSpecialObjectValue(reflectValue.Interface())
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
func ExtractMultipleObjectsFields(object interface{}, userType string, action uint) interface{} {
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
func ExtractMapObjectsFields(object interface{}, userType string, action uint) interface{} {
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

// extractField extracts the field and returns a map from string to interface
func extractField(field reflect.StructField, value reflect.Value, userType string, action uint) map[string]interface{} {
	resultField := map[string]interface{}{}

	if field.PkgPath != "" { // Field is exported or not
		return resultField
	}

	if !hasPermission(field.Tag.Get(PermissionTag), userType, action) {
		return resultField
	}

	// Get the field name
	fieldName := getJSONFieldName(field.Tag.Get("json"))
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

// getSpecialObjectValue returns the value of a special object
func getSpecialObjectValue(object interface{}) interface{} {
	switch object.(type) {
	case time.Time:
		return object.(time.Time).String()
	case sql.NullBool:
		value := object.(sql.NullBool).Bool
		if !object.(sql.NullBool).Valid {
			return nil
		}
		return value
	case sql.NullFloat64:
		value := object.(sql.NullFloat64).Float64
		if !object.(sql.NullFloat64).Valid {
			return nil
		}
		return value
	case sql.NullInt64:
		value := object.(sql.NullInt64).Int64
		if !object.(sql.NullInt64).Valid {
			return nil
		}
		return value
	case sql.NullString:
		value := object.(sql.NullString).String
		if !object.(sql.NullString).Valid {
			return nil
		}
		return value
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

// getJSONFieldName returns the field name given a JSON tag
func getJSONFieldName(jsonTag string) string {
	if jsonTag == "" {
		return ""
	}

	return strings.Split(jsonTag, ",")[0]
}

// hasPermission checks if a certain user type has permission for a given action.
// It returns false if the permission for that user is not defined, the user does not have permission
// for that action or the action is invalid. Returns true if the permission tag is not defined or the user
// has permission for that action.
func hasPermission(permissionTag string, userType string, action uint) bool {
	// Get permissions tag
	if permissionTag == "" {
		return true
	}

	permissions := mapPermissions(permissionTag)

	// Check if user type permission is defined
	permission, ok := permissions[userType]
	if !ok {
		return false
	}

	// Check permissions
	if action == ActionRead {
		return strings.Contains(permission, PermissionRead)
	} else if action == ActionWrite {
		return strings.Contains(permission, PermissionWrite)
	}

	return false
}

// mapPermissions converts a permission tag into a map from user type to permission
func mapPermissions(permissionTag string) map[string]string {
	// Create permissions map
	permissions := make(map[string]string)
	for _, permission := range strings.Split(permissionTag, ",") {
		pair := strings.Split(permission, ":")
		permissions[pair[0]] = pair[1]
	}

	return permissions
}
