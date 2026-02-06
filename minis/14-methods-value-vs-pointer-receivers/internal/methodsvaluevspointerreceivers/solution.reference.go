//go:build reference

package methodsvaluevspointerreceivers

func (b *BankAccount) DepositSolution(amount int) {
	if b == nil {
		return
	}
	b.balance += amount
}

func (b *BankAccount) BalanceSolution() int {
	if b == nil {
		return 0
	}
	return b.balance
}

func (b *BankAccount) WithdrawSolution(amount int) {
	if b == nil {
		return
	}
	b.balance -= amount
}

func (r Rectangle) AreaSolution() float64 {
	return r.Width * r.Height
}

func (c *Circle) AreaSolution() float64 {
	if c == nil {
		return 0
	}
	return 3.14159 * c.Radius * c.Radius
}

func TotalAreaSolution(shapes []Shape) float64 {
	total := 0.0
	for _, shape := range shapes {
		total += shape.Area()
	}
	return total
}

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

func (l *StringList) ContainsSolution(value string) bool {
	if l == nil {
		return false
	}
	if l.value == value {
		return true
	}
	return l.next.ContainsSolution(value)
}

func (l *StringList) FirstSolution() string {
	if l == nil {
		return ""
	}
	return l.value
}

func (c SmallConfig) ValidateSolution() bool {
	return c.ID > 0 && c.Name != ""
}

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

func (u *User) SetNameSolution(name string) {
	if u == nil {
		return
	}
	u.Name = name
}

func (u *User) SetEmailSolution(email string) {
	if u == nil {
		return
	}
	u.Email = email
}

func (u *User) GetNameSolution() string {
	if u == nil {
		return ""
	}
	return u.Name
}

func (u *User) IsAdultSolution() bool {
	if u == nil {
		return false
	}
	return u.Age >= 18
}

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

func NewSafeCounterMapSolution() SafeCounterMap {
	return SafeCounterMap{counters: make(map[string]int)}
}

func (m *SafeCounterMap) IncrementSolution(key string) {
	if m == nil {
		return
	}
	if m.counters == nil {
		m.counters = make(map[string]int)
	}
	m.counters[key]++
}

func (m *SafeCounterMap) GetSolution(key string) int {
	if m == nil || m.counters == nil {
		return 0
	}
	return m.counters[key]
}

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
