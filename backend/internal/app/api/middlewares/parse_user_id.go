package middlewares

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

const UserIdContextKey = "userId"

func ParseUserIDMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		userId := c.Request().Header.Get("X-UserId")
		if userId == "" {
			return c.NoContent(http.StatusUnauthorized)
		}

		c.Set("userId", userId)

		return next(c)
	}
}
