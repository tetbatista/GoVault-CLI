package auth

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("govault-local-secret-2024")

type Claims struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func sessionFilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".govault_session"), nil
}

func CreateSession(userID, username string) error {
	claims := Claims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(30 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(jwtSecret)
	if err != nil {
		return err
	}

	path, err := sessionFilePath()
	if err != nil {
		return err
	}

	return os.WriteFile(path, []byte(token), 0600)
}

func GetSession() (*Claims, error) {
	path, err := sessionFilePath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, errors.New("no active session — please login first")
	}

	claims := &Claims{}

	token, err := jwt.ParseWithClaims(string(data), claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return jwtSecret, nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("invalid or expired session — please login again")
	}

	return claims, nil
}

func DestroySession() error {
	path, err := sessionFilePath()
	if err != nil {
		return err
	}
	return os.Remove(path)
}

func RequireAuth(fn func(claims *Claims)) {
	claims, err := GetSession()
	if err != nil {
		fmt.Println("❌", err.Error())
		return
	}
	fn(claims)
}
