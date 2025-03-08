package env

import (
	"bufio"
	"os"
	"strings"
)

// LoadEnv reads key-value pairs from the file at the given filename and sets
// them as environment variables. Empty lines and lines starting with '#' are
// ignored. If a line cannot be split into a key-value pair (i.e. it does not
// contain a '='), it is also ignored.
func LoadEnv(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		// ignore empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		// split key-calue pairs
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		// set environment variable
		os.Setenv(key, value)
	}
	return scanner.Err()
}
