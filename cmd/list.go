package cmd

import (
	"sort"

	"github.com/spf13/cobra"

	"github.com/niktheblak/teleport/pkg/warppoint"
)

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Args:    cobra.NoArgs,
	Short:   "Lists warp points",
	RunE: func(cmd *cobra.Command, args []string) error {
		f, err := warpPointsFile()
		if err != nil {
			return err
		}
		wps, err := warppoint.ReadFromFile(f)
		if err != nil {
			return err
		}
		var maxLen int
		var sorted []string
		for key := range wps {
			sorted = append(sorted, key)
			if len(key) > maxLen {
				maxLen = len(key)
			}
		}
		sort.Strings(sorted)
		for _, wp := range sorted {
			cmd.Printf("%-*s %s\n", maxLen, wp, wps[wp])
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
