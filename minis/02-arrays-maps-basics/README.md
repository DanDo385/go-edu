# Project 02: arrays-maps-basics

## What Is This Project About?

Imagine you have a long document and you want to know which words appear most frequently. Maybe you're analyzing customer feedback, studying a book, or processing survey responses. This project teaches you how to count things efficiently using Go's built-in data structures.

You'll build a word frequency analyzer that:
1. Reads text input (one word per line) from a file
2. Counts how many times each word appears
3. Finds the most common word

## The Fundamental Problem: Counting Things Efficiently

### First Principles: What Is Counting?

At the simplest level, counting means keeping track of how many times you see something. If you have a list of fruits:

```
apple
banana
apple
orange
apple
banana
```

You could count them by hand:
- apple: 3
- banana: 2
- orange: 1

But how do we make a computer do this efficiently?

### The Naive Approach (What NOT to Do)

You might think: "Let's use a list and search through it every time!"

```
words = ["apple", "banana", "apple"]
counts = [1, 1, 2]  // apple appears twice
```

When you see "orange", you search the entire `words` list to see if it's there. If not, add it.

**Problem**: Searching takes time proportional to the list length. For 1 million words, that's slow!

### The Smart Approach: Maps (Hash Tables)

A **map** is a data structure that links **keys** to **values**, like a real dictionary:

```
"apple" → 3
"banana" → 2
"orange" → 1
```

The magic: Looking up a key in a map is nearly instant (O(1) time complexity), regardless of how many items are in the map!

In Go, maps are written as: `map[KeyType]ValueType`

For our problem: `map[string]int` means "map from word (string) to count (integer)"

## Breaking Down the Problem (First Principles)

Let's solve this step-by-step, as if explaining to someone who's never programmed before.

### Step 1: Reading Text Line-by-Line

**The Human Process**:
Imagine reading a physical document one line at a time with your finger. You move your finger down, read a line, process it, move to the next line, repeat until you reach the end.

**The Computer Process**:
We do exactly the same! But we use a **scanner** (like a cursor) that reads through the file:

```go
scanner := bufio.NewScanner(file)
for scanner.Scan() {  // Move to next line, return true if there's a line
    line := scanner.Text()  // Get the text of current line
    // Process the line...
}
```

**Why not read the entire file at once?**
For a 1GB log file with millions of words, reading it all into memory would crash your program! Scanning line-by-line uses constant memory.

### Step 2: Normalizing Words

**The Problem**: Is "Hello" the same as "hello"? For counting, we usually want them to be the same.

**The Solution**: Convert everything to lowercase:
```
"Hello" → strings.ToLower() → "hello"
"HELLO" → strings.ToLower() → "hello"
```

Also, remove extra whitespace:
```
"  hello  " → strings.TrimSpace() → "hello"
```

### Step 3: Counting with a Map

Here's where maps shine. The pattern is beautifully simple:

```go
freq := make(map[string]int)  // Create empty map

for each word {
    freq[word]++  // Increment count for this word
}
```

**Wait, what if the word isn't in the map yet?**

This is Go's magic: When you access a map key that doesn't exist, Go returns the **zero value** for that type. For `int`, the zero value is `0`.

So:
```go
freq["apple"]++  // First time: 0 + 1 = 1
freq["apple"]++  // Second time: 1 + 1 = 2
freq["apple"]++  // Third time: 2 + 1 = 3
```

No need to check "if word exists, increment, else set to 1"—Go handles it!

### Step 4: Finding the Maximum

Once we have our frequency map:
```
"apple" → 3
"banana" → 2
"orange" → 1
```

How do we find which word has the highest count?

**The Algorithm**:
```
1. Start with maxCount = 0 and maxWord = ""
2. For each word and its count in the map:
   - If count > maxCount:
     - Update maxCount = count
     - Update maxWord = word
3. Return maxWord
```

In code:
```go
maxCount := 0
maxWord := ""

for word, count := range freq {
    if count > maxCount {
        maxCount = count
        maxWord = word
    }
}
```

## The Complete Solution (Explained)

Let's walk through the entire function:

```go
func FreqFromReader(r io.Reader) (map[string]int, string, error) {
    // Create empty frequency map
    freq := make(map[string]int)

    // Create scanner to read line-by-line
    scanner := bufio.NewScanner(r)

    // Process each line
    for scanner.Scan() {
        line := scanner.Text()

        // Normalize: lowercase and trim whitespace
        word := strings.ToLower(strings.TrimSpace(line))

        // Skip blank lines
        if word == "" {
            continue
        }

        // Increment count (works even if word not yet in map!)
        freq[word]++
    }

    // Check for reading errors
    if err := scanner.Err(); err != nil {
        return nil, "", err
    }

    // Find most common word
    maxCount := 0
    maxWord := ""

    for word, count := range freq {
        if count > maxCount {
            maxCount = count
            maxWord = word
        }
    }

    return freq, maxWord, nil
}
```

