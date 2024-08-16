package cmd

import (
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
		keys, maxLen := keysAndMaxLen(wps)
		for _, wp := range keys {
			cmd.Printf("%-*s %s\n", maxLen, wp, wps[wp])
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
