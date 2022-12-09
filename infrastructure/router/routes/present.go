package routes

import (
	"net/http"

	"github.com/hidenari-yuda/paychan/domain/entity"
	"github.com/hidenari-yuda/paychan/infrastructure/database"
	"github.com/hidenari-yuda/paychan/infrastructure/di"
	"github.com/hidenari-yuda/paychan/usecase"
	"github.com/labstack/echo/v4"
)

type PresentRoutes struct{}

// 	Create(db *database.DB, firebase usecase.Firebase) func(c echo.Context) error

func (r *PresentRoutes) Create(db *database.DB, firebase usecase.Firebase) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			param entity.PresentParam
		)

		err := bindAndValidate(c, &param)
		if err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

		tx, _ := db.Begin()
		h := di.InitializePresentHandler(tx, firebase)
		presenter, err := h.Create(&param)
		if err != nil {
			tx.Rollback()
			return c.JSON(http.StatusInternalServerError, err)
		}
		tx.Commit()
		c.JSON(http.StatusOK, presenter)
		return nil
	}
}
