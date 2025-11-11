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
        s.Count++
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
