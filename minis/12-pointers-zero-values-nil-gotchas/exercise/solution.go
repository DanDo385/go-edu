package exercise

// SafeDeref safely dereferences a pointer.
func SafeDeref(p *int, defaultValue int) int {
	if p == nil {
		return defaultValue
	}
	return *p
}

// Swap exchanges two integer values.
func Swap(a, b *int) {
	temp := *a
	*a = *b
	*b = temp
}

// InitializeMap creates a map if nil.
func InitializeMap(m map[string]int) map[string]int {
	if m == nil {
		return make(map[string]int)
	}
	return m
}

// AppendNode appends to a linked list.
func AppendNode(head *Node, value int) *Node {
	newNode := &Node{Value: value}

	if head == nil {
		return newNode
	}

	current := head
	for current.Next != nil {
		current = current.Next
	}
	current.Next = newNode

	return head
}

// ListLength counts nodes in a list.
func ListLength(head *Node) int {
	count := 0
	current := head
	for current != nil {
		count++
		current = current.Next
	}
	return count
}
