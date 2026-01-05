# 05: Command-Line Todo App with File Persistence

This project marks a major milestone: you will build your first **complete, interactive application**. You will assemble all the skills you've learned—strings, slices, structs, and JSON processing—to create a command-line (CLI) todo manager. The central challenge is **persistence**: making your application remember data between runs by saving it to a file.

## Table of Contents

- [Learning Objectives](#learning-objectives)
- [The Big Picture: Building a Complete Application](#the-big-picture-building-a-complete-application)
- [Application Architecture: Separation of Concerns](#application-architecture-separation-of-concerns)
- [First Principles: CLIs, Persistence, and CRUD](#first-principles-clis-persistence-and-crud)
- [Project Structure](#project-structure)
- [Key Go Concepts in This Project](#key-go-concepts-in-this-project)
  - [Parsing User Input: The `flag` Package](#parsing-user-input-the-flag-package)
  - [Interfaces for Abstraction: The `Store`](#interfaces-for-abstraction-the-store)
  - [JSON for Data Serialization](#json-for-data-serialization)
- [Progression: Tying It All Together](#progression-tying-it-all-together)
- [How to Run and Test](#how-to-run-and-test)
- [Key Takeaways](#key-takeaways)
- [Further Reading](#further-reading)

## Learning Objectives

By the end of this project, you will be able to:

-   **Build a complete, interactive CLI application** from scratch.
-   **Parse command-line arguments and flags** using Go's standard `flag` package.
-   **Implement persistence** by saving and loading application state to a JSON file.
-   **Design a clean application architecture** using interfaces to separate concerns (UI vs. data).
-   **Implement the Dependency Inversion Principle** to create testable and flexible code.
-   **Apply CRUD (Create, Read, Update, Delete) principles** in a real application.

## The Big Picture: Building a Complete Application

So far, you've been building libraries—small, reusable pieces of code that process data. Now, you'll assemble those skills to build your first **complete, interactive application**: a command-line (CLI) todo manager.

This project is a major step. It introduces the fundamental challenge of **persistence**: how does your application remember data between runs? You'll solve this by creating a data "store" that saves your todo items to a JSON file, introducing you to application architecture, interfaces, and the full lifecycle of data management.

## Application Architecture: Separation of Concerns

This project introduces a clean, layered architecture. Understanding this separation of concerns is key to building maintainable software.

```
+--------------------------------+
|           CLI Layer            |  <-- Knows what the user wants.
|           (main.go)            |      (e.g., "add an item", "list items")
|                                |      Knows NOTHING about files or JSON.
+--------------------------------+
               |
               | Calls methods on the Store interface.
               | Example: store.Add("Buy milk")
               v
+--------------------------------+
|      Abstraction Layer         |  <-- The "Contract".
|      (Store interface)         |      Defines WHAT can be done, not HOW.
|                                |      (e.g., "a Store must have an Add method")
+--------------------------------+
               ^
               | The fileStore struct IMPLEMENTS the Store interface.
               | It promises to fulfill the contract.
               |
+--------------------------------+
|     Data Persistence Layer     |  <-- Knows how to handle data.
|       (fileStore struct)       |      (e.g., read/write a JSON file)
|                                |      Knows NOTHING about CLI flags.
+--------------------------------+
```

-   The **CLI Layer** is the user-facing part.
-   The **Data Persistence Layer** is the data-handling part.
-   The **`Store` interface** is the glue that allows them to be developed and tested independently.

## First Principles: CLIs, Persistence, and CRUD

1.  **Command-Line Interface (CLI)**: A way for users to interact with a program by typing text commands, arguments, and flags (e.g., `-add "My Task"`).
2.  **Persistence**: The ability of an application to save its **state** (data) so that it survives when the application closes. Our choice, a simple JSON file, is a very common and effective persistence strategy.
3.  **CRUD Operations**: A foundational concept in data management, representing the four basic functions you can perform on data: **C**reate, **R**ead, **U**pdate, **D**elete.

## Project Structure

```
.
├── cmd/
│   └── app/
│       └── main.go       # Entry point. Parses flags and drives the app.
├── internal/
│   ├── store.go      # Defines the Store interface and fileStore implementation.
│   └── item.go       # Defines the Item struct.
└── testdata/
    └── todos.json    # The file where your todo items will be saved.
```

-   **`cmd/app`**: The main application package. Its only job is to interpret user commands and call the appropriate `Store` methods.
-   **`internal/`**: The core business logic and data persistence implementation.

## Key Go Concepts in This Project

### Parsing User Input: The `flag` Package

The `flag` package is the idiomatic Go way to define and parse command-line flags. It's more robust than parsing `os.Args` manually, handles type conversions, and automatically generates help text.

```go
add := flag.Bool("add", false, "Add a new todo item")
list := flag.Bool("list", false, "List all todo items")
flag.Parse() // Parses the user's input and populates the variables.

if *add {
    // ... logic to add an item
}
```

### Interfaces for Abstraction: The `Store`

This is the most important architectural concept in this project. We define a `Store` **interface** that describes *what* our data store must be able to do.

```go
type Store interface {
    Load() error
    Save() error
    Add(text string)
    Toggle(id int) error
    List(pendingOnly bool) []Item
}
```
We then create a `fileStore` **struct** that *implements* this interface.

**Why is this so powerful? Decoupling and Dependency Injection.**

The `main` function of our application depends only on the `Store` interface, not the concrete `fileStore` implementation. This means we can "inject" any type that satisfies the interface.

```go
// In production, we inject the real file-based store:
var s Store = NewFileStore("path/to/todos.json")
runApp(s)

// In our tests, we can inject a completely different, in-memory mock store:
var mock Store = NewMockStore()
// Add some test data
mock.Add("A test item")
// Now we can test the `runApp` logic without touching the file system!
runApp(mock)
```
This pattern of "programming to an interface" is one of the most important principles of modern software design. It makes your code flexible, modular, and easy to test.

### JSON for Data Serialization

This project puts your JSON skills from Project 04 to practical use.
*   `json.MarshalIndent`: To save the `[]Item` slice to the file, you'll "marshal" it into a nicely formatted JSON byte slice.
*   `json.Unmarshal`: To load the app, you'll read the file and "unmarshal" the JSON data back into the `[]Item` slice.

This process of converting in-memory objects to a storable format is called **serialization**.

## Progression: Tying It All Together

This project is a synthesis of everything you've learned so far. You are moving from writing data-processing pipelines to building interactive, stateful applications. The architectural patterns here—separating concerns with interfaces, a clear persistence layer—are the same ones used to build large-scale systems.

## How to Run and Test

This project is designed to be run from your terminal.

```bash
# Build the application binary
go build -o todo ./cmd/app

# --- Examples of how to use it ---

# Add a new item (the text is the first argument after flags)
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

## Key Takeaways

-   **Separate your application's concerns** into layers (e.g., UI, business logic, data storage).
-   **Use interfaces as contracts** to decouple these layers.
-   This "dependency injection" makes your code **testable and flexible**.
-   The `flag` package is the standard way to parse CLI flags.
-   **Persistence** is the act of saving application state; `json.Marshal` and `os.WriteFile` is a simple and effective way to do it.

## Further Reading

-   [**Go by Example: Command-Line Flags**](https://gobyexample.com/command-line-flags): A guide to using the `flag` package.
-   [**Go by Example: Interfaces**](https://gobyexample.com/interfaces): A concise introduction to interfaces in Go.
-   [**SOLID Design Principles**](https://en.wikipedia.org/wiki/SOLID): The "D" in SOLID stands for Dependency Inversion, which is exactly the principle our `Store` interface demonstrates.
-   [**Twelve-Factor App**](https://12factor.net/): A methodology for building robust, scalable applications. This project touches on several factors, including treating backing services (our JSON file) as attached resources.