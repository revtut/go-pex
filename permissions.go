package go_pex

import (
	"reflect"
	"strings"
)

// CleanObject removes all the fields that a given user does not have access and
// returns a string to interface map where each key is the field of the object
// and the value, the value of that field.
// It uses the json tag to get the field name, it it is not defined uses the field
// name of the struct.
func CleanObject(object interface{}, userType uint, action uint) map[string]interface{} {
	// TODO: Check if slice because it might receive slice pointers

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

		if field.Kind() == reflect.Slice {
			if reflect.TypeOf(field.Interface()).Elem().Kind() == reflect.Struct { // Type of the slice
				resultObject[fieldName] = CleanSlice(field.Interface().([]interface{}), userType, action)
			} else {
				resultObject[fieldName] = field.Interface()
			}
		} else if field.Kind() == reflect.Ptr || field.Kind() == reflect.Struct || field.Kind() == reflect.Interface {
			subObject := CleanObject(field.Interface(), userType, action)
			if reflectType.Field(i).Anonymous && subObject != nil {
				for key, value := range subObject {
					resultObject[key] = value
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
// returns an array of string to interface map where each key is the field of the object
// and the value, the value of that field.
// It uses the json tag to get the field name, it it is not defined uses the field
// name of the struct.
func CleanSlice(object []interface{}, userType uint, action uint) []map[string]interface{} {
	resultObjects := make([]map[string]interface{}, len(object))
	for i := 0; i < len(object); i++ {
		resultObjects[i] = CleanObject(object[i], userType, action)
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

	/*
	TODO: Consider this case
	if fieldName == "-" {
		return ""
	}
	*/

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
