package utils

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
)

type CustomDateFormat time.Time

// OTP character set
const otpChars = "0123456789"

// MarshalJSON implements the json.Marshaler interface
func (c CustomDateFormat) MarshalJSON() ([]byte, error) {
	t := time.Time(c)
	formatted := t.Format("2006-01-02") // Format the time as "YYYY-MM-DD"
	return json.Marshal(formatted)
}

func GenerateUniqueNumber(length int) (string, error) {
	otp := make([]byte, length)
	for i := range otp {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(otpChars))))
		if err != nil {
			return "", fmt.Errorf("failed to generate OTP: %v", err)
		}
		otp[i] = otpChars[num.Int64()]
	}
	return string(otp), nil
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

// GenerateOTP generates a random OTP of the given length using numeric characters.
func GenerateOTP(length int) (string, error) {
	otp := make([]byte, length)
	for i := range otp {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(otpChars))))
		if err != nil {
			return "", fmt.Errorf("failed to generate OTP: %v", err)
		}
		otp[i] = otpChars[num.Int64()]
	}
	return string(otp), nil
}

// GenerateOrderCode generates a random code of the given length using the provided character set.
func GenerateOrderCode(length int, charset string) (string, error) {
	code := make([]byte, length)
	for i := range code {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", fmt.Errorf("failed to generate code: %v", err)
		}
		code[i] = charset[num.Int64()]
	}
	return string(code), nil
}

// StringToIntEnv retrieves an environment variable as an integer, returning a default value if the conversion fails.
func StringToIntEnv(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}

// Key type should match the one used in the middleware.
type key int

const (
	ClaimsKey key = iota
)

// GetAuthUser retrieves the claims from the context.

func GetAuthUser(ctx context.Context) (map[string]interface{}, error) {
	claims := ctx.Value(ClaimsKey)
	if claims == nil {
		return nil, fmt.Errorf("claims not found in context")
	}
	user, err := claims.(jwt.MapClaims)
	if !err {
		return nil, fmt.Errorf("failed to get user from claims")
	}
	return user, nil
}

type Bool bool

// UnmarshalBool is a reusable function that unmarshals both boolean and string representations
func UnmarshalBool(data []byte) (Bool, error) {
	var b Bool

	// Try to unmarshal into a bool
	var boolValue bool
	if err := json.Unmarshal(data, &boolValue); err == nil {
		b = Bool(boolValue)
		return b, nil
	}

	// Try to unmarshal into a string, then convert to bool
	var strValue string
	if err := json.Unmarshal(data, &strValue); err == nil {
		boolValue, err := strconv.ParseBool(strValue)
		if err != nil {
			return b, fmt.Errorf("invalid boolean value: '%s'. Expected 'true' or 'false'", strValue)
		}
		b = Bool(boolValue)
		return b, nil
	}

	return b, fmt.Errorf("invalid boolean value")
}

// UnmarshalJSON customizes the unmarshaling for the Bool type
func (b *Bool) UnmarshalJSON(data []byte) error {
	parsedBool, err := UnmarshalBool(data)
	if err != nil {
		return err
	}
	*b = parsedBool
	return nil
}
