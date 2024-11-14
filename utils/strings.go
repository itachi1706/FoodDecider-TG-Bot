package utils

import "strings"

func SplitString(message string) []string {
    return SplitStringWithDelimiter(message, " ")
}

func SplitStringWithDelimiter(message string, delimiter string) []string {
    return strings.Split(message, delimiter)
}
