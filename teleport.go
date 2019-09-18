/*
 * Copyright 2017 Niko Korhonen
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

// main executable for the teleport command
package main

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/niktheblak/teleport/pkg/warppoint"
	"github.com/urfave/cli"
)

var commands = []string{"help", "add", "remove", "rm", "list", "ls"}

var (
	teleportHome     = ""
	teleportFileName = ".tp"
)

func main() {
	app := cli.NewApp()
	app.Name = "teleport"
	app.Usage = "Tool for rapidly switching between directories"
	app.Version = "1.0.0"
	app.Compiled = time.Now()
	app.Authors = []cli.Author{
		{
			Name:  "Niko Korhonen",
			Email: "niko@bitnik.fi",
		},
	}
	app.Copyright = "(c) 2018 Niko Korhonen"
	app.Commands = []cli.Command{
		{
			Name:      "add",
			Aliases:   []string{"a"},
			Usage:     "adds a warp point to the current directory or to the specified directory",
			ArgsUsage: "{warp point} [dir]",
			Action: func(c *cli.Context) error {
				switch c.NArg() {
				case 0:
					return fmt.Errorf("warp point name is required")
				case 1:
					key := c.Args().First()
					if !isValidKey(key) {
						return fmt.Errorf("%s cannot be used as warp point key\n", key)
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
					key := c.Args().First()
					if !isValidKey(key) {
						return fmt.Errorf("%s cannot be used as warp point key\n", key)
					}
					dir, err := filepath.Abs(c.Args().Get(1))
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
		},
		{
			Name:      "remove",
			Aliases:   []string{"rm"},
			Usage:     "removes warp point pointing to the current directory or the specified warp point",
			ArgsUsage: "[dir]",
			Action: func(c *cli.Context) error {
				switch c.NArg() {
				case 0:
					return removeCurrentDirWarpPoint()
				case 1:
					key := c.Args().First()
					return removeWarpPoint(key)
				}
				return nil
			},
		},
		{
			Name:    "list",
			Aliases: []string{"ls"},
			Usage:   "lists warp points",
			Action: func(c *cli.Context) error {
				return listWarpPoints()
			},
		},
		{
			Name:      "warp",
			Aliases:   []string{"w"},
			Usage:     "teleport to the specified warp point",
			ArgsUsage: "{warp point}",
			Action: func(c *cli.Context) error {
				return warpTo(c.Args().First())
			},
		},
	}
	tpHome, ok := os.LookupEnv("TELEPORT_HOME")
	if ok {
		teleportHome = tpHome
	}
	tpFile, ok := os.LookupEnv("TELEPORT_FILE")
	if ok {
		teleportFileName = tpFile
	}
	if err := app.Run(os.Args); err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}
}

func isValidKey(key string) bool {
	return !isCommand(key) && !strings.Contains(key, "=")
}

func isCommand(s string) bool {
	for _, c := range commands {
		if s == c {
			return true
		}
	}
	return false
}

func warpTo(key string) error {
	f, err := warpPointsFile()
	if err != nil {
		return err
	}
	wps, err := warppoint.ReadFromFile(f)
	if err != nil {
		return err
	}
	dir, ok := wps[key]
	if !ok {
		return fmt.Errorf("warp point %s does not exist", key)
	}
	fmt.Println(dir)
	return nil
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

func listWarpPoints() error {
	f, err := warpPointsFile()
	if err != nil {
		return err
	}
	wps, err := warppoint.ReadFromFile(f)
	if err != nil {
		return err
	}
	var sorted []string
	for key := range wps {
		sorted = append(sorted, key)
	}
	sort.Strings(sorted)
	for _, wp := range sorted {
		fmt.Printf("%s\t%s\n", wp, wps[wp])
	}
	return nil
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
	for key, val := range wps {
		if val == dir {
			delete(wps, key)
			break
		}
	}
	return warppoint.WriteToFile(f, wps)
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
