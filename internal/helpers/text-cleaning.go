package helpers

import (
	"fmt"
	"strings"
)

func FormatDate(date string) string {
	nums := strings.Split(strings.Split(date, "T")[0], "-")
	return fmt.Sprintf("%v-%v-%v", nums[2], nums[1], nums[0])
}

func CleanName(name string) string {
	unallowedChars := []string{"/", "\\", ":", "*", "?", "\"", "<", ">", "|"}

	for i := range unallowedChars {
		name = strings.ReplaceAll(name, unallowedChars[i], " ")
	}
	return name
}
