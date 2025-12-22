# Project 03: csv-stats - Understanding Structured Data Processing

## What Is This Project About?

Imagine you work for a company and you need to analyze thousands of credit card transactions stored in a spreadsheet. You want to know: "How much did we spend on groceries this month? What was the average transaction amount per category?"

This project teaches you how to:
1. Read structured data from CSV (Comma-Separated Values) files
2. Parse and validate data row by row
3. Compute aggregate statistics (count, sum, average) per category
4. Handle malformed data gracefully

## First Principles: What Is Structured Data?

### From Unstructured to Structured

**Unstructured** data is like a paragraph of text—there's no predefined format:
```
We bought groceries for $12.50, then books for $10, and more groceries for $7.50.
```

**Structured** data organizes information into rows and columns, like a table:

| ID | Category   | Amount |
|----|-----------|--------|
| 1  | groceries | 12.50  |
| 2  | books     | 10.00  |
| 3  | groceries | 7.50   |

CSV is one of the simplest formats for structured data.

### What Is CSV?

CSV (Comma-Separated Values) is a text file where:
- Each line is a **row**
- Values in a row are separated by **commas**
- The first row usually contains **headers** (column names)

Example `transactions.csv`:
```csv
id,category,amount
1,groceries,12.50
2,books,10.00
3,groceries,7.50
```

When you open this in Excel or Google Sheets, it appears as a table. But it's just plain text!

## The Problem We're Solving

**Input**: A CSV file with transactions
**Output**: Statistics per category (count, sum, average)

Example:
```
Input CSV:
id,category,amount
1,groceries,12.50
2,groceries,7.50
3,books,10.00

Output:
groceries: Count=2, Sum=$20.00, Avg=$10.00
books: Count=1, Sum=$10.00, Avg=$10.00
```

## Breaking Down the Solution (Step by Step)

### Step 1: Understanding CSV Structure

A CSV file is just text with a specific pattern:
```
header1,header2,header3
value1,value2,value3
value1,value2,value3
```

To process it, we need to:
1. Read it line by line
2. Split each line by commas
3. Extract values from the split pieces

### Step 2: Validating the Header

The first line tells us the column names. We expect:
```
id,category,amount
```

If we get something else (like `foo,bar,baz`), the file is invalid. We should return an error immediately rather than processing garbage data.

**Why validate early?**
- Prevents processing wrong files
- Gives clear error messages ("expected 'id,category,amount', got 'foo,bar,baz'")
- Follows the "fail fast" principle

### Step 3: Processing Each Row

For each data row:
```
1,groceries,12.50
```

We need to:
1. **Split by comma** → `["1", "groceries", "12.50"]`
2. **Extract fields**:
   - `id` = "1" (we don't use this, but could validate it's a number)
   - `category` = "groceries"
   - `amount` = "12.50"
3. **Convert amount to number**:
   - String `"12.50"` → Float `12.50`
   - If conversion fails (like `"invalid"`), return error
4. **Validate**:
   - Check category isn't empty
   - Check amount is valid number

### Step 4: Aggregating Statistics

We need to track, for each category:
- **Count**: How many transactions?
- **Sum**: Total amount?
- **Avg**: Average amount?

**Data Structure**: A map from category name to statistics:
```go
map[string]Stat{
    "groceries": {Count: 2, Sum: 20.00, Avg: 10.00},
    "books": {Count: 1, Sum: 10.00, Avg: 10.00},
}
```

**Algorithm**:
```
For each row:
    1. Get current stats for this category (or create new stats)
    2. Increment count
    3. Add amount to sum
    4. (Don't compute average yet—do it at the end)

After all rows:
    For each category:
        Compute avg = sum / count
```

## The Complete Solution (Explained Line by Line)

