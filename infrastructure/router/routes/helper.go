package routes

import (
	"fmt"
	"strconv"

	"github.com/hidenari-yuda/paychan-server/interfaces/presenter"
	"github.com/labstack/echo/v4"
	"gopkg.in/go-playground/validator.v9"
)

func GetFirebaseToken(c echo.Context) string {
	token := c.Request().Header.Get("FirebaseAuthorization")
	return token
}

func bindAndValidate(c echo.Context, obj interface{}) (err error) {
	err = c.Bind(obj)
	if err != nil {
		return
	}
	return validator.New().Struct(obj)
}

func renderJSON(c echo.Context, p presenter.Presenter) {
	c.JSON(p.StatusCode(), p.Data())
}

// ファイル送信用
func renderFile(c echo.Context, filePath string) error {
	err := c.File(filePath)
	if err != nil {
		return err
	}
	return nil
}

// stringToint
func stringToint(s string) (int, error) {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}

	if i < 1 {
		return 0, fmt.Errorf("invalid value: %d", i)
	}

	return int(i), nil
}
