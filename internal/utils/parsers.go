package utils

import (
	"fmt"
	"strings"
)

func StringParseBoolean(str string) (bool, error) {
	lowercaseStr := strings.ToLower(str)

	switch lowercaseStr {
	case "true", "yes", "1":
		return true, nil
	case "false", "no", "0":
		return false, nil
	default:
		return false, fmt.Errorf("invalid value to cast bool: %s", str)
	}
}
