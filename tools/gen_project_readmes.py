#!/usr/bin/env python3
"""
Generate per-project README.md files for minis/ and geth/.

Existing READMEs are left untouched (not overwritten).
"""

from __future__ import annotations

from pathlib import Path


ROOT = Path(__file__).resolve().parents[1]
MINIS = ROOT / "minis"
GETH = ROOT / "geth"


def title_from_dir(dirname: str) -> tuple[str, str]:
    """
    Given e.g. '07-eth-call' -> ('07', 'Eth Call')
    """
    if "-" not in dirname:
        return ("", dirname)
    num, slug = dirname.split("-", 1)
    words = [w for w in slug.replace("_", "-").split("-") if w]
    name = " ".join(w.upper() if w in {"rpc", "eip1559", "grpc", "http", "jwt", "lru", "pprof", "tcp", "udp"} else w.capitalize() for w in words)
    # Fix common Ethereum acronyms
    name = name.replace("Eip1559", "EIP-1559").replace("Geth", "Geth").replace("Abi", "ABI").replace("Erc20", "ERC20")
    return (num, name)


def readme_template(track: str, dirname: str) -> str:
    num, name = title_from_dir(dirname)

    header = f"# {num}: {name}" if num else f"# {name}"

    if track == "minis":
        about = (
            "This project is a focused, self-contained mini that teaches a single Go concept by making you implement and test it in isolation.\n\n"
            "You’ll work in the `internal/` exercise package (where the tests live), then use `cmd/dev` to step through deterministic examples with a debugger, and `cmd/app` to run the same idea with real CLI arguments."
        )
        why = (
            "Go rewards building strong intuition around its core primitives (types, interfaces, concurrency, I/O, and the standard library). "
            "These minis are designed to build that intuition quickly by combining tight exercises with fast feedback from tests and a debug harness."
        )
        problems = [
            "Turning a language feature into reliable, testable code",
            "Designing small, composable APIs and data structures",
            "Debugging correctness and performance issues with repeatable inputs",
        ]
        prereq = ["Basic Go syntax (functions, structs, slices/maps)", "Comfort running `go test` and `go run`"]
    else:
        about = (
            "This project is part of the `geth/` track and teaches an Ethereum concept using Go + go-ethereum.\n\n"
            "You’ll connect to an RPC node, query chain data, and learn the underlying primitives (blocks, transactions, calls, logs) that most higher-level tooling is built on."
        )
        why = (
            "Ethereum development is ultimately about understanding what the node exposes over JSON-RPC and how those APIs map to on-chain state. "
            "Learning the concepts in small, targeted projects makes later contract interaction and infra work dramatically easier."
        )
        problems = [
            "Building reliable RPC-based tooling (indexers, monitors, analyzers)",
            "Debugging on-chain behavior by inspecting canonical node data structures",
            "Bridging console/RPC intuition into production Go services",
        ]

        prereq = ["Completion of earlier geth modules in sequence (recommended)"]
        if dirname.startswith("07-"):
            prereq.insert(0, "Completion of `geth/06-smart-contracts` (console intuition for call vs transaction)")
        if dirname.startswith("08-"):
            prereq = [
                "Completion of `geth/06-smart-contracts`",
                "Completion of `geth/07-eth-call`",
            ] + prereq
        if dirname.startswith("09-"):
            prereq = [
                "Completion of `geth/06-smart-contracts`",
                "Completion of `geth/07-eth-call`",
                "Completion of `geth/08-abigen`",
            ] + prereq

    concepts = [
        f"{name}: the core idea behind this module",
        "Debugging with deterministic inputs (`cmd/dev`) and tests (`go test`)",
        "Reading and reasoning about results + common failure modes",
    ]

    structure = f"""```text
{dirname}/
  cmd/
    app/  # Application entry point (CLI arguments)
    dev/  # Debug harness (fixed inputs)
  internal/
    <package>/  # Exercise implementation
      exercise.go
      exercise_test.go
      solution.reference.go
      solution_no_err.reference.go
  .vscode/
    launch.json  # Debug configurations
```"""

    return f"""{header}

## What Is This Project About?

{about}

## Why Is This Important?

{why}

## Real-World Problems This Solves

- {problems[0]}
- {problems[1]}
- {problems[2]}

## Key Concepts You’ll Learn

- {concepts[0]}
- {concepts[1]}
- {concepts[2]}

## Prerequisites

{chr(10).join(f"- {p}" for p in prereq)}

## Project Structure

{structure}

## How to Run

### Using `cmd/app/main.go` (CLI arguments)

```bash
# from this project directory
go run ./cmd/app
```

### Using `cmd/dev/main.go` (debug harness)

```bash
# from this project directory
go run ./cmd/dev
```

### How to Debug

- Set breakpoints at `// BREAKPOINT:` comments
- Press **F5** and select:
  - **Debug: cmd/app (with RPC_URL argument support)** (geth) / the closest matching config (minis)
  - **Debug: cmd/dev (Debug Harness)**
  - **Test: Run All Tests** / **Test: Current Test Function**

## Testing

```bash
go test ./...
go test -v ./...
go test -v -run TestName ./...
go test -tags=reference -v ./...
```

## Exercises

- Implement the functions in `internal/<package>/exercise.go` until tests pass.

## Additional Resources

- Go testing: `go help test`
- go-ethereum docs (geth track): `github.com/ethereum/go-ethereum`
"""


def main() -> None:
    created = 0

    # minis/*
    for p in sorted(MINIS.iterdir()):
        if not p.is_dir() or p.name.startswith("."):
            continue
        readme = p / "README.md"
        if readme.exists():
            continue
        readme.write_text(readme_template("minis", p.name), encoding="utf-8")
        created += 1

    # geth/*
    for p in sorted(GETH.iterdir()):
        if not p.is_dir() or p.name.startswith("."):
            continue
        if p.name == "06-smart-contracts":
            continue  # custom tutorial README already exists
        readme = p / "README.md"
        if readme.exists():
            continue
        readme.write_text(readme_template("geth", p.name), encoding="utf-8")
        created += 1

    print(f"created {created} README.md files")


if __name__ == "__main__":
    main()
