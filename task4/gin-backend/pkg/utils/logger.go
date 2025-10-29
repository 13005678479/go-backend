package utils

// 字符串截断（用于文章列表显示摘要）
func TruncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

// uint转字符串
func UintToString(u uint) string {
	return string(rune(u))
}