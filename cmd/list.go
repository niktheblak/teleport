package cmd

import (
	"fmt"
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
		format := fmt.Sprintf("%%-%ds %%s\n", maxLen)
		for _, wp := range sorted {
			cmd.Printf(format, wp, wps[wp])
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
