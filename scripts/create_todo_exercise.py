#!/usr/bin/env python3
"""
Convert solution.reference.go files to TODO-based exercise.go files.
Preserves structure while removing implementation details.
"""

import sys
import re
from typing import List, Tuple

def extract_function_signature(lines: List[str], start_idx: int) -> Tuple[List[str], int]:
    """Extract complete function signature."""
    sig_lines = []
    i = start_idx
    
    # Skip initial comment lines before function
    while i < len(lines) and lines[i].strip().startswith('//'):
        i += 1
    
    # Get the function signature (might span multiple lines)
    # Keep reading until we find a line that ends with ' {' or starts with '{'
    while i < len(lines):
        line = lines[i].rstrip()
        sig_lines.append(line)
        
        # Check if this line has the function body opening brace
        # It's either at the end of the signature line or on the next line
        stripped = line.strip()
        if stripped.endswith('{') or stripped == '{':
            break
        i += 1
    
    return sig_lines, i

def parse_return_type(sig_lines: List[str]) -> str:
    """Parse the return type from function signature."""
    full_sig = ' '.join(sig_lines)
    
    # Match function signature pattern
    match = re.search(r'\)\s*([^{]+)\s*\{', full_sig)
    if match:
        return_part = match.group(1).strip()
        return return_part if return_part else ""
    return ""

def generate_return_statement(return_type: str) -> str:
    """Generate appropriate return statement for the return type."""
    if not return_type:
        return ""
    
    # Handle multiple return values
    if ',' in return_type:
        parts = [p.strip() for p in return_type.split(',')]
        returns = []
        for part in parts:
            if part == 'error':
                returns.append('nil')
            elif part in ['string', 'int', 'uint8', 'uint16', 'uint32', 'uint64', 'int8', 'int16', 'int32', 'int64']:
                zero_vals = {'string': '""', 'int': '0', 'uint8': '0', 'uint16': '0', 'uint32': '0', 'uint64': '0',
                            'int8': '0', 'int16': '0', 'int32': '0', 'int64': '0'}
                returns.append(zero_vals.get(part, '0'))
            elif part == 'bool':
                returns.append('false')
            elif part in ['float32', 'float64']:
                returns.append('0.0')
            else:
                returns.append('nil')
        return f"\treturn {', '.join(returns)}\n"
    
    # Single return value
    if return_type == 'error':
        return '\treturn nil\n'
    elif return_type in ['string']:
        return '\treturn ""\n'
    elif return_type in ['int', 'uint8', 'uint16', 'uint32', 'uint64', 'int8', 'int16', 'int32', 'int64']:
        return '\treturn 0\n'
    elif return_type == 'bool':
        return '\treturn false\n'
    elif return_type in ['float32', 'float64']:
        return '\treturn 0.0\n'
    else:
        return '\treturn nil\n'

