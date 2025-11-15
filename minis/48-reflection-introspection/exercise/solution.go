//go:build solution
// +build solution

// Package exercise contains complete solutions for reflection exercises.
package exercise

import (
	"fmt"
	"reflect"
)

// ============================================================================
// EXERCISE 1: Type Inspection
// ============================================================================

// GetTypeName returns the name of the type of the given value.
func GetTypeName(v interface{}) string {
	t := reflect.TypeOf(v)
	return t.Name()
}

// GetKind returns the kind of the type (the underlying category).
func GetKind(v interface{}) string {
	t := reflect.TypeOf(v)
	return t.Kind().String()
}

// CountFields returns the number of fields in a struct.
func CountFields(v interface{}) int {
	t := reflect.TypeOf(v)

	// Handle pointer to struct
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	// Check if it's a struct
	if t.Kind() != reflect.Struct {
		return 0
	}

	return t.NumField()
}

// ============================================================================
// EXERCISE 2: Struct Tags
// ============================================================================

// GetJSONTag returns the json tag for a named field in a struct.
func GetJSONTag(v interface{}, fieldName string) string {
	t := reflect.TypeOf(v)

	// Handle pointer to struct
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	// Check if it's a struct
	if t.Kind() != reflect.Struct {
		return ""
	}

	// Find the field
	field, ok := t.FieldByName(fieldName)
	if !ok {
		return ""
	}

	// Return the json tag
	return field.Tag.Get("json")
}

// GetAllTags returns all struct tags for a named field.
func GetAllTags(v interface{}, fieldName string) string {
	t := reflect.TypeOf(v)

	// Handle pointer to struct
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	// Check if it's a struct
	if t.Kind() != reflect.Struct {
		return ""
	}

	// Find the field
	field, ok := t.FieldByName(fieldName)
	if !ok {
		return ""
	}

	// Return all tags
	return string(field.Tag)
}

// ============================================================================
// EXERCISE 3: Value Inspection
// ============================================================================

// GetFieldValue returns the value of a named field in a struct.
func GetFieldValue(val interface{}, fieldName string) interface{} {
	v := reflect.ValueOf(val)

	// Handle pointer
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	// Check if it's a struct
	if v.Kind() != reflect.Struct {
		return nil
	}

	// Get the field
	field := v.FieldByName(fieldName)
	if !field.IsValid() {
		return nil
	}

	return field.Interface()
}

// GetFieldValues returns a map of field names to their values for a struct.
func GetFieldValues(val interface{}) map[string]interface{} {
	result := make(map[string]interface{})

	v := reflect.ValueOf(val)
	t := reflect.TypeOf(val)

	// Handle pointer
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
		t = t.Elem()
	}

	// Check if it's a struct
	if v.Kind() != reflect.Struct {
		return result
	}

	// Iterate over fields
	for i := 0; i < v.NumField(); i++ {
		fieldType := t.Field(i)
		fieldValue := v.Field(i)

		// Only include exported fields
		if fieldType.IsExported() {
			result[fieldType.Name] = fieldValue.Interface()
		}
	}

	return result
}

// ============================================================================
// EXERCISE 4: Value Modification
// ============================================================================

// SetFieldValue sets the value of a named field in a struct.
func SetFieldValue(val interface{}, fieldName string, newValue interface{}) error {
	v := reflect.ValueOf(val)

	// Check if it's a pointer
	if v.Kind() != reflect.Ptr {
		return fmt.Errorf("value must be a pointer")
	}

	// Get the element being pointed to
	elem := v.Elem()

	// Check if it's a struct
	if elem.Kind() != reflect.Struct {
		return fmt.Errorf("value must point to a struct")
	}

	// Get the field
	field := elem.FieldByName(fieldName)
	if !field.IsValid() {
		return fmt.Errorf("field %s does not exist", fieldName)
	}

	// Check if we can set it
	if !field.CanSet() {
		return fmt.Errorf("field %s cannot be set", fieldName)
	}

	// Set the value
	newVal := reflect.ValueOf(newValue)
	if field.Type() != newVal.Type() {
		return fmt.Errorf("type mismatch: field is %s, value is %s", field.Type(), newVal.Type())
	}

	field.Set(newVal)
	return nil
}

