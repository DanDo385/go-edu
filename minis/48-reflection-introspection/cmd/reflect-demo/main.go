// Package main demonstrates Go's reflection and introspection capabilities.
//
// MACRO-COMMENT: What This Program Teaches
// =========================================
// This program provides hands-on examples of:
// 1. Type inspection (TypeOf, Kind, Fields, Methods)
// 2. Struct tags (reading metadata from struct fields)
// 3. Value inspection (ValueOf, getting field values)
// 4. Value modification (setting fields through reflection)
// 5. Dynamic method calls (calling methods by name)
// 6. Performance implications (reflection vs direct access)
// 7. Common pitfalls (unaddressable values, unexported fields)
//
// COMPILER BEHAVIOR: Reflection and Type Information
// ====================================================
// The Go compiler embeds type metadata in the binary for all types.
// This metadata includes:
// - Type names and kinds
// - Struct field names and types
// - Method names and signatures
// - Struct tags
//
// The reflect package accesses this metadata at runtime.
// This adds to binary size but enables powerful runtime introspection.
//
// RUNTIME BEHAVIOR: Performance Cost
// ===================================
// Reflection operations are 100-300x slower than direct access because:
// 1. Runtime type checking (compiler can't verify types)
// 2. Interface conversions (values wrapped in interfaces)
// 3. No inlining (compiler can't optimize reflection code)
// 4. Heap allocations (reflection often escapes to heap)
//
// Use reflection sparingly, typically only in library code or startup code.

package main

import (
	"fmt"
	"reflect"
	"strings"
	"time"
)

// ============================================================================
// SECTION 1: Type Inspection Basics
// ============================================================================

// Person demonstrates struct reflection with tags.
//
// MACRO-COMMENT: Struct Tags
// Struct tags are string literals that store metadata about fields.
// They're commonly used for:
// - JSON/XML marshaling (json:"name")
// - Database mapping (db:"user_name")
// - Validation rules (validate:"required,min=1")
// - ORM relationships (gorm:"foreignKey:UserID")
//
// Format: `key:"value" key2:"value2"`
// Access: field.Tag.Get("key")
type Person struct {
	Name  string `json:"name" db:"user_name" validate:"required"`
	Age   int    `json:"age" db:"age" validate:"min=0,max=150"`
	Email string `json:"email,omitempty" db:"email" validate:"email"`
}

// demonstrateTypeInspection shows how to inspect types at runtime.
func demonstrateTypeInspection() {
	fmt.Println("=== Type Inspection: Understanding Types at Runtime ===")

	// MICRO-COMMENT: Create a value to inspect
	p := Person{Name: "Alice", Age: 30, Email: "alice@example.com"}

	// MACRO-COMMENT: reflect.TypeOf - The Gateway to Type Information
	// TypeOf returns a reflect.Type that describes the type of the value.
	// This includes:
	// - The type's name ("Person", "int", "string")
	// - The type's kind (struct, int, slice, etc.)
	// - For structs: field information
	// - For interfaces: method information
	t := reflect.TypeOf(p)

	// MICRO-COMMENT: Basic type information
	fmt.Printf("Type name: %s\n", t.Name())       // "Person"
	fmt.Printf("Package path: %s\n", t.PkgPath()) // "main"
	fmt.Printf("Kind: %s\n", t.Kind())            // "struct"

	// MACRO-COMMENT: Type vs Kind - A Critical Distinction
	// - Type is the specific type: "Person", "time.Duration", "MyInt"
	// - Kind is the underlying category: struct, int, slice, map, etc.
	//
	// Example:
	//   type MyInt int
	//   x := MyInt(42)
	//   TypeOf(x).Name() → "MyInt"  (Type)
	//   TypeOf(x).Kind() → "int"    (Kind)
	type MyInt int
	mi := MyInt(42)
	miType := reflect.TypeOf(mi)
	fmt.Printf("\nMyInt Type: %s, Kind: %s\n", miType.Name(), miType.Kind())

	// MACRO-COMMENT: Inspecting Struct Fields
	// For struct types, we can iterate over fields and access metadata.
	// This is how JSON marshalers, ORMs, and validators work.
	fmt.Printf("\nNumber of fields: %d\n", t.NumField())

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		// MICRO-COMMENT: Field provides complete information about a struct field:
		// - Name: the field's identifier
		// - Type: the field's type (as reflect.Type)
		// - Tag: the struct tag (for metadata)
		// - Offset: memory offset (for unsafe operations)
		// - Anonymous: whether it's an embedded field
		fmt.Printf("\nField %d:\n", i)
		fmt.Printf("  Name: %s\n", field.Name)
		fmt.Printf("  Type: %s\n", field.Type)
		fmt.Printf("  Tag: %s\n", field.Tag)

		// MACRO-COMMENT: Reading Struct Tags
		// Tags are parsed as space-separated key:"value" pairs.
		// Use Tag.Get(key) to extract a specific tag.
		// Returns empty string if the tag doesn't exist.
		jsonTag := field.Tag.Get("json")
		dbTag := field.Tag.Get("db")
		validateTag := field.Tag.Get("validate")

		if jsonTag != "" {
			fmt.Printf("  JSON tag: %s\n", jsonTag)
		}
		if dbTag != "" {
			fmt.Printf("  DB tag: %s\n", dbTag)
		}
		if validateTag != "" {
			fmt.Printf("  Validate tag: %s\n", validateTag)
		}
	}

	// MICRO-COMMENT: You can also access fields by name
	nameField, found := t.FieldByName("Name")
	if found {
		fmt.Printf("\nDirect field lookup: %s (%s)\n", nameField.Name, nameField.Type)
	}

	fmt.Println()
}

