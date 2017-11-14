package go_pex

import (
	"reflect"
	"strings"
)

// CleanObject removes all the fields that a given user does not have access and
// returns a JSON interface of that object either its an array, slice or struct.
// It uses the json tag to get the field name, it it is not defined uses the field
// name of the struct.
func CleanObject(object interface{}, userType uint, action uint) interface{} {
	reflectType := reflect.TypeOf(object)
	switch reflectType.Kind() {
	case reflect.Slice:
		return CleanSlice(object, userType, action)
	case reflect.Array:
		return object
	default:
		return CleanSingleObject(object, userType, action)
	}
}

// CleanSingleObject removes all the fields that a given user does not have access and
// returns a JSON interface of that object.
// It uses the json tag to get the field name, it it is not defined uses the field
// name of the struct.
func CleanSingleObject(object interface{}, userType uint, action uint) interface{} {
	// TODO: Check for time.Time also

	// Get the value of the object
	reflectValue := getReflectValue(object)
	if reflectValue == nil {
		return nil
	}

	// Iterate through all the fields
	reflectType := reflect.TypeOf(reflectValue.Interface())
	resultObject := map[string]interface{}{}
	for i := 0; i < reflectValue.NumField(); i++ {
		field := reflectValue.Field(i)
		tags := reflectType.Field(i).Tag

		// Check permission
		if !HasPermission(tags.Get(PermissionTag), userType, action) {
			continue
		}

		// Get the field name
		fieldName := GetJSONFieldName(tags.Get("json"))
		if fieldName == "" {
			fieldName = reflectType.Field(i).Name
		}

		if field.Kind() == reflect.Slice || field.Kind() == reflect.Ptr || field.Kind() == reflect.Interface {
			resultObject[fieldName] = CleanObject(field.Interface(), userType, action)
		} else if field.Kind() == reflect.Struct {
			subObject := CleanObject(field.Interface(), userType, action)
			if reflectType.Field(i).Anonymous && subObject != nil {
				subObjectMap, ok := subObject.(map[string]interface{})
				if ok {
					for key, value := range subObjectMap {
						resultObject[key] = value
					}
				} else {
					resultObject[fieldName] = subObject
				}
			} else {
				resultObject[fieldName] = subObject
			}
		} else {
			resultObject[fieldName] = field.Interface()
		}
	}

	return resultObject
}

// CleanSlice removes all the fields that a given user does not have access and
// returns a JSON interface of an array of objects.
// It uses the json tag to get the field name of each of the objects,
// it it is not defined uses the field name of the struct.
func CleanSlice(object interface{}, userType uint, action uint) interface{} {
	// Slice of builtin types, then no need to iterate
	if reflect.TypeOf(object).Elem().Kind() != reflect.Struct {
		return object
	}

	// Iterate through each single object in the slice
	reflectValue := getReflectValue(object)
	if reflectValue == nil {
		return nil
	}

	resultObjects := make([]interface{}, reflectValue.Len())
	for i := 0; i < reflectValue.Len(); i++ {
		resultObjects[i] = CleanObject(reflectValue.Index(i), userType, action)
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

// getReflectValue returns the reflect value of an interface it is exists and its valid
func getReflectValue(object interface{}) *reflect.Value {
	// Get the structure of the object
	reflectValue := reflect.ValueOf(object)
	for reflectValue.Kind() == reflect.Ptr || reflectValue.Kind() == reflect.Interface {
		reflectValue = reflectValue.Elem()
	}
	if !reflectValue.IsValid() {
		return nil
	}

	return &reflectValue
}
