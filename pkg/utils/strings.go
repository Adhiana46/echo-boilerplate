package utils

import "regexp"

func ContainAlpha(text string) bool {
	return regexp.MustCompile(`[a-zA-Z]+`).MatchString(text)
}

func ContainAlphaLower(text string) bool {
	return regexp.MustCompile(`[a-z]+`).MatchString(text)
}

func ContainAlphaUpper(text string) bool {
	return regexp.MustCompile(`[A-Z]+`).MatchString(text)
}

func ContainNumeric(text string) bool {
	return regexp.MustCompile(`[0-9]+`).MatchString(text)
}
