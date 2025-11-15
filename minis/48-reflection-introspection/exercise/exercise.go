//go:build !solution
// +build !solution

// Package exercise contains hands-on exercises for understanding Go reflection.
//
// LEARNING OBJECTIVES:
// - Inspect types and values at runtime
// - Read and parse struct tags
// - Modify values through reflection
// - Call methods dynamically
// - Understand reflection performance trade-offs
// - Avoid common reflection pitfalls

package exercise

// TODO: Import the "reflect" package and other packages as needed.

// TODO: Implement these functions according to the specifications in the tests.
// Each function tests a different aspect of reflection.

// ============================================================================
// EXERCISE 1: Type Inspection
// ============================================================================

// GetTypeName returns the name of the type of the given value.
//
// REQUIREMENTS:
// - Use reflect.TypeOf to get the type
// - Return the type's Name() (e.g., "User", "int", "string")
//
// EXAMPLES:
//   GetTypeName(User{}) → "User"
//   GetTypeName(42) → "int"
//   GetTypeName("hello") → "string"
//
// HINT: t := reflect.TypeOf(v); return t.Name()
func GetTypeName(v interface{}) string {
	// TODO: Implement this function
	return ""
}

// GetKind returns the kind of the type (the underlying category).
//
// REQUIREMENTS:
// - Use reflect.TypeOf to get the type
// - Return the Kind as a string (use .String() on the Kind)
//
// EXAMPLES:
//   GetKind(User{}) → "struct"
//   GetKind(42) → "int"
//   GetKind(&User{}) → "ptr"
//
// HINT: t := reflect.TypeOf(v); return t.Kind().String()
func GetKind(v interface{}) string {
	// TODO: Implement this function
	return ""
}

// CountFields returns the number of fields in a struct.
//
// REQUIREMENTS:
// - Check if the value is a struct (Kind() == reflect.Struct)
// - If not a struct, return 0
// - If a struct, return the number of fields
//
// HINT: Use NumField() on the Type
func CountFields(v interface{}) int {
	// TODO: Implement this function
	return 0
}

// ============================================================================
// EXERCISE 2: Struct Tags
// ============================================================================

// GetJSONTag returns the json tag for a named field in a struct.
//
// REQUIREMENTS:
// - Get the type of the value
// - Find the field by name using FieldByName
// - If the field doesn't exist, return ""
// - If it exists, return the "json" tag value
//
// EXAMPLES:
//   GetJSONTag(User{}, "Name") → "name"
//   GetJSONTag(User{}, "Email") → "email"
//   GetJSONTag(User{}, "NoSuchField") → ""
//
// HINT: field, ok := t.FieldByName(fieldName)
//       if ok { return field.Tag.Get("json") }
func GetJSONTag(v interface{}, fieldName string) string {
	// TODO: Implement this function
	return ""
}

// GetAllTags returns all struct tags for a named field.
//
// REQUIREMENTS:
// - Get the type of the value
// - Find the field by name
// - If the field doesn't exist, return ""
// - If it exists, return the complete tag string (all tags, not just one)
//
// EXAMPLES:
//   GetAllTags(User{}, "Name") → `json:"name" validate:"required"`
//   GetAllTags(Product{}, "Name") → `db:"product_name" json:"name"`
//
// HINT: Use field.Tag (not field.Tag.Get(...))
//       Convert to string with string(field.Tag)
func GetAllTags(v interface{}, fieldName string) string {
	// TODO: Implement this function
	return ""
}

// ============================================================================
// EXERCISE 3: Value Inspection
// ============================================================================

// GetFieldValue returns the value of a named field in a struct.
//
// REQUIREMENTS:
// - Use reflect.ValueOf to get the value
// - Get the field by name using FieldByName
// - Return the field value as interface{} using .Interface()
// - If the field doesn't exist or value is not a struct, return nil
//
// EXAMPLES:
//   user := User{Name: "Alice", Email: "alice@example.com", Age: 30}
//   GetFieldValue(user, "Name") → "Alice"
//   GetFieldValue(user, "Age") → 30
//
// HINT: v := reflect.ValueOf(val); field := v.FieldByName(fieldName)
//       Check field.IsValid() before calling .Interface()
func GetFieldValue(val interface{}, fieldName string) interface{} {
	// TODO: Implement this function
	return nil
}