// ============================================================================
// SECTION 2: Value Inspection
// ============================================================================

// demonstrateValueInspection shows how to read values through reflection.
func demonstrateValueInspection() {
	fmt.Println("=== Value Inspection: Reading Values at Runtime ===")

	// MICRO-COMMENT: Create a person with specific values
	p := Person{Name: "Bob", Age: 25, Email: "bob@example.com"}

	// MACRO-COMMENT: reflect.ValueOf - Accessing Values
	// ValueOf returns a reflect.Value that holds the value and knows its type.
	// You can extract the actual value using type-specific methods:
	// - v.Int() for integers
	// - v.String() for strings
	// - v.Bool() for booleans
	// - v.Interface() to get as interface{}
	v := reflect.ValueOf(p)

	// MICRO-COMMENT: You can get the type from a Value too
	fmt.Printf("Type: %s\n", v.Type())
	fmt.Printf("Kind: %s\n", v.Kind())

	// MACRO-COMMENT: Iterating Over Field Values
	// For struct values, we can access individual field values.
	// This is how JSON marshalers read your struct data.
	fmt.Println("\nField values:")

	for i := 0; i < v.NumField(); i++ {
		fieldValue := v.Field(i)
		fieldType := v.Type().Field(i)

		// MICRO-COMMENT: Use .Interface() to get the value as interface{}
		// Then you can type assert or print it
		fmt.Printf("  %s = %v (type: %s)\n",
			fieldType.Name,
			fieldValue.Interface(),
			fieldValue.Type())
	}

	// MACRO-COMMENT: Type-Specific Value Extraction
	// For known types, use the typed getters for safety and performance.
	nameField := v.FieldByName("Name")
	if nameField.Kind() == reflect.String {
		fmt.Printf("\nName (string method): %s\n", nameField.String())
	}

	ageField := v.FieldByName("Age")
	if ageField.Kind() == reflect.Int {
		fmt.Printf("Age (int method): %d\n", ageField.Int())
	}

	// MICRO-COMMENT: Checking if a value is zero
	if v.IsZero() {
		fmt.Println("Value is zero (default)")
	} else {
		fmt.Println("Value is non-zero")
	}

	fmt.Println()
}

// ============================================================================
// SECTION 3: Value Modification
// ============================================================================