```go
func SummarizeCSV(r io.Reader) (map[string]Stat, error) {
    // Step 1: Create CSV reader
    // The csv.Reader handles:
    // - Splitting rows by newlines
    // - Splitting columns by commas
    // - Handling quoted values (like "value,with,commas")
    csvReader := csv.NewReader(r)

    // Step 2: Read and validate header
    headers, err := csvReader.Read()
    if err != nil {
        if err == io.EOF {
            return nil, fmt.Errorf("empty CSV file")
        }
        return nil, fmt.Errorf("reading header: %w", err)
    }

    // Expect exactly: id, category, amount
    if len(headers) != 3 || 
       headers[0] != "id" || 
       headers[1] != "category" || 
       headers[2] != "amount" {
        return nil, fmt.Errorf("invalid header: expected [id,category,amount], got %v", headers)
    }

    // Step 3: Initialize aggregation map
    stats := make(map[string]Stat)

    rowNum := 2  // Row 1 is header, data starts at row 2

    // Step 4: Process each row
    for {
        record, err := csvReader.Read()
        if err == io.EOF {
            break  // No more rows
        }
        if err != nil {
            return nil, fmt.Errorf("row %d: %w", rowNum, err)
        }

        // Validate row has 3 columns
        if len(record) != 3 {
            return nil, fmt.Errorf("row %d: expected 3 fields, got %d", rowNum, len(record))
        }

        // Extract fields
        category := record[1]
        amountStr := record[2]

        // Validate category
        if category == "" {
            return nil, fmt.Errorf("row %d: empty category", rowNum)
        }

        // Parse amount as float
        amount, err := strconv.ParseFloat(amountStr, 64)
        if err != nil {
            return nil, fmt.Errorf("row %d: invalid amount %q: %w", rowNum, amountStr, err)
        }

        // Update statistics
        s := stats[category]  // Get current stats (zero value if new)
        s.Count++ // 
        s.Sum += amount
        stats[category] = s  // Write back (required because Stat is a value type)

        rowNum++
    }

    // Step 5: Compute averages
    for category, s := range stats {
        if s.Count > 0 {
            s.Avg = s.Sum / float64(s.Count)
            stats[category] = s
        }
    }

    return stats, nil
}
```

## Key Concepts Explained

### Why `io.Reader` Instead of a File Path?

Our function accepts `io.Reader`, an interface that represents "anything you can read bytes from":
- Files
- Network connections
- In-memory strings
- Compressed data
- **Test data** (critical for testing!)

This makes the function flexible and testable.

### Why Check for `io.EOF`?

`EOF` (End Of File) is **not an error**—it's the normal way to signal "no more data". We need to distinguish between:
- `io.EOF`: Success, just finished reading
- Other errors: Actual problems (file corrupt, disk error, etc.)

### The Map Update Pattern

Go maps store **values**, not references. When we do:
```go
s := stats[category]  // Get copy of stats
s.Count++             // Modify copy
stats[category] = s   // Write copy back
```

We must write back because `s` is a **copy** of the struct, not a reference to it.

**Alternative with pointers**:
```go
stats := make(map[string]*Stat)  // Pointers as values
s := stats[category]
if s == nil {
    s = &Stat{}
    stats[category] = s
}
s.Count++  // Modifies original (no write-back needed)
```

### Floating-Point Precision

We use `float64` for amounts. Be aware:
- `0.1 + 0.2 = 0.30000000000000004` (in binary floating point!)
- For financial apps, consider using integer cents or decimal libraries

## Common Patterns You're Learning

### Pattern 1: CSV Reading with Validation
```go
csvReader := csv.NewReader(r)
headers, _ := csvReader.Read()
// Validate headers...
for {
    record, err := csvReader.Read()
    if err == io.EOF { break }
    // Process record...
}
```

### Pattern 2: String to Number Conversion
```go
amount, err := strconv.ParseFloat(amountStr, 64)
if err != nil {
    return fmt.Errorf("invalid amount: %w", err)
}
```

### Pattern 3: Accumulation in Map
```go
stats := make(map[string]Stat)
for each row {
    s := stats[key]
    s.Count++
    s.Sum += value
    stats[key] = s
}
```

