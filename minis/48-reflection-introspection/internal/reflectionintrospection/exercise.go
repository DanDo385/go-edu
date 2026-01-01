//go:build !solution && !reference

// Package exercise contains hands-on exercises for understanding Go reflection.
//
// LEARNING OBJECTIVES:
// - Inspect types and values at runtime
// - Read and parse struct tags
// - Modify values through reflection
// - Call methods dynamically
// - Understand reflection performance trade-offs
// - Avoid common reflection pitfalls

package reflectionintrospection

// TODO: Import the "reflect" package and other packages as needed.

// TODO: Implement these functions according to the specifications in the tests.
// Each function tests a different aspect of reflection.

// ============================================================================
// EXERCISE 1: Type Inspection
// ============================================================================

// GetTypeName returns the name of the type of the given value.
func GetTypeName(v interface{}) string {
	// TODO: Implement this function.
	// - Use `reflect.TypeOf(v)` to get the `reflect.Type` of the value.
	// - The `Name()` method of a `reflect.Type` returns the type's name within its package.
	return reflect.TypeOf(v).Name()
}

// GetKind returns the kind of the type (the underlying category).
func GetKind(v interface{}) string {
	// TODO: Implement this function.
	// - `Kind` is different from `Name`. For a struct `type User struct{}`, the Name is "User" but the Kind is "struct".
	// - Use `reflect.TypeOf(v).Kind()` to get the `reflect.Kind`.
	// - Call the `.String()` method on the `Kind` to get its string representation.
	return reflect.TypeOf(v).Kind().String()
}

// CountFields returns the number of fields in a struct.
func CountFields(v interface{}) int {
	// TODO: Implement this function.
	// - Get the `reflect.Type` of `v`.
	// - If the type is a pointer to a struct, you need to get the element it points to: `t = t.Elem()`.
	// - Check if the `Kind` is `reflect.Struct`. If not, return 0.
	// - If it is a struct, use the `NumField()` method to get the number of fields.
	t := reflect.TypeOf(v)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
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
	// TODO: Implement this function.
	// - Get the `reflect.Type`.
	// - Handle pointers with `t.Elem()`.
	// - Get the specific field with `t.FieldByName(fieldName)`. This returns a `StructField` and a boolean.
	// - If the field is found, use the `Tag.Get("json")` method on the `StructField` to look up the value associated with the "json" key.
	t := reflect.TypeOf(v)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return ""
	}
	field, ok := t.FieldByName(fieldName)
	if !ok {
		return ""
	}
	return field.Tag.Get("json")
}

// GetAllTags returns all struct tags for a named field.
func GetAllTags(v interface{}, fieldName string) string {
	// TODO: Implement this function.
	// - This is similar to `GetJSONTag`, but instead of getting a specific key, you'll return the whole tag string.
	// - After finding the `StructField`, access its `Tag` property and convert it to a `string`.
	t := reflect.TypeOf(v)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return ""
	}
	field, ok := t.FieldByName(fieldName)
	if !ok {
		return ""
	}
	return string(field.Tag)
}

// ============================================================================
// EXERCISE 3: Value Inspection
// ============================================================================

// GetFieldValue returns the value of a named field in a struct.
func GetFieldValue(val interface{}, fieldName string) interface{} {
	// TODO: Implement this function.
	// - Use `reflect.ValueOf(val)` to get the `reflect.Value`.
	// - Handle pointers with `v.Elem()`.
	// - Get the field with `v.FieldByName(fieldName)`.
	// - **Important**: Before accessing the value, check if the field is valid (i.e., it was found) using `field.IsValid()`.
	// - If it's valid, return its value by calling `.Interface()`. This converts the `reflect.Value` back into a standard `interface{}`.
	v := reflect.ValueOf(val)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return nil
	}
	field := v.FieldByName(fieldName)
	if !field.IsValid() {
		return nil
	}
	return field.Interface()
}

// GetFieldValues returns a map of field names to their values for a struct.
func GetFieldValues(val interface{}) map[string]interface{} {
	// TODO: Implement this function.
	// - Get both the `reflect.Value` and `reflect.Type` of the input.
	// - Handle pointers using `.Elem()`.
	// - Loop from `i = 0` to `v.NumField() - 1`.
	// - In the loop, get the `reflect.StructField` (for the name and to check if it's exported) from the type: `t.Field(i)`.
	// - Also get the `reflect.Value` for the field: `v.Field(i)`.
	// - If the field `IsExported()`, add its name and value (`.Interface()`) to the results map.
	result := make(map[string]interface{})
	v := reflect.ValueOf(val)
	t := reflect.TypeOf(val)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
		t = t.Elem()
	}
	if v.Kind() != reflect.Struct {
		return result
	}
	for i := 0; i < v.NumField(); i++ {
		fieldType := t.Field(i)
		fieldValue := v.Field(i)
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
	// TODO: Implement this function.

	// Modifying values with reflection is more complex and requires care.

	// Step 1: Get the `reflect.Value` of `val`.
	// Step 2: Check that the value is a pointer. You can only modify things that are addressable.
	// - `if v.Kind() != reflect.Ptr { return error }`
	// Step 3: Get the element the pointer points to.
	// - `elem := v.Elem()`
	// Step 4: Find the field on the element.
	// - `field := elem.FieldByName(fieldName)`
	// - Check if the field is valid.
	// Step 5: Check if the field can be set.
	// - `if !field.CanSet() { return error }`
	// - A field can only be set if it is exported (starts with an uppercase letter).
	// Step 6: Get the `reflect.Value` of the `newValue`.
	// - `newVal := reflect.ValueOf(newValue)`
	// Step 7: Check that the types match.
	// - `if field.Type() != newVal.Type() { return error }`
	// Step 8: Set the value.
	// - `field.Set(newVal)`
	return nil
}

