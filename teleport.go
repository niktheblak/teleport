package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/user"
	"sort"
	"strings"
)

var commands = []string{"add", "remove", "rm", "list", "ls"}

const WarpPointsFile = ".tp"

func main() {
	if len(os.Args) == 1 {
		printUsage()
		os.Exit(1)
	}
	cmd := os.Args[1]
	args := os.Args[2:]
	switch cmd {
	case "add":
		switch len(args) {
		case 0:
			printUsage()
			os.Exit(1)
		case 1:
			key := args[0]
			if isCommand(key) {
				fmt.Fprintf(os.Stderr, "%s cannot be used as warp point key\n", key)
				os.Exit(1)
			}
			dir, err := os.Getwd()
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(2)
			}
			err = addWarpPoint(key, dir)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(2)
			}
		case 2:
			key := args[0]
			if isCommand(key) {
				fmt.Fprintf(os.Stderr, "%s cannot be used as warp point key\n", key)
				os.Exit(1)
			}
			dir := args[1]
			err := addWarpPoint(key, dir)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(2)
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
			os.Exit(2)
		}
	case "list", "ls":
		err := listWarpPoints()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(2)
		}
	default:
		target := cmd
		err := warpTo(target)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(2)
		}
	}
}

func printUsage() {
	fmt.Println(`Usage: tp {command} [args]

Supported commands are:
goto [key]
	changes current directory to to the warp point with the given key
add [key] [dir]
	adds warp point to current directory or the specified directory
remove
	removes current directory from warp points
list
	lists warp points`)
}

func isCommand(s string) bool {
	for _, c := range commands {
		if s == c {
			return true
		}
	}
	return false
}

func warpTo(target string) error {
	wps, err := loadWarpPoints()
	if err != nil {
		return err
	}
	dir, ok := wps[target]
	if !ok {
		return fmt.Errorf("Warp point %s does not exist", target)
	}
	fmt.Println(dir)
	return nil
}

func addWarpPoint(key, dir string) error {
	wps, err := loadWarpPoints()
	if err != nil {
		return err
	}
	wps[key] = dir
	return saveWarpPoints(wps)
}

func listWarpPoints() error {
	wps, err := loadWarpPoints()
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
	wps, err := loadWarpPoints()
	if err != nil {
		return err
	}
	delete(wps, key)
	return saveWarpPoints(wps)
}

func removeCurrentDirWarpPoint() error {
	dir, err := os.Getwd()
	if err != nil {
		return err
	}
	wps, err := loadWarpPoints()
	if err != nil {
		return err
	}
	for key, val := range wps {
		if val == dir {
			delete(wps, key)
			break
		}
	}
	return saveWarpPoints(wps)
}

func loadWarpPoints() (map[string]string, error) {
	fileName, err := warpPointsFile()
	if err != nil {
		return nil, err
	}
	f, err := os.Open(fileName)
	switch err.(type) {
	case *os.PathError:
		return make(map[string]string), nil
	case nil:
	default:
		return nil, err
	}
	defer f.Close()
	return readWarpPointsFile(f)
}

func saveWarpPoints(warpPoints map[string]string) error {
	fileName, err := warpPointsFile()
	if err != nil {
		return err
	}
	f, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer f.Close()
	return writeWarpPointsFile(f, warpPoints)
}

func warpPointsFile() (string, error) {
	u, err := user.Current()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s/%s", u.HomeDir, WarpPointsFile), nil
}

func readWarpPointsFile(r io.Reader) (map[string]string, error) {
	wps := make(map[string]string)
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, ";") || strings.HasPrefix(line, "#") {
			continue
		}
		if strings.HasPrefix(line, "[") {
			continue
		}
		tokens := strings.SplitN(line, "=", 2)
		if len(tokens) != 2 {
			return nil, fmt.Errorf("Invalid warp point: %s", line)
		}
		key := strings.TrimSpace(tokens[0])
		dir := strings.TrimSpace(tokens[1])
		wps[key] = dir
	}

	if err := scanner.Err(); err != nil {
		return wps, err
	}
	return wps, nil
}

func writeWarpPointsFile(w io.Writer, warpPoints map[string]string) error {
	var keys []string
	for key := range warpPoints {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		_, err := fmt.Fprintf(w, "%s = %s\n", key, warpPoints[key])
		if err != nil {
			return err
		}
	}
	return nil
}
