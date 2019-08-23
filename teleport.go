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
	"sort"
	"strings"

	"github.com/niktheblak/teleport/pkg/warppoint"
)

var commands = []string{"help", "add", "remove", "rm", "list", "ls"}

var (
	teleportHome     = ""
	teleportFileName = ".tp"
)

func main() {
	if len(os.Args) == 1 {
		printUsage()
		os.Exit(3)
	}
	tpHome, ok := os.LookupEnv("TELEPORT_HOME")
	if ok {
		teleportHome = tpHome
	}
	tpFile, ok := os.LookupEnv("TELEPORT_FILE")
	if ok {
		teleportFileName = tpFile
	}
	cmd := os.Args[1]
	args := os.Args[2:]
	switch cmd {
	case "help", "-h", "--help":
		printUsage()
		os.Exit(0)
	case "add":
		switch len(args) {
		case 0:
			printUsage()
			os.Exit(3)
		case 1:
			key := args[0]
			if !isValidKey(key) {
				fmt.Fprintf(os.Stderr, "%s cannot be used as warp point key\n", key)
				os.Exit(3)
			}
			dir, err := os.Getwd()
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(4)
			}
			err = addWarpPoint(key, dir)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(4)
			}
		case 2:
			key := args[0]
			if !isValidKey(key) {
				fmt.Fprintf(os.Stderr, "%s cannot be used as warp point key\n", key)
				os.Exit(3)
			}
			dir := args[1]
			err := addWarpPoint(key, dir)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(4)
			}
		}
	case "remove", "rm":
		var err error
		switch len(args) {
		case 0:
			err = removeCurrentDirWarpPoint()
		case 1:
			key := args[0]
			err = removeWarpPoint(key)
		}
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(4)
		}
	case "list", "ls":
		err := listWarpPoints()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(4)
		}
	default:
		target := cmd
		err := warpTo(target)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(4)
		}
	}
}

func printUsage() {
	fmt.Println(`Usage: tp {command} [args]

Supported commands are:
[key]
	changes current directory to to the warp point with the given key
add [key]
	adds warp point to the current directory
add [key] [dir]
	adds warp point to the specified directory
remove [key]
	removes key from warp points
list
	lists warp points`)
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
