package cmd

import (
	"os"
	"os/user"
	"path"
	"slices"
	"sort"
	"strings"
	"unicode"

	"github.com/spf13/cobra"
)

var (
	teleportHome     = ""
	teleportFileName = ".tp"
)

var rootCmd = &cobra.Command{
	Use:   "teleport",
	Short: "Tool for rapidly switching between directories",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		tpHome, ok := os.LookupEnv("TELEPORT_HOME")
		if ok {
			teleportHome = tpHome
		}
		tpFile, ok := os.LookupEnv("TELEPORT_FILE")
		if ok {
			teleportFileName = tpFile
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func isCommand(cmd *cobra.Command, s string) bool {
	return slices.ContainsFunc(cmd.Commands(), func(c *cobra.Command) bool {
		if s == c.Name() {
			return true
		}
		if slices.Contains(c.Aliases, s) {
			return true
		}
		return false
	})
}

func isValidKey(cmd *cobra.Command, key string) bool {
	if isCommand(cmd, key) || strings.Contains(key, "=") {
		return false
	}
	for _, r := range key {
		if unicode.IsSpace(r) || !unicode.IsPrint(r) {
			return false
		}
	}
	return true
}

func warpPointsFile() (string, error) {
	if teleportHome != "" {
		return path.Join(teleportHome, teleportFileName), nil
	}
	u, err := user.Current()
	if err != nil {
		return "", err
	}
	return path.Join(u.HomeDir, teleportFileName), nil
}

func keysAndMaxLen(wps map[string]string) ([]string, int) {
	var maxLen int
	var sorted []string
	for key := range wps {
		sorted = append(sorted, key)
		if len(key) > maxLen {
			maxLen = len(key)
		}
	}
	sort.Strings(sorted)
	return sorted, maxLen
}