## Real-World Applications

1. **Financial Analysis**: Tracking spending by category (Mint, YNAB)
2. **E-commerce**: Sales reports by product category
3. **Log Analysis**: Counting events by type/severity
4. **Scientific Data**: Aggregating measurements by experiment/condition
5. **Business Intelligence**: Any group-by aggregation query

## How to Run

```bash
# Prepare the project
cd minis/03-csv-stats/exercise
mv solution.go solution.go.reference

# Look at the test data
cat ../testdata/transactions.csv

# Run tests
go test -v

# Implement your solution in exercise.go
# Then test again
go test

# Run the demo program
cd ../..
make run P=03-csv-stats
```

## Common Mistakes to Avoid

1. **Not validating headers**: Always check the first row matches expectations
2. **Not handling blank categories**: Check `if category == ""`
3. **Forgetting to check field count**: CSV rows can have wrong number of columns
4. **Computing average in the loop**: Wait until end to divide sum by count
5. **Ignoring EOF vs real errors**: `io.EOF` is success, not failure!

## Stretch Goals

1. **Add median calculation**: Track all amounts per category to compute median
2. **Support different CSV formats**: Accept column order via configuration
3. **Add date filtering**: Include a date column, filter by date range
4. **Output JSON**: Marshal results to JSON for APIs
5. **Handle currency**: Parse amounts like "$12.50" (strip $ before parsing)

---

## Deep Dive: Advanced Implementation Details

This section provides a comprehensive analysis of the implementation, including complexity analysis, language comparisons, memory layout details, and execution traces.

### Problem Statement

Given a CSV with columns (id, category, amount), we need to:
1. Parse the CSV line-by-line (streaming for memory efficiency)
2. Group transactions by category
3. Compute count, sum, and average for each category
4. Handle malformed data gracefully

Constraints:
- CSV has a header row that must be validated
- Amounts are decimal numbers (use float64)
- Missing or invalid amounts should cause an error (fail-fast)
- Empty categories should be treated as an error

### Complexity Analysis

**Time Complexity: O(n)** where n = number of rows
- Single pass through the CSV (streaming)
- Each row: O(1) parse + O(1) map lookup/insert (amortized)
- Final average computation: O(c) where c = unique categories
- Total: O(n + c), but c ≤ n, so O(n)

**Space Complexity: O(c)** where c = number of unique categories
- We DON'T store all rows in memory (streaming)
- We only store one Stat struct per unique category
- If you have 1 million rows but 10 categories, we use ~10 × sizeof(Stat)

**Memory Layout of Stat struct:**
```go
type Stat struct {
    Count int      // 8 bytes on 64-bit systems
    Sum   float64  // 8 bytes
    Avg   float64  // 8 bytes
}
// Total: 24 bytes per category (no padding needed, fields are aligned)
```

### Why Go is Well-Suited for This Problem

