package middlewares

import (
	"fmt"
	"net/http"

	"gocommerce/core"
	"gocommerce/helper"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func SetClientJWTmiddlewares(g *echo.Group, role string) {
	env := core.App.Config

	g.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningMethod: "HS256",
		SigningKey:    []byte(env.JWT_SECRET),
	}))

	g.Use(ValidateJWTLogin)
}

func ValidateJWTLogin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user := c.Get("user")
		token := user.(*jwt.Token)

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			if claims["id"] != 0 {
				return next(c)
			} else {
				return helper.Response(http.StatusUnauthorized, "", fmt.Sprintf("%s", "Please Sign In \n Woops! Gonna sign in first\n Only a click away and you can continue to enjoy"))

			}
		}
		return helper.Response(http.StatusUnauthorized, "", fmt.Sprintf("%s", "Please Sign In \n Woops! Gonna sign in first\n Only a click away and you can continue to enjoy"))
	}
}
