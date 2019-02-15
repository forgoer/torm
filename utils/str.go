package utils

import "strings"

var (
	// The cache of snake-cased words.
	snakeCache = make(map[string]string)

	// The cache of studly-cased words.
	studlyCache = make(map[string]string)
)

// SnakeCase Convert a string to snake case, XxYy to xx_yy , XxYY to xx_yy
func SnakeCase(s string) string {
	if v, ok := snakeCache[s]; ok {
		return v
	}
	data := make([]byte, 0, len(s)*2)
	flag := false
	num := len(s)
	for i := 0; i < num; i++ {
		d := s[i]
		if i > 0 && d >= 'A' && d <= 'Z' && !flag {
			data = append(data, '_')
			flag = true
		} else {
			flag = false
		}
		data = append(data, d)
	}
	v := strings.ToLower(string(data[:]))
	return v
}

// StudlyCase Convert a value to studly caps case, xx_yy to XxYy
func StudlyCase(s string) string {
	if v, ok := studlyCache[s]; ok {
		return v
	}
	data := make([]byte, 0, len(s))
	flag, num := true, len(s)-1
	for i := 0; i <= num; i++ {
		d := s[i]
		if d == '_' {
			flag = true
			continue
		} else if flag {
			if d >= 'a' && d <= 'z' {
				d = d - 32
			}
			flag = false
		}
		data = append(data, d)
	}
	v := string(data[:])
	return v
}
