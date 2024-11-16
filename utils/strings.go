package utils

import "strings"

func SplitString(message string) []string {
	return SplitStringWithDelimiter(message, " ")
}

func SplitStringWithDelimiter(message string, delimiter string) []string {
	return strings.Split(message, delimiter)
}

func EscapeMarkdown(message string) string {
	message = strings.ReplaceAll(message, "_", "\\_")
	message = strings.ReplaceAll(message, "*", "\\*")
	message = strings.ReplaceAll(message, "[", "\\[")
	message = strings.ReplaceAll(message, "]", "\\]")
	message = strings.ReplaceAll(message, "~", "\\~")
	message = strings.ReplaceAll(message, "`", "\\`")
	message = strings.ReplaceAll(message, ">", "\\>")
	message = strings.ReplaceAll(message, "#", "\\#")
	message = strings.ReplaceAll(message, "+", "\\+")
	message = strings.ReplaceAll(message, "-", "\\-")
	message = strings.ReplaceAll(message, "=", "\\=")
	message = strings.ReplaceAll(message, "|", "\\|")
	message = strings.ReplaceAll(message, "{", "\\{")
	message = strings.ReplaceAll(message, "}", "\\}")
	message = strings.ReplaceAll(message, ".", "\\.")
	message = strings.ReplaceAll(message, "!", "\\!")
	return message
}
