package warppoint

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
)

var commentPrefixes = []string{
	";",
	"#",
	"[",
}

// ReadFromFile reads a collection of warp points from a file
func ReadFromFile(fileName string) (map[string]string, error) {
	f, err := os.Open(fileName)
	switch err.(type) {
	case *os.PathError:
		// File does not exist, meaning there aren't any stored warp points.
		// Return an empty map.
		return make(map[string]string), nil
	case nil:
	default:
		return nil, err
	}
	defer f.Close()
	return Read(f)
}

// WriteToFile writes a collection of warp points to a file
func WriteToFile(fileName string, warpPoints map[string]string) error {
	f, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer f.Close()
	return Write(f, warpPoints)
}

// Read reads a collection of warp points from an io.Reader.
// Warp points are serialized as one key-value pair per line where the
// key and the value (directory) are separated with an equals sign (=).
// Empty lines and commented lines that begin with ;, # or [ are ignored.
func Read(r io.Reader) (map[string]string, error) {
	wps := make(map[string]string)
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || isComment(line) {
			continue
		}
		tokens := strings.SplitN(line, "=", 2)
		if len(tokens) != 2 {
			return nil, fmt.Errorf("invalid warp point: %s", line)
		}
		key := strings.TrimSpace(tokens[0])
		dir := strings.TrimSpace(tokens[1])
		wps[key] = dir
	}

	return wps, scanner.Err()
}

// Write writes a collection of warp points to an io.Writer.
// Warp points are written one per line, sorted alphabetically by key
func Write(w io.Writer, warpPoints map[string]string) error {
	var keys []string
	for key := range warpPoints {
		keys = append(keys, key)
	}
	// Always write warp points sorted by key
	sort.Strings(keys)
	for _, key := range keys {
		_, err := fmt.Fprintf(w, "%s = %s\n", key, warpPoints[key])
		if err != nil {
			return err
		}
	}
	return nil
}

func isComment(line string) bool {
	for _, p := range commentPrefixes {
		if strings.HasPrefix(line, p) {
			return true
		}
	}
	return false
}
