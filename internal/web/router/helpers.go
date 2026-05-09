package router

import "strings"

func join(parts ...string) string {
	var cleaned []string

	for _, p := range parts {
		if p == "" {
			continue
		}

		cleaned = append(cleaned, strings.Trim(p, "/"))
	}

	if len(cleaned) == 0 {
		return "/"
	}

	result := "/" + strings.Join(cleaned, "/")

	// Preserve trailing slash if the last non-empty part had one
	for i := len(parts) - 1; i >= 0; i-- {
		if parts[i] != "" {
			if strings.HasSuffix(parts[i], "/") && parts[i] != "/" {
				result += "/"
			}
			break
		}
	}

	return result
}
