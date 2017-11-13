package go_pex

import (
	"reflect"
	"strings"
	"fmt"
)

// CleanObject removes all the fields that a given user does not have access and
// returns a JSON interface of that object.
// It uses the json tag to get the field name, it it is not defined uses the field
// name of the struct.
func CleanObject(object interface{}, userType uint, action uint) interface{} {
	reflectType := reflect.TypeOf(object)
	switch reflectType.Kind() {
	case reflect.Slice:
		return cleanSlice(object, userType, action)
	case reflect.Array:
		return object
	default:
		return cleanSingleObject(object, userType, action)
	}
}

func cleanSingleObject(object interface{}, userType uint, action uint) interface{} {
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
		if !hasPermission(tags, userType, action) {
			continue
		}

		// Get the field name
		fieldName := getJSONFieldName(tags)
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

func cleanSlice(object interface{}, userType uint, action uint) interface{} {
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

func getJSONFieldName(tags reflect.StructTag) string {
	fieldName := tags.Get("json")
	if fieldName == "" {
		return ""
	}

	if strings.HasSuffix(fieldName, ",omitempty") {
		fieldName = fieldName[0: len(fieldName)-len(",omitempty")]
	}

	return fieldName
}

func hasPermission(tags reflect.StructTag, userType uint, action uint) bool {
	// Get permissions tag
	permissionTag := tags.Get(PERMISSION_TAG)
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
