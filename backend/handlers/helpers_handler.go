package handlers

import (
	"strings"
	"unicode"
)

// คงชื่อให้ตรงกับ handler ที่เรียกใช้
func OnlyDigits(s string) string {
	var b strings.Builder
	for _, r := range s {
		if unicode.IsDigit(r) {
			b.WriteRune(r)
		}
	}
	return b.String()
}

func Norm(s string) string {
	return strings.TrimSpace(s)
}
