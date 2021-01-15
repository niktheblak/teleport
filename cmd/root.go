package cmd

import (
	"fmt"
	"os"
	"os/user"
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
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func isCommand(cmd *cobra.Command, s string) bool {
	for _, c := range cmd.Commands() {
		if s == c.Name() {
			return true
		}
		for _, a := range c.Aliases {
			if s == a {
				return true
			}
		}
	}
	return false
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
		return fmt.Sprintf("%s/%s", teleportHome, teleportFileName), nil
	}
	u, err := user.Current()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s/%s", u.HomeDir, teleportFileName), nil
}
