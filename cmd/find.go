package cmd

import (
	"fmt"
	"os"

	"github.com/atian25/projj-go/internal/cache"
	"github.com/atian25/projj-go/internal/picker"
	"github.com/spf13/cobra"
)

func newFindCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "find [query]",
		Short: "查找仓库路径",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			query := ""
			if len(args) > 0 {
				query = args[0]
			}

			c, err := cache.Load(cache.DefaultCachePath())
			if err != nil {
				return err
			}

			keys := c.FindKeys(query)
			if len(keys) == 0 {
				fmt.Fprintf(os.Stderr, "未找到匹配的仓库: %q\n", query)
				os.Exit(1)
			}

			if len(keys) == 1 {
				fmt.Println(c[keys[0]])
				return nil
			}

			pathByKey := make(map[string]string, len(keys))
			for _, k := range keys {
				pathByKey[k] = c[k]
			}

			selected, err := picker.Pick(keys, pathByKey)
			if err != nil {
				if err == picker.ErrCancelled {
					os.Exit(1)
				}
				return err
			}

			fmt.Println(selected)
			return nil
		},
	}
}
