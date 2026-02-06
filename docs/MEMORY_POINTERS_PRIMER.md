# Memory and Pointers Primer

This primer is intentionally detailed. Use it whenever a lesson introduces `*` or `&`.

## At a High Level

In Go:

1. Every variable stores a value.
2. Some values are plain data (`int`, `struct` copies).
3. Some values contain addresses (`*T`, slices, maps, channels, functions, interfaces).
4. Function arguments are always passed by value.

The key detail: a copied value may itself contain an address. That is where aliasing comes from.

## Under the Hood: `*` Has Multiple Meanings

`*` can mean different things depending on context.

### 1. Pointer Type: `*T`

Example:

```go
var p *int
```

Meaning:

1. `p` is a variable whose value is either an address of an `int` or `nil`.
2. `*` here is part of a type declaration.
3. No dereference happens yet.

Memory view:

```text
stack:
  p = nil   // p exists, but points to nothing
```

Common misconception:

1. "`*` means read data." Not here. In `*int`, it defines pointer type only.

### 2. Dereference: `*p`

Example:

```go
x := 10
p := &x
*p = 20
```

Meaning:

1. `&x` computes the address of `x`.
2. `p` stores that address value.
3. `*p` means "go to the address stored in `p`, then access that cell."
4. `*p = 20` mutates the same `x` cell.

Memory before assignment:

```text
stack:
  x = 10
  p = &x
```

Memory after `*p = 20`:

```text
stack:
  x = 20
  p = &x
```

Common misconception:

1. "`*p = 20` creates a copy." It does not. It writes into the original location.

### 3. Multiplication: `a * b`

Example:

```go
area := width * height
```

Meaning:

1. `*` is arithmetic multiplication only.
2. No pointer semantics are involved.

Rule of thumb:

1. If both sides are numeric expressions, `*` is multiplication.
2. If `*` appears before an expression (`*p`) or inside a type (`*T`), it is pointer-related.

## Under the Hood: `&` (Address-Of)

Example:

```go
x := 7
p := &x
```

Meaning:

1. `&x` computes an address value.
2. `p` receives that address as its value.
3. No data copy of `x` is needed to create `p`; only the address is copied into `p`.

Memory:

```text
stack:
  x = 7
  p = 0xABC...   // address where x is stored
```

Common misconception:

1. "`&` means pass-by-reference." Go remains pass-by-value.
2. If you pass `p`, Go copies the pointer value (the address), not the pointed-to `int`.

## Pointer Receivers: `func (s *State) Update()`

Example:

```go
type Counter struct{ n int }

func (c *Counter) Inc() {
	c.n++
}
```

Meaning:

1. Receiver type is `*Counter`, so method receives a copied pointer value.
2. That pointer value still addresses the original `Counter`, so mutations are visible to callers.

Memory model:

```text
caller stack:
  ctr = Counter{n: 0}

call frame:
  c = &ctr   // pointer value copied into method frame

method writes through c -> ctr.n becomes 1
```

Common misconception:

1. "Pointer receiver means object must be on heap." Not necessarily.
2. Heap vs stack is decided by escape analysis, not by receiver syntax alone.

## Escape Analysis and Allocation Placement

The compiler can keep values on stack when safe. Values tend to escape to heap when:

1. Their address is returned.
2. Their address outlives the current frame.
3. They are captured by goroutines/closures with longer lifetime.

Example:

```go
func makePtr() *int {
	x := 1
	return &x // x must outlive function frame, so compiler allocates appropriately
}
```

Key point:

1. `&x` does not force heap in all cases.
2. `&x` may still remain stack-local if lifetime does not escape.

## Aliasing: Why Bugs Happen

Aliasing means two program paths can reach the same underlying memory.

Common sources:

1. Pointers (`*T`).
2. Slices sharing an array.
3. Maps/channels shared between goroutines.

When debugging aliasing bugs ask:

1. Which variable owns mutation rights?
2. Which variables share the same underlying storage?
3. Is synchronization explicit when concurrent access exists?

## Minimal Mental Model

1. A pointer is just a value containing an address.
2. Dereference follows that address to data.
3. `&` produces an address value.
4. Go is pass-by-value always.
5. Aliasing and lifetime determine behavior, performance, and safety.
