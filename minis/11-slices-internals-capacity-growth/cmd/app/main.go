package main

import "fmt"

func main() {
	// Create a nil slice. Note that a nil slice has a len and cap of 0.
	var s []int

	fmt.Println("Appending to a slice and observing its capacity growth.")
	fmt.Println("----------------------------------------------------")
	fmt.Printf("Initial state: len=%d, cap=%d\n", len(s), cap(s))
	fmt.Println("----------------------------------------------------")

	// Loop 100 times, appending a new element in each iteration.
	for i := 0; i < 100; i++ {
		// Store the capacity *before* appending.
		oldCap := cap(s)

		// Append an element.
		s = append(s, i)

		// Get the new capacity.
		newCap := cap(s)

		// If the capacity changed, print a message showing the growth.
		if newCap > oldCap {
			fmt.Printf("Append item %2d: len=%-3d, cap=%-4d -> cap grew from %d to %d\n", i, len(s), newCap, oldCap, newCap)
		}
	}

	fmt.Println("----------------------------------------------------")
	fmt.Println("Notice how the capacity doesn't grow one by one.")
	fmt.Println("It grows in larger and larger chunks to reduce the number of reallocations.")
	fmt.Println("This is Go's slice growth algorithm in action.")
}