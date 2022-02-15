package util

import "strings"

func GetSuffix(filename string) string {
	return filename[strings.LastIndex(filename, ".")+1:]
}
