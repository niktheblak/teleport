package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/niktheblak/teleport/pkg/warppoint"
)

var addCmd = &cobra.Command{
	Use:   "add name",
	Args:  cobra.RangeArgs(1, 2),
	Short: "Adds a warp point to the current directory or to the specified directory",
	RunE: func(cmd *cobra.Command, args []string) error {
		switch len(args) {
		case 1:
			key := strings.TrimSpace(args[0])
			if !isValidKey(cmd.Root(), key) {
				return fmt.Errorf("%s cannot be used as warp point key", key)
			}
			dir, err := os.Getwd()
			if err != nil {
				return err
			}
			err = addWarpPoint(key, dir)
			if err != nil {
				return err
			}
		case 2:
			key := strings.TrimSpace(args[0])
			if !isValidKey(cmd.Root(), key) {
				return fmt.Errorf("%s cannot be used as warp point key", key)
			}
			dir, err := filepath.Abs(args[1])
			if err != nil {
				return err
			}
			err = addWarpPoint(key, dir)
			if err != nil {
				return err
			}
		}
		return nil
	},
}

func addWarpPoint(key, dir string) error {
	f, err := warpPointsFile()
	if err != nil {
		return err
	}
	wps, err := warppoint.ReadFromFile(f)
	if err != nil {
		return err
	}
	wps[key] = dir
	return warppoint.WriteToFile(f, wps)
}

func init() {
	rootCmd.AddCommand(addCmd)
}
