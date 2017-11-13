package go_pex

import (
	"reflect"
	"encoding/json"
	"strings"
	"time"
)

// CleanStruct removes the fields that the user can't read / write
func CleanStruct(obj interface{}, userType uint, action uint) {
	// Remove fields without permission to change
	tmp := RemoveFields(obj, userType, action)
	p := reflect.ValueOf(obj).Elem()
	p.Set(reflect.Zero(p.Type()))
	InterfaceToStructure(tmp, &obj)
}

// InterfaceToStructure is a method that converts a map string to value into a
// interface
func InterfaceToStructure(data interface{}, object interface{}) error {
	result, err := json.Marshal(data)
	if err != nil {
		return err
	}
	json.Unmarshal(result, &object)

	return nil
}

// RemoveFields removes the fields that the user can't view, write or export from the interface
func RemoveFields(obj interface{}, userType uint, action uint) interface{} {
	// Slice
	if reflect.TypeOf(obj).Kind() == reflect.Slice {
		values := reflect.ValueOf(obj)

		result := make([]interface{}, values.Len())
		for i := 0; i < values.Len(); i++ {
			sliceObject := values.Index(i).Interface()
			result[i] = RemoveFields(sliceObject, userType, action)
		}

		return result
	}

	// Single Object
	objectValue := reflect.ValueOf(obj)
	if reflect.TypeOf(obj).Kind() == reflect.Ptr {
		objectValue = objectValue.Elem()
	}
	if !objectValue.IsValid() {
		return nil
	}

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
