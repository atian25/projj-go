package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "projj",
	Short: "本地代码仓库管理工具",
	Long:  "projj 按 $BASE/{host}/{owner}/{repo} 层级结构管理 Git 仓库。",
}

// Execute runs the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(
		newInitCmd(),
		newAddCmd(),
		newInternalAddCmd(),
		newFindCmd(),
		newListCmd(),
		newImportCmd(),
		newSyncCmd(),
		newRunCmd(),
		newRunAllCmd(),
		newShellInitCmd(),
	)
}
