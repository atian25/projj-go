package cmd

import (
	"fmt"

	"github.com/atian25/projj-go/internal/cache"
	"github.com/spf13/cobra"
)

func newListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "列出所有缓存的仓库",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := cache.Load(cache.DefaultCachePath())
			if err != nil {
				return err
			}
			for key := range c {
				fmt.Println(key)
			}
			return nil
		},
	}
}
