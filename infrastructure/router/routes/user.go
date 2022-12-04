package routes

import (
	"fmt"
	"net/http"

	"github.com/hidenari-yuda/umerun-resume/domain/entity"
	"github.com/hidenari-yuda/umerun-resume/infrastructure/database"
	"github.com/hidenari-yuda/umerun-resume/infrastructure/di"
	"github.com/hidenari-yuda/umerun-resume/usecase"
	"github.com/labstack/echo/v4"
	"github.com/line/line-bot-sdk-go/v7/linebot"
)

type UserRouteFunc interface {
	SignUp(db *database.DB, firebase usecase.Firebase) func(c echo.Context) error
	SignIn(db *database.DB, firebase usecase.Firebase) func(c echo.Context) error
	GetByFirebaseToken(db *database.DB, firebase usecase.Firebase) func(c echo.Context) error

	//line
	GetLineWebHook(db *database.DB, firebase usecase.Firebase) func(c echo.Context) error

	// resume
	UploadResume(db *database.DB, firebase usecase.Firebase) func(c echo.Context) error
}

type UserRoutes struct {
	UserRouteFunc
}

func (r *UserRoutes) SignUp(db *database.DB, firebase usecase.Firebase) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			param entity.SignUpParam
		)

		err := bindAndValidate(c, &param)
		if err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

		tx, _ := db.Begin()
		h := di.InitializeUserHandler(tx, firebase)
		presenter, err := h.SignUp(&param)
		if err != nil {
			tx.Rollback()
			return c.JSON(http.StatusInternalServerError, err)
		}
		tx.Commit()
		c.JSON(http.StatusOK, presenter)
		return nil
	}
}

func (r *UserRoutes) SignIn(db *database.DB, firebase usecase.Firebase) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			param entity.SignInParam
		)

		err := bindAndValidate(c, &param)
		if err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

		h := di.InitializeUserHandler(db, firebase)
		presenter, err := h.SignIn(&param)
		if err != nil {
			err = fmt.Errorf("サインインエラー: %s:%w", err.Error(), entity.ErrRequestError)
			fmt.Println(err)
			return c.JSON(http.StatusInternalServerError, err)
		}

		renderJSON(c, presenter)
		return nil
	}
}

func (r *UserRoutes) GetByFirebaseToken(db *database.DB, firebase usecase.Firebase) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			firebaseToken = GetFirebaseToken(c)
		)

		h := di.InitializeUserHandler(db, firebase)
		presenter, err := h.GetByFirebaseToken(firebaseToken)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}

		c.JSON(http.StatusOK, presenter)
		return nil
	}
}

// func (r *UserRoutes) GetByLineUserId(db *database.DB, firebase usecase.Firebase) func(c echo.Context) error {
// 	return func(c echo.Context) error {
// 		var (
// 			lineUserId = c.Param("lineUserId")
// 		)

// 		h := di.InitializeUserHandler(db, firebase)
// 		presenter, err := h.GetByFirebaseToken(firebaseToken)
// 		if err != nil {
// 			return c.JSON(http.StatusInternalServerError, err)
// 		}

// 		c.JSON(http.StatusOK, presenter)
// 		return nil
// 	}
// }

/************************ line関連 ************************/
func (r *UserRoutes) GetLineWebHook(db *database.DB, firebase usecase.Firebase) func(c echo.Context) error {
	return func(c echo.Context) error {
		req := c.Request()

		bot, err := linebot.New("LINE_SECRET", "LINE_ACCESS_TOKEN")
		if err != nil {
			fmt.Println(err)
			return err
		}

		// Webhookイベントオブジェクトの取得
		// 参考: https://developers.line.biz/ja/reference/messaging-api/#webhook-event-objects
		events, err := bot.ParseRequest(req)
		if err != nil {
			if err == linebot.ErrInvalidSignature {
				fmt.Println("c.Response().WriteHeader(400)")
				c.Response().WriteHeader(400)
				return err
			} else {
				fmt.Println("c.Response().WriteHeader(500)")
				c.Response().WriteHeader(500)
				return err
			}
		}

		param := &entity.LineWebHookParam{
			Bot:    bot,
			Events: events,
		}

		h := di.InitializeUserHandler(db, firebase)
		presenter, err := h.GetLineWebHook(param)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}

		c.JSON(http.StatusOK, presenter)
		return nil
	}
}
