// middleware/jwt_test.go
package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestAdminMiddleware(t *testing.T) {
	e := echo.New()

	// Apply JWT Middleware
	e.Use(JWTMiddleware())

	// Create a valid token with admin role
	claims := &YourCustomClaims{
		ID:   1,
		Role: 2, // Admin role
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 5)), // Set token expiry
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte("your-secret-key"))
	if err != nil {
		t.Fatalf("Failed to sign token: %v", err)
	}

	// Create a request with the token
	req := httptest.NewRequest(http.MethodGet, "/admin", nil)
	req.Header.Set("Authorization", "Bearer "+signedToken)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Apply Admin Middleware
	mw := AdminMiddleware(func(c echo.Context) error {
		return c.String(http.StatusOK, "Admin Access")
	})

	// Execute request
	err = mw(c)
	if err != nil {
		t.Fatalf("Middleware returned an error: %v", err)
	}

	// Check results
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "Admin Access", rec.Body.String())
}

// func TestUserMiddleware(t *testing.T) {
// 	e := echo.New()

// 	claims := &YourCustomClaims{
// 		ID:   1,
// 		Role: 1, // User role
// 		RegisteredClaims: jwt.RegisteredClaims{
// 			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 5)),
// 		},
// 	}
// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
// 	signedToken, err := token.SignedString([]byte("your-secret-key"))
// 	if err != nil {
// 		t.Fatalf("Failed to sign token: %v", err)
// 	}

// 	req := httptest.NewRequest(http.MethodGet, "/user", nil)
// 	req.Header.Set("Authorization", "Bearer "+signedToken)
// 	rec := httptest.NewRecorder()
// 	c := e.NewContext(req, rec)

// 	// Apply JWT Middleware
// 	e.Use(JWTMiddleware())

// 	// Apply User Middleware
// 	mw := UserMiddleware(func(c echo.Context) error {
// 		return c.String(http.StatusOK, "User Access")
// 	})

// 	// Execute request
// 	err = mw(c)
// 	if err != nil {
// 		t.Fatalf("Middleware returned an error: %v", err)
// 	}

// 	// Check results
// 	assert.NoError(t, err)
// 	assert.Equal(t, http.StatusOK, rec.Code)
// 	assert.Equal(t, "User Access", rec.Body.String())
// }
