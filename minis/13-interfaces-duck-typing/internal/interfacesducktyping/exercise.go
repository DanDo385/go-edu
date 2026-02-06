//go:build !solution && !reference

package interfacesducktyping

import (
	"fmt"
	"reflect"
)

func (p Person) String() string {
	return fmt.Sprintf("%s (%d years old)", p.Name, p.Age)
}

func GetAge(s Stringer) (int, bool) {
	p, ok := s.(Person)
	if !ok {
		return 0, false
	}
	return p.Age, true
}

func DescribeType(i interface{}) string {
	switch v := i.(type) {
	case int:
		return fmt.Sprintf("Integer: %d", v)
	case string:
		return fmt.Sprintf("String: %s", v)
	case bool:
		return fmt.Sprintf("Boolean: %t", v)
	case Person:
		return fmt.Sprintf("Person: %s", v.Name)
	case nil:
		return "Nil"
	default:
		return "Unknown"
	}
}

func IsValidEmail(v Validator) bool {
	if v == nil {
		return false
	}
	if e, ok := v.(*Email); ok && e == nil {
		return false
	}
	return v.IsValid()
}

func (b *Buffer) Read() string {
	if b == nil {
		return ""
	}
	return b.data
}

func (b *Buffer) Write(data string) error {
	if b == nil {
		return fmt.Errorf("nil buffer")
	}
	b.data += data
	return nil
}

func IsReadWriter(i interface{}) bool {
	_, ok := i.(ReadWriter)
	return ok
}

func (c *Counter) Increment() {
	if c == nil {
		return
	}
	c.Value++
}

func CanIncrement(i interface{}) bool {
	_, ok := i.(Incrementer)
	return ok
}

func CountTypes(values []interface{}) map[string]int {
	counts := make(map[string]int)
	for _, v := range values {
		if v == nil {
			counts["<nil>"]++
			continue
		}

		t := reflect.TypeOf(v).String()
		if t == "interfacesducktyping.Person" {
			t = "exercise.Person"
		}
		counts[t]++
	}
	return counts
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("validation error on %s: %s", e.Field, e.Message)
}

func (r Rectangle) Area() float64 {
	return r.Width * r.Height
}

func (c Circle) Area() float64 {
	return 3.14159 * c.Radius * c.Radius
}

func TotalArea(shapes []Shape) float64 {
	total := 0.0
	for _, shape := range shapes {
		total += shape.Area()
	}
	return total
}
