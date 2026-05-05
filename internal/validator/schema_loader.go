package validator

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

// LoadSchema parses a simple schema definition file into a Schema.
//
// Each non-blank, non-comment line has the form:
//
//	KEY [required] [pattern=<regex>] [allowed=val1,val2,...]
//
// Example:
//
//	APP_ENV required allowed=development,staging,production
//	PORT    required pattern=^\d+$
//	LOG_LEVEL allowed=debug,info,warn,error
func LoadSchema(path string) (Schema, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("schema_loader: open %q: %w", path, err)
	}
	defer f.Close()

	schema := make(Schema)
	scanner := bufio.NewScanner(f)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.Fields(line)
		if len(parts) == 0 {
			continue
		}

		key := parts[0]
		rule := Rule{}

		for _, token := range parts[1:] {
			switch {
			case token == "required":
				rule.Required = true

			case strings.HasPrefix(token, "pattern="):
				raw := strings.TrimPrefix(token, "pattern=")
				re, err := regexp.Compile(raw)
				if err != nil {
					return nil, fmt.Errorf("schema_loader: line %d: invalid pattern %q: %w", lineNum, raw, err)
				}
				rule.Pattern = re

			case strings.HasPrefix(token, "allowed="):
				raw := strings.TrimPrefix(token, "allowed=")
				rule.AllowedValues = strings.Split(raw, ",")

			default:
				return nil, fmt.Errorf("schema_loader: line %d: unknown token %q", lineNum, token)
			}
		}

		schema[key] = rule
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("schema_loader: scanning %q: %w", path, err)
	}

	return schema, nil
}
