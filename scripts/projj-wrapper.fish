# projj shell wrapper for fish
# This wrapper enables automatic directory changing after 'projj add'

function projj
    set projj_cmd "projj"
    
    # Check if projj binary exists in PATH
    if not command -v $projj_cmd > /dev/null 2>&1
        # Try to find projj-go binary
        if command -v "projj-go" > /dev/null 2>&1
            set projj_cmd "projj-go"
        else
            echo "Error: projj command not found in PATH" >&2
            return 1
        end
    end
    
    # Execute the original projj command and capture output
    set output (eval $projj_cmd $argv 2>&1)
    set exit_code $status
    
    # Print the output
    echo $output
    
    # Check if this was an 'add' or 'find' command and if change_directory is enabled
    if test $exit_code -eq 0 -a \( "$argv[1]" = "add" -o "$argv[1]" = "find" \)
        # Look for the special PROJJ_CHANGE_DIRECTORY line in output
        set change_dir (echo $output | grep "^PROJJ_CHANGE_DIRECTORY=" | cut -d'=' -f2-)
        
        if test -n "$change_dir" -a -d "$change_dir"
            echo "Changing directory to: $change_dir"
            cd "$change_dir"
            if test $status -ne 0
                echo "Error: Failed to change directory to $change_dir" >&2
                return 1
            end
        end
    end
    
    return $exit_code
end