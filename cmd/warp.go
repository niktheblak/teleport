package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/niktheblak/teleport/pkg/warppoint"
)

var warpCmd = &cobra.Command{
	Use:     "warp",
	Aliases: []string{"w"},
	Args:    cobra.ExactArgs(1),
	Short:   "Teleport to the specified warp point",
	RunE: func(cmd *cobra.Command, args []string) error {
		f, err := warpPointsFile()
		if err != nil {
			return err
		}
		wps, err := warppoint.ReadFromFile(f)
		if err != nil {
			return err
		}
		key := args[0]
		dir, ok := wps[key]
		if !ok {
			return fmt.Errorf("warp point %s does not exist", key)
		}
		fmt.Println(dir)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(warpCmd)
}