def create_todo_version(solution_file: str, exercise_file: str):
    """Create TODO-based exercise file from solution file."""
    with open(solution_file, 'r') as f:
        lines = f.readlines()
    
    output = []
    i = 0
    problem_statement_done = False
    
    while i < len(lines):
        line = lines[i].rstrip()
        
        # Handle build tags
        if line.startswith('//go:build'):
            output.append('//go:build !solution && !reference\n\n')
            i += 1
            continue
        
        # Handle package declaration
        if line.startswith('package '):
            output.append(line + '\n\n')
            i += 1
            continue
        
        # Handle imports
        if line.startswith('import'):
            output.append(line + '\n')
            i += 1
            if '(' in line:
                while i < len(lines) and ')' not in lines[i-1]:
                    output.append(lines[i])
                    i += 1
            output.append('\n')
            continue
        
        # Handle const declarations
        if line.startswith('const '):
            output.append(line + '\n')
            i += 1
            continue
        
        # Handle top-level comment blocks (problem statements)
        if line == '/*' and not problem_statement_done:
            comment_block = []
            start_i = i
            while i < len(lines) and '*/' not in lines[i]:
                comment_block.append(lines[i])
                i += 1
            if i < len(lines):
                comment_block.append(lines[i])
                i += 1
            
            # Extract concise problem statement
            output.append('/*\n')
            keep_next = False
            for cline in comment_block[1:-1]:
                stripped = cline.strip()
                
                # Keep key sections
                if any(stripped.startswith(x) for x in ['Problem:', 'Requirements:', 'Algorithm:', 'Constraints:', 'Time/Space']):
                    output.append(cline)
                    keep_next = True
                elif keep_next and stripped and not any(x in stripped for x in ['Why', 'How it works', 'DEBUGGING', 'Computer science', 'Go Concepts', 'BREAKPOINT', 'DEBUG:', 'Building on']):
                    if stripped.startswith('-') or stripped.startswith('1.') or stripped.startswith('2.') or stripped.startswith('3.'):
                        output.append(cline)
                    else:
                        keep_next = False
                else:
                    keep_next = False
            
            output.append('*/\n\n')
            problem_statement_done = True
            continue
        
        # Handle type definitions
        if re.match(r'^type \w+', line):
            # Keep type definitions
            output.append(line + '\n')
            i += 1
            # If it's a struct, keep the struct definition
            if '{' in line:
                brace_count = line.count('{') - line.count('}')
                while i < len(lines) and brace_count > 0:
                    output.append(lines[i])
                    brace_count += lines[i].count('{') - lines[i].count('}')
                    i += 1
            output.append('\n')
            continue
        
        # Handle function definitions
        if re.match(r'^func ', line) or re.match(r'^func \(', line):
            # Extract function signature
            sig_lines, end_idx = extract_function_signature(lines, i)
            
            # Add concise comment
            func_name_match = re.search(r'func (?:\([^)]+\) )?(\w+)', line)
            if func_name_match:
                func_name = func_name_match.group(1)
                output.append(f'// {func_name} - TODO: implement this function\n')
            
            # Add signature - reconstruct it properly
            # Need to handle the opening brace correctly
            for idx, sig_line in enumerate(sig_lines):
                if '{' in sig_line:
                    # This is the line with the opening brace
                    # Find the position of the last { which is the function body opening
                    # But we need to avoid splitting interface{} 
                    # Strategy: find the rightmost { that's not part of interface{}
                    
                    # Simple heuristic: if line ends with {, split from the end
                    if sig_line.rstrip().endswith('{'):
                        clean_line = sig_line.rstrip()[:-1].rstrip()
                        if clean_line:
                            output.append(clean_line + ' {\n')
                        else:
                            output.append('{\n')
                    else:
                        # { is in the middle, likely interface{} - keep the whole line
                        output.append(sig_line + '\n')
                else:
                    output.append(sig_line + '\n')
            
            # Add TODO comments
            output.append('\t// TODO: Implement this function\n')
            output.append('\t// Refer to solution.reference.go for the complete implementation with detailed explanations\n')
            
            # Generate appropriate return statement
            return_type = parse_return_type(sig_lines)
            output.append(generate_return_statement(return_type))
            
            output.append('}\n\n')
            
            # Skip to the end of the function in the original file
            i = end_idx + 1
            brace_count = 1
            while i < len(lines) and brace_count > 0:
                brace_count += lines[i].count('{') - lines[i].count('}')
                i += 1
            continue
        
        # Skip other comment blocks
        if line.startswith('//'):
            i += 1
            continue
        
        i += 1
    
    # Write output
    with open(exercise_file, 'w') as f:
        f.write(''.join(output))

if __name__ == '__main__':
    if len(sys.argv) != 3:
        print("Usage: create_todo_exercise.py <solution_file> <exercise_file>")
        sys.exit(1)
    
    create_todo_version(sys.argv[1], sys.argv[2])