// GetFieldValues returns a map of field names to their values for a struct.
//
// REQUIREMENTS:
// - Return a map[string]interface{} with field names as keys
// - Only include exported (public) fields
// - If the value is not a struct, return an empty map
//
// EXAMPLES:
//   user := User{Name: "Bob", Email: "bob@example.com", Age: 25}
//   GetFieldValues(user) → map[string]interface{}{
//     "Name": "Bob",
//     "Email": "bob@example.com",
//     "Age": 25,
//   }
//
// HINT: Use NumField() and Field(i) to iterate
//       Check field.IsExported() (from the Type)
func GetFieldValues(val interface{}) map[string]interface{} {
	// TODO: Implement this function
	return nil
}

// ============================================================================
// EXERCISE 4: Value Modification
// ============================================================================

// SetFieldValue sets the value of a named field in a struct.
//
// REQUIREMENTS:
// - The input val must be a POINTER to a struct
// - Use .Elem() to get the struct being pointed to
// - Find the field by name
// - Set the field value using Set(reflect.ValueOf(newValue))
// - Return nil on success, error on failure
// - Return an error if:
//   - val is not a pointer
//   - val doesn't point to a struct
//   - field doesn't exist
//   - field is not settable
//
// EXAMPLES:
//   user := &User{Name: "Alice"}
//   SetFieldValue(user, "Name", "Bob") → nil (and user.Name is now "Bob")
//   SetFieldValue(User{}, "Name", "Bob") → error (not a pointer)
//
// HINT: v := reflect.ValueOf(val)
//       Check v.Kind() == reflect.Ptr
//       elem := v.Elem()
//       field := elem.FieldByName(fieldName)
//       Check field.CanSet()
func SetFieldValue(val interface{}, fieldName string, newValue interface{}) error {
	// TODO: Implement this function
	return nil
}

// ============================================================================
// EXERCISE 5: Dynamic Method Calls
// ============================================================================

// CallMethod calls a method by name with the given arguments.
//
// REQUIREMENTS:
// - Use MethodByName to get the method
// - Convert arguments to []reflect.Value
// - Call the method
// - Return the results as []interface{}
// - If the method doesn't exist, return nil
//
// EXAMPLES:
//   calc := Calculator{}
//   CallMethod(calc, "Add", 5, 3) → []interface{}{8}
//   CallMethod(calc, "Multiply", 4, 7) → []interface{}{28}
//
// HINT: method := v.MethodByName(methodName)
//       Check method.IsValid()
//       Build args as []reflect.Value using reflect.ValueOf for each arg
//       results := method.Call(args)
//       Convert results back to []interface{}
func CallMethod(obj interface{}, methodName string, args ...interface{}) []interface{} {
	// TODO: Implement this function
	return nil
}

// HasMethod checks if a value has a method with the given name.
//
// REQUIREMENTS:
// - Use MethodByName to look up the method
// - Return true if the method exists, false otherwise
//
// EXAMPLES:
//   HasMethod(Calculator{}, "Add") → true
//   HasMethod(Calculator{}, "Divide") → false
//
// HINT: method := v.MethodByName(methodName)
//       return method.IsValid()
func HasMethod(obj interface{}, methodName string) bool {
	// TODO: Implement this function
	return false
}

// ============================================================================
// EXERCISE 6: Type Comparison
// ============================================================================

// SameType checks if two values have the same type.
//
// REQUIREMENTS:
// - Return true if both values have the exact same type
// - Return false otherwise
//
// EXAMPLES:
//   SameType(42, 10) → true (both int)
//   SameType(42, "hello") → false (int vs string)
//   SameType(User{}, User{Name: "Alice"}) → true (both User)
//
// HINT: Use reflect.TypeOf for both values and compare with ==
func SameType(a, b interface{}) bool {
	// TODO: Implement this function
	return false
}

