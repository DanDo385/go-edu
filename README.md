# Go Educational Repository

Welcome to the comprehensive Go exercises repository. This project is designed to take you from Go basics to advanced systems programming, including Ethereum internals.

## Repository Structure

Every project in this repository follows a consistent, production-grade structure:

*   **`cmd/app/`**: The entry point for the "real" application. Runs the logic as a complete program.
*   **`cmd/dev/`**: A development harness. Often simplifies inputs or setup for easier debugging and iteration.
*   **`internal/<pkg>/`**: The core logic library. This is where you work.
    *   `exercise.go`: The file you edit. Contains `TODO`s and stubbed functions.
    *   `exercise_test.go`: Tests to verify your implementation.
    *   `solution.reference.go`: The reference implementation (hidden by build tags).

## How to Work on Exercises

1.  **Choose a Project**: Navigate to `minis/` or `geth/` and pick a numbered directory.
2.  **Open in VS Code**: Open the repository root or the specific project folder.
3.  **Read the Exercise**: Open `internal/<pkg>/exercise.go`. Read the comments and `TODO`s.
4.  **Implement**: Write your code in `exercise.go`.
5.  **Test**: Run tests from the terminal:
    ```bash
    go test ./...
    ```
6.  **Debug**: Use the "Run and Debug" tab in VS Code. We provide pre-configured launch configurations:
    *   **Debug: cmd/dev**: Steps through your code with fixed inputs.
    *   **Test: Run All Tests**: debugs the test suite.

## Build Tags Explained

We use Go build tags to manage exercises and solutions:

*   **`!solution && !reference`**: This is the default tag for your environment. It compiles your `exercise.go` file.
*   **`reference`**: This tag is used for the reference solution. These files are excluded from your normal build so they don't conflict with your code. To run tests against the reference solution (to see expected behavior), use:
    ```bash
    go test -tags=reference ./...
    ```

## Project Index

### Minis (Go Foundations & Systems)
Explore the `minis/` directory for 50+ projects covering:
*   Language Basics (Strings, Maps, Slices)
*   Concurrency (Goroutines, Channels, Worker Pools)
*   Systems (HTTP, gRPC, Files)
*   Patterns (Middleware, Options, State Machines)
*   Cryptography (Hashing, Merkle Trees)

### Geth (Ethereum Internals)
Explore the `geth/` directory for Ethereum implementation exercises:
*   [01-stack](geth/01-stack) - RPC & Chain initialization
*   [02-rpc-basics](geth/02-rpc-basics) - Client interaction
*   ... and more.