// demonstrateValueModification shows how to change values through reflection.
func demonstrateValueModification() {
	fmt.Println("=== Value Modification: Changing Values at Runtime ===")

	// MACRO-COMMENT: The Addressability Problem
	// You can only modify values that are "addressable".
	// A value is addressable if you can take its address with &.
	//
	// WRONG: v := reflect.ValueOf(person)
	// This creates a reflect.Value of a COPY of person.
	// Modifying the copy doesn't affect the original.
	//
	// RIGHT: v := reflect.ValueOf(&person).Elem()
	// This creates a reflect.Value of a POINTER to person,
	// then .Elem() gives us the addressable value being pointed to.

	p := Person{Name: "Charlie", Age: 35, Email: "charlie@example.com"}
	fmt.Printf("Before: %+v\n", p)

	// MICRO-COMMENT: Get addressable value
	// Step 1: Pass pointer to ValueOf
	// Step 2: Use Elem() to get the value being pointed to
	v := reflect.ValueOf(&p).Elem()

	// MACRO-COMMENT: Checking Addressability
	// Always check CanSet() before attempting to modify.
	// This prevents panics from trying to set unaddressable values.
	if !v.CanSet() {
		fmt.Println("ERROR: Value is not settable!")
		return
	}

	// MICRO-COMMENT: Modify the Name field
	nameField := v.FieldByName("Name")
	if nameField.CanSet() && nameField.Kind() == reflect.String {
		nameField.SetString("David")
		fmt.Printf("After changing name: %+v\n", p)
	}

	// MICRO-COMMENT: Modify the Age field
	ageField := v.FieldByName("Age")
	if ageField.CanSet() && ageField.Kind() == reflect.Int {
		ageField.SetInt(40)
		fmt.Printf("After changing age: %+v\n", p)
	}

	// MACRO-COMMENT: Set Methods by Type
	// Different types have different Set methods:
	// - SetBool(bool)
	// - SetInt(int64)
	// - SetUint(uint64)
	// - SetFloat(float64)
	// - SetString(string)
	// - Set(reflect.Value) - for any type
	//
	// Using the wrong one causes a panic!

	// MACRO-COMMENT: Setting with interface{}
	// For generic setting, use Set() with ValueOf(value)
	emailField := v.FieldByName("Email")
	if emailField.CanSet() {
		emailField.Set(reflect.ValueOf("david@example.com"))
		fmt.Printf("After changing email: %+v\n", p)
	}

	fmt.Println()
}

// ============================================================================
// SECTION 4: Unexported Fields
// ============================================================================

// Config demonstrates unexported field reflection.
type Config struct {
	PublicField  string // Exported - can read and set
	privateField string // Unexported - can read but NOT set
}

// demonstrateUnexportedFields shows the limitations with unexported fields.
func demonstrateUnexportedFields() {
	fmt.Println("=== Unexported Fields: Reflection Respects Visibility ===")

	cfg := Config{
		PublicField:  "public",
		privateField: "private",
	}

	v := reflect.ValueOf(&cfg).Elem()

	// MACRO-COMMENT: Reading Unexported Fields
	// You CAN see unexported fields through reflection,
	// but you CANNOT call .Interface() on them (it will panic).
	// You can still inspect their type and check if they're settable.
	publicField := v.FieldByName("PublicField")
	privateField := v.FieldByName("privateField")

	fmt.Printf("Public field value: %v\n", publicField.Interface())

	// MICRO-COMMENT: For unexported fields, we can't call .Interface()
	// But we can still see the field exists and check its properties
	fmt.Printf("Private field exists: %v\n", privateField.IsValid())
	fmt.Printf("Private field type: %v\n", privateField.Type())
	// privateField.Interface() would PANIC here!

	// MACRO-COMMENT: Setting Unexported Fields
	// You CANNOT set unexported fields, even with reflection.
	// This preserves Go's encapsulation rules.
	fmt.Printf("Can set public field? %v\n", publicField.CanSet())   // true
	fmt.Printf("Can set private field? %v\n", privateField.CanSet()) // false

	if publicField.CanSet() {
		publicField.SetString("modified")
		fmt.Printf("After modification: %+v\n", cfg)
	}

	// MICRO-COMMENT: Attempting to set a private field would panic:
	// privateField.SetString("won't work")  // PANIC!

	fmt.Println()
}

// ============================================================================
// SECTION 5: Method Inspection and Dynamic Calls
// ============================================================================

// Counter demonstrates method reflection.
type Counter struct {
	count int
}

// Increment adds 1 to the counter.
func (c *Counter) Increment() {
	c.count++
}

// Add adds a specific amount to the counter.
func (c *Counter) Add(amount int) {
	c.count += amount
}

// Get returns the current count.
func (c *Counter) Get() int {
	return c.count
}

