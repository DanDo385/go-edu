#!/bin/bash
# Script to reset exercise.go files to TODO list format
# Usage:
#   ./scripts/reset-exercises.sh               # Reset all exercises
#   ./scripts/reset-exercises.sh minis         # Reset only minis exercises
#   ./scripts/reset-exercises.sh geth          # Reset only geth exercises
#   ./scripts/reset-exercises.sh minis/01-*    # Reset specific exercise

set -e

CYAN='\033[0;36m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
NC='\033[0m' # No Color

# Function to create TODO version from solution.reference.go
create_todo_exercise() {
    local solution_file="$1"
    local exercise_file="$2"
    
    if [ ! -f "$solution_file" ]; then
        echo -e "${YELLOW}Warning: Solution file not found: $solution_file${NC}"
        return 1
    fi
    
    echo -e "${CYAN}Resetting: $exercise_file${NC}"
    
    # Use the Python script to create TODO version
    python3 "$(dirname "$0")/create_todo_exercise.py" "$solution_file" "$exercise_file"
    
    echo -e "${GREEN}✓ Reset complete${NC}"
}

# Function to process a directory
process_directory() {
    local dir="$1"
    local count=0
    
    echo -e "${CYAN}Finding exercises in $dir...${NC}"
    
    # Find all solution.reference.go files in the directory
    find "$dir" -name "solution.reference.go" -type f | while IFS= read -r solution_file; do
        exercise_file="${solution_file/solution.reference.go/exercise.go}"
        create_todo_exercise "$solution_file" "$exercise_file" && ((count++)) || true
    done
    
    # Count files for reporting
    local total=$(find "$dir" -name "solution.reference.go" -type f | wc -l)
    echo -e "${GREEN}✓ Reset $total exercise(s) in $dir${NC}"
}

# Main logic
if [ $# -eq 0 ]; then
    # Reset all exercises
    echo -e "${CYAN}Resetting all exercises...${NC}"
    process_directory "minis"
    process_directory "geth"
elif [ "$1" == "minis" ]; then
    echo -e "${CYAN}Resetting minis exercises...${NC}"
    process_directory "minis"
elif [ "$1" == "geth" ]; then
    echo -e "${CYAN}Resetting geth exercises...${NC}"
    process_directory "geth"
elif [ -d "$1" ]; then
    # Specific directory provided
    echo -e "${CYAN}Resetting exercises in $1...${NC}"
    process_directory "$1"
else
    echo -e "${YELLOW}Error: Invalid argument '$1'${NC}"
    echo "Usage:"
    echo "  $0               # Reset all exercises"
    echo "  $0 minis         # Reset only minis exercises"
    echo "  $0 geth          # Reset only geth exercises"
    echo "  $0 minis/01-*    # Reset specific exercise"
    exit 1
fi

echo -e "${GREEN}✓ All exercises reset successfully${NC}"
