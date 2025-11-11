#!/bin/bash

# Test Runner for Go Mini Projects
# This script helps you test your solutions for each mini project.

set -e

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Project list
PROJECTS=(
    "01-hello-strings"
    "02-arrays-maps-basics"
    "03-csv-stats"
    "04-jsonl-log-filter"
    "05-cli-todo-files"
    "06-worker-pool-wordcount"
    "07-generic-lru-cache"
    "08-http-client-retries"
    "09-http-server-graceful"
    "10-grpc-telemetry-service"
)

# Usage information
usage() {
    echo "Usage: $0 [project_number|all] [options]"
    echo ""
    echo "Examples:"
    echo "  $0 1           # Test project 01-hello-strings"
    echo "  $0 06          # Test project 06-worker-pool-wordcount"
    echo "  $0 all         # Test all projects"
    echo "  $0 1 -v        # Test with verbose output"
    echo "  $0 1 --solution # Test reference solution instead of your code"
    echo "  $0 all --bench # Run benchmarks (for projects that have them)"
    echo ""
    echo "Options:"
    echo "  -v, --verbose   Show verbose test output"
    echo "  --bench         Run benchmarks"
    echo "  --solution      Test the reference solution instead of exercise code"
    echo "  -h, --help      Show this help message"
    exit 1
}

# Test a single project
test_project() {
    local project_num=$1
    local verbose=$2
    local bench=$3
    local use_solution=$4

    # Pad project number with zero if needed
    if [[ $project_num =~ ^[0-9]$ ]]; then
        project_num="0$project_num"
    fi

    # Find the project directory
    local project_dir=""
    for proj in "${PROJECTS[@]}"; do
        if [[ $proj == $project_num-* ]]; then
            project_dir="minis/$proj/exercise"
            break
        fi
    done

    if [[ -z $project_dir ]]; then
        echo -e "${RED}✗ Project $project_num not found${NC}"
        return 1
    fi

    if [[ ! -d $project_dir ]]; then
        echo -e "${RED}✗ Directory $project_dir does not exist${NC}"
        return 1
    fi

    if [[ $use_solution == "true" ]]; then
        echo -e "${BLUE}Testing ${project_dir} ${YELLOW}(reference solution)${NC}"
    else
        echo -e "${BLUE}Testing ${project_dir}${NC}"
    fi

    cd "$project_dir"

    # Build test command
    local test_cmd="go test"
    if [[ $use_solution == "true" ]]; then
        test_cmd="$test_cmd -tags=solution"
    fi
    if [[ $verbose == "true" ]]; then
        test_cmd="$test_cmd -v"
    fi

    # Run tests
    if eval "$test_cmd"; then
        echo -e "${GREEN}✓ Tests passed${NC}\n"
    else
        echo -e "${RED}✗ Tests failed${NC}\n"
        cd - > /dev/null
        return 1
    fi

    # Run benchmarks if requested
    if [[ $bench == "true" ]] && ls *bench*.go &>/dev/null; then
        echo -e "${YELLOW}Running benchmarks...${NC}"
        local bench_cmd="go test -bench=. -benchmem"
        if [[ $use_solution == "true" ]]; then
            bench_cmd="go test -tags=solution -bench=. -benchmem"
        fi
        eval "$bench_cmd"
        echo ""
    fi

    cd - > /dev/null
    return 0
}

# Test all projects
test_all() {
    local verbose=$1
    local bench=$2
    local use_solution=$3
    local failed=0
    local passed=0

    echo -e "${BLUE}========================================${NC}"
    echo -e "${BLUE}Testing All Projects${NC}"
    echo -e "${BLUE}========================================${NC}\n"

    for i in "${!PROJECTS[@]}"; do
        project_num=$(printf "%02d" $((i + 1)))
        if test_project "$project_num" "$verbose" "$bench" "$use_solution"; then
            ((passed++))
        else
            ((failed++))
        fi
    done

    echo -e "${BLUE}========================================${NC}"
    echo -e "${BLUE}Summary${NC}"
    echo -e "${BLUE}========================================${NC}"
    echo -e "${GREEN}Passed: $passed${NC}"
    if [[ $failed -gt 0 ]]; then
        echo -e "${RED}Failed: $failed${NC}"
        return 1
    else
        echo -e "${GREEN}All tests passed!${NC}"
        return 0
    fi
}

# Main script
main() {
    # Parse arguments
    local project=""
    local verbose="false"
    local bench="false"
    local use_solution="false"

    if [[ $# -eq 0 ]]; then
        usage
    fi

    for arg in "$@"; do
        case $arg in
            -h|--help)
                usage
                ;;
            -v|--verbose)
                verbose="true"
                ;;
            --bench)
                bench="true"
                ;;
            --solution)
                use_solution="true"
                ;;
            *)
                if [[ -z $project ]]; then
                    project="$arg"
                fi
                ;;
        esac
    done

    # Change to script directory
    cd "$(dirname "$0")"

    if [[ $project == "all" ]]; then
        test_all "$verbose" "$bench" "$use_solution"
    else
        test_project "$project" "$verbose" "$bench" "$use_solution"
    fi
}

main "$@"
