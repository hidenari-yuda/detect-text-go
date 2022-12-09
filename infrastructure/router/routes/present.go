package routes

import (
	"net/http"

	"github.com/hidenari-yuda/paychan-server/domain/entity"
	"github.com/hidenari-yuda/paychan-server/infrastructure/database"
	"github.com/hidenari-yuda/paychan-server/infrastructure/di"
	"github.com/hidenari-yuda/paychan-server/usecase"
	"github.com/labstack/echo/v4"
)

type PresentRoutes struct{}

// Create(param *entity.Present) (presenter.Presenter, error)
// Update(param *entity.Present) (presenter.Presenter, error)
// GetById(id uint) (presenter.Presenter, error)
// GetByLineUserId(lineUserId string) (presenter.Presenter, error)
// GetAll() (presenter.Presenter, error)
// DeleteByExpired() (presenter.Presenter, error)

// 	Create(db *database.DB, firebase usecase.Firebase) func(c echo.Context) error
// 	Update(db *database.DB, firebase usecase.Firebase) func(c echo.Context) error
// 	GetById(db *database.DB, firebase usecase.Firebase) func(c echo.Context) error
// 	GetByLineUserId(db *database.DB, firebase usecase.Firebase) func(c echo.Context) error
// 	GetAll(db *database.DB, firebase usecase.Firebase) func(c echo.Context) error
// 	DeleteByExpired(db *database.DB, firebase usecase.Firebase) func(c echo.Context) error

func (r *PresentRoutes) Create(db *database.DB, firebase usecase.Firebase) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			param entity.Present
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

func (r *PresentRoutes) Update(db *database.DB, firebase usecase.Firebase) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			param entity.Present
		)

		err := bindAndValidate(c, &param)
		if err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

		tx, _ := db.Begin()
		h := di.InitializePresentHandler(tx, firebase)
		presenter, err := h.Update(&param)
		if err != nil {
			tx.Rollback()
			return c.JSON(http.StatusInternalServerError, err)
		}
		tx.Commit()
		c.JSON(http.StatusOK, presenter)
		return nil
	}
}

func (r *PresentRoutes) GetById(db *database.DB, firebase usecase.Firebase) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			idStr = c.Param("id")
		)

		id, err := stringToUint(idStr)
		if err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

		h := di.InitializePresentHandler(db, firebase)
		presenter, err := h.GetById(id)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}

		c.JSON(http.StatusOK, presenter)
		return nil
	}
}

func (r *PresentRoutes) GetByLineUserId(db *database.DB, firebase usecase.Firebase) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			lineUserId = c.Param("lineUserId")
		)

		if lineUserId == "" {
			return c.JSON(http.StatusBadRequest, "lineUserId is required")
		}

		h := di.InitializePresentHandler(db, firebase)
		presenter, err := h.GetByLineUserId(lineUserId)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}

		c.JSON(http.StatusOK, presenter)
		return nil
	}
}

func (r *PresentRoutes) GetAll(db *database.DB, firebase usecase.Firebase) func(c echo.Context) error {
	return func(c echo.Context) error {
		h := di.InitializePresentHandler(db, firebase)
		presenter, err := h.GetAll()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}

		c.JSON(http.StatusOK, presenter)
		return nil
	}
}

func (r *PresentRoutes) DeleteByExpired(db *database.DB, firebase usecase.Firebase) func(c echo.Context) error {
	return func(c echo.Context) error {
		h := di.InitializePresentHandler(db, firebase)
		presenter, err := h.DeleteByExpired()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}

		c.JSON(http.StatusOK, presenter)
		return nil
	}
}
