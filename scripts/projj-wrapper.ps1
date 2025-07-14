# projj shell wrapper for PowerShell
# This wrapper enables automatic directory changing after 'projj add'

function projj {
    param(
        [Parameter(ValueFromRemainingArguments=$true)]
        [string[]]$Arguments
    )
    
    $projjCmd = "projj"
    
    # Check if projj binary exists in PATH
    if (-not (Get-Command $projjCmd -ErrorAction SilentlyContinue)) {
        # Try to find projj-go binary
        if (Get-Command "projj-go" -ErrorAction SilentlyContinue) {
            $projjCmd = "projj-go"
        } else {
            Write-Error "Error: projj command not found in PATH"
            return 1
        }
    }
    
    # Execute the original projj command and capture output
    try {
        $output = & $projjCmd @Arguments 2>&1
        $exitCode = $LASTEXITCODE
    } catch {
        Write-Error "Error executing projj: $_"
        return 1
    }
    
    # Print the output
    $output | ForEach-Object { Write-Host $_ }
    
    # Check if this was an 'add' or 'find' command and if change_directory is enabled
    if ($exitCode -eq 0 -and $Arguments.Length -gt 0 -and ($Arguments[0] -eq "add" -or $Arguments[0] -eq "find")) {
        # Look for the special PROJJ_CHANGE_DIRECTORY line in output
        $changeDirLine = $output | Where-Object { $_ -match "^PROJJ_CHANGE_DIRECTORY=" }
        
        if ($changeDirLine) {
            $changeDir = $changeDirLine -replace "^PROJJ_CHANGE_DIRECTORY=", ""
            
            if ($changeDir -and (Test-Path $changeDir -PathType Container)) {
                Write-Host "Changing directory to: $changeDir"
                try {
                    Set-Location $changeDir
                } catch {
                    Write-Error "Error: Failed to change directory to $changeDir"
                    return 1
                }
            }
        }
    }
    
    return $exitCode
}

# Create an alias for easier access
Set-Alias -Name projj -Value projj -Force