//go:build !solution && !reference

// TODO:
// - Read the tests in exercise_test.go to understand expected behavior.
// - Implement the exported API in this file.
// - Compare with the fully-commented reference in solution.reference.go (go test -tags=reference ./...).
package escapeanalysisinlining
// TODO: implement SumIntsOptimizedSolution.
func SumIntsOptimizedSolution(values []int) int { panic("TODO: implement") }
// TODO: implement CalculateAreaOptimizedSolution.
func CalculateAreaOptimizedSolution(width, height float64) float64 { panic("TODO: implement") }
// TODO: implement JoinStringsOptimizedSolution.
func JoinStringsOptimizedSolution(parts []string, separator string) string { panic("TODO: implement") }
// TODO: implement AreaValueReceiverSolution.
func (r Rectangle) AreaValueReceiverSolution() float64 { panic("TODO: implement") }
// TODO: implement ProcessItemsOptimizedSolution.
func ProcessItemsOptimizedSolution(items []string) [][]byte { panic("TODO: implement") }
// TODO: implement FormatIntOptimizedSolution.
func FormatIntOptimizedSolution(prefix string, value int) string { panic("TODO: implement") }
// TODO: implement FormatIntOptimizedManual.
func FormatIntOptimizedManual(prefix string, value int) string { panic("TODO: implement") }
// TODO: implement FilterPositiveOptimizedSolution.
func FilterPositiveOptimizedSolution(numbers []int) []int { panic("TODO: implement") }
// TODO: implement FilterPositiveOptimizedEstimate.
func FilterPositiveOptimizedEstimate(numbers []int) []int { panic("TODO: implement") }
// TODO: implement GetConfigOptimizedSolution.
func GetConfigOptimizedSolution() Config { panic("TODO: implement") }
// TODO: implement BuildStringNoAlloc.
func BuildStringNoAlloc(parts []string, separator string) string { panic("TODO: implement") }

type Point struct {
	X, Y float64
}
// TODO: implement NewPoint.
func NewPoint(x, y float64) Point { panic("TODO: implement") }
// TODO: implement SumArrayNoBoundsCheck.
func SumArrayNoBoundsCheck(arr *[1000]int) int { panic("TODO: implement") }
// TODO: implement EscapingAllocation.
func EscapingAllocation() *int { panic("TODO: implement") }
// TODO: implement NonEscapingAllocation.
func NonEscapingAllocation() int { panic("TODO: implement") }
// TODO: implement InlinableFunction.
func InlinableFunction(a, b int) int { panic("TODO: implement") }
// TODO: implement NonInlinableFunction.
func NonInlinableFunction(a, b, c, d, e, f int) int { panic("TODO: implement") }

type SmallStruct struct{ A, B int }
// TODO: implement Good.
func Good() SmallStruct { panic("TODO: implement") }
// TODO: implement Bad.
func Bad() *SmallStruct { panic("TODO: implement") }
// TODO: implement GoodLocal.
func GoodLocal() int { panic("TODO: implement") }
// TODO: implement BadLocal.
func BadLocal() []int { panic("TODO: implement") }
// TODO: implement GoodClosure.
func GoodClosure() func() int { panic("TODO: implement") }
// TODO: implement BetterNoClosure.
func BetterNoClosure(x int) int { panic("TODO: implement") }
// TODO: implement GoodConcrete.
func GoodConcrete(x int) int { panic("TODO: implement") }
// TODO: implement BadInterface.
func BadInterface(x interface{}) interface{} { panic("TODO: implement") }
