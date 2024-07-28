package typedfetch

import (
	"sort"
	"strings"
)

// Capitalize the first letter, and make the rest lower case
func pascalize(s string) string {
	if len(s) < 1 {
		return s
	}

	return strings.ToUpper(s[:1]) + strings.ToLower(s[1:])
}

// Capitalize the first letter
func capitalize(s string) string {
	if len(s) < 1 {
		return s
	}

	return strings.ToUpper(s[:1]) + s[1:]
}

// Convert /path/to/{var} to PathToVar
func pathToVar(path string) string {
	// remove invalid variable characters: -{}.~%
	path = strings.ReplaceAll(path, "-", "")
	path = strings.ReplaceAll(path, "{", "")
	path = strings.ReplaceAll(path, "}", "")
	path = strings.ReplaceAll(path, ".", "")
	path = strings.ReplaceAll(path, "~", "")
	path = strings.ReplaceAll(path, "%", "")

	// Split on /, capitalize each part, and join
	parts := strings.Split(path, "/")
	for i, part := range parts {
		if part == "" {
			continue
		}
		parts[i] = capitalize(part)
	}

	return strings.Join(parts, "")
}

func itemInSlice[item comparable](slice []item, i item) bool {
	for _, s := range slice {
		if s == i {
			return true
		}
	}
	return false
}

func isValidJsonType(t string) bool {
	switch t {
	case "string", "number", "integer", "boolean", "array", "object":
		return true
	default:
		return false
	}
}

func sortedMapKeys[T any](m map[string]T) []string {
	keys := []string{}
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
