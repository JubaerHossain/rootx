package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"time"
)

type CustomDateFormat time.Time

// MarshalJSON implements the json.Marshaler interface
func (c CustomDateFormat) MarshalJSON() ([]byte, error) {
	t := time.Time(c)
	formatted := t.Format("2006-01-02") // Format the time as "YYYY-MM-DD"
	return json.Marshal(formatted)
}

func GenerateUniqueNumber(length int) string {
	var buffer bytes.Buffer
	for i := 0; i < length; i++ {
		b := rand.Intn(10)                       // Generate a random digit (0-9)
		buffer.WriteString(fmt.Sprintf("%d", b)) // Convert digit to string and append
	}
	return buffer.String()
}

func Slugify(s string) string {
	// Replace all non-alphanumeric characters with a space
	s = replaceNonAlphanumeric(s, ' ')

	// Replace all spaces with a hyphen
	s = replaceSpaces(s, '-')

	// Convert to lowercase
	s = toLower(s)

	return s
}

func replaceNonAlphanumeric(s string, r rune) string {
	var buffer bytes.Buffer
	for _, c := range s {
		if isAlphanumeric(c) {
			buffer.WriteRune(c)
		} else {
			buffer.WriteRune(r)
		}
	}
	return buffer.String()
}

func replaceSpaces(s string, r rune) string {
	var buffer bytes.Buffer
	for _, c := range s {
		if c == ' ' {
			buffer.WriteRune(r)
		} else {
			buffer.WriteRune(c)
		}
	}
	return buffer.String()
}

func toLower(s string) string {
	var buffer bytes.Buffer
	for _, c := range s {
		if 'A' <= c && c <= 'Z' {
			buffer.WriteRune(c + 32)
		} else {
			buffer.WriteRune(c)
		}
	}
	return buffer.String()
}

func isAlphanumeric(c rune) bool {
	return ('a' <= c && c <= 'z') || ('A' <= c && c <= 'Z') || ('0' <= c && c <= '9')
}
