package router

import (
	"github.com/hidenari-yuda/paychan-server/domain/entity"
	"github.com/hidenari-yuda/paychan-server/infrastructure/router/routes"
	"github.com/labstack/echo/v4"
)

func validateFirebaseToken(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if routes.GetFirebaseToken(c) == "" {
			return entity.ErrFirebaseEmptyToken
		}
		if err := next(c); err != nil {
			return err
		}
		return nil
	}
}