## Key Concepts Explained

### What is `io.Reader`?

Instead of requiring a specific file type, we accept anything that implements the `io.Reader` **interface**. This means our function works with:
- Files (`os.Open()`)
- Strings (`strings.NewReader()`)
- Network connections
- Compressed data streams
- Any custom reader!

This is called **dependency injection**—the caller decides what to read from.

### What is `bufio.Scanner`?

Think of it as a smart wrapper around a reader that:
- Automatically handles line endings (`\n`, `\r\n`, etc.)
- Buffers data for efficiency (reads chunks, not byte-by-byte)
- Provides a simple `Scan()` → `Text()` API

### Why Return Multiple Values?

Go functions can return multiple values:
```go
func DoSomething() (result int, success bool, err error)
```

Our function returns:
1. `map[string]int`: The full frequency map (for detailed analysis)
2. `string`: The most common word (quick answer)
3. `error`: Any errors that occurred (nil if successful)

This is idiomatic Go—errors are values, not exceptions!

## Common Patterns You're Learning

### Pattern 1: Map Increment Pattern
```go
m := make(map[string]int)
m[key]++  // Safe even if key doesn't exist!
```

### Pattern 2: Scanner Pattern
```go
scanner := bufio.NewScanner(r)
for scanner.Scan() {
    line := scanner.Text()
    // Process line...
}
if err := scanner.Err(); err != nil {
    // Handle error
}
```

### Pattern 3: Finding Maximum in Map
```go
maxKey, maxVal := "", 0
for k, v := range m {
    if v > maxVal {
        maxKey, maxVal = k, v
    }
}
```

## Why This Approach Is Efficient

### Time Complexity Analysis

For N words:
- **Reading**: O(N) — must read each word once
- **Inserting into map**: O(1) per word → O(N) total
- **Finding max**: O(U) where U = unique words (usually U << N)

**Total: O(N)** — linear time, the best possible!

Compare to naive approach with a list:
- Search list each time: O(N) per word → O(N²) total
- For 1 million words: 1 trillion operations vs 1 million!

### Space Complexity

We store one entry per unique word: **O(U)** where U = unique words.

For most text, unique words are a small fraction of total words (maybe 5-10% for English).

## Real-World Applications

1. **Search Engines**: Term frequency is used in ranking algorithms (TF-IDF)
2. **Spam Detection**: Analyzing word frequency patterns to detect spam
3. **Text Summarization**: Identifying key terms by frequency
4. **Customer Feedback**: Finding most mentioned products/issues
5. **Log Analysis**: Counting error types in server logs

## How to Run

```bash
# Prepare the project
cd minis/02-arrays-maps-basics/exercise
mv solution.go solution.go.reference

# Look at the test data
cat ../testdata/input.txt

# Run tests
go test

# Implement your solution in exercise.go
# Then test again
go test

# Run the demo program
cd ../..
make run P=02-arrays-maps-basics
```

## Testing Strategy

The tests use `strings.NewReader()` to create an in-memory reader from a string:

```go
input := "hello\nworld\nhello\n"
r := strings.NewReader(input)
freq, mostCommon, err := FreqFromReader(r)
```

This is faster than creating real files and makes tests deterministic!

## Common Mistakes to Avoid

1. **Forgetting to normalize**: "Hello" and "hello" count separately if you don't use `ToLower()`
2. **Not handling blank lines**: Always check `if word == ""` after trimming
3. **Reading entire file into memory**: Use `bufio.Scanner`, not `ioutil.ReadAll()` for large files
4. **Panicking on missing keys**: Go's zero-value semantics handle this, but other languages don't!
5. **Iterating map multiple times**: You can find max in one pass through the map

## Stretch Goals

1. **Top N Words**: Instead of just the most common, return the top 5 words
   - Hint: You'll need to sort the map by value (requires converting to a slice of key-value pairs)

2. **Ignore Stop Words**: Filter out common words like "the", "a", "is", "are"
   - Create a set (map[string]bool) of stop words to exclude

3. **Case-Sensitive Mode**: Add a parameter to make counting case-sensitive
   - Compare "Hello" vs "hello" as different words

4. **Word Length Statistics**: Also compute average word length, min, and max
   - Process length alongside frequency

5. **Benchmark**: Compare different approaches
   - `bufio.Scanner` vs `ioutil.ReadAll()` + `strings.Split()`
   - Map vs slice of structs for storage
