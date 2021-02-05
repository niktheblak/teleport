package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/niktheblak/teleport/pkg/warppoint"
)

var pruneCmd = &cobra.Command{
	Use:   "prune",
	Args:  cobra.NoArgs,
	Short: "Removes warp points whose target directory no longer exist",
	RunE: func(cmd *cobra.Command, args []string) error {
		f, err := warpPointsFile()
		if err != nil {
			return err
		}
		wps, err := warppoint.ReadFromFile(f)
		if err != nil {
			return err
		}
		var pruned []string
		for wp, path := range wps {
			target, err := os.Open(path)
			if err != nil {
				if os.IsNotExist(err) {
					pruned = append(pruned, wp)
					continue
				} else {
					return err
				}
			}
			info, err := target.Stat()
			if err != nil {
				pruned = append(pruned, wp)
				continue
			}
			if !info.IsDir() {
				pruned = append(pruned, wp)
				continue
			}
		}
		for _, wp := range pruned {
			cmd.Printf("Removing warp point %s to %s\n", wp, wps[wp])
			delete(wps, wp)
		}
		return warppoint.WriteToFile(f, wps)
	},
}

func init() {
	rootCmd.AddCommand(pruneCmd)
}
