# Project 01: hello-strings

## What Is This Project About?

Imagine you're writing a program that needs to work with text from around the world—English, Spanish, Japanese, emoji, and more. This project teaches you how to manipulate strings (text) in a way that respects all these different writing systems.

You'll build three utilities:
1. **TitleCase**: Convert text to title case (First Letter Of Each Word Capitalized)
2. **Reverse**: Flip text backwards while keeping special characters intact
3. **RuneLen**: Count the actual number of characters (not bytes) in text

## The Fundamental Problem: Text Isn't Simple

### First Principles: What Is Text in a Computer?

When you type "Hello" on your keyboard, your computer doesn't store the word "Hello" directly. Instead, it stores **numbers**. Each letter gets converted to a number:
- 'H' → 72
- 'e' → 101
- 'l' → 108
- 'l' → 108
- 'o' → 111

This conversion system is called **character encoding**. The most common encoding today is called **UTF-8**.

### The Critical Insight: Bytes vs Characters

Here's where it gets tricky. Simple English letters take **1 byte** (one number) to store. But many characters from other languages or symbols like emoji take **multiple bytes**:

- 'A' = 1 byte
- 'é' (e with accent) = 2 bytes
- '日' (Japanese kanji) = 3 bytes
- '😀' (emoji) = 4 bytes

**This is the most important concept in this project**: In Go, when you count the length of a string with `len()`, you get the **number of bytes**, not the **number of characters**.

Example:
```
"Hello"   → len() = 5 (5 bytes, 5 characters) ✓
"café"    → len() = 5 (5 bytes, but only 4 characters!) ✗
"Hello👋" → len() = 9 (9 bytes, but only 6 characters!) ✗
```

In Go terminology:
- **Byte**: A single number (0-255) representing part of a character
- **Rune**: A complete character (what humans think of as "one character")

## Problem 1: TitleCase - Making The First Letter Of Each Word Capital

### The Human-Level Problem

You have text like: `"hello world"` and you want: `"Hello World"`

Seems simple, right? But consider edge cases:
- Multiple spaces: `"hello    world"` → `"Hello World"` (collapse spaces)
- Already capitalized: `"Hello World"` → `"Hello World"` (don't break it)
- Special characters: `"café résumé"` → `"Café Résumé"` (handle accents)
- Emoji: `"hello 👋 world"` → `"Hello 👋 World"` (leave emoji alone)

### Breaking Down the Solution (First Principles)

**Step 1: Split the text into words**
How do we find where one word ends and another begins? By looking for **whitespace** (spaces, tabs, newlines).

`"hello world"` → `["hello", "world"]`

In Go, we use `strings.Fields()` which automatically:
- Splits on ANY whitespace (not just spaces)
- Removes empty strings from the result
- Handles multiple consecutive spaces

**Step 2: Capitalize the first character of each word**
For each word like "hello", we need to:
1. Get the first character (but remember—it might be multiple bytes!)
2. Make it uppercase
3. Keep the rest of the word as-is

**The Rune Approach**: Convert the word to a slice of runes (characters), then we can safely access the first rune:
```
"hello" → ['h', 'e', 'l', 'l', 'o'] (as runes)
         → ['H', 'e', 'l', 'l', 'o'] (capitalize first)
         → "Hello" (convert back)
```

**Step 3: Join the words back together**
Put the words back together with single spaces between them.

### Translating to Code

```go
func TitleCase(s string) string {
    // Step 1: Split into words
    words := strings.Fields(s)

    // Step 2: Process each word
    for i, word := range words {
        // Convert to runes so we can work with characters, not bytes
        runes := []rune(word)

        if len(runes) > 0 {
            // Capitalize first rune
            runes[0] = unicode.ToUpper(runes[0])
        }

        // Convert back to string
        words[i] = string(runes)
    }

    // Step 3: Join with spaces
    return strings.Join(words, " ")
}
```

## Problem 2: Reverse - Flipping Text Backwards

### The Human-Level Problem

You want to reverse "Hello" to get "olleH". Easy for English. But what about "Hello👋"?

**Wrong approach** (reversing bytes):
```
"Hello👋" (bytes: H e l l o [4 bytes for 👋])
Reversed bytes: [4 bytes reversed] o l l e H
Result: "���olleH" (BROKEN EMOJI!)
```

**Right approach** (reversing characters):
```
"Hello👋" (characters: H e l l o 👋)
Reversed characters: 👋 o l l e H
Result: "👋olleH" (EMOJI INTACT!)
```

### Breaking Down the Solution

**Step 1: Convert string to runes (characters)**
This separates the string into actual characters, where each emoji is ONE rune, not four bytes.

`"Hi👋"` → `['H', 'i', '👋']` (3 runes)

**Step 2: Reverse the rune slice**
Use the classic "two-pointer" algorithm:
- Start with a pointer at the beginning and one at the end
- Swap the elements they point to
- Move the pointers toward each other
- Stop when they meet in the middle

```
['H', 'i', '👋']
 ↑           ↑     Swap H and 👋
['👋', 'i', 'H']
      ↑↑          Pointers meet, done!
```

**Step 3: Convert back to string**

### Translating to Code

```go
func Reverse(s string) string {
    // Step 1: Convert to runes
    runes := []rune(s)

    // Step 2: Reverse with two pointers
    for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
        runes[i], runes[j] = runes[j], runes[i]  // Swap
    }

    // Step 3: Convert back
    return string(runes)
}
```

**Why the condition `i < j`?**
When `i` and `j` meet or cross, we've swapped all pairs. If we continued, we'd reverse the reverse (undoing our work)!

## Problem 3: RuneLen - Counting Actual Characters

### The Human-Level Problem

When a human looks at "café", they see **4 characters**: c, a, f, é.
But Go's `len()` returns **5** because 'é' is stored as 2 bytes.

We need a function that counts what humans see as characters.

### The Core Concept: What Is a Rune?

In Go, a **rune** is an integer that represents a single Unicode character. It's defined as `type rune = int32`.

Unicode is a giant table that assigns a unique number to every character in every language:
- 'A' = 65
- '日' = 26085
- '😀' = 128512

So when we count runes, we're counting these unique character codes, which matches human perception.

### Breaking Down the Solution

We need to iterate through the string and count each complete character (rune), not each byte.

**Approach 1: Convert to []rune and get length**
```go
return len([]rune(s))
```
This works but allocates a new slice in memory (potentially wasteful for large strings).

**Approach 2: Use the standard library**
Go provides `utf8.RuneCountInString()` which:
- Scans through the string
- Recognizes multi-byte sequences
- Counts complete characters
- Doesn't allocate extra memory

### Translating to Code

```go
func RuneLen(s string) int {
    return utf8.RuneCountInString(s)
}
```

That's it! The standard library does the heavy lifting.

**Alternative (educational) implementation**:
```go
func RuneLen(s string) int {
    count := 0
    for range s {  // 'for range' on a string iterates over RUNES, not bytes!
        count++
    }
    return count
}
```

The `for range` loop over a string automatically decodes UTF-8 and gives you runes!

## The Bigger Picture: Why This Matters

### Real-World Applications

1. **Internationalization (i18n)**: Apps that work in multiple languages MUST handle UTF-8 correctly
2. **Social Media**: User names and posts contain emoji—reversed/truncated incorrectly looks broken
3. **Text Analysis**: Word count, character limits (Twitter), search—all need proper character handling
4. **Security**: Incorrectly handling UTF-8 can lead to vulnerabilities

### What You're Learning

- **Go's explicit byte/rune distinction** forces you to think about encoding (prevents bugs)
- **Standard library utilities** (`strings`, `unicode`, `utf8`) are your friends
- **Rune slices** are the key to safe character-level manipulation
- **Immutability** of strings means we create new strings, not modify existing ones

## How to Run

```bash
# Prepare the project (rename solution so it doesn't conflict with your code)
cd minis/01-hello-strings/exercise
mv solution.go solution.go.reference

# Run tests to see what's expected
go test

# Implement your solution in exercise.go, then test again
# Edit exercise.go...
go test

# Run the demo program
cd ../..
make run P=01-hello-strings
```

## Common Mistakes to Avoid

1. **Using `len()` to count characters**: `len()` counts bytes, use `utf8.RuneCountInString()`
2. **Reversing bytes instead of runes**: Always convert to `[]rune` first
3. **Forgetting `[]rune` allocates**: For performance-critical code, consider alternatives
4. **Not testing with emoji and accented characters**: Your tests should include these edge cases

## Stretch Goals

Once you've completed the basics:

1. **Palindrome Checker**: `IsPalindrome("racecar")` → `true`
   - Hint: Use your `Reverse` function and compare with the original

2. **Word Truncation**: `TruncateWords("Hello world from Go", 3)` → `"Hello world from..."`
   - Handle edge cases: what if the text is already short?

3. **Vowel Counter**: `CountVowels("Hello")` → `2` (e, o)
   - Remember to check both uppercase and lowercase vowels

4. **Performance Benchmark**: Compare byte-based vs rune-based reversal for ASCII-only text
   - Learn to use Go's `testing.B` benchmarking tools
