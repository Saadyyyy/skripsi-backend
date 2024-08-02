package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

type YourCustomClaims struct {
	ID   int64 `json:"user_id"`
	Role int64 `json:"role"`
	jwt.RegisteredClaims
}

func JWTMiddleware() echo.MiddlewareFunc {
	return echo.MiddlewareFunc(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			tokenString := c.Request().Header.Get("Authorization")
			if tokenString == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, "No token found")
			}

			// Menghapus skema Bearer jika ada
			if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
				tokenString = tokenString[7:]
			}

			claims := &YourCustomClaims{}
			token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
				// Pastikan metode penandatanganan yang digunakan sesuai
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, errors.New("unexpected signing method")
				}
				return []byte(os.Getenv("JWT_SECRET")), nil
			})

			if err != nil || !token.Valid {
				return echo.NewHTTPError(http.StatusUnauthorized, "Invalid token")
			}

			// Simpan klaim ke dalam context
			c.Set("user", claims)
			return next(c)
		}
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

func AdminMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		claims, ok := c.Get("user").(*YourCustomClaims)
		if !ok || claims == nil {
			return echo.NewHTTPError(http.StatusUnauthorized, "No token found")
		}

		fmt.Println("role", claims.Role)

		if claims.Role == 1 { // Periksa apakah role adalah admin
			return echo.NewHTTPError(http.StatusForbidden, "Access denied: Only admins can access this resource")
		}

		return next(c)
	}
}

func SetTokenCookie(e echo.Context, token string) {
	cookie := new(http.Cookie)
	cookie.Name = "token"
	cookie.Value = token
	cookie.Path = "/"

	e.SetCookie(cookie)
}
