package loader

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// EnvMap represents a set of key-value environment variable pairs.
type EnvMap map[string]string

// LoadFile reads a .env file and returns its key-value pairs.
// Lines starting with '#' are treated as comments and ignored.
// Empty lines are also ignored.
func LoadFile(path string) (EnvMap, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("loader: could not open file %q: %w", path, err)
	}
	defer f.Close()

	env := make(EnvMap)
	scanner := bufio.NewScanner(f)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		key, value, err := parseLine(line)
		if err != nil {
			return nil, fmt.Errorf("loader: parse error in %q at line %d: %w", path, lineNum, err)
		}

		env[key] = value
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("loader: error reading file %q: %w", path, err)
	}

	return env, nil
}

// parseLine splits a line into a key and value on the first '=' character.
// It trims surrounding whitespace and strips optional surrounding quotes from values.
func parseLine(line string) (string, string, error) {
	idx := strings.IndexByte(line, '=')
	if idx < 0 {
		return "", "", fmt.Errorf("missing '=' in line %q", line)
	}

	key := strings.TrimSpace(line[:idx])
	if key == "" {
		return "", "", fmt.Errorf("empty key in line %q", line)
	}

	value := strings.TrimSpace(line[idx+1:])
	value = stripQuotes(value)

	return key, value, nil
}

// stripQuotes removes matching surrounding double or single quotes from a value.
func stripQuotes(s string) string {
	if len(s) >= 2 {
		if (s[0] == '"' && s[len(s)-1] == '"') ||
			(s[0] == '\'' && s[len(s)-1] == '\'') {
			return s[1 : len(s)-1]
		}
	}
	return s
}
