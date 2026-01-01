//go:build !solution && !reference

// TODO:
// - Read the tests in exercise_test.go to understand expected behavior.
// - Implement the exported API in this file.
// - Compare with the fully-commented reference in solution.reference.go (go test -tags=reference ./...).
package atomiccountersvsmutex
// TODO: implement NewAtomicCounter.
func NewAtomicCounter() *AtomicCounter { panic("TODO: implement") }
// TODO: implement Increment.
func (c *AtomicCounter) Increment() { panic("TODO: implement") }
// TODO: implement Decrement.
func (c *AtomicCounter) Decrement() { panic("TODO: implement") }
// TODO: implement Add.
func (c *AtomicCounter) Add(delta int64) { panic("TODO: implement") }
// TODO: implement Value.
func (c *AtomicCounter) Value() int64 { panic("TODO: implement") }
// TODO: implement Reset.
func (c *AtomicCounter) Reset() int64 { panic("TODO: implement") }
// TODO: implement NewAtomicFlag.
func NewAtomicFlag() *AtomicFlag { panic("TODO: implement") }
// TODO: implement Set.
func (f *AtomicFlag) Set() { panic("TODO: implement") }
// TODO: implement Clear.
func (f *AtomicFlag) Clear() { panic("TODO: implement") }
// TODO: implement IsSet.
func (f *AtomicFlag) IsSet() bool { panic("TODO: implement") }
// TODO: implement TestAndSet.
func (f *AtomicFlag) TestAndSet() bool { panic("TODO: implement") }
// TODO: implement NewRateLimiter.
func NewRateLimiter(capacity, tokensPerSecond int64) *RateLimiter { panic("TODO: implement") }
// TODO: implement Allow.
func (rl *RateLimiter) Allow() bool { panic("TODO: implement") }
// TODO: implement NewAtomicMax.
func NewAtomicMax() *AtomicMax { panic("TODO: implement") }
// TODO: implement Update.
func (am *AtomicMax) Update(value int64) { panic("TODO: implement") }
// TODO: implement Max.
func (am *AtomicMax) Max() int64 { panic("TODO: implement") }
// TODO: implement NewSpinLock.
func NewSpinLock() *SpinLock { panic("TODO: implement") }
// TODO: implement Lock.
func (sl *SpinLock) Lock() { panic("TODO: implement") }
// TODO: implement Unlock.
func (sl *SpinLock) Unlock() { panic("TODO: implement") }
// TODO: implement NewAtomicState.
func NewAtomicState() *AtomicState { panic("TODO: implement") }
// TODO: implement CurrentState.
func (as *AtomicState) CurrentState() int64 { panic("TODO: implement") }
// TODO: implement Transition.
func (as *AtomicState) Transition(expectedCurrent, newState int64) bool { panic("TODO: implement") }
// TODO: implement NewReferenceCounter.
func NewReferenceCounter() *ReferenceCounter { panic("TODO: implement") }
// TODO: implement Acquire.
func (rc *ReferenceCounter) Acquire() { panic("TODO: implement") }
// TODO: implement Release.
func (rc *ReferenceCounter) Release() bool { panic("TODO: implement") }
// TODO: implement Count.
func (rc *ReferenceCounter) Count() int64 { panic("TODO: implement") }
// TODO: implement NewConfigManager.
func NewConfigManager() *ConfigManager { panic("TODO: implement") }
// TODO: implement Update.
func (cm *ConfigManager) Update(newConfig *Config) { panic("TODO: implement") }
// TODO: implement Get.
func (cm *ConfigManager) Get() *Config { panic("TODO: implement") }
// TODO: implement NewLoadBalancer.
func NewLoadBalancer(numWorkers int) *LoadBalancer { panic("TODO: implement") }
// TODO: implement NextWorker.
func (lb *LoadBalancer) NextWorker() int64 { panic("TODO: implement") }
// TODO: implement NewAtomicBitmap.
func NewAtomicBitmap() *AtomicBitmap { panic("TODO: implement") }
// TODO: implement SetBit.
func (ab *AtomicBitmap) SetBit(bitIndex int) { panic("TODO: implement") }
// TODO: implement ClearBit.
func (ab *AtomicBitmap) ClearBit(bitIndex int) { panic("TODO: implement") }
// TODO: implement TestBit.
func (ab *AtomicBitmap) TestBit(bitIndex int) bool { panic("TODO: implement") }
// TODO: implement IncrementCounterNonAtomic.
func IncrementCounterNonAtomic(counter *int64) { panic("TODO: implement") }
// TODO: implement IncrementCounterAtomic.
func IncrementCounterAtomic(counter *int64) { panic("TODO: implement") }
// TODO: implement NewCircularBuffer.
func NewCircularBuffer(capacity int) *CircularBuffer { panic("TODO: implement") }
// TODO: implement Push.
func (cb *CircularBuffer) Push(value int64) bool { panic("TODO: implement") }
// TODO: implement Pop.
func (cb *CircularBuffer) Pop() (int64, bool) { panic("TODO: implement") }
