package healper

import (
	"errors"
	"os"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CompareHash(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func DecodeJWTToken(tokenString string) (string, error) {
	// Parse token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Token verification key (harus sesuai dengan key yang digunakan untuk sign token)
		secretKey := os.Getenv("JWT_SECRET")
		if secretKey == "" {
			return nil, errors.New("secret key not found in environment variable")
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		return "", err
	}

	// Memeriksa apakah token valid
	if !token.Valid {
		return "", errors.New("token tidak valid")
	}

	// Memeriksa apakah token memiliki klaim "sub"
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("gagal mendapatkan claims dari token")
	}

	userID, ok := claims["sub"].(string)
	if !ok || userID == "" {
		return "", errors.New("tidak dapat mendapatkan user ID dari token")
	}

	return userID, nil
}
