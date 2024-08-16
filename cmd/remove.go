package cmd

import (
	"maps"
	"os"

	"github.com/spf13/cobra"

	"github.com/niktheblak/teleport/pkg/warppoint"
)

var removeCmd = &cobra.Command{
	Use:     "remove",
	Aliases: []string{"rm"},
	Short:   "Removes warp point pointing to the current directory or the specified warp points",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return removeCurrentDirWarpPoint()
		}
		for _, key := range args {
			if err := removeWarpPoint(key); err != nil {
				return err
			}
		}
		return nil
	},
}

func removeCurrentDirWarpPoint() error {
	dir, err := os.Getwd()
	if err != nil {
		return err
	}
	f, err := warpPointsFile()
	if err != nil {
		return err
	}
	wps, err := warppoint.ReadFromFile(f)
	if err != nil {
		return err
	}
	maps.DeleteFunc(wps, func(key string, val string) bool {
		return val == dir
	})
	return warppoint.WriteToFile(f, wps)
}

func removeWarpPoint(key string) error {
	f, err := warpPointsFile()
	if err != nil {
		return err
	}
	wps, err := warppoint.ReadFromFile(f)
	if err != nil {
		return err
	}
	delete(wps, key)
	return warppoint.WriteToFile(f, wps)
}

func init() {
	rootCmd.AddCommand(removeCmd)
}
