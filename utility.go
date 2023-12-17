package main

import (
	"encoding/base64"
	"regexp"
)

// checks if a secret is base64 encoded
func IsBase64Encoded(val string) bool {
	base64Regex, _ := regexp.MatchString(Base64Pattern, val)
	return base64Regex
}

// decode a string in case it is base 64 encoded
func DecodeBase64(encodedString string) (string, error) {
	decoded, err := base64.StdEncoding.DecodeString(encodedString)
	if err != nil {
		return "", err
	}
	return string(decoded), nil
}