// demonstrateMethodCalls shows dynamic method invocation.
func demonstrateMethodCalls() {
	fmt.Println("=== Dynamic Method Calls: Calling Methods by Name ===")

	counter := &Counter{count: 0}
	v := reflect.ValueOf(counter)
	t := reflect.TypeOf(counter)

	// MACRO-COMMENT: Discovering Methods
	// Methods are discovered at runtime, similar to fields.
	// The method set depends on whether you have a value or pointer.
	// Methods with pointer receivers are ONLY in the pointer's method set.
	fmt.Printf("Number of methods: %d\n", t.NumMethod())
	fmt.Println("Methods:")

	for i := 0; i < t.NumMethod(); i++ {
		method := t.Method(i)
		fmt.Printf("  %d. %s\n", i+1, method.Name)
	}

	// MACRO-COMMENT: Calling Methods Dynamically
	// Use MethodByName to get a method, then Call it with arguments.
	// Arguments must be provided as []reflect.Value.
	// Returns are also []reflect.Value.

	// MICRO-COMMENT: Call Increment (no arguments)
	fmt.Println("\nCalling Increment()...")
	incrementMethod := v.MethodByName("Increment")
	if incrementMethod.IsValid() {
		// Call with no arguments (empty slice)
		incrementMethod.Call([]reflect.Value{})
	}

	// MICRO-COMMENT: Call Add (with argument)
	fmt.Println("Calling Add(5)...")
	addMethod := v.MethodByName("Add")
	if addMethod.IsValid() {
		// Prepare arguments as []reflect.Value
		args := []reflect.Value{
			reflect.ValueOf(5),
		}
		addMethod.Call(args)
	}

	// MICRO-COMMENT: Call Get (returns a value)
	fmt.Println("Calling Get()...")
	getMethod := v.MethodByName("Get")
	if getMethod.IsValid() {
		// Call returns []reflect.Value
		results := getMethod.Call([]reflect.Value{})

		// Extract the first (only) return value
		if len(results) > 0 {
			count := results[0].Int()
			fmt.Printf("Counter value: %d\n", count)
		}
	}

	// MACRO-COMMENT: Checking Method Signatures
	// Before calling a method, you should validate:
	// 1. It exists (IsValid)
	// 2. It has the right number of parameters
	// 3. The parameter types match
	// 4. The return types match your expectations
	fmt.Println("\nMethod signature inspection:")
	addType := addMethod.Type()
	fmt.Printf("Add method:\n")
	fmt.Printf("  Input parameters: %d\n", addType.NumIn())
	fmt.Printf("  Output parameters: %d\n", addType.NumOut())

	if addType.NumIn() > 0 {
		fmt.Printf("  First parameter type: %s\n", addType.In(0))
	}

	fmt.Println()
}

// ============================================================================
// SECTION 6: Working with Slices and Maps
// ============================================================================

// demonstrateCollections shows reflection with slices and maps.
func demonstrateCollections() {
	fmt.Println("=== Collections: Slices and Maps ===")

	// MACRO-COMMENT: Slice Reflection
	// Slices can be inspected and modified through reflection.
	numbers := []int{1, 2, 3, 4, 5}
	v := reflect.ValueOf(numbers)

	fmt.Printf("Slice type: %s\n", v.Type())
	fmt.Printf("Slice kind: %s\n", v.Kind())
	fmt.Printf("Slice length: %d\n", v.Len())
	fmt.Printf("Slice capacity: %d\n", v.Cap())

	// MICRO-COMMENT: Access elements by index
	fmt.Println("\nSlice elements:")
	for i := 0; i < v.Len(); i++ {
		elem := v.Index(i)
		fmt.Printf("  [%d] = %v\n", i, elem.Interface())
	}

	// MACRO-COMMENT: Map Reflection
	// Maps require special handling because their keys can be any comparable type.
	users := map[string]int{
		"Alice": 30,
		"Bob":   25,
		"Carol": 35,
	}

	mv := reflect.ValueOf(users)
	fmt.Printf("\nMap type: %s\n", mv.Type())
	fmt.Printf("Map kind: %s\n", mv.Kind())
	fmt.Printf("Map length: %d\n", mv.Len())

	// MICRO-COMMENT: Iterate over map keys
	// MapKeys() returns all keys as []reflect.Value
	fmt.Println("\nMap entries:")
	for _, key := range mv.MapKeys() {
		value := mv.MapIndex(key)
		fmt.Printf("  %v = %v\n", key.Interface(), value.Interface())
	}

	// MACRO-COMMENT: Setting Map Values
	// Maps obtained through ValueOf are not addressable,
	// but you CAN modify the map itself (add/remove/update entries)
	// because maps are reference types.
	mv.SetMapIndex(reflect.ValueOf("David"), reflect.ValueOf(40))
	fmt.Printf("\nAfter adding David: %v\n", users)

	// MICRO-COMMENT: Delete a map entry by setting to zero Value
	mv.SetMapIndex(reflect.ValueOf("Bob"), reflect.Value{})
	fmt.Printf("After deleting Bob: %v\n", users)

	fmt.Println()
}

