# Exercise File Management

This document explains how exercise files are structured and how to reset them to TODO format.

## File Structure

Each exercise in both `minis/` and `geth/` directories follows this structure:

### exercise.go
- **Purpose**: Student-facing file with TODO lists
- **Content**: 
  - Concise problem statement
  - Key constraints and requirements
  - Function signatures with return types
  - TODO comments indicating what needs to be implemented
  - References to `solution.reference.go` for details
- **Build Tag**: `//go:build !solution && !reference`
- **Size**: ~20-50 lines (vs 150-700 lines in solution files)

### solution.reference.go
- **Purpose**: Complete implementation with extensive educational comments
- **Content**:
  - Full working implementation
  - Detailed step-by-step explanations
  - Computer science principles
  - Debugging tips and breakpoint suggestions
  - Memory layout explanations
  - Algorithm descriptions
- **Build Tag**: `//go:build reference`
- **Size**: Full implementation with comprehensive commentary

## Resetting Exercises

You can reset exercise files to TODO format using the provided Makefile commands:

### Reset All Exercises
```bash
make reset-exercises
```
Resets all exercises in both `minis/` and `geth/` directories.

### Reset Minis Only
```bash
make reset-exercises-minis
```
Resets only exercises in the `minis/` directory.

### Reset Geth Only
```bash
make reset-exercises-geth
```
Resets only exercises in the `geth/` directory.

### Reset Specific Exercise
```bash
./scripts/reset-exercises.sh minis/01-hello-strings
```
Resets a specific exercise directory.

## How It Works

The reset process:

1. **Finds Solution Files**: Locates all `solution.reference.go` files
2. **Extracts Structure**: Parses package, imports, types, and function signatures
3. **Creates TODOs**: Generates concise exercise.go with:
   - Minimal problem statement (no verbose explanations)
   - Complete function signatures
   - TODO comments for implementation
   - Appropriate zero-value returns
4. **Preserves Solution**: Never modifies `solution.reference.go` files

## Script Details

### Main Script
`scripts/reset-exercises.sh` - Bash script that orchestrates the reset process

### Python Generator
`scripts/create_todo_exercise.py` - Python script that:
- Parses Go source files
- Extracts function signatures correctly (handles `interface{}`, variadic parameters, etc.)
- Generates appropriate return statements based on return types
- Preserves essential type definitions and constants
- Creates concise problem statements

## Example Transformation

### Before (exercise.go with full implementation)
```go
func Run(ctx context.Context, client RPCClient, cfg Config) (*Result, error) {
    // ============================================================================
    // STEP 1: Input Validation - Defensive Programming Pattern
    // ============================================================================
    // Why validate inputs? This function is a library API that will be called by
    // other code. We can't trust callers to always pass valid inputs...
    // [150+ lines of implementation and comments]
}
```

### After (exercise.go with TODO)
```go
// Run - TODO: implement this function
func Run(ctx context.Context, client RPCClient, cfg Config) (*Result, error) {
    // TODO: Implement this function
    // Refer to solution.reference.go for the complete implementation with detailed explanations
    return nil, nil
}
```

## Benefits

1. **Clear Learning Path**: Students start with clean TODO lists
2. **Complete Reference**: All implementation details preserved in solution files
3. **Easy Reset**: Teachers/students can reset to fresh state anytime
4. **Consistent Structure**: All 75 exercises follow the same pattern
5. **Build System Integration**: Works seamlessly with Go build tags

## Statistics

- **Total Exercises**: 75 (49 minis + 26 geth)
- **Average Reduction**: ~95% fewer lines in exercise files
- **Solution Files**: Unchanged, maintain full educational content
- **Time to Reset**: ~30 seconds for all exercises

## Maintenance

To modify the reset behavior:

1. Edit `scripts/create_todo_exercise.py` for generation logic
2. Edit `scripts/reset-exercises.sh` for orchestration
3. Test on a single exercise before running on all files

## Build Tags

The project uses Go build tags to manage which files are compiled:

- **Default build** (`go build` or `go test`): Uses `exercise.go` files
- **Solution build** (`go build -tags solution`): Uses `solution.go` files
- **Reference build** (`go build -tags reference`): Uses `solution.reference.go` files

This ensures students work with exercise files by default, while reference implementations remain available for review.