// ============================================================================
// EXERCISE 5: Dynamic Method Calls
// ============================================================================

// CallMethod calls a method by name with the given arguments.
func CallMethod(obj interface{}, methodName string, args ...interface{}) []interface{} {
	// TODO: Implement this function.

	// Step 1: Get the `reflect.Value` of the object.
	// Step 2: Get the method by its name.
	// - `method := v.MethodByName(methodName)`
	// - Check if `method.IsValid()`. If not, the method doesn't exist, so return `nil`.

	// Step 3: Prepare the arguments for the call.
	// - The `Call` method expects a `[]reflect.Value`.
	// - Create a slice of `reflect.Value` with the same length as `args`.
	// - Loop through your `args` and convert each one to a `reflect.Value` using `reflect.ValueOf(arg)`.

	// Step 4: Call the method.
	// - `results := method.Call(reflectArgs)`
	// - This returns a `[]reflect.Value`.

	// Step 5: Convert the results back to `[]interface{}`.
	// - Create a new slice of `interface{}`.
	// - Loop through the `results` and convert each `reflect.Value` back using `.Interface()`.

	return nil
}

// HasMethod checks if a value has a method with the given name.
func HasMethod(obj interface{}, methodName string) bool {
	// TODO: Implement this function.
	// - Get the `reflect.Value` of the object.
	// - Get the method by name using `.MethodByName(methodName)`.
	// - The `IsValid()` method on the result will be `true` if the method exists and `false` otherwise.
	return false
}

// ============================================================================
// EXERCISE 6: Type Comparison
// ============================================================================

// SameType checks if two values have the same type.
func SameType(a, b interface{}) bool {
	// TODO: Implement this function.
	// - Use `reflect.TypeOf()` on both `a` and `b`.
	// - The resulting `reflect.Type` values can be compared directly with `==`.
	return false
}

// IsPointer checks if a value is a pointer.
func IsPointer(v interface{}) bool {
	// TODO: Implement this function.
	// - Get the `reflect.Type` of `v`.
	// - Get the `Kind` of the type.
	// - Compare the `Kind` to the constant `reflect.Ptr`.
	return false
}

// ============================================================================
// EXERCISE 7: Slice and Map Operations
// ============================================================================

// SliceLength returns the length of a slice using reflection.
func SliceLength(slice interface{}) int {
	// TODO: Implement this function.
	// - Get the `reflect.Value` of the input.
	// - Check if its `Kind` is `reflect.Slice`. If not, return -1.
	// - If it is a slice, use the `.Len()` method on the `reflect.Value` to get its length.
	return -1
}

// MapKeys returns all keys from a map as []interface{}.
func MapKeys(m interface{}) []interface{} {
	// TODO: Implement this function.
	// - Get the `reflect.Value` of the input.
	// - Check if its `Kind` is `reflect.Map`. If not, return `nil`.
	// - If it is a map, use the `.MapKeys()` method. This returns a `[]reflect.Value`.
	// - You need to convert this `[]reflect.Value` into an `[]interface{}` by looping through it and calling `.Interface()` on each key.
	return nil
}

// ============================================================================
// EXERCISE 8: Creating Values
// ============================================================================

// NewInstance creates a new instance of the same type as the given value.
func NewInstance(v interface{}) interface{} {
	// TODO: Implement this function.
	// - `reflect.New(t)` creates a new zero value of type `t` and returns a `reflect.Value` representing a pointer to it.
	// - Get the type of `v` using `reflect.TypeOf(v)`.
	// - Call `reflect.New()` with that type.
	// - Call `.Interface()` on the result to get the pointer back as an `interface{}`.
	return nil
}

// ============================================================================
// EXERCISE 9: Advanced - Struct Field Names
// ============================================================================

// GetFieldNames returns the names of all fields in a struct.
func GetFieldNames(v interface{}) []string {
	// TODO: Implement this function.
	// - Get the `reflect.Type` and handle pointers.
	// - Check if it's a struct.
	// - Loop from `i = 0` to `t.NumField() - 1`.
	// - In the loop, get the `StructField` with `t.Field(i)`.
	// - The `StructField` type has a `Name` property and an `IsExported()` method.
	// - If the field is exported, add its name to your result slice.
	return nil
}

// ============================================================================
// EXERCISE 10: Advanced - Deep Copy
// ============================================================================

// DeepCopy creates a deep copy of a struct using reflection.
func DeepCopy(v interface{}) interface{} {
	// TODO: Implement this function.
	// This function should create a new struct of the same type and copy all the field values.

	// Step 1: Get the `reflect.Value` and `reflect.Type` of the input, handling pointers.
	// Step 2: Check if it's a struct.
	// Step 3: Create a new instance of the struct type.
	// - `newVal := reflect.New(typ)` gives you a pointer to the new struct.
	// - `newElem := newVal.Elem()` gives you the struct itself, which you can set fields on.
	// Step 4: Loop through the fields of the original value (`val`).
	// - `for i := 0; i < val.NumField(); i++ { ... }`
	// Step 5: In the loop, get the source field (`val.Field(i)`) and the destination field (`newElem.Field(i)`).
	// Step 6: Check if the destination field `.CanSet()`.
	// Step 7: If it can be set, set its value from the source field: `dstField.Set(srcField)`.
	// Step 8: Return the new pointer: `return newVal.Interface()`.
	return nil
}
