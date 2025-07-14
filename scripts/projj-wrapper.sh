#!/bin/bash
# projj shell wrapper for bash/zsh
# This wrapper enables automatic directory changing after 'projj add'

projj() {
    local projj_binary="./projj-go"
    
    # Check if local projj-go binary exists
    if [[ ! -f "$projj_binary" ]]; then
        # Try to find projj-go binary in PATH
        if command -v "projj-go" &> /dev/null; then
            projj_binary="projj-go"
        # Try to find projj binary in PATH (but avoid recursion)
        elif command -v "projj" &> /dev/null && [[ "$(command -v projj)" != *"projj()"* ]]; then
            projj_binary="projj"
        else
            echo "Error: projj command not found in PATH" >&2
            return 1
        fi
    fi
    
    # Execute the original projj command and capture output
    local output
    output=$("$projj_binary" "$@" 2>&1)
    local exit_code=$?
    
    # Print the output
    echo "$output"
    
    # Check if this was an 'add' command and if change_directory is enabled
    if [[ $exit_code -eq 0 && "$1" == "add" ]]; then
        # Look for the special PROJJ_CHANGE_DIRECTORY line in output
        local change_dir
        change_dir=$(echo "$output" | grep "^PROJJ_CHANGE_DIRECTORY=" | cut -d'=' -f2-)
        
        if [[ -n "$change_dir" && -d "$change_dir" ]]; then
            echo "Changing directory to: $change_dir"
            cd "$change_dir" || {
                echo "Error: Failed to change directory to $change_dir" >&2
                return 1
            }
        fi
    fi
    
    return $exit_code
}

# Export the function so it's available in subshells
export -f projj