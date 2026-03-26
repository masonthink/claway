package middleware

import (
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

// RequireAuth is a middleware that validates JWT from the Authorization header
// and sets user_id in the echo context.
func RequireAuth(jwtSecret string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "missing authorization header"})
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid authorization format"})
			}

			tokenStr := parts[1]

			token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
				if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, jwt.ErrSignatureInvalid
				}
				return []byte(jwtSecret), nil
			},
				jwt.WithIssuer("claway"),
				jwt.WithAudience("claway-api"),
				jwt.WithValidMethods([]string{"HS256"}),
			)
			if err != nil || !token.Valid {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid or expired token"})
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid token claims"})
			}

			// Extract user_id from claims (stored as float64 by default in JSON)
			userIDFloat, ok := claims["user_id"].(float64)
			if !ok || userIDFloat <= 0 {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid user_id in token"})
			}

			c.Set("user_id", int64(userIDFloat))
			return next(c)
		}
	}
}
