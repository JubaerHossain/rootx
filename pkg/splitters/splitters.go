package splitters

import (
	"errors"
	"strings"
)

// SplitArgs splits a string based on the given separator while handling quotes and optional comment support.
func SplitArgs(input, separator string, keepQuotes bool) ([]string, error) {
	return splitArgs(input, separator, keepQuotes, "")
}

// SplitSQL splits SQL statements, supporting custom comment markers and preserving quotes.
func SplitSQL(input, separator string, keepQuotes bool) ([]string, error) {
	return splitArgs(input, separator, keepQuotes, "--")
}

// splitArgs handles string splitting, optionally preserving quotes and skipping comments.
func splitArgs(input, separator string, keepQuotes bool, commentSign string) ([]string, error) {
	// Ensure the separator is not empty
	if separator == "" {
		return nil, errors.New("separator cannot be empty")
	}

	singleQuoteOpen := false
	doubleQuoteOpen := false
	commentSignOpen := false
	commentSignIndex := 0
	separatorIndex := 0

	var tokenBuffer strings.Builder
	var result []string

	// Iterate over each character in the input string
	for _, char := range input {
		inputChar := string(char)

		// Handle newlines to reset comment state
		if inputChar == "\n" {
			commentSignIndex = 0
			commentSignOpen = false
			tokenBuffer.WriteString(inputChar)
			continue
		}

		// Skip over characters in comment mode
		if commentSignOpen {
			tokenBuffer.WriteString(inputChar)
			continue
		}

		// Handle single and double quotes
		if inputChar == "'" && !doubleQuoteOpen && !commentSignOpen {
			if keepQuotes {
				tokenBuffer.WriteString(inputChar)
			}
			singleQuoteOpen = !singleQuoteOpen
			continue
		} else if inputChar == `"` && !singleQuoteOpen && !commentSignOpen {
			if keepQuotes {
				tokenBuffer.WriteString(inputChar)
			}
			doubleQuoteOpen = !doubleQuoteOpen
			continue
		}

		// Handle comment markers
		if !singleQuoteOpen && !doubleQuoteOpen && !commentSignOpen {
			if commentSign != "" && inputChar == string(commentSign[commentSignIndex]) {
				// Comment marker detected
				commentSignIndex++
				if commentSignIndex == len(commentSign) {
					commentSignOpen = true
					commentSignIndex = 0
				}
				tokenBuffer.WriteString(inputChar)
				continue
			} else {
				commentSignIndex = 0
			}
		}

		// Handle separator logic outside quotes and comments
		if !singleQuoteOpen && !doubleQuoteOpen && !commentSignOpen {
			if inputChar == string(separator[separatorIndex]) {
				separatorIndex++
				if separatorIndex == len(separator) {
					// Separator fully matched
					result = append(result, tokenBuffer.String())
					tokenBuffer.Reset()
					separatorIndex = 0
					continue
				}
			} else {
				// Reset separator matching state
				separatorIndex = 0
			}
		}
		tokenBuffer.WriteString(inputChar)
	}

	// Add remaining buffer content to the result
	if tokenBuffer.Len() > 0 {
		result = append(result, tokenBuffer.String())
	}

	return result, nil
}
