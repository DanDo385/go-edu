//go:build !solution && !reference

// TODO:
// - Read the tests in exercise_test.go to understand expected behavior.
// - Implement the exported API in this file.
// - Compare with the fully-commented reference in solution.reference.go (go test -tags=reference ./...).
package interfacesducktyping
// TODO: implement String.
func (p Person) String() string { panic("TODO: implement") }
// TODO: implement GetAge.
func GetAge(s Stringer) (int, bool) { panic("TODO: implement") }
// TODO: implement DescribeType.
func DescribeType(i interface{}) string { panic("TODO: implement") }
// TODO: implement IsValidEmail.
func IsValidEmail(v Validator) bool { panic("TODO: implement") }
// TODO: implement Read.
func (b *Buffer) Read() string { panic("TODO: implement") }
// TODO: implement Write.
func (b *Buffer) Write(data string) error { panic("TODO: implement") }
// TODO: implement IsReadWriter.
func IsReadWriter(i interface{}) bool { panic("TODO: implement") }
// TODO: implement Increment.
func (c *Counter) Increment() { panic("TODO: implement") }
// TODO: implement CanIncrement.
func CanIncrement(i interface{}) bool { panic("TODO: implement") }
// TODO: implement CountTypes.
func CountTypes(values []interface{}) map[string]int { panic("TODO: implement") }
// TODO: implement Error.
func (e ValidationError) Error() string { panic("TODO: implement") }
// TODO: implement Area.
func (r Rectangle) Area() float64 { panic("TODO: implement") }
// TODO: implement Area.
func (c Circle) Area() float64 { panic("TODO: implement") }
// TODO: implement TotalArea.
func TotalArea(shapes []Shape) float64 { panic("TODO: implement") }
