package auth

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/JubaerHossain/rootx/pkg/core/app"
	"github.com/golang-jwt/jwt/v5"
)

// Errors used in the function
var (
	ErrMissingJwtExpiration   = errors.New("missing configuration: JwtExpiration")
	ErrMissingJwtSecretKey    = errors.New("missing configuration: JwtSecretKey")
	ErrInvalidAccessTokenExp  = errors.New("invalid access token expiration value")
	ErrInvalidRefreshTokenExp = errors.New("invalid refresh token expiration value")
	ErrFailedToGenerateToken  = errors.New("failed to generate token")
)

// CreateTokens generates access and refresh tokens
func CreateTokens(userId uint, role string, app *app.App, refreshExp uint) (string, string, error) {

	jwtTime := app.Config.JwtExpiration
	if jwtTime == 0 {
		jwtTime = 24 * time.Hour
	}

	if refreshExp == 0 {
		refreshExp = 24 * 7
	}

	// Define token expiration times
	accessTokenExpirationTime := time.Now().Add(jwtTime).Unix()
	refreshTokenExpirationTime := time.Now().Add(time.Hour * time.Duration(refreshExp)).Unix()

	// refreshExpInt := int(refreshExp) // Ensure refreshExp is an integer

	// Define token expiration times
	// accessTokenExpirationTime := time.Now().Add(time.Minute * time.Duration(accessExp)).Unix()
	// refreshTokenExpirationTime := time.Now().Add(time.Hour * time.Duration(refreshExpInt)).Unix()

	secretKey := app.Config.JwtSecretKey
	if secretKey == "" {
		secretKey = "your-secret"
	}
	// Create access token claims
	accessClaims := jwt.MapClaims{
		"sub":  userId,
		"iat":  time.Now().Unix(),
		"exp":  accessTokenExpirationTime,
		"role": role,
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString([]byte(secretKey))
	if err != nil {
		return "", "", fmt.Errorf("%w: %v", ErrFailedToGenerateToken, err)
	}

	// Create refresh token claims
	refreshClaims := jwt.MapClaims{
		"sub": userId,
		"exp": refreshTokenExpirationTime,
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString([]byte(secretKey))
	if err != nil {
		return "", "", fmt.Errorf("%w: %v", ErrFailedToGenerateToken, err)
	}

	return accessTokenString, refreshTokenString, nil
}

func ValidateToken(tokenString string) (jwt.MapClaims, error) {
	secretKey := os.Getenv("JWT_SECRET_KEY")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	return claims, nil
}
