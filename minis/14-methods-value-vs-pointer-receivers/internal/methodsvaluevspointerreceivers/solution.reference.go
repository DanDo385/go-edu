//go:build reference

package methodsvaluevspointerreceivers

/*
Reference Solution - Value vs Pointer Receivers
==============================================

Deep explanation of (b *BankAccount) vs (r Rectangle) (per .cursorrules):

(b *BankAccount): b is a POINTER to a BankAccount. When we write b.balance += amount,
Go automatically dereferences: b.balance is (*b).balance — we read/write the
balance in the struct that b points to. The caller's BankAccount is mutated.
If b is nil, b.balance would panic (dereference nil). Hence the nil check.

(r Rectangle): r is a COPY of the Rectangle. The method receives a full copy.
r.Width * r.Height reads from our copy. Mutating r would not affect the caller's
original. Value receiver = pass-by-value; pointer receiver = pass-by-reference.

Before call: caller has acct (a BankAccount). acct.Deposit(10) passes &acct.
Inside Deposit, b holds that address. b.balance += 10 writes to acct's balance.
After: acct.balance has changed. The same memory location was modified.
*/

func (b *BankAccount) DepositSolution(amount int) {
	if b == nil {
		return
	}
	b.balance += amount
}

// BalanceSolution - Pointer receiver for consistency with Deposit/Withdraw; nil-safe.
func (b *BankAccount) BalanceSolution() int {
	if b == nil {
		return 0
	}
	return b.balance
}

// WithdrawSolution - Pointer receiver: mutates balance. Nil check prevents panic.
func (b *BankAccount) WithdrawSolution(amount int) {
	if b == nil {
		return
	}
	b.balance -= amount
}

// AreaSolution - Value receiver: Rectangle is small, read-only. No mutation.
func (r Rectangle) AreaSolution() float64 {
	return r.Width * r.Height
}

// AreaSolution - Pointer receiver for Circle (may be part of interface); nil-safe.
func (c *Circle) AreaSolution() float64 {
	if c == nil {
		return 0
	}
	return 3.14159 * c.Radius * c.Radius
}

// TotalAreaSolution - Uses Shape interface; both (Rectangle) and (*Circle) satisfy it.
func TotalAreaSolution(shapes []Shape) float64 {
	total := 0.0
	for _, shape := range shapes {
		total += shape.Area()
	}
	return total
}

// AppendSolution - Pointer receiver: mutates list. Nil list: create new head.
func (l *StringList) AppendSolution(value string) *StringList {
	if l == nil {
		return &StringList{value: value}
	}
	if l.next == nil {
		l.next = &StringList{value: value}
		return l
	}
	l.next = l.next.AppendSolution(value)
	return l
}

// ContainsSolution - Pointer receiver for linked list; recursive with nil base case.
func (l *StringList) ContainsSolution(value string) bool {
	if l == nil {
		return false
	}
	if l.value == value {
		return true
	}
	return l.next.ContainsSolution(value)
}

// FirstSolution - Read-only but pointer receiver for consistency with list API.
func (l *StringList) FirstSolution() string {
	if l == nil {
		return ""
	}
	return l.value
}

// ValidateSolution - Value receiver: SmallConfig is small, read-only validation.
func (c SmallConfig) ValidateSolution() bool {
	return c.ID > 0 && c.Name != ""
}

// SumSolution - Pointer receiver: LargeConfig has slice, avoid copying. Nil-safe.
func (l *LargeConfig) SumSolution() int {
	if l == nil {
		return 0
	}
	total := 0
	for _, v := range l.Data {
		total += v
	}
	return total
}

// SetNameSolution - Pointer receiver: mutates User.Name.
func (u *User) SetNameSolution(name string) {
	if u == nil {
		return
	}
	u.Name = name
}

// SetEmailSolution - Pointer receiver: mutates User.Email.
func (u *User) SetEmailSolution(email string) {
	if u == nil {
		return
	}
	u.Email = email
}

// GetNameSolution - Read-only but pointer for consistency with setters; nil-safe.
func (u *User) GetNameSolution() string {
	if u == nil {
		return ""
	}
	return u.Name
}

// IsAdultSolution - Read-only; pointer for API consistency.
func (u *User) IsAdultSolution() bool {
	if u == nil {
		return false
	}
	return u.Age >= 18
}

// EqualsSolution - Value receiver; Point is small. Type switch for Point vs *Point.
func (p Point) EqualsSolution(other Comparable) bool {
	switch v := other.(type) {
	case Point:
		return p.X == v.X && p.Y == v.Y
	case *Point:
		if v == nil {
			return false
		}
		return p.X == v.X && p.Y == v.Y
	default:
		return false
	}
}

// NewSafeCounterMapSolution - Returns value (struct); mutex/counters embedded.
func NewSafeCounterMapSolution() SafeCounterMap {
	return SafeCounterMap{counters: make(map[string]int)}
}

// IncrementSolution - Pointer receiver: mutates map. Nil-safe, lazy-init map.
func (m *SafeCounterMap) IncrementSolution(key string) {
	if m == nil {
		return
	}
	if m.counters == nil {
		m.counters = make(map[string]int)
	}
	m.counters[key]++
}

// GetSolution - Pointer receiver for consistency; nil-safe.
func (m *SafeCounterMap) GetSolution(key string) int {
	if m == nil || m.counters == nil {
		return 0
	}
	return m.counters[key]
}

// AppendIterative - Same as AppendSolution but iterative (no recursion).
func (l *StringList) AppendIterative(value string) *StringList {
	if l == nil {
		return &StringList{value: value}
	}
	cur := l
	for cur.next != nil {
		cur = cur.next
	}
	cur.next = &StringList{value: value}
	return l
}
