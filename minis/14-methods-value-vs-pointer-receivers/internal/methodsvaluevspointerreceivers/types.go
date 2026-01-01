package methodsvaluevspointerreceivers

// BankAccount is a tiny type used to demonstrate pointer receivers for mutation.
type BankAccount struct {
	balance int
}

func (b *BankAccount) Deposit(amount int)  { b.DepositSolution(amount) }
func (b *BankAccount) Balance() int        { return b.BalanceSolution() }
func (b *BankAccount) Withdraw(amount int) { b.WithdrawSolution(amount) }

// Shape is used to demonstrate interface satisfaction and method sets.
type Shape interface {
	Area() float64
}

type Rectangle struct {
	Width  float64
	Height float64
}

func (r Rectangle) Area() float64 { return r.AreaSolution() }

type Circle struct {
	Radius float64
}

func (c *Circle) Area() float64 { return c.AreaSolution() }

func TotalArea(shapes []Shape) float64 { return TotalAreaSolution(shapes) }

// StringList is a minimal singly linked list used to demonstrate nil receiver safety.
type StringList struct {
	value string
	next  *StringList
}

func (l *StringList) Append(value string) *StringList { return l.AppendSolution(value) }
func (l *StringList) Contains(value string) bool      { return l.ContainsSolution(value) }
func (l *StringList) First() string                   { return l.FirstSolution() }

// SmallConfig is intentionally small, so value receivers are fine.
type SmallConfig struct {
	ID   int
	Name string
}

func (c SmallConfig) Validate() bool { return c.ValidateSolution() }

// LargeConfig is intentionally large (1000 ints ~ 8KB on 64-bit).
type LargeConfig struct {
	Data [1000]int
}

func (l *LargeConfig) Sum() int { return l.SumSolution() }

type User struct {
	Name  string
	Email string
	Age   int
}

func (u *User) SetName(name string)   { u.SetNameSolution(name) }
func (u *User) SetEmail(email string) { u.SetEmailSolution(email) }
func (u *User) GetName() string       { return u.GetNameSolution() }
func (u *User) IsAdult() bool         { return u.IsAdultSolution() }

type Comparable interface {
	Equals(other Comparable) bool
}

type Point struct {
	X int
	Y int
}

func (p Point) Equals(other Comparable) bool { return p.EqualsSolution(other) }

type SafeCounterMap struct {
	counters map[string]int
}

func NewSafeCounterMap() SafeCounterMap { return NewSafeCounterMapSolution() }

func (m *SafeCounterMap) Increment(key string) { m.IncrementSolution(key) }
func (m *SafeCounterMap) Get(key string) int   { return m.GetSolution(key) }