// ============================================================================
// SECTION 7: Creating New Values
// ============================================================================

// demonstrateNewValues shows how to create values dynamically.
func demonstrateNewValues() {
	fmt.Println("=== Creating New Values: Dynamic Construction ===")

	// MACRO-COMMENT: Creating Values by Type
	// You can create new instances of types dynamically,
	// which is useful for deserializers, dependency injection, etc.

	// MICRO-COMMENT: Create a new Person
	personType := reflect.TypeOf(Person{})
	newPersonValue := reflect.New(personType) // Returns *Person as reflect.Value

	fmt.Printf("Created type: %s (kind: %s)\n",
		newPersonValue.Type(),
		newPersonValue.Kind()) // "ptr"

	// MICRO-COMMENT: Access the struct being pointed to
	personValue := newPersonValue.Elem()
	fmt.Printf("Dereferenced type: %s (kind: %s)\n",
		personValue.Type(),
		personValue.Kind()) // "struct"

	// MICRO-COMMENT: Set fields on the new value
	personValue.FieldByName("Name").SetString("Eve")
	personValue.FieldByName("Age").SetInt(28)
	personValue.FieldByName("Email").SetString("eve@example.com")

	// MICRO-COMMENT: Convert back to concrete type
	person := newPersonValue.Interface().(*Person)
	fmt.Printf("Created person: %+v\n", person)

	// MACRO-COMMENT: Creating Slices
	// reflect.MakeSlice creates a new slice with specified length and capacity
	intSliceType := reflect.TypeOf([]int{})
	newSlice := reflect.MakeSlice(intSliceType, 0, 5)

	// Append elements
	newSlice = reflect.Append(newSlice, reflect.ValueOf(10))
	newSlice = reflect.Append(newSlice, reflect.ValueOf(20))
	newSlice = reflect.Append(newSlice, reflect.ValueOf(30))

	fmt.Printf("\nCreated slice: %v\n", newSlice.Interface())

	// MACRO-COMMENT: Creating Maps
	// reflect.MakeMap creates a new map
	mapType := reflect.TypeOf(map[string]int{})
	newMap := reflect.MakeMap(mapType)

	// Set entries
	newMap.SetMapIndex(reflect.ValueOf("x"), reflect.ValueOf(100))
	newMap.SetMapIndex(reflect.ValueOf("y"), reflect.ValueOf(200))

	fmt.Printf("Created map: %v\n", newMap.Interface())

	fmt.Println()
}

// ============================================================================
// SECTION 8: Performance Comparison
// ============================================================================

// demonstratePerformance shows the cost of reflection.
func demonstratePerformance() {
	fmt.Println("=== Performance: The Cost of Flexibility ===")

	p := Person{Name: "Test", Age: 30, Email: "test@example.com"}

	// MACRO-COMMENT: Benchmarking Direct Access vs Reflection
	// This demonstrates why reflection should be used sparingly.
	iterations := 1000000

	// MICRO-COMMENT: Direct field access (fast)
	start := time.Now()
	for i := 0; i < iterations; i++ {
		_ = p.Name
		_ = p.Age
	}
	directDuration := time.Since(start)

	// MICRO-COMMENT: Reflection-based access (slow)
	v := reflect.ValueOf(p)
	start = time.Now()
	for i := 0; i < iterations; i++ {
		_ = v.FieldByName("Name").String()
		_ = v.FieldByName("Age").Int()
	}
	reflectDuration := time.Since(start)

	fmt.Printf("Direct access (%d iterations): %v\n", iterations, directDuration)
	fmt.Printf("Reflection access (%d iterations): %v\n", iterations, reflectDuration)

	slowdown := float64(reflectDuration) / float64(directDuration)
	fmt.Printf("Reflection is ~%.0fx slower\n", slowdown)

	// MACRO-COMMENT: Why the Difference?
	fmt.Println("\nWhy reflection is slower:")
	fmt.Println("  - Runtime type checking (compiler can't verify)")
	fmt.Println("  - Interface conversions (wrapping/unwrapping)")
	fmt.Println("  - No inlining (compiler can't optimize)")
	fmt.Println("  - Potential heap allocations (escape analysis)")
	fmt.Println("\nWhen to use reflection anyway:")
	fmt.Println("  - JSON/XML marshaling (I/O bound, not CPU bound)")
	fmt.Println("  - Database ORMs (network latency >> reflection cost)")
	fmt.Println("  - Plugin systems (flexibility > performance)")
	fmt.Println("  - Testing frameworks (not in hot path)")

	fmt.Println()
}

