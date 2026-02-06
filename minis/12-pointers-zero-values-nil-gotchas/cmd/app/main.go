package main

import "fmt"

// MyError is a custom error type for demonstrating the typed nil interface gotcha.
type MyError struct{}

func (e *MyError) Error() string {
	return "this is my custom error"
}

// getError returns a nil pointer to MyError, but as an `error` interface.
func getError() error {
	var err *MyError = nil
	return err
}

func main() {
	fmt.Println("--- 1. Pointers (& and *)")
	x := 100
	p := &x // p is a pointer that holds the memory address of x.

	fmt.Printf("x (value)       = %d\n", x)
	fmt.Printf("&x (address)    = %p\n", &x)
	fmt.Printf("p (pointer)     = %p\n", p)
	fmt.Printf("*p (dereference) = %d\n", *p)

	// Modify the value AT the address p.
	*p = 200
	fmt.Printf("x is now        = %d (changed via pointer)\n", x)
	fmt.Println()

	fmt.Println("--- 2. Zero Values")
	var z_int int
	var z_bool bool
	var z_str string
	var z_ptr *int
	var z_slice []int
	var z_map map[string]int
	fmt.Printf("int zero value: %d\n", z_int)
	fmt.Printf("bool zero value: %v\n", z_bool)
	fmt.Printf("string zero value: %q\n", z_str)
	fmt.Printf("pointer zero value: %v\n", z_ptr)
	fmt.Printf("slice zero value: %v (is nil: %v)\n", z_slice, z_slice == nil)
	fmt.Printf("map zero value: %v (is nil: %v)\n", z_map, z_map == nil)
	fmt.Println()

	fmt.Println("--- 3. Gotcha #1: Nil Slice vs. Nil Map")
	// A nil slice is usable.
	var nilSlice []int
	fmt.Printf("Appending to a nil slice... (len=%d, cap=%d)\n", len(nilSlice), cap(nilSlice))
	nilSlice = append(nilSlice, 42)
	fmt.Printf("Success! Slice is now: %v\n", nilSlice)

	// A nil map will panic on write.
	var nilMap map[string]int
	fmt.Println("Attempting to write to a nil map...")
	// The following line would cause a panic. We use a recover to demonstrate.
	func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("Success! Caught expected panic: %v\n", r)
			}
		}()
		// This is the line that panics.
		nilMap["key"] = 1
	}()
	fmt.Println()

	fmt.Println("--- 4. Gotcha #2: The Typed Nil Interface")
	err := getError()
	fmt.Printf("Value of err: %v\n", err)
	fmt.Printf("Type of err: %T\n", err)

	// The `err` variable holds a `*MyError` type, even though the value is nil.
	// Because the interface has a type, it is NOT considered nil itself.
	if err != nil {
		fmt.Println("`err != nil` is TRUE, which can be surprising.")
		fmt.Println("This is because the interface has a type (*MyError), even though its value is nil.")
	} else {
		fmt.Println("`err != nil` is FALSE.")
	}
}