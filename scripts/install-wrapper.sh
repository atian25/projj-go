#!/bin/bash
# Install script for projj shell wrapper
# This script helps users set up the shell wrapper for automatic directory changing

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Get the directory where this script is located
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

echo -e "${BLUE}Projj Shell Wrapper Installer${NC}"
echo "=============================="
echo

# Detect current shell
CURRENT_SHELL=$(basename "$SHELL")
echo -e "Detected shell: ${GREEN}$CURRENT_SHELL${NC}"
echo

# Function to add wrapper to shell config
add_to_shell_config() {
    local config_file="$1"
    local wrapper_file="$2"
    local shell_name="$3"
    
    if [[ -f "$config_file" ]]; then
        # Check if wrapper is already sourced
        if grep -q "projj-wrapper" "$config_file"; then
            echo -e "${YELLOW}Projj wrapper already configured in $config_file${NC}"
            return 0
        fi
        
        echo -e "${GREEN}Adding projj wrapper to $config_file${NC}"
        echo "" >> "$config_file"
        echo "# Projj shell wrapper for automatic directory changing" >> "$config_file"
        echo "source \"$wrapper_file\"" >> "$config_file"
        echo -e "${GREEN}✓ Added to $config_file${NC}"
        return 0
    else
        echo -e "${YELLOW}$config_file not found, creating it...${NC}"
        mkdir -p "$(dirname "$config_file")"
        echo "# Projj shell wrapper for automatic directory changing" > "$config_file"
        echo "source \"$wrapper_file\"" >> "$config_file"
        echo -e "${GREEN}✓ Created and configured $config_file${NC}"
        return 0
    fi
}

# Function to add fish wrapper
add_fish_wrapper() {
    local fish_config_dir="$HOME/.config/fish"
    local fish_functions_dir="$fish_config_dir/functions"
    local wrapper_file="$SCRIPT_DIR/projj-wrapper.fish"
    
    if [[ ! -d "$fish_functions_dir" ]]; then
        echo -e "${YELLOW}Creating fish functions directory...${NC}"
        mkdir -p "$fish_functions_dir"
    fi
    
    local projj_function_file="$fish_functions_dir/projj.fish"
    
    if [[ -f "$projj_function_file" ]]; then
        echo -e "${YELLOW}Projj function already exists in $projj_function_file${NC}"
        return 0
    fi
    
    echo -e "${GREEN}Installing projj function for fish...${NC}"
    cp "$wrapper_file" "$projj_function_file"
    echo -e "${GREEN}✓ Installed fish function${NC}"
}

# Install based on detected shell
case "$CURRENT_SHELL" in
    "bash")
        add_to_shell_config "$HOME/.bashrc" "$SCRIPT_DIR/projj-wrapper.sh" "bash"
        if [[ "$OSTYPE" == "darwin"* ]]; then
            # On macOS, also add to .bash_profile
            add_to_shell_config "$HOME/.bash_profile" "$SCRIPT_DIR/projj-wrapper.sh" "bash"
        fi
        ;;
    "zsh")
        add_to_shell_config "$HOME/.zshrc" "$SCRIPT_DIR/projj-wrapper.sh" "zsh"
        ;;
    "fish")
        add_fish_wrapper
        ;;
    *)
        echo -e "${YELLOW}Unsupported shell: $CURRENT_SHELL${NC}"
        echo "Please manually source the appropriate wrapper script:"
        echo "  - For bash/zsh: source $SCRIPT_DIR/projj-wrapper.sh"
        echo "  - For fish: copy $SCRIPT_DIR/projj-wrapper.fish to ~/.config/fish/functions/projj.fish"
        echo "  - For PowerShell: source $SCRIPT_DIR/projj-wrapper.ps1"
        exit 1
        ;;
esac

echo
echo -e "${GREEN}Installation completed!${NC}"
echo
echo -e "${BLUE}Next steps:${NC}"
echo "1. Restart your terminal or run: source ~/.${CURRENT_SHELL}rc"
echo "2. Ensure 'change_directory' is set to true in your projj config:"
echo "   projj config set -k change_directory -v true"
echo "3. Test with: projj add <repository-url>"
echo
echo -e "${YELLOW}Note: The wrapper will automatically change to the repository directory after 'projj add'${NC}"