1. **`encoding/csv` in stdlib**: No external dependencies for CSV parsing
   - Handles edge cases: quoted fields, escaped quotes, different line endings
   - Streaming by default (doesn't load entire file into memory)

2. **Strong typing**: Compile-time detection of struct field mismatches
   - If you try `s.count` instead of `s.Count`, compiler catches it
   - If you try `s.Count = "five"`, compiler catches it

3. **Explicit error handling**: No silent data corruption
   - Every function that can fail returns an error
   - You MUST handle it (or explicitly ignore with `_`)
   - Compare to Python: `int("invalid")` throws at runtime, often unhandled

4. **Predictable memory**: No hidden allocations, no GC surprises for this workload
   - We know exactly when allocations happen (make, new, string parsing)
   - For a streaming CSV, memory stays constant regardless of file size

### Language Comparison

#### Python (pandas)
```python
df = pd.read_csv('transactions.csv')
stats = df.groupby('category')['amount'].agg(['count','sum','mean'])
```

**Pros**: One-liner, powerful analytics (median, std dev, percentiles)

**Cons**:
- Loads ENTIRE file into memory (100MB CSV → 500MB+ RAM)
- DataFrame has overhead: index, column metadata, dtypes
- Slower for large files (pandas overhead, Python interpreter)
- Dynamic typing hides schema errors until runtime

#### JavaScript (Node.js)
```javascript
const csv = require('csv-parser');
const stats = {};
fs.createReadStream('transactions.csv')
  .pipe(csv())
  .on('data', (row) => { ... })
  .on('end', () => { ... });
```

**Pros**: Streaming support, widely known syntax

**Cons**:
- Async complexity (callbacks/promises/async-await)
- Requires external library (csv-parser)
- Dynamic typing misses schema errors
- Event-driven model harder to reason about for sequential processing

#### Rust
```rust
use csv::ReaderBuilder;
let mut stats: HashMap<String, (usize, f64)> = HashMap::new();
for result in rdr.records() {
    let record = result?;
    let entry = stats.entry(record[1].to_string()).or_insert((0, 0.0));
    entry.0 += 1;
    entry.1 += record[2].parse::<f64>()?;
}
```

**Pros**:
- Zero-copy parsing possible (csv crate can return &str slices)
- Faster execution (no GC, better optimization)
- `?` operator for ergonomic error handling

**Cons**:
- Ownership rules more complex (borrow checker learning curve)
- `entry()` API less intuitive than Go's read-modify-write

#### SQL
```sql
SELECT category,
       COUNT(*) as count,
       SUM(amount) as sum,
       AVG(amount) as avg
FROM transactions
GROUP BY category;
```

**Pros**:
- Declarative (say WHAT you want, not HOW)
- Query planner optimizes execution
- Scales to billions of rows (with proper indexing)

**Cons**:
- Requires database setup/infrastructure
- Less portable than a single CSV file
- Overkill for ad-hoc analysis

**Go's sweet spot**: Single-binary deployment, "fast enough" for millions of rows, readable code, strong guarantees at compile time.

### Data Structures Deep Dive

#### Why struct instead of separate variables?

```go
type Stat struct {
    Count int
    Sum   float64
    Avg   float64
}
```

**Benefits**:
- **Logical grouping**: These 3 values always belong together
- **Single map lookup**: `stats[category]` returns all 3 values at once
- **Cache locality**: When CPU loads Count, Sum and Avg are likely in same cache line
- **Type safety**: Can't accidentally mix up Count from one category with Sum from another

#### Why exported fields (Capital letters)?

- Allows callers to read the results: `stat.Count`, `stat.Sum`, `stat.Avg`
- If we used lowercase, only this package could access them
- For a return value, we WANT callers to read the fields

### Why io.Reader Instead of *os.File?

The `io.Reader` interface is defined as:
```go
type Reader interface {
    Read(p []byte) (n int, err error)
}
```

ANY type that has a Read method with this signature satisfies io.Reader.

This is **DEPENDENCY INJECTION via interfaces**:
- Function depends on BEHAVIOR (Read), not IMPLEMENTATION (*os.File)
- Caller decides what to inject: file, network, buffer, mock

**Concrete types that implement io.Reader**:
- `*os.File` - Read from filesystem
- `*bytes.Buffer` - Read from in-memory buffer (great for testing!)
- `*strings.Reader` - Read from string (great for testing!)
- `http.Response.Body` - Read from HTTP response
- `*gzip.Reader` - Read from compressed stream
- `io.LimitReader` - Read with size limit

**Testing benefit**:
```go
// Production code:
file, _ := os.Open("transactions.csv")
stats, err := SummarizeCSV(file)

// Test code (no file needed!):
testData := "id,category,amount\n1,groceries,10.00\n"
stats, err := SummarizeCSV(strings.NewReader(testData))
```

### Execution Trace: Three Scenarios

#### Scenario 1: Valid CSV (Happy Path)

**Input:**
```csv
id,category,amount
1,groceries,12.50
2,groceries,7.50
3,books,10.00
```

**Step-by-step execution:**

1. **`csv.NewReader(r)`**
   - Memory: Allocates csv.Reader struct on heap (~200 bytes)
   - Contains: bufio.Reader (4096 byte buffer), field buffer, line buffer

2. **`csvReader.Read()` → headers = ["id", "category", "amount"]**
   - Memory: Allocates []string with 3 elements, plus 3 string headers
   - The strings point into csvReader's internal buffer (no copy yet)

3. **Header validation passes**

4. **`stats := make(map[string]Stat)`**
   - Memory: Allocates empty hash map structure (~48 bytes initially)
   - Contains: pointer to buckets array, count, hash seed

5. **First iteration: `csvReader.Read()` → ["1", "groceries", "12.50"]**
   - `category = "groceries"`
   - `strconv.ParseFloat("12.50", 64)` → 12.50 (no allocation, returns float64)
   - `s := stats["groceries"]` → s = Stat{0, 0.0, 0.0} (zero value, key doesn't exist)
   - `s.Count++` → s = Stat{1, 0.0, 0.0}
   - `s.Sum += 12.50` → s = Stat{1, 12.50, 0.0}
   - `stats["groceries"] = s` → Map allocates bucket, stores copy of s
   - Map state: `{"groceries": {Count:1, Sum:12.50, Avg:0.0}}`

6. **Second iteration: `csvReader.Read()` → ["2", "groceries", "7.50"]**
   - `category = "groceries"`
   - `amount = 7.50`
   - `s := stats["groceries"]` → s = Stat{1, 12.50, 0.0} (COPY from map)
   - `s.Count++` → s = Stat{2, 12.50, 0.0}
   - `s.Sum += 7.50` → s = Stat{2, 20.00, 0.0}
   - `stats["groceries"] = s` → Overwrites existing entry
   - Map state: `{"groceries": {Count:2, Sum:20.00, Avg:0.0}}`

7. **Third iteration: `csvReader.Read()` → ["3", "books", "10.00"]**
   - `category = "books"`
   - `amount = 10.00`
   - `s := stats["books"]` → s = Stat{0, 0.0, 0.0} (zero value, new key)
   - `s.Count++` → s = Stat{1, 0.0, 0.0}
   - `s.Sum += 10.00` → s = Stat{1, 10.00, 0.0}
   - `stats["books"] = s` → Map adds new entry
   - Map state: `{"groceries": {Count:2, Sum:20.00, Avg:0.0}, "books": {Count:1, Sum:10.00, Avg:0.0}}`

8. **Fourth iteration: `csvReader.Read()` → io.EOF**
   - Loop exits

9. **Average computation pass:**
   - for "groceries": `s.Avg = 20.00 / 2 = 10.00`; `stats["groceries"] = s`
   - for "books": `s.Avg = 10.00 / 1 = 10.00`; `stats["books"] = s`
   - **Final**: `{"groceries": {2, 20.00, 10.00}, "books": {1, 10.00, 10.00}}`

#### Scenario 2: Empty CSV (Edge Case)

**Input:**
```csv
id,category,amount
(no data rows)
```

**Execution:**
1. Read header → ["id", "category", "amount"]
2. Read next row → io.EOF immediately
3. Loop never executes
4. Average pass: range over empty map → no iterations
5. Return empty `map{}`, nil

This is **VALID behavior**: empty file = empty results, not an error.

#### Scenario 3: Malformed Amount (Failure Case)

**Input:**
```csv
id,category,amount
1,groceries,12.50
2,books,invalid
```

**Execution:**
1. Read header → OK
2. First row → OK, `stats = {"groceries": {1, 12.50, 0.0}}`
3. Second row: `strconv.ParseFloat("invalid", 64)` → ERROR
4. Return `nil, error("row 3: invalid amount \"invalid\": ...")`

**FAIL-FAST behavior**: We stop immediately on first error.

### The Map Aggregation Pattern (Critical!)

This is the **CRITICAL section** to understand.

**FUNDAMENTAL RULE**: Go maps store VALUES, not references. When you write `s := stats[category]`, you get a COPY of the Stat.

**Step-by-step memory trace:**

1. **`s := stats[category]`**

   Memory BEFORE (first time seeing "groceries"):
   ```
   stats (map) → {"books": Stat{...}} (no "groceries" key)
   ```

   What happens:
   - Map performs hash lookup for "groceries"
   - Key not found → return ZERO VALUE of Stat
   - `s = Stat{Count: 0, Sum: 0.0, Avg: 0.0}`
   - s is a LOCAL VARIABLE on the STACK (24 bytes)

   Memory AFTER:
   ```
   stats (map) → {"books": Stat{...}} (unchanged!)
   s (stack)   → Stat{0, 0.0, 0.0}
   ```

2. **`s.Count++`**

   Memory BEFORE:
   ```
   s (stack) → Stat{0, 0.0, 0.0}
   ```

   What happens:
   - Increment the Count field of the LOCAL copy
   - Map is NOT modified

   Memory AFTER:
   ```
   stats (map) → {"books": Stat{...}} (still unchanged!)
   s (stack)   → Stat{1, 0.0, 0.0}
   ```

3. **`s.Sum += amount`**

   Memory BEFORE:
   ```
   s (stack) → Stat{1, 0.0, 0.0}
   amount = 12.50
   ```

   What happens:
   - Add amount to the Sum field of the LOCAL copy

   Memory AFTER:
   ```
   stats (map) → {"books": Stat{...}} (still unchanged!)
   s (stack)   → Stat{1, 12.50, 0.0}
   ```

4. **`stats[category] = s`**

   Memory BEFORE:
   ```
   stats (map) → {"books": Stat{...}}
   s (stack)   → Stat{1, 12.50, 0.0}
   ```

   What happens:
   - Map performs hash lookup for "groceries"
   - Key not found → allocate new bucket entry
   - COPY s into the map's value storage
   - (If key existed, OVERWRITE the existing value)

   Memory AFTER:
   ```
   stats (map) → {"books": Stat{...}, "groceries": Stat{1, 12.50, 0.0}}
   s (stack)   → Stat{1, 12.50, 0.0} (will be reused next iteration)
   ```

#### Why Can't We Do `stats[category].Count++`?

Go prohibits this! You'll get a compile error:
```
cannot assign to stats[category].Count
```

**Why?** Map values are NOT ADDRESSABLE.

The Go runtime might RELOCATE map entries when:
- Map grows (load factor exceeded)
- Buckets are reorganized

If you could take the address of `stats[category]`, it would become a DANGLING POINTER after relocation. Go prevents this at compile time.

#### Alternative: Use *Stat as value type

```go
stats := make(map[string]*Stat)

s := stats[category]
if s == nil {
    s = &Stat{}          // Allocate new Stat on HEAP
    stats[category] = s  // Store pointer in map
}
s.Count++  // Now this works! We're modifying the heap object
s.Sum += amount
// No write-back needed—we modified the object in place
```

**Trade-offs:**
- ✅ Modify in-place (no write-back)
- ❌ More allocations (one `&Stat{}` per category, on heap)
- ❌ Nil checks required
- ❌ Pointer indirection (cache miss when accessing s.Count)

For small structs like Stat (24 bytes), value semantics are usually better.

### Why Compute Averages in a Separate Pass?

**Option A (compute in loop):**
```go
s := stats[category]
s.Count++
s.Sum += amount
s.Avg = s.Sum / float64(s.Count)  // Redundant! Computed every row
stats[category] = s
```

**Option B (separate pass):**
```go
// After loop:
for k, s := range stats {
    s.Avg = s.Sum / float64(s.Count)  // Computed once per category
    stats[k] = s
}
```

Option A does division on EVERY row.
Option B does division once per CATEGORY.

**If you have 1 million rows and 10 categories:**
- Option A: 1,000,000 divisions
- Option B: 10 divisions

The difference is negligible in practice, but Option B is cleaner and makes the intent clear: "Average is a derived value."

### Why rowNum Starts at 2

This is about the difference between **programming index conventions** and **user-friendly row numbers** in CSV files.

**CSV File Structure (Human Perspective):**
```
Row 1: header1,header2,header3    ← Row 1 = header
Row 2: value1,value2,value3        ← Row 2 = first data row
Row 3: value1,value2,value3        ← Row 3 = second data row
```

**Why Not Start at 0?**
- Index 0 doesn't match how humans think about rows in spreadsheets/text editors
- If there's an error, saying "error on row 0" is confusing

**Why Not Start at 1?**
- Index 1 would be the header row
- The first **data** row comes after the header, which is row 2
- When your CSV reader skips the header, you're already past row 1

**The Goal: Error Messages**

This is for error reporting. If there's a problem with the first data row:
```go
return nil, fmt.Errorf("invalid value on row %d", rowNum)
// → "invalid value on row 2"
```

This matches what users see in Excel, Google Sheets, or their text editor (line 2 of the file), making debugging much easier!

### Alternatives & Trade-offs

#### 1. Pointer values in map

```go
stats := make(map[string]*Stat)
s := stats[category]
if s == nil {
    s = &Stat{}
    stats[category] = s
}
s.Count++
s.Sum += amount
```

**Pros**: Modify in-place (no write-back needed)

**Cons**:
- Extra heap allocation per category
- Nil checks required
- Pointer indirection (potential cache miss)
- Less idiomatic for small structs

#### 2. Accumulate errors instead of failing fast

```go
var errs []error
for {
    // ... on parse error:
    errs = append(errs, fmt.Errorf("row %d: %w", rowNum, err))
    continue  // Don't return, keep going
}
if len(errs) > 0 {
    return nil, errors.Join(errs...)  // Go 1.20+
}
```

**Pros**: Process entire file, report all errors at once

**Cons**:
- More complex error handling
- May hide systemic issues (e.g., all rows are malformed)
- Partial results: should we return the valid rows?

#### 3. Use integer cents instead of float64

Floating-point has rounding errors:
```
0.1 + 0.2 = 0.30000000000000004 (in binary!)
```

For financial data, store as integer cents:
```go
type Stat struct {
    Count    int
    SumCents int64
    AvgCents int64
}

// Parse "12.50" → 1250 cents
amountCents := int64(amount * 100)
s.SumCents += amountCents

// Convert back for display
fmt.Printf("$%.2f", float64(s.SumCents) / 100.0)
```

**Pros**: Exact arithmetic, no rounding errors

**Cons**: More code, still need float64 for display, need to handle rounding in conversion

#### 4. Concurrent processing

For VERY large files, we could split into chunks and process in parallel:

```go
// Pseudo-code
chunks := splitFileIntoChunks(file, numCPU)
results := make(chan map[string]Stat, numCPU)
for _, chunk := range chunks {
    go func(c Chunk) {
        results <- processChunk(c)
    }(chunk)
}
finalStats := mergeResults(results)
```

**Pros**: Faster on multi-core for huge files (>100MB)

**Cons**:
- Much more complex code
- Need to handle chunk boundaries (rows split across chunks)
- Memory overhead of multiple maps
- Overkill for typical CSV sizes

#### 5. Using csv.Reader.ReuseRecord

```go
csvReader.ReuseRecord = true
```

**What it does**:
- Reuses the same []string slice for each Read()
- Strings inside the slice point to internal buffer

**Pros**: Fewer allocations per row

**Cons**:
- `record` becomes invalid after next Read()
- Category would need to be copied manually: `category := strings.Clone(record[1])`
- Error-prone if you forget to clone

For our use case, the map key assignment already copies the string, so ReuseRecord would have minimal benefit.