// ============================================================================
// SECTION 9: Practical Example - Simple JSON Encoder
// ============================================================================

// SimpleJSONEncode demonstrates a basic JSON encoder using reflection.
//
// MACRO-COMMENT: Real-World Application
// This shows how libraries like encoding/json work internally.
// The actual JSON encoder is more complex (handles edge cases, escaping, etc.),
// but the core idea is the same: use reflection to discover struct fields
// and their tags, then encode accordingly.
func SimpleJSONEncode(v interface{}) string {
	val := reflect.ValueOf(v)
	typ := val.Type()

	// MICRO-COMMENT: Only handle structs in this simple example
	if typ.Kind() != reflect.Struct {
		return fmt.Sprintf("%v", v)
	}

	var fields []string

	// MICRO-COMMENT: Iterate over fields
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		fieldValue := val.Field(i)

		// MICRO-COMMENT: Skip unexported fields
		if !field.IsExported() {
			continue
		}

		// MICRO-COMMENT: Read the json tag
		jsonName := field.Tag.Get("json")
		if jsonName == "" {
			jsonName = field.Name
		}

		// MICRO-COMMENT: Handle omitempty
		// In real JSON encoder, this would check if value is zero
		if strings.Contains(jsonName, ",omitempty") {
			parts := strings.Split(jsonName, ",")
			jsonName = parts[0]
			if fieldValue.IsZero() {
				continue // Skip zero values with omitempty
			}
		}

		// MICRO-COMMENT: Encode the value based on kind
		var encoded string
		switch fieldValue.Kind() {
		case reflect.String:
			encoded = fmt.Sprintf(`"%s"`, fieldValue.String())
		case reflect.Int, reflect.Int64:
			encoded = fmt.Sprintf("%d", fieldValue.Int())
		case reflect.Bool:
			encoded = fmt.Sprintf("%t", fieldValue.Bool())
		default:
			encoded = fmt.Sprintf(`"%v"`, fieldValue.Interface())
		}

		fields = append(fields, fmt.Sprintf(`"%s": %s`, jsonName, encoded))
	}

	return "{" + strings.Join(fields, ", ") + "}"
}

// demonstrateJSONEncoder shows our simple JSON encoder in action.
func demonstrateJSONEncoder() {
	fmt.Println("=== Practical Example: Simple JSON Encoder ===")

	// MICRO-COMMENT: Test with a Person struct
	p := Person{
		Name:  "Alice",
		Age:   30,
		Email: "alice@example.com",
	}

	json := SimpleJSONEncode(p)
	fmt.Printf("JSON encoded: %s\n", json)

	// MICRO-COMMENT: Test with omitempty
	p2 := Person{
		Name: "Bob",
		Age:  25,
		// Email omitted - should be skipped due to omitempty tag
	}

	json2 := SimpleJSONEncode(p2)
	fmt.Printf("JSON with omitempty: %s\n", json2)

	fmt.Println()
}

// ============================================================================
// SECTION 10: Common Pitfalls and How to Avoid Them
// ============================================================================