// IsPointer checks if a value is a pointer.
//
// REQUIREMENTS:
// - Return true if the value's kind is reflect.Ptr
// - Return false otherwise
//
// EXAMPLES:
//   x := 42
//   IsPointer(x) → false
//   IsPointer(&x) → true
//
// HINT: Use Kind() and compare with reflect.Ptr
func IsPointer(v interface{}) bool {
	// TODO: Implement this function
	return false
}

// ============================================================================
// EXERCISE 7: Slice and Map Operations
// ============================================================================

// SliceLength returns the length of a slice using reflection.
//
// REQUIREMENTS:
// - Check if the value is a slice (Kind() == reflect.Slice)
// - Return the length if it's a slice
// - Return -1 if it's not a slice
//
// EXAMPLES:
//   SliceLength([]int{1, 2, 3}) → 3
//   SliceLength("not a slice") → -1
//
// HINT: v := reflect.ValueOf(slice)
//       Check v.Kind() == reflect.Slice
//       return v.Len()
func SliceLength(slice interface{}) int {
	// TODO: Implement this function
	return -1
}

// MapKeys returns all keys from a map as []interface{}.
//
// REQUIREMENTS:
// - Check if the value is a map (Kind() == reflect.Map)
// - Return all keys as []interface{}
// - Return nil if not a map
//
// EXAMPLES:
//   m := map[string]int{"a": 1, "b": 2}
//   MapKeys(m) → []interface{}{"a", "b"} (order may vary)
//
// HINT: v := reflect.ValueOf(m)
//       Check v.Kind() == reflect.Map
//       keys := v.MapKeys()
//       Convert each key.Interface() to []interface{}
func MapKeys(m interface{}) []interface{} {
	// TODO: Implement this function
	return nil
}

// ============================================================================
// EXERCISE 8: Creating Values
// ============================================================================

// NewInstance creates a new instance of the same type as the given value.
//
// REQUIREMENTS:
// - Use reflect.New to create a pointer to a new zero value
// - Return the pointer as interface{}
//
// EXAMPLES:
//   u := User{}
//   newU := NewInstance(u).(*User)  // newU is a *User with zero values
//
// HINT: t := reflect.TypeOf(v)
//       ptr := reflect.New(t)
//       return ptr.Interface()
func NewInstance(v interface{}) interface{} {
	// TODO: Implement this function
	return nil
}

// ============================================================================
// EXERCISE 9: Advanced - Struct Field Names
// ============================================================================

// GetFieldNames returns the names of all fields in a struct.
//
// REQUIREMENTS:
// - Return a slice of field names (strings)
// - Only include exported (public) fields
// - Return nil if the value is not a struct
//
// EXAMPLES:
//   GetFieldNames(User{}) → []string{"Name", "Email", "Age"}
//   GetFieldNames(42) → nil
//
// HINT: Iterate with NumField() and Field(i).Name
//       Check Field(i).IsExported()
func GetFieldNames(v interface{}) []string {
	// TODO: Implement this function
	return nil
}

// ============================================================================
// EXERCISE 10: Advanced - Deep Copy
// ============================================================================

// DeepCopy creates a deep copy of a struct using reflection.
//
// REQUIREMENTS:
// - Create a new instance of the same type
// - Copy all field values from the original to the new instance
// - Return a pointer to the new instance
// - Only handle simple types (int, string, bool, float64)
// - Return nil if not a struct
//
// EXAMPLES:
//   original := User{Name: "Alice", Age: 30}
//   copy := DeepCopy(original).(*User)
//   copy.Name → "Alice"
//   &copy != &original (different instances)
//
// HINT: Create new instance with reflect.New
//       Iterate fields and copy values using Field(i)
//       Use Set to copy each field value
func DeepCopy(v interface{}) interface{} {
	// TODO: Implement this function
	return nil
}
