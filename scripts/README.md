# Projj Shell Wrapper Scripts

These scripts enable automatic directory changing after `projj add` commands, similar to the Node.js version of projj.

## Overview

When `change_directory` is enabled in your projj configuration, these wrapper scripts will automatically change your current working directory to the newly added repository after a successful `projj add` command.

## Supported Shells

- **Bash** (`projj-wrapper.sh`)
- **Zsh** (`projj-wrapper.sh`)
- **Fish** (`projj-wrapper.fish`)
- **PowerShell** (`projj-wrapper.ps1`)

## Quick Installation

Run the installation script to automatically set up the wrapper for your current shell:

```bash
./scripts/install-wrapper.sh
```

This script will:
1. Detect your current shell
2. Add the appropriate wrapper to your shell configuration
3. Provide instructions for activation

## Manual Installation

### Bash/Zsh

Add this line to your `~/.bashrc` or `~/.zshrc`:

```bash
source "/path/to/projj-go/scripts/projj-wrapper.sh"
```

### Fish

Copy the fish wrapper to your functions directory:

```bash
cp scripts/projj-wrapper.fish ~/.config/fish/functions/projj.fish
```

### PowerShell

Add this line to your PowerShell profile (`$PROFILE`):

```powershell
. "/path/to/projj-go/scripts/projj-wrapper.ps1"
```

## Configuration

Ensure that `change_directory` is enabled in your projj configuration:

```bash
projj config set -k change_directory -v true
```

You can verify the setting with:

```bash
projj config get -k change_directory
```

## How It Works

1. The wrapper intercepts `projj` commands
2. When you run `projj add <repository-url>` or `projj find <query>`, it executes the original command
3. If the command succeeds and `change_directory` is enabled, projj outputs a special line: `PROJJ_CHANGE_DIRECTORY=/path/to/repo`
4. The wrapper detects this line and automatically changes to that directory

### Supported Commands

- **`projj add`**: Always changes directory after successfully adding a repository
- **`projj find`**: Changes directory when exactly one repository matches the query

## Example Usage

```bash
# Before: you're in ~/Documents
$ pwd
/Users/username/Documents

# Add a repository
$ projj add golang/go
正在克隆 https://github.com/golang/go.git 到 /Users/username/Workspaces/coding/github.com/golang/go...
仓库添加成功: /Users/username/Workspaces/coding/github.com/golang/go
Changing directory to: /Users/username/Workspaces/coding/github.com/golang/go

# After: you're automatically in the new repository
$ pwd
/Users/username/Workspaces/coding/github.com/golang/go
```

## Troubleshooting

### Wrapper Not Working

1. Ensure you've restarted your terminal or sourced your shell configuration
2. Check that `change_directory` is set to `true` in your projj config
3. Verify the wrapper is properly loaded by running `type projj`

### Binary Not Found

The wrapper will automatically try to find:
1. `projj` command in PATH
2. `projj-go` command in PATH

If neither is found, ensure your projj binary is in your PATH or create a symlink:

```bash
# Create a symlink (adjust paths as needed)
ln -s /path/to/projj-go /usr/local/bin/projj
```

### Fish Shell Issues

If the fish function doesn't work, try:

```bash
# Reload fish functions
fish -c "source ~/.config/fish/functions/projj.fish"
```

## Cross-Platform Compatibility

These scripts are designed to work across different platforms:

- **macOS**: Bash, Zsh, Fish
- **Linux**: Bash, Zsh, Fish
- **Windows**: PowerShell, WSL with Bash/Zsh/Fish

## Technical Details

The wrapper works by:

1. Intercepting the `projj` command
2. Executing the original command and capturing its output
3. Parsing the output for the special `PROJJ_CHANGE_DIRECTORY=` line
4. Using the shell's `cd` command to change directories

This approach is necessary because external programs cannot directly change the working directory of their parent shell process.