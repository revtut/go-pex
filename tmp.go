package go_pex

/*
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
}*/
