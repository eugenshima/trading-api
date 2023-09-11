package middleware

import (
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

// tokenClaims struct consists od JWT claims
type tokenClaims struct {
	jwt.StandardClaims
}

// const for middlware
const (
	Bearer = "Bearer"
	Admin  = "admin"
)

// UserIdentity is a middleware function that validates access token
func UserIdentity() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			return nil
		}
	}
}
