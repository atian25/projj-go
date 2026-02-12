package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

const bashTemplate = `projj() {
  if [ "$1" = "add" ]; then
    local dir
    dir=$(command projj _add "${@:2}") && cd "$dir"
  elif [ "$1" = "find" ]; then
    local dir
    dir=$(command projj find "${@:2}") && cd "$dir"
  else
    command projj "$@"
  fi
}`

func newShellInitCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "shell-init",
		Short: "输出 shell 集成代码（在 .zshrc/.bashrc 中执行 eval \"$(projj shell-init)\"）",
		RunE: func(cmd *cobra.Command, args []string) error {
			shell := detectShell()
			switch shell {
			case "bash", "zsh":
				fmt.Println(bashTemplate)
			default:
				fmt.Fprintf(os.Stderr, "不支持的 shell: %s（仅支持 bash 和 zsh）\n", shell)
				os.Exit(1)
			}
			return nil
		},
	}
}

func detectShell() string {
	shell := os.Getenv("SHELL")
	return filepath.Base(shell)
}
