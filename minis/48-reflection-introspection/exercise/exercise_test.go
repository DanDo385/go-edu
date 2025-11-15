package exercise

import (
	"reflect"
	"sort"
	"testing"
)

// ============================================================================
// EXERCISE 1: Type Inspection Tests
// ============================================================================

func TestGetTypeName(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected string
	}{
		{"User struct", User{}, "User"},
		{"int", 42, "int"},
		{"string", "hello", "string"},
		{"bool", true, "bool"},
		{"Product struct", Product{}, "Product"},
		{"Point struct", Point{}, "Point"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetTypeName(tt.input)
			if result != tt.expected {
				t.Errorf("GetTypeName(%v) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestGetKind(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected string
	}{
		{"struct", User{}, "struct"},
		{"int", 42, "int"},
		{"string", "hello", "string"},
		{"pointer", &User{}, "ptr"},
		{"slice", []int{1, 2, 3}, "slice"},
		{"map", map[string]int{}, "map"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetKind(tt.input)
			if result != tt.expected {
				t.Errorf("GetKind(%v) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestCountFields(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected int
	}{
		{"User has 3 fields", User{}, 3},
		{"Product has 3 fields", Product{}, 3},
		{"Point has 2 fields", Point{}, 2},
		{"int has 0 fields", 42, 0},
		{"string has 0 fields", "hello", 0},
		{"pointer to User", &User{}, 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CountFields(tt.input)
			if result != tt.expected {
				t.Errorf("CountFields(%T) = %d, want %d", tt.input, result, tt.expected)
			}
		})
	}
}

// ============================================================================
// EXERCISE 2: Struct Tags Tests
// ============================================================================

func TestGetJSONTag(t *testing.T) {
	tests := []struct {
		name      string
		input     interface{}
		fieldName string
		expected  string
	}{
		{"User.Name", User{}, "Name", "name"},
		{"User.Email", User{}, "Email", "email"},
		{"User.Age", User{}, "Age", "age"},
		{"Product.ID", Product{}, "ID", "id"},
		{"Product.Name", Product{}, "Name", "name"},
		{"non-existent field", User{}, "NoSuchField", ""},
		{"non-struct", 42, "Field", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetJSONTag(tt.input, tt.fieldName)
			if result != tt.expected {
				t.Errorf("GetJSONTag(%T, %q) = %q, want %q",
					tt.input, tt.fieldName, result, tt.expected)
			}
		})
	}
}

func TestGetAllTags(t *testing.T) {
	tests := []struct {
		name      string
		input     interface{}
		fieldName string
		expected  string
	}{
		{"User.Name", User{}, "Name", `json:"name" validate:"required"`},
		{"User.Email", User{}, "Email", `json:"email" validate:"email"`},
		{"Product.Name", Product{}, "Name", `db:"product_name" json:"name"`},
		{"non-existent", User{}, "NoSuchField", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetAllTags(tt.input, tt.fieldName)
			if result != tt.expected {
				t.Errorf("GetAllTags(%T, %q) = %q, want %q",
					tt.input, tt.fieldName, result, tt.expected)
			}
		})
	}
}

// ============================================================================
// EXERCISE 3: Value Inspection Tests
// ============================================================================

func TestGetFieldValue(t *testing.T) {
	user := User{
		Name:  "Alice",
		Email: "alice@example.com",
		Age:   30,
	}

	tests := []struct {
		name      string
		input     interface{}
		fieldName string
		expected  interface{}
	}{
		{"User.Name", user, "Name", "Alice"},
		{"User.Email", user, "Email", "alice@example.com"},
		{"User.Age", user, "Age", 30},
		{"non-existent", user, "NoSuchField", nil},
		{"non-struct", 42, "Field", nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetFieldValue(tt.input, tt.fieldName)
			if result != tt.expected {
				t.Errorf("GetFieldValue(%v, %q) = %v, want %v",
					tt.input, tt.fieldName, result, tt.expected)
			}
		})
	}
}

func TestGetFieldValues(t *testing.T) {
	user := User{
		Name:  "Bob",
		Email: "bob@example.com",
		Age:   25,
	}

	result := GetFieldValues(user)

	expected := map[string]interface{}{
		"Name":  "Bob",
		"Email": "bob@example.com",
		"Age":   25,
	}

	if len(result) != len(expected) {
		t.Errorf("GetFieldValues returned %d fields, want %d", len(result), len(expected))
	}

	for key, expectedValue := range expected {
		if result[key] != expectedValue {
			t.Errorf("GetFieldValues()[%q] = %v, want %v", key, result[key], expectedValue)
		}
	}
}

func TestGetFieldValuesNonStruct(t *testing.T) {
	result := GetFieldValues(42)
	if len(result) != 0 {
		t.Errorf("GetFieldValues(42) should return empty map, got %v", result)
	}
}

// ============================================================================
// EXERCISE 4: Value Modification Tests
// ============================================================================

func TestSetFieldValue(t *testing.T) {
	t.Run("set string field", func(t *testing.T) {
		user := &User{Name: "Alice", Email: "alice@example.com", Age: 30}
		err := SetFieldValue(user, "Name", "Bob")

		if err != nil {
			t.Errorf("SetFieldValue failed: %v", err)
		}

		if user.Name != "Bob" {
			t.Errorf("Name = %q, want %q", user.Name, "Bob")
		}
	})

	t.Run("set int field", func(t *testing.T) {
		user := &User{Name: "Alice", Email: "alice@example.com", Age: 30}
		err := SetFieldValue(user, "Age", 35)

		if err != nil {
			t.Errorf("SetFieldValue failed: %v", err)
		}

		if user.Age != 35 {
			t.Errorf("Age = %d, want %d", user.Age, 35)
		}
	})

	t.Run("error on non-pointer", func(t *testing.T) {
		user := User{Name: "Alice"}
		err := SetFieldValue(user, "Name", "Bob")

		if err == nil {
			t.Error("SetFieldValue should fail on non-pointer")
		}
	})

	t.Run("error on non-existent field", func(t *testing.T) {
		user := &User{}
		err := SetFieldValue(user, "NoSuchField", "value")

		if err == nil {
			t.Error("SetFieldValue should fail on non-existent field")
		}
	})
}

// ============================================================================
// EXERCISE 5: Dynamic Method Calls Tests
// ============================================================================

func TestCallMethod(t *testing.T) {
	t.Run("Calculator.Add", func(t *testing.T) {
		calc := Calculator{}
		results := CallMethod(calc, "Add", 5, 3)

		if len(results) != 1 {
			t.Fatalf("Expected 1 result, got %d", len(results))
		}

		if results[0] != 8 {
			t.Errorf("Add(5, 3) = %v, want 8", results[0])
		}
	})

	t.Run("Calculator.Multiply", func(t *testing.T) {
		calc := Calculator{}
		results := CallMethod(calc, "Multiply", 4, 7)

		if len(results) != 1 {
			t.Fatalf("Expected 1 result, got %d", len(results))
		}

		if results[0] != 28 {
			t.Errorf("Multiply(4, 7) = %v, want 28", results[0])
		}
	})

	t.Run("non-existent method", func(t *testing.T) {
		calc := Calculator{}
		results := CallMethod(calc, "Divide", 10, 2)

		if results != nil {
			t.Errorf("CallMethod on non-existent method should return nil, got %v", results)
		}
	})

	t.Run("Counter.Increment", func(t *testing.T) {
		counter := &Counter{Value: 0}
		CallMethod(counter, "Increment")

		if counter.Value != 1 {
			t.Errorf("After Increment, Value = %d, want 1", counter.Value)
		}
	})
}

func TestHasMethod(t *testing.T) {
	tests := []struct {
		name       string
		obj        interface{}
		methodName string
		expected   bool
	}{
		{"Calculator has Add", Calculator{}, "Add", true},
		{"Calculator has Multiply", Calculator{}, "Multiply", true},
		{"Calculator doesn't have Divide", Calculator{}, "Divide", false},
		{"Counter has Increment", &Counter{}, "Increment", true},
		{"Counter has Decrement", &Counter{}, "Decrement", true},
		{"int doesn't have Add", 42, "Add", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := HasMethod(tt.obj, tt.methodName)
			if result != tt.expected {
				t.Errorf("HasMethod(%T, %q) = %v, want %v",
					tt.obj, tt.methodName, result, tt.expected)
			}
		})
	}
}

// ============================================================================
// EXERCISE 6: Type Comparison Tests
// ============================================================================

func TestSameType(t *testing.T) {
	tests := []struct {
		name     string
		a        interface{}
		b        interface{}
		expected bool
	}{
		{"both int", 42, 10, true},
		{"int and string", 42, "hello", false},
		{"both User", User{}, User{Name: "Alice"}, true},
		{"User and Product", User{}, Product{}, false},
		{"int and *int", 42, new(int), false},
		{"both *User", &User{}, &User{Name: "Bob"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SameType(tt.a, tt.b)
			if result != tt.expected {
				t.Errorf("SameType(%T, %T) = %v, want %v",
					tt.a, tt.b, result, tt.expected)
			}
		})
	}
}

func TestIsPointer(t *testing.T) {
	x := 42
	user := User{}

	tests := []struct {
		name     string
		input    interface{}
		expected bool
	}{
		{"int value", 42, false},
		{"int pointer", &x, true},
		{"User value", user, false},
		{"User pointer", &user, true},
		{"string value", "hello", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsPointer(tt.input)
			if result != tt.expected {
				t.Errorf("IsPointer(%T) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

// ============================================================================
// EXERCISE 7: Slice and Map Operations Tests
// ============================================================================

func TestSliceLength(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected int
	}{
		{"int slice length 3", []int{1, 2, 3}, 3},
		{"empty slice", []string{}, 0},
		{"string slice", []string{"a", "b", "c", "d"}, 4},
		{"not a slice", 42, -1},
		{"string is not slice", "hello", -1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SliceLength(tt.input)
			if result != tt.expected {
				t.Errorf("SliceLength(%v) = %d, want %d", tt.input, result, tt.expected)
			}
		})
	}
}

func TestMapKeys(t *testing.T) {
	t.Run("string keys", func(t *testing.T) {
		m := map[string]int{
			"a": 1,
			"b": 2,
			"c": 3,
		}

		keys := MapKeys(m)

		if len(keys) != 3 {
			t.Fatalf("Expected 3 keys, got %d", len(keys))
		}

		// Convert to string slice for easier comparison
		var strKeys []string
		for _, k := range keys {
			strKeys = append(strKeys, k.(string))
		}
		sort.Strings(strKeys)

		expected := []string{"a", "b", "c"}
		for i, k := range expected {
			if strKeys[i] != k {
				t.Errorf("Key %d: got %q, want %q", i, strKeys[i], k)
			}
		}
	})

	t.Run("empty map", func(t *testing.T) {
		m := map[string]int{}
		keys := MapKeys(m)

		if len(keys) != 0 {
			t.Errorf("Expected 0 keys for empty map, got %d", len(keys))
		}
	})

	t.Run("not a map", func(t *testing.T) {
		keys := MapKeys(42)
		if keys != nil {
			t.Errorf("MapKeys on non-map should return nil, got %v", keys)
		}
	})
}

// ============================================================================
// EXERCISE 8: Creating Values Tests
// ============================================================================

func TestNewInstance(t *testing.T) {
	t.Run("User instance", func(t *testing.T) {
		original := User{Name: "Alice", Age: 30}
		newInst := NewInstance(original)

		// Should be a pointer to User
		userPtr, ok := newInst.(*User)
		if !ok {
			t.Fatalf("NewInstance should return *User, got %T", newInst)
		}

		// Should be zero-valued (not a copy of original)
		if userPtr.Name != "" || userPtr.Age != 0 {
			t.Errorf("NewInstance should return zero-valued struct, got %+v", userPtr)
		}

		// Should be a different instance
		if userPtr == &original {
			t.Error("NewInstance should return a different instance")
		}
	})

	t.Run("Point instance", func(t *testing.T) {
		original := Point{X: 10, Y: 20}
		newInst := NewInstance(original)

		pointPtr, ok := newInst.(*Point)
		if !ok {
			t.Fatalf("NewInstance should return *Point, got %T", newInst)
		}

		if pointPtr.X != 0 || pointPtr.Y != 0 {
			t.Errorf("NewInstance should return zero-valued struct, got %+v", pointPtr)
		}
	})
}

// ============================================================================
// EXERCISE 9: Advanced - Struct Field Names Tests
// ============================================================================

func TestGetFieldNames(t *testing.T) {
	t.Run("User fields", func(t *testing.T) {
		result := GetFieldNames(User{})
		expected := []string{"Name", "Email", "Age"}

		if len(result) != len(expected) {
			t.Fatalf("Expected %d fields, got %d", len(expected), len(result))
		}

		for i, name := range expected {
			if result[i] != name {
				t.Errorf("Field %d: got %q, want %q", i, result[i], name)
			}
		}
	})

	t.Run("Product fields", func(t *testing.T) {
		result := GetFieldNames(Product{})
		expected := []string{"ID", "Name", "Price"}

		if len(result) != len(expected) {
			t.Fatalf("Expected %d fields, got %d", len(expected), len(result))
		}

		for i, name := range expected {
			if result[i] != name {
				t.Errorf("Field %d: got %q, want %q", i, result[i], name)
			}
		}
	})

	t.Run("non-struct", func(t *testing.T) {
		result := GetFieldNames(42)
		if result != nil {
			t.Errorf("GetFieldNames on non-struct should return nil, got %v", result)
		}
	})

	t.Run("pointer to struct", func(t *testing.T) {
		result := GetFieldNames(&User{})
		expected := []string{"Name", "Email", "Age"}

		if len(result) != len(expected) {
			t.Fatalf("Expected %d fields, got %d", len(expected), len(result))
		}
	})
}

// ============================================================================
// EXERCISE 10: Advanced - Deep Copy Tests
// ============================================================================

func TestDeepCopy(t *testing.T) {
	t.Run("User copy", func(t *testing.T) {
		original := User{
			Name:  "Alice",
			Email: "alice@example.com",
			Age:   30,
		}

		copyPtr := DeepCopy(original)
		if copyPtr == nil {
			t.Fatal("DeepCopy returned nil")
		}

		copy, ok := copyPtr.(*User)
		if !ok {
			t.Fatalf("DeepCopy should return *User, got %T", copyPtr)
		}

		// Values should match
		if copy.Name != original.Name {
			t.Errorf("Name: got %q, want %q", copy.Name, original.Name)
		}
		if copy.Email != original.Email {
			t.Errorf("Email: got %q, want %q", copy.Email, original.Email)
		}
		if copy.Age != original.Age {
			t.Errorf("Age: got %d, want %d", copy.Age, original.Age)
		}

		// Should be different instances
		if copy == &original {
			t.Error("DeepCopy should return a different instance")
		}

		// Modifying copy shouldn't affect original
		copy.Name = "Bob"
		if original.Name == "Bob" {
			t.Error("Modifying copy affected original")
		}
	})

	t.Run("Point copy", func(t *testing.T) {
		original := Point{X: 10, Y: 20}
		copyPtr := DeepCopy(original)

		copy, ok := copyPtr.(*Point)
		if !ok {
			t.Fatalf("DeepCopy should return *Point, got %T", copyPtr)
		}

		if copy.X != original.X || copy.Y != original.Y {
			t.Errorf("Copy values don't match: got %+v, want %+v", copy, original)
		}
	})

	t.Run("non-struct", func(t *testing.T) {
		result := DeepCopy(42)
		if result != nil {
			t.Errorf("DeepCopy on non-struct should return nil, got %v", result)
		}
	})
}

// ============================================================================
// Benchmark Tests - Demonstrating Reflection Overhead
// ============================================================================

func BenchmarkDirectFieldAccess(b *testing.B) {
	user := User{Name: "Alice", Email: "alice@example.com", Age: 30}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = user.Name
		_ = user.Email
		_ = user.Age
	}
}

func BenchmarkReflectionFieldAccess(b *testing.B) {
	user := User{Name: "Alice", Email: "alice@example.com", Age: 30}
	v := reflect.ValueOf(user)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = v.FieldByName("Name").String()
		_ = v.FieldByName("Email").String()
		_ = v.FieldByName("Age").Int()
	}
}

func BenchmarkDirectMethodCall(b *testing.B) {
	calc := Calculator{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = calc.Add(5, 3)
	}
}

func BenchmarkReflectionMethodCall(b *testing.B) {
	calc := Calculator{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = CallMethod(calc, "Add", 5, 3)
	}
}