// ============================================================================
// EXERCISE 5: Dynamic Method Calls
// ============================================================================

// CallMethod calls a method by name with the given arguments.
func CallMethod(obj interface{}, methodName string, args ...interface{}) []interface{} {
	v := reflect.ValueOf(obj)

	// Get the method
	method := v.MethodByName(methodName)
	if !method.IsValid() {
		return nil
	}

	// Convert arguments to []reflect.Value
	reflectArgs := make([]reflect.Value, len(args))
	for i, arg := range args {
		reflectArgs[i] = reflect.ValueOf(arg)
	}

	// Call the method
	results := method.Call(reflectArgs)

	// Convert results to []interface{}
	interfaceResults := make([]interface{}, len(results))
	for i, result := range results {
		interfaceResults[i] = result.Interface()
	}

	return interfaceResults
}

// HasMethod checks if a value has a method with the given name.
func HasMethod(obj interface{}, methodName string) bool {
	v := reflect.ValueOf(obj)
	method := v.MethodByName(methodName)
	return method.IsValid()
}

// ============================================================================
// EXERCISE 6: Type Comparison
// ============================================================================

// SameType checks if two values have the same type.
func SameType(a, b interface{}) bool {
	return reflect.TypeOf(a) == reflect.TypeOf(b)
}

// IsPointer checks if a value is a pointer.
func IsPointer(v interface{}) bool {
	return reflect.TypeOf(v).Kind() == reflect.Ptr
}

// ============================================================================
// EXERCISE 7: Slice and Map Operations
// ============================================================================

// SliceLength returns the length of a slice using reflection.
func SliceLength(slice interface{}) int {
	v := reflect.ValueOf(slice)

	if v.Kind() != reflect.Slice {
		return -1
	}

	return v.Len()
}

// MapKeys returns all keys from a map as []interface{}.
func MapKeys(m interface{}) []interface{} {
	v := reflect.ValueOf(m)

	if v.Kind() != reflect.Map {
		return nil
	}

	keys := v.MapKeys()
	result := make([]interface{}, len(keys))

	for i, key := range keys {
		result[i] = key.Interface()
	}

	return result
}

// ============================================================================
// EXERCISE 8: Creating Values
// ============================================================================

// NewInstance creates a new instance of the same type as the given value.
func NewInstance(v interface{}) interface{} {
	t := reflect.TypeOf(v)
	ptr := reflect.New(t)
	return ptr.Interface()
}

// ============================================================================
// EXERCISE 9: Advanced - Struct Field Names
// ============================================================================

// GetFieldNames returns the names of all fields in a struct.
func GetFieldNames(v interface{}) []string {
	t := reflect.TypeOf(v)

	// Handle pointer
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	// Check if it's a struct
	if t.Kind() != reflect.Struct {
		return nil
	}

	var names []string
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		// Only include exported fields
		if field.IsExported() {
			names = append(names, field.Name)
		}
	}

	return names
}

// ============================================================================
// EXERCISE 10: Advanced - Deep Copy
// ============================================================================

// DeepCopy creates a deep copy of a struct using reflection.
func DeepCopy(v interface{}) interface{} {
	val := reflect.ValueOf(v)
	typ := reflect.TypeOf(v)

	// Handle pointer
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
		typ = typ.Elem()
	}

	// Check if it's a struct
	if val.Kind() != reflect.Struct {
		return nil
	}

	// Create new instance
	newVal := reflect.New(typ)
	newElem := newVal.Elem()

	// Copy all fields
	for i := 0; i < val.NumField(); i++ {
		srcField := val.Field(i)
		dstField := newElem.Field(i)

		// Only set if we can (exported fields)
		if dstField.CanSet() {
			dstField.Set(srcField)
		}
	}

	return newVal.Interface()
}
