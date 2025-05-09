package middlewares

import (
	"TelegaFeed/internal/pkg/core/abstractions/infrastructure/repositories"
	"github.com/labstack/echo/v4"
	"net/http"
)

func UserExistsMiddleware(usersRepository abstractrepositories.UsersRepository) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			userId := c.Get(UserIdContextKey).(string)

			user, err := usersRepository.GetUserById(c.Request().Context(), userId)

			if err != nil {
				return c.NoContent(http.StatusInternalServerError)
			}

			if user == nil {
				return c.NoContent(http.StatusUnauthorized)
			}

			return next(c)
		}
	}
}
