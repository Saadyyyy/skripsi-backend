package middleware

import (
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func JWTMiddleware() echo.MiddlewareFunc {
	godotenv.Load()
	return middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey: []byte(os.Getenv("JWT_SECRET")),
	})
}

func CreateToken(id string, role int64) (string, error) {
	godotenv.Load()
	claims := jwt.MapClaims{}
	claims["id"] = id
	claims["role"] = role
	claims["exp"] = time.Now().Add(time.Hour * 5).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

func SetTokenCookie(c echo.Context, token string) {
	cookie := new(http.Cookie)
	cookie.Name = "token"
	cookie.Value = token
	cookie.Path = "/"

	c.SetCookie(cookie)
}

func ExtractToken(c echo.Context) (string, string, error) {
	user := c.Get("user").(*jwt.Token)
	if user.Valid {
		claims := user.Claims.(jwt.MapClaims)
		Id := claims["id"].(string)
		Role := claims["role"].(string)
		return Id, Role, nil
	}
	return "", "", errors.New("invalid token")
}

func CreateTokenVerifikasi(email string) (string, error) {
	godotenv.Load()
	claims := jwt.MapClaims{}
	claims["email"] = email
	claims["exp"] = time.Now().Add(time.Minute * 10).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

func ExtractTokenVerifikasi(c echo.Context) (string, error) {
	user := c.Get("user").(*jwt.Token)
	if user.Valid {
		claims := user.Claims.(jwt.MapClaims)
		email := claims["email"].(string)

		return email, nil
	}
	return "", errors.New("invalid token")
}
