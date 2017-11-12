package go_pex

import (
	"reflect"
	"strings"
	"time"
	"encoding/json"
	"fmt"
)

// CleanObject removes all the fields that a given user does not have access and
// returns a string to interface map where each key is the field of the object
// and the value, the value of that field.
// It uses the json tag to get the field name, it it is not defined uses the field
// name of the struct.
func CleanObject(object interface{}, userType uint, action uint) map[string]interface{} {
	if action == ActionWrite {
		return nil
	} else if action == ActionRead {
		return nil
	} else {
		panic("Invalid value for action.")
	}
}

// CleanSlice removes all the fields that a given user does not have access and
// returns an array of string to interface map where each key is the field of the object
// and the value, the value of that field.
// It uses the json tag to get the field name, it it is not defined uses the field
// name of the struct.
func CleanSlice(object []interface{}, userType uint, action uint) []map[string]interface{} {
	objects := make([]map[string]interface{}, len(object))
	for i := 0; i < len(object); i++ {
		objects[i] = CleanObject(object[i], userType, action)
	}
	return objects
}

func cleanObjectRead(object interface{}, userType uint, action uint) map[string]interface{} {
	// Get the structure of the object
	reflectValue := reflect.ValueOf(object)
	for reflectValue.Kind() == reflect.Ptr || reflectValue.Kind() == reflect.Interface {
		reflectValue = reflectValue.Elem()
	}
	if !reflectValue.IsValid() {
		return nil
	}
}

// RemoveFields removes the fields that the user can't view, write or export from the interface
func RemoveFields(obj interface{}, userType uint, action uint) interface{} {

	objectType := reflect.TypeOf(objectValue.Interface())
	resultObject := map[string]interface{}{}

	for i := 0; i < objectValue.NumField(); i++ {
		field := objectValue.Field(i)
		if field.Kind() == reflect.Slice && action == ActionWrite {
			continue
		}

		fieldName := objectType.Field(i).Name
		tags := objectType.Field(i).Tag

		// Check the gorm tag
		gormTag := tags.Get("gorm")
		if gormTag == "-" && action == ActionWrite {
			continue
		}

		// Get permissions tag
		uptecTag := tags.Get("uptec")
		permission := PermissionReadWrite
		if uptecTag != "" {
			permission = int(uptecTag[userType-1] - '0')
		}

		// Check permissions
		if action == ActionRead || action == ActionExport {
			if permission != PermissionRead && permission != PermissionReadWrite {
				continue
			}
		} else if action == ActionWrite {
			if permission != PermissionReadWrite && permission != PermissionWrite {
				continue
			}
		}

		// Tag Key
		var tag string
		if action == ActionExport {
			tag = tags.Get("export")
			if tag == "" {
				tag = fieldName
			} else if tag == "-" {
				continue
			}
		} else {
			tag = tags.Get("json")
			if tag == "" {
				tag = fieldName
			} else {
				if strings.HasSuffix(tag, ",omitempty") {
					tag = tag[0 : len(tag)-len(",omitempty")]
				}

				if tag == "-" {
					continue
				}
			}
		}

		// Check if it is a struct or pointer
		if field.Kind() == reflect.Ptr || field.Kind() == reflect.Struct || field.Kind() == reflect.Slice {
			switch field.Interface().(type) {
			case time.Time, *time.Time:
				if action == ActionExport {
					resultObject[tag] = field.Interface().(time.Time).Format("02-01-2006")
				} else {
					resultObject[tag] = field.Interface()
				}
			case null.Int, null.Float, null.String, *null.Int, *null.Float, *null.String:
				resultObject[tag] = field.Interface()
			case null.Time, *null.Time:
				if action == ActionExport {
					nullTime := field.Interface().(null.Time)
					if nullTime.Valid {
						resultObject[tag] = nullTime.Time.Format("02-01-2006")
					} else {
						resultObject[tag] = field.Interface()
					}
				} else {
					resultObject[tag] = field.Interface()
				}
			case null.Bool, *null.Bool:
				if action == ActionExport {
					if field.Interface().(bool) {
						resultObject[tag] = "Yes"
					} else {
						resultObject[tag] = "No"
					}
				} else {
					resultObject[tag] = field.Interface()
				}
			case []uint, []string, []int, []int32, []int64,
			[]bool, []float32, []float64, []byte:
				resultObject[tag] = field.Interface()
			default:
				subObject := RemoveFields(field.Interface(), userType, action)

				if objectType.Field(i).Anonymous && subObject != nil { // Inheritance
					subObjectMap, ok := subObject.(map[string]interface{})
					if ok {
						for key, value := range subObjectMap {
							resultObject[key] = value
						}
					} else {
						resultObject[tag] = subObject
					}
				} else {
					resultObject[tag] = subObject
				}
			}
		} else {
			if field.Interface() == nil {
				continue
			}

			switch field.Interface().(type) {
			case bool:
				if action == ActionExport {
					if field.Interface().(bool) {
						resultObject[tag] = "Yes"
					} else {
						resultObject[tag] = "No"
					}
				} else {
					resultObject[tag] = field.Interface()
				}
			default:
				resultObject[tag] = field.Interface()
			}
		}
	}

	return resultObject
}
