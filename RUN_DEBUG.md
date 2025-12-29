# Run & Debug Guide: Visualizing Variable State and Memory Changes

## Table of Contents

1. [What is Run & Debug?](#what-is-run--debug)
2. [Getting Started - Your First Debug Session](#getting-started---your-first-debug-session)
3. [Understanding the Debug Interface](#understanding-the-debug-interface)
4. [Essential Keyboard Shortcuts](#essential-keyboard-shortcuts)
5. [Understanding the Debug Panels](#understanding-the-debug-panels)
6. [Setting Breakpoints](#setting-breakpoints)
7. [Stepping Through Code](#stepping-through-code)
8. [Tracking Variables Across Functions](#tracking-variables-across-functions)
9. [Inspecting Memory Composition](#inspecting-memory-composition)
10. [Using Watch Expressions](#using-watch-expressions)
11. [Debugging Concurrent Code](#debugging-concurrent-code)
12. [Advanced Techniques](#advanced-techniques)
13. [Common Debugging Workflows](#common-debugging-workflows)
14. [Troubleshooting](#troubleshooting)
15. [Practice Exercises](#practice-exercises)

---

## What is Run & Debug?

Run & Debug is a powerful tool in Cursor (and VS Code) that lets you:

- **Pause execution** at specific points (breakpoints)
- **Step through code** line by line
- **Inspect variables** to see their current values
- **Watch expressions** to monitor specific values
- **See how data flows** through your program
- **Understand memory** - how slices, maps, channels change
- **Debug concurrent code** - see multiple goroutines running

Think of it like a **time machine for your code** - you can pause, rewind, and examine what's happening at any moment.

---

## Getting Started - Your First Debug Session

### Step 1: Open a Project

1. Open Cursor
2. Open a project folder (e.g., `minis/06-worker-pool-wordcount`)
3. You should see the project files in the file explorer

### Step 2: Open a File to Debug

1. Open `main.go` or `solution.go`
2. You'll see the code with line numbers on the left

### Step 3: Set Your First Breakpoint

1. **Find a line** you want to pause at (look for `// BREAKPOINT:` comments)
2. **Click in the gutter** (the area left of the line numbers)
3. A **red dot** appears - this is your breakpoint
4. The program will pause when it reaches this line

**Example:**

```go
func main() {
    // BREAKPOINT: Set breakpoint here
    workers := flag.Int("workers", 3, "Number of workers")
    // ↑ Click here (left of line number) to set breakpoint
```

### Step 4: Start Debugging

- Press `F5` OR
- Click the Run and Debug icon in the left sidebar (looks like a play button with a bug)
- Click the green play button at the top OR
- Press `F5` again

### Step 5: Select Debug Configuration

A dropdown appears at the top asking which configuration to use:

- **"Debug: Run main.go (Default Args)"** - Start here!

Select it and press Enter.

### Step 6: Code Executes and Pauses

- Your program starts running
- When it hits your breakpoint, execution pauses
- The line with the breakpoint is highlighted in yellow
- You're now in "debug mode"

**What you'll see:**

- The line you're paused at is highlighted
- Debug panels appear at the bottom
- You can now inspect variables and step through code

### Step 7: Inspect Variables

- Look at the **Variables panel** (usually at bottom left)
- You'll see all variables in the current scope
- Expand variables by clicking the `>` arrow
- See their current values

**Example:**

```
Variables:
  > workers (*int): 0x... (pointer address)
    > *workers (int): 3
```

### Step 8: Step Through Code

- Press `F10` (Step Over) to execute the current line
- Watch the highlighted line move to the next line
- Watch Variables panel - values may change
- Continue pressing `F10` to step through

### Step 9: Stop Debugging

- Press `Shift+F5` to stop
- OR click the stop button (square icon)

**Congratulations! You've completed your first debug session.**

---

## Understanding the Debug Interface

When you start debugging, Cursor shows several panels:

### Top Bar (Debug Controls)

```
[Continue] [Step Over] [Step Into] [Step Out] [Restart] [Stop]
   F5         F10         F11        Shift+F11  Ctrl+Shift+F5  Shift+F5
```

- **Continue (F5)**: Resume execution until next breakpoint
- **Step Over (F10)**: Execute current line, don't enter functions
- **Step Into (F11)**: Enter function calls to see internals
- **Step Out (Shift+F11)**: Exit current function, return to caller
- **Restart**: Start debugging again from beginning
- **Stop**: End debugging session

### Left Sidebar - Debug Panels

When debugging, you'll see these panels:

- **Variables** - Shows all variables in current scope
- **Watch** - Custom expressions you want to monitor
- **Call Stack** - Shows function call hierarchy
- **Breakpoints** - Lists all your breakpoints
- **Threads** - Shows all goroutines (for concurrent code)

### Bottom Panel - Debug Console

- Shows debug output and logs
- Can evaluate expressions
- Type variable names to see their values

---

## Essential Keyboard Shortcuts

Memorize these - they'll make debugging much faster:

| Shortcut | Action | When to Use |
|----------|--------|-------------|
| `F5` | Start/Continue | Start debugging or continue to next breakpoint |
| `F9` | Toggle Breakpoint | Add/remove breakpoint at current line |
| `F10` | Step Over | Execute current line, don't enter functions |
| `F11` | Step Into | Enter function calls to see internals |
| `Shift+F11` | Step Out | Exit current function |
| `Ctrl+Shift+F5` | Restart | Start debugging from beginning |
| `Shift+F5` | Stop | End debugging session |

**Pro Tip:** Keep your hand on `F10` - you'll use it constantly!

---

## Understanding the Debug Panels

### Variables Panel

**Location:** Usually bottom-left when debugging

**What it shows:**

- All variables in the current scope (local variables, function parameters, package-level variables)
- Variable name, type, and current value
- Expandable structures (structs, slices, maps, channels)

**How to use:**

1. **View variables:**
   - Variables automatically appear when you pause
   - They're organized by scope (Local, Function parameters, Package)

2. **Expand structures:**
   - Click the `>` arrow next to variable name
   - See nested values (map entries, slice elements, struct fields)

3. **Copy values:**
   - Right-click variable → "Copy Value"
   - Useful for comparing or documenting

4. **Set values (advanced):**
   - Right-click variable → "Set Value"
   - Change variable value during debugging
   - Useful for testing different scenarios

**Example - Inspecting a Map:**

```
Variables Panel shows:
  counts (map[string]int)
    > len (int): 3
    > "hello" (int): 5
    > "world" (int): 3
    > "go" (int): 1
```

**What this tells you:**

- Map has 3 entries (len=3)
- Key "hello" has value 5
- Key "world" has value 3
- Key "go" has value 1

**Example - Inspecting a Slice:**

```
Variables Panel shows:
  urls ([]string)
    > len (int): 3
    > cap (int): 4
    > [0] (string): "http://example.com/1"
    > [1] (string): "http://example.com/2"
    > [2] (string): "http://example.com/3"
```

**What this tells you:**

- Slice has 3 elements (len=3)
- Underlying array can hold 4 elements (cap=4)
- Elements are indexed [0], [1], [2]

### Watch Panel

**Location:** Usually below Variables panel

**What it shows:**

- Custom expressions you want to monitor
- Updates automatically as you step through code
- Can evaluate any valid Go expression

**How to use:**

1. **Add watch expression:**
   - Click the `+` icon in Watch panel
   - Type an expression (e.g., `len(mySlice)`)
   - Press Enter

2. **Expression updates automatically:**
   - As you step through code (F10), expression re-evaluates
   - See how value changes at each step

3. **Remove watch expression:**
   - Click the `-` icon next to expression
   - OR right-click → Remove

**Useful Watch Expressions:**

**For Slices:**
```go
len(urls)                    // Current length
cap(urls)                    // Current capacity
urls[0]                      // First element
urls[len(urls)-1]            // Last element (if len > 0)
```

**For Maps:**
```go
len(counts)                  // Number of entries
counts["hello"]              // Value for key "hello"
```

**For Channels:**
```go
len(jobs)                    // Number of items in buffer
cap(jobs)                    // Channel capacity
```

### Call Stack Panel

**Location:** Usually below Watch panel

**What it shows:**

- Function call hierarchy (current function at top)
- Shows how you got to current location
- Can navigate to any function in the stack

**Example:**

```
Call Stack:
  > tokenizeAndCount (current)    ← You are here
    fetchAndCount                 ← Called from here
    WordCount (worker goroutine)  ← Called from here
    main                          ← Started here
```

### Threads Panel (For Concurrent Code)

**Location:** Usually below Call Stack

**What it shows:**

- All running goroutines
- Current location of each goroutine
- Can switch between goroutines

**Example:**

```
Threads:
  > [1] main.main (main goroutine)     ← Currently selected
    [2] worker (goroutine #1)          ← Click to switch
    [3] worker (goroutine #2)          ← Click to switch
    [4] worker (goroutine #3)          ← Click to switch
```

---

## Setting Breakpoints

### Basic Breakpoints

**How to set:**

1. Click in the gutter (left of line numbers)
2. Red dot appears
3. Program pauses when it reaches this line

**How to remove:**

1. Click the red dot again
2. OR right-click → Remove Breakpoint

### Conditional Breakpoints

**How to set:**

1. Right-click on breakpoint (red dot)
2. Select "Edit Breakpoint"
3. Enter condition (e.g., `len(urls) > 5`)
4. Breakpoint only triggers when condition is true

**Use cases:**

- Break only on specific iterations: `i == 5`
- Break when variable reaches value: `count > 100`
- Break on error conditions: `err != nil`

**Example:**

```go
for i, url := range urls {
    // Conditional breakpoint: i == 5
    // Only breaks on 6th iteration (when i=5)
    process(url)
}
```

### Logpoints

**How to set:**

1. Right-click on line (not on breakpoint)
2. Select "Add Logpoint"
3. Enter expression: `{url} processed, count: {len(counts)}`
4. Expression is logged without stopping execution

**Example:**

```go
for url := range urls {
    // Logpoint: {url} - count: {len(counts)}
    // Outputs: "http://example.com - count: 5"
    process(url)
}
```

---

## Stepping Through Code

### Step Over (F10)

**What it does:**

- Executes the current line
- Moves to the next line
- Does NOT enter function calls

**When to use:**

- You understand what the function does
- You want to skip function internals
- You're focusing on the current function

**Example:**

```go
counts := make(map[string]int)  // ← You are here, press F10
counts["hello"]++                // ← Moves here (skips make() internals)
fmt.Println(counts)              // ← Then here
```

### Step Into (F11)

**What it does:**

- Executes the current line
- If line contains function call, ENTERS that function
- Shows you function internals

**When to use:**

- You want to see how a function works
- Function behavior is unclear
- You're debugging a specific function

**Example:**

```go
result := processURL(url)  // ← You are here, press F11
// Now you're INSIDE processURL function
// You can see how it processes the URL
```

**Pro Tip:** Use F11 to "dive deep" into functions you don't understand.

### Step Out (Shift+F11)

**What it does:**

- Executes rest of current function
- Returns to the caller
- Useful when you've seen enough of a function

**When to use:**

- You've seen what you need in current function
- You want to return to caller quickly
- Function is long and you're done with it

### Continue (F5)

**What it does:**

- Resumes execution
- Continues until next breakpoint
- Or until program ends

**When to use:**

- You've inspected what you need
- You want to skip to next breakpoint
- You're done with current section

---

## Tracking Variables Across Functions

### Technique 1: Step Through Function Calls

**Complete Workflow:**

1. **Set breakpoint at function call site:**

```go
func main() {
    urls := []string{"url1", "url2", "url3"}
    // BREAKPOINT: Set here
    counts, err := WordCount(ctx, urls, 3)
}
```

2. **Before calling function:**
   - Inspect parameters in Variables panel
   - See: `urls` (slice with 3 elements), `ctx`, `workers=3`
   - Add watch expressions: `len(urls)`, `workers`

3. **Step Into function (F11):**
   - Now you're inside `WordCount` function
   - Variables panel shows function parameters
   - See: `urls`, `workers` (same values, but in function scope)

4. **Step through function:**
   - Press F10 to step through function
   - Watch Variables panel - see new variables created
   - Watch: jobs channel created, results channel created
   - Watch: finalCounts map created and populated

5. **At function exit:**
   - Inspect return value before returning
   - See: finalCounts map with all word counts

6. **Step Out (Shift+F11) to return to caller**

7. **Back in caller:**
   - See return value assigned to `counts` variable
   - Inspect `counts` - should match what you saw in function
   - See how return value affects caller's code

**What you learned:**

- How parameters flow into function
- How function processes data
- How return value flows back to caller
- How data transforms through function

---

## Inspecting Memory Composition

### Understanding Go Memory Layout

**Basic Types:**

- `int`, `string`, `bool`: Value types (stored directly)
- `[]T`: Slice (header + underlying array)
- `map[K]V`: Map (hash table structure)
- `chan T`: Channel (buffer + synchronization)
- `*T`: Pointer (address to memory)

### Inspecting Slices

**What Variables Panel Shows:**

```
Variables Panel:
  words ([]string)
    > len (int): 5           // Current length
    > cap (int): 8           // Current capacity
    > [0] (string): "hello"  // Elements
    > [1] (string): "world"
    > [2] (string): "go"
    > [3] (string): "lang"
    > [4] (string): "rocks"
```

**What to Watch:**

1. **Length (len):**
   - Shows how many elements are currently used
   - Increases as you append items
   - Decreases if you slice the slice

2. **Capacity (cap):**
   - Shows underlying array size
   - When len == cap, next append will allocate new array
   - Usually doubles when it grows

3. **Elements:**
   - Expand slice to see individual elements
   - Indexed [0], [1], [2], etc.
   - See actual values stored

**Watch Expressions:**

```go
len(words)                   // Current length
cap(words)                   // Current capacity
words[0]                     // First element
words[len(words)-1]          // Last element
```

**Example Workflow - Watching Slice Grow:**

```go
// Set breakpoint at slice creation:
words := make([]string, 0, 4)   // BREAKPOINT: Here
// Inspect initial state: len=0, cap=4

words = append(words, "hello")   // BREAKPOINT: Here
// Variables: len=1, cap=4 (capacity not exceeded)

words = append(words, "world", "go", "lang")   // BREAKPOINT: Here
// Variables: len=4, cap=4 (at capacity!)

words = append(words, "rocks")   // BREAKPOINT: Here
// Variables: len=5, cap=8 (capacity doubled!)
```

### Inspecting Maps

**What Variables Panel Shows:**

```
Variables Panel:
  counts (map[string]int)
    > len (int): 3           // Number of entries
    > "hello" (int): 5       // Key-value pairs
    > "world" (int): 3
    > "go" (int): 1
```

**Example Workflow - Watching Map Grow:**

```go
// Step 1: Breakpoint at map creation
counts := make(map[string]int)   // BREAKPOINT: Here
// Variables: counts (map[string]int): {}

// Step 2: Add entries
counts["hello"]++   // BREAKPOINT: Here
// Variables: counts has 1 entry
//   > "hello" (int): 1

counts["world"]++   // BREAKPOINT: Here
// Variables: counts has 2 entries
//   > "hello" (int): 1
//   > "world" (int): 1

counts["hello"]++   // BREAKPOINT: Here
// Variables: counts still has 2 entries (no new key)
//   > "hello" (int): 2 (value updated!)
//   > "world" (int): 1
```

### Inspecting Channels

**What Variables Panel Shows:**

```
Variables Panel:
  jobs (chan string)
    > len (int): 2           // Items in buffer
    > cap (int): 5           // Buffer capacity
```

**Example Workflow:**

```go
jobs := make(chan string, 5)   // BREAKPOINT: Here
// Variables: jobs (chan string): len=0, cap=5

jobs <- "url1"   // BREAKPOINT: Here
// Variables: len=1, cap=5 (1 item in buffer)

jobs <- "url2"   // BREAKPOINT: Here
// Variables: len=2, cap=5 (2 items in buffer)

url := <-jobs   // BREAKPOINT: Here
// Variables: len=1, cap=5 (1 item removed)
```

---

## Using Watch Expressions

### Adding Watch Expressions

**Step-by-step:**

1. Start debugging (F5)
2. Pause at breakpoint
3. Open Watch panel (usually below Variables)
4. Click `+` icon
5. Type expression (e.g., `len(urls)`)
6. Press Enter

**Expression appears in Watch panel:**

```
Watch:
  > len(urls): 3
```

### Useful Watch Expressions by Type

**Slices:**
```go
len(mySlice)                 // Current length
cap(mySlice)                 // Current capacity
mySlice[0]                   // First element
mySlice[len(mySlice)-1]      // Last element
```

**Maps:**
```go
len(myMap)                   // Number of entries
myMap["key"]                 // Value for key
```

**Channels:**
```go
len(myChannel)               // Items in buffer
cap(myChannel)               // Channel capacity
```

**Complex Expressions:**
```go
len(mySlice) > 0             // Boolean
myMap["a"] + myMap["b"]      // Arithmetic
len(myChannel) == cap(myChannel)  // Comparison
```

---

## Debugging Concurrent Code

### Setting Up for Concurrent Debugging

**Set breakpoints in goroutines:**

```go
go func() {
    // BREAKPOINT: Set here
    process(url)
}()
```

### Using Threads Panel

**What you see:**

```
Threads:
  > [1] main.main (main goroutine)
    [2] worker (goroutine #1)
    [3] worker (goroutine #2)
    [4] worker (goroutine #3)
```

**How to use:**

1. Click on goroutine to switch context
2. Variables panel updates to show that goroutine's variables
3. Set breakpoints in different goroutines
4. Switch between them to see each state

### Inspecting Shared Resources

**Channels:**
- All goroutines see same channel
- Watch `len(channel)` to see items
- See how goroutines coordinate through channels

**Shared Maps:**
- Multiple goroutines can access
- Use race detector configuration
- Watch for concurrent modifications

---

## Advanced Techniques

### Conditional Breakpoints

**How to set:**

1. Right-click on breakpoint (red dot)
2. Select "Edit Breakpoint"
3. Enter condition: `len(urls) > 5`
4. Breakpoint only triggers when condition is true

**Use cases:**

- Break only on specific iterations: `i == 5`
- Break when variable reaches value: `count > 100`
- Break on error conditions: `err != nil`

### Logpoints

**How to set:**

1. Right-click on line (not breakpoint)
2. Select "Add Logpoint"
3. Enter: `{url} processed, count: {len(counts)}`
4. Expression logs without stopping

### Data Breakpoints (Memory Watchpoints)

**How to set:**

1. Right-click variable in Variables panel
2. Select "Break on Value Change"
3. Breaks whenever variable is modified

---

## Common Debugging Workflows

### Workflow 1: Understanding Function Behavior

**Goal:** Understand how a function works

**Steps:**

1. Set breakpoint at function call site
2. Inspect parameters before call
3. Step Into (F11) function
4. Step through function (F10), watching variables
5. Inspect return value
6. Step Out (Shift+F11) to caller
7. See how return value affects caller

### Workflow 2: Tracking Data Transformation

**Goal:** See how data changes through transformations

**Steps:**

1. Set breakpoint at data creation
2. Add watch expressions for key variables
3. Step through transformations
4. Watch values change at each step
5. Inspect final result

### Workflow 3: Debugging Concurrent Code

**Goal:** Understand concurrent execution

**Steps:**

1. Set breakpoints in multiple goroutines
2. Start debugging
3. When breakpoint hits, check Threads panel
4. Switch between goroutines
5. Inspect variables for each
6. Watch shared channels/maps
7. Use race detector configuration

### Workflow 4: Finding Bugs

**Goal:** Locate where bug occurs

**Steps:**

1. Reproduce bug
2. Set breakpoint before bug occurs
3. Step through code carefully
4. Watch variables at each step
5. Identify where values become incorrect
6. Add watch expressions for suspicious variables
7. Use conditional breakpoints to narrow down

---

## Troubleshooting

### Common Issues

**Problem:** Breakpoint not hitting

**Solutions:**
- Verify code is actually executed
- Check build tags (solution vs exercise)
- Ensure debug configuration is correct
- Rebuild: `go build`

**Problem:** Variables not showing

**Solutions:**
- Ensure you're paused at breakpoint
- Check variable scope
- Try expanding parent structures
- Disable optimizations: `-gcflags='-N -l'`

**Problem:** Can't see goroutines

**Solutions:**
- Open Threads panel
- Set breakpoints inside goroutines
- Ensure goroutines are running
- Check for deadlocks

---

## Practice Exercises

### Exercise 1: Your First Debug Session

**Goal:** Get comfortable with basic debugging

**Steps:**

1. Open `minis/01-hello-strings`
2. Open `solution.go`
3. Set breakpoint at function entry
4. Press F5, select "Debug: Test solution.go"
5. When paused, inspect variables
6. Press F10 to step through 5 lines
7. Press F5 to continue
8. Press Shift+F5 to stop

### Exercise 2: Track a Variable

**Goal:** Watch how a variable changes

**Steps:**

1. Open `minis/02-arrays-maps-basics`
2. Open `solution.go`
3. Find a map being built
4. Set breakpoint before map creation
5. Add watch: `len(myMap)`
6. Step through code (F10)
7. Watch the length increase

### Exercise 3: Debug a Function Call

**Goal:** Follow data through a function

**Steps:**

1. Open `minis/06-worker-pool-wordcount`
2. Open `solution.go`
3. Find function call
4. Set breakpoint at call site
5. Step Into (F11) the function
6. Inspect parameters
7. Step through function
8. Step Out (Shift+F11)
9. Inspect return value

### Exercise 4: Debug Concurrent Code

**Goal:** Debug multiple goroutines

**Steps:**

1. Open `minis/06-worker-pool-wordcount`
2. Open `solution.go`
3. Set breakpoint inside worker goroutine
4. Start debugging
5. When paused, open Threads panel
6. Click different goroutines
7. See their different states

---

## Additional Resources

- **VS Code Debugging Docs:** https://code.visualstudio.com/docs/editor/debugging
- **Go Debugging in VS Code:** https://github.com/golang/vscode-go/wiki/debugging
- **Delve Debugger:** https://github.com/go-delve/delve

---

**Happy Debugging! Remember: The debugger is your best friend for understanding code. Use it often!**