// demonstratePitfalls shows common reflection mistakes.
func demonstratePitfalls() {
	fmt.Println("=== Common Pitfalls: Learning from Mistakes ===")

	// PITFALL 1: Forgetting to check CanSet
	fmt.Println("Pitfall 1: Not checking CanSet")
	x := 42
	v := reflect.ValueOf(x) // Value, not pointer!
	fmt.Printf("  CanSet: %v (need to pass pointer!)\n", v.CanSet())

	v2 := reflect.ValueOf(&x).Elem()
	fmt.Printf("  With pointer.Elem(), CanSet: %v\n", v2.CanSet())

	// PITFALL 2: Using wrong type method
	fmt.Println("\nPitfall 2: Using wrong type accessor")
	s := "hello"
	sv := reflect.ValueOf(s)
	fmt.Printf("  String value: %s\n", sv.String()) // ✅ Correct

	// sv.Int() would panic!
	fmt.Printf("  Calling .Int() on string would panic!\n")

	// PITFALL 3: Nil interface checks
	fmt.Println("\nPitfall 3: Nil pointer in interface")
	var p *Person = nil
	v3 := reflect.ValueOf(p)
	fmt.Printf("  Value.IsValid(): %v\n", v3.IsValid()) // true (it's a valid *Person)
	fmt.Printf("  Value.IsNil(): %v\n", v3.IsNil())     // true (but the pointer is nil)

	// v3.Elem() would panic because p is nil!
	fmt.Printf("  Must check IsNil() before calling Elem()!\n")

	// PITFALL 4: Forgetting FieldByName can fail
	fmt.Println("\nPitfall 4: FieldByName returns zero Value if not found")
	p4 := Person{Name: "Test", Age: 30}
	v4 := reflect.ValueOf(&p4).Elem()
	badField := v4.FieldByName("NonExistent")
	fmt.Printf("  Field.IsValid(): %v (field doesn't exist)\n", badField.IsValid())

	if badField.IsValid() {
		fmt.Println("  Field found")
	} else {
		fmt.Println("  Field NOT found - must check IsValid()!")
	}

	fmt.Println()
}

// ============================================================================
// MAIN FUNCTION: Orchestrates All Demonstrations
// ============================================================================

// main executes all demonstration functions in order.
//
// MACRO-COMMENT: Learning Progression
// ====================================
// The demos build on each other:
// 1. Type inspection (understanding types)
// 2. Value inspection (reading values)
// 3. Value modification (changing values)
// 4. Unexported fields (visibility limits)
// 5. Method calls (dynamic invocation)
// 6. Collections (slices and maps)
// 7. Creating values (dynamic construction)
// 8. Performance (understanding the cost)
// 9. JSON encoder (practical application)
// 10. Pitfalls (avoiding common mistakes)
//
// AFTER RUNNING THIS:
// You should understand:
// - How JSON marshalers work
// - How ORMs map structs to databases
// - When reflection is appropriate
// - How to avoid reflection pitfalls
// - The performance implications
func main() {
	demonstrateTypeInspection()
	demonstrateValueInspection()
	demonstrateValueModification()
	demonstrateUnexportedFields()
	demonstrateMethodCalls()
	demonstrateCollections()
	demonstrateNewValues()
	demonstrateJSONEncoder()
	demonstratePerformance()
	demonstratePitfalls()

	// MACRO-COMMENT: Key Insights
	// ============================
	fmt.Println("=== Summary: When to Use Reflection ===")
	fmt.Println("\n✅ Good use cases:")
	fmt.Println("  - Serialization/deserialization (JSON, XML, etc.)")
	fmt.Println("  - Database ORMs (mapping structs to tables)")
	fmt.Println("  - Dependency injection frameworks")
	fmt.Println("  - Testing frameworks (test discovery)")
	fmt.Println("  - RPC frameworks (method dispatch)")
	fmt.Println("  - Validation frameworks (checking constraints)")

	fmt.Println("\n❌ Avoid reflection for:")
	fmt.Println("  - Hot paths (inner loops)")
	fmt.Println("  - Simple configuration (use direct access)")
	fmt.Println("  - Type switches (use interface type assertions)")
	fmt.Println("  - Generic code (use Go 1.18+ generics instead)")

	fmt.Println("\n💡 Best practices:")
	fmt.Println("  - Cache reflection results (don't reflect in loops)")
	fmt.Println("  - Always check CanSet, IsValid, IsNil")
	fmt.Println("  - Validate types before type assertions")
	fmt.Println("  - Consider code generation as an alternative")
	fmt.Println("  - Profile before optimizing (measure, don't guess)")
}
