package routes

import (
	"fmt"

	"github.com/hidenari-yuda/paychan-server/domain/config"
	"github.com/hidenari-yuda/paychan-server/domain/entity"
	"github.com/hidenari-yuda/paychan-server/infrastructure/database"
	"github.com/hidenari-yuda/paychan-server/infrastructure/di"
	"github.com/hidenari-yuda/paychan-server/interfaces/presenter"
	"github.com/labstack/echo/v4"
)

type AdminAuthorizeParam struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// AdminAuthorize is
//
func AdminAuthorize(db *database.DB, appConfig config.App) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			param = new(AdminAuthorizeParam)
		)
		if err := bindAndValidate(c, param); err != nil {
			wrapped := fmt.Errorf("%s:%w", err.Error(), entity.ErrRequestError)
			renderJSON(c, presenter.NewErrorJsonPresenter(wrapped))
			return err
		}

		h := di.InitializeAdminHandler(db, appConfig)
		p, err := h.Authorize(param.Username, param.Password)
		if err != nil {
			renderJSON(c, presenter.NewErrorJsonPresenter(err))
		}

		renderJSON(c, p)
		return nil
	}
}
