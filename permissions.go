package go_pex

import (
	"reflect"
	"strings"
	"time"
)

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
	case reflect.Slice, reflect.Array:
		return ExtractMultipleObjectsFields(object, userType, action)
	case reflect.Struct:
		return ExtractSingleObjectFields(object, userType, action)
	case reflect.Map:
		// TODO: Clean maps
		fallthrough
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
		field := reflectType.Field(i)

		if field.PkgPath != "" { // Field is exported or not
			continue
		}

		if !HasPermission(field.Tag.Get(PermissionTag), userType, action) {
			continue
		}

		// Get the field name
		fieldName := GetJSONFieldName(field.Tag.Get("json"))
		if fieldName == "" {
			fieldName = field.Name
		}

		// Anonymous fields
		cleanedField := ExtractFields(reflectValue.Field(i).Interface(), userType, action)
		if field.Anonymous {
			subObjectMap, ok := cleanedField.(map[string]interface{})
			if ok {
				for key, value := range subObjectMap {
					resultObject[key] = value
				}
			} else {
				resultObject[fieldName] = cleanedField
			}
		} else {
			resultObject[fieldName] = cleanedField
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

// isSpecialObject returns true if the given object is from a certain type
func isSpecialObject(object interface{}) bool {
	switch object.(type) {
	case time.Time:
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
