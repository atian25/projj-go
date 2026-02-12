package cmd

import (
	"fmt"
	"os"

	"github.com/atian25/projj-go/internal/cache"
	"github.com/spf13/cobra"
)

func newSyncCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "sync",
		Short: "校验缓存中的仓库是否存在，清理失效项",
		RunE: func(cmd *cobra.Command, args []string) error {
			cachePath := cache.DefaultCachePath()
			c, err := cache.Load(cachePath)
			if err != nil {
				return err
			}

			removed := 0
			for key, path := range c {
				if _, err := os.Stat(path); os.IsNotExist(err) {
					fmt.Printf("移除失效项: %s (%s)\n", key, path)
					delete(c, key)
					removed++
				}
			}

			if removed > 0 {
				if err := cache.Save(cachePath, c); err != nil {
					return err
				}
				fmt.Printf("共移除 %d 个失效项\n", removed)
			} else {
				fmt.Println("缓存状态正常，无失效项")
			}
			return nil
		},
	}
}
