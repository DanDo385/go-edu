# Go Educational Projects

Welcome to the Go Educational Repository! This workspace is a collection of self-contained projects designed to teach Go fundamentals (`minis/`) and Ethereum development (`geth/`).

## 📂 Repository Structure

The repository is organized into two main tracks:

*   **`minis/`**: Go language fundamentals. These small projects focus on specific Go concepts like concurrency, interfaces, and standard library packages.
*   **`geth/`**: Ethereum development. These projects leverage `go-ethereum` (geth) to teach blockchain interaction, from basic connectivity to smart contracts and indexing.

### Project Layout

Each project follows a standard Go project layout:

```text
category/project-name/
├── cmd/
│   ├── app/
│   │   └── main.go          # CLI application entry point (accepts arguments)
│   └── dev/
│       └── main.go          # Debug harness (fixed inputs for debugging)
├── internal/
│   └── package_name/
│       ├── exercise.go      # TODO: Implement your solution here
│       ├── exercise_test.go # Tests to verify your implementation
│       ├── solution.reference.go # Reference solution (for learning/checking)
│       └── types.go         # (Optional) Type definitions
└── .vscode/                 # (Optional) Project-specific configurations
```

## 🚀 Getting Started

1.  **Initialize the workspace**:
    ```bash
    make setup
    ```
2.  **List available projects**:
    ```bash
    make list
    ```
3.  **Choose a project** (e.g., `minis/01-hello-strings`).

## 📝 How to Complete Exercises

1.  **Navigate** to the project directory (e.g., `cd minis/01-hello-strings`).
2.  **Open** the `internal/*/exercise.go` file.
3.  **Find the `TODO` comments**. These comments provide step-by-step instructions.
    *   *Note: Commentary in `exercise.go` is minimal. For detailed explanations, refer to `solution.reference.go` or the project's README.*
4.  **Implement** the code to satisfy the TODOs.
5.  **Run Tests** to verify your solution:
    ```bash
    make test P=minis/01-hello-strings
    # OR from within the directory:
    go test ./...
    ```

## 🛠️ Debugging & Running

We provide two entry points for every project:

### 1. The Debug Harness (`cmd/dev/main.go`)
*   **Purpose**: Learning and debugging with breakpoints.
*   **Behavior**: Uses hardcoded inputs so you don't need to pass CLI arguments.
*   **How to use**:
    *   Open `cmd/dev/main.go`.
    *   Set breakpoints in `exercise.go` or `main.go`.
    *   Press **F5** in VS Code (select "Debug: cmd/dev").
    *   Step through the code to understand the flow.

### 2. The Application (`cmd/app/main.go`)
*   **Purpose**: The "real" application.
*   **Behavior**: Accepts command-line arguments (e.g., RPC URLs, input strings).
*   **How to use**:
    ```bash
    go run ./cmd/app/main.go [arguments]
    ```
    *   Refer to the specific project's README for required arguments.

### VS Code Configuration
The `.vscode/launch.json` file in the root directory contains configurations to help you debug:
*   **Debug: cmd/dev**: Runs the debug harness for the project containing the active file.
*   **Debug: cmd/app**: Runs the CLI app (you may need to edit `args` in `launch.json`).
*   **Test: Current Package**: Debugs tests in the currently open package.

## 🧹 Make Commands

We provide a simplified `Makefile` to manage the repository.

### Resetting Exercises (`make todo`)
You can reset `exercise.go` files to their initial TODO state at any time.

*   **Reset Everything**:
    ```bash
    make todo all
    ```
*   **Reset Contextually**:
    *   From `geth/` directory: `make todo all` (Resets only `geth/` projects)
    *   From a specific project (e.g., `geth/01-stack/`): `make todo all` (Resets only that project)
*   **Reset Specific Path**:
    ```bash
    make todo minis/01-hello-strings
    ```

### Other Commands
*   `make setup`: Download dependencies and verify builds.
*   `make list`: Show all projects.
*   `make test`: Run all tests.
*   `make test P=<path>`: Run tests for a specific project.
*   `make run P=<path>`: Run the project's `cmd/app`.

## 📜 License
See [LICENSE](./LICENSE).
