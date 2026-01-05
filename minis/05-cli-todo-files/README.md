# 05: Command-Line Todo App with File Persistence

## The Big Picture: Building a Complete Application

So far, you've been building libraries—small, reusable pieces of code that process data. Now, you'll assemble those skills to build your first **complete, interactive application**: a command-line (CLI) todo manager.

This project is a major step. It introduces the fundamental challenge of **persistence**: how does your application remember data between runs? You'll solve this by creating a data "store" that saves your todo items to a JSON file, introducing you to application architecture, interfaces, and the full lifecycle of data management.

## First Principles: CLIs, Persistence, and CRUD

1.  **Command-Line Interface (CLI)**: A way for users to interact with a program by typing text commands. A CLI takes input via **arguments** (positional inputs) and **flags** (named, optional inputs like `-verbose` or `-user=admin`). Designing a good CLI is a UI/UX challenge: it should be intuitive, helpful, and "do what the user expects."

2.  **Persistence**: The ability of an application to save its **state** (data) so that it survives when the application closes.
    *   **In-Memory Storage**: Data is stored in RAM (e.g., in a slice or map). It is fast but **volatile**—it disappears when the program ends.
    *   **Persistent Storage**: Data is stored on a non-volatile medium like a hard drive. Our choice here, a simple JSON file, is a very common and effective persistence strategy.

3.  **CRUD Operations**: A foundational concept in data management, representing the four basic functions you can perform on data:
    *   **C**reate: Add a new todo item.
    *   **R**ead: List the existing todo items.
    *   **U**pdate: Mark a todo item as complete or not complete.
    *   **D**elete: Remove a todo item (in our case, we'll focus on the 'Update' part).

This project requires you to implement a full CRUD lifecycle for your todo items.

## Key Go Concepts in This Project

### Parsing User Input: `os.Args` and the `flag` Package

Go gives you two main ways to read from the command line:
*   `os.Args`: A raw slice of strings (`[]string`) containing everything the user typed. `os.Args[0]` is the program name, `os.Args[1]` is the first argument, and so on. It's simple but requires manual parsing and validation.
*   `flag` package: The idiomatic Go way to define and parse command-line flags. It's more robust, handles type conversions (e.g., parsing strings into `bool` or `int`), and automatically generates help text. This is the preferred method for any real application.

### Persistence with `os.ReadFile` and `os.WriteFile`

You'll use two key functions from the `os` package:
*   `os.ReadFile(path)`: Reads an entire file into a byte slice (`[]byte`). It's simple and perfect for configuration files or small data files like our todo list.
*   `os.WriteFile(path, data, perms)`: Writes a byte slice to a file with specified permissions (e.g., `0644`).

**Pro Tip: Atomic Writes**. A simple `os.WriteFile` can be risky. If your program crashes mid-write, the file can become corrupted and empty. A more robust "atomic" pattern is to write to a temporary file first, and only if that succeeds, rename the temporary file to the final destination. This ensures you never have a partially written file.

### Interfaces for Abstraction: The `Store`

This is the most important architectural concept in this project. We define a `Store` **interface** that describes *what* our data store must be able to do (Load, Save, Add, etc.).

```go
type Store interface {
    Load() error
    Save() error
    Add(text string) Item
    // ...
}
```

We then create a `fileStore` **struct** that *implements* this interface.

**Why is this so powerful?**
1.  **Testability**: In our tests, we can create a fake "mock" store that doesn't touch the filesystem, making tests fast and reliable.
2.  **Flexibility**: The `main` part of our application only knows about the `Store` interface. It doesn't care if the data is stored in a JSON file, a database, or a cloud service. We can completely change the persistence layer later without changing the CLI logic, simply by creating a new struct that satisfies the `Store` interface. This is a form of **dependency injection**.

### JSON for Data Serialization

This project puts your JSON skills from Project 04 to practical use.
*   `json.MarshalIndent`: To save the `[]Item` slice to the file, you'll "marshal" it into a nicely formatted JSON byte slice.
*   `json.Unmarshal`: To load the app, you'll read the file and "unmarshal" the JSON data back into the `[]Item` slice.

This process of converting in-memory objects to a storable format is called **serialization**.

## Progression: Tying It All Together

This project is a synthesis of everything you've learned:
*   You'll manage text with **strings** (Project 01).
*   Your core data structure will be a **slice** of `Item` structs (Project 02).
*   You'll use the `encoding/json` package to **serialize** your data for persistence (Project 04).
*   You'll parse data from a file, a skill you practiced with CSV and JSONL (Projects 03 & 04).

You are moving from writing data-processing pipelines to building interactive, stateful applications. The architectural patterns here—separating concerns with interfaces, a clear persistence layer—are the same ones used to build large-scale systems.

## How to Run and Test

This project is designed to be run from your terminal to see the CLI in action. You'll build the binary and then interact with it.

```bash
# Build the application binary
go build -o todo ./cmd/app

# Examples of how to use it:
# Add a new item
./todo -add "Write a great README"

# List all items
./todo -list

# Mark item with ID 1 as complete
./todo -toggle 1

# List only pending (incomplete) items
./todo -list -pending

# Run the provided tests
go test -v ./...
```