package warppoint

import (
	"io"
	"bufio"
	"strings"
	"fmt"
	"sort"
	"os"
)

func ReadFromFile(fileName string) (map[string]string, error) {
	f, err := os.Open(fileName)
	switch err.(type) {
	case *os.PathError:
		// File does not exist, return empty map
		return make(map[string]string), nil
	case nil:
	default:
		return nil, err
	}
	defer f.Close()
	return Read(f)
}

func WriteToFile(fileName string, warpPoints map[string]string) error {
	f, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer f.Close()
	return Write(f, warpPoints)
}

func Read(r io.Reader) (map[string]string, error) {
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
			return nil, fmt.Errorf("invalid warp point: %s", line)
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

func Write(w io.Writer, warpPoints map[string]string) error {
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
