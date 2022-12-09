package routes

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"

	"github.com/hidenari-yuda/paychan-server/domain/config"
	"github.com/hidenari-yuda/paychan-server/domain/entity"
	"github.com/hidenari-yuda/paychan-server/infrastructure/database"
	"github.com/hidenari-yuda/paychan-server/infrastructure/di"
	"github.com/hidenari-yuda/paychan-server/usecase"
	"github.com/labstack/echo/v4"
	"github.com/line/line-bot-sdk-go/v7/linebot"
)

type UserRoutes struct{}

// 	SignUp(db *database.DB, firebase usecase.Firebase) func(c echo.Context) error
// 	SignIn(db *database.DB, firebase usecase.Firebase) func(c echo.Context) error
// 	GetByFirebaseToken(db *database.DB, firebase usecase.Firebase) func(c echo.Context) error
// 	GetLineWebHook(db *database.DB, firebase usecase.Firebase) func(c echo.Context) error

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

func (r *UserRoutes) GetByLineUserId(db *database.DB, firebase usecase.Firebase) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			lineUserId = c.Param("lineUserId")
		)

		h := di.InitializeUserHandler(db, firebase)
		presenter, err := h.GetByLineUserId(lineUserId)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}

		c.JSON(http.StatusOK, presenter)
		return nil
	}
}

/************************ line関連 ************************/
func (r *UserRoutes) GetLineWebHook(db *database.DB, firebase usecase.Firebase) func(c echo.Context) error {
	return func(c echo.Context) error {
		cfg, err := config.New()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}

		bot, err := linebot.New(cfg.Line.ChannelSecret, cfg.Line.ChannelAccessToken)
		if err != nil {
			fmt.Println(err)
			return err
		}

		req := c.Request() // リクエストの取得
		defer req.Body.Close()

		// Webhookイベントオブジェクトの処理
		body, err := io.ReadAll(req.Body)
		if err != nil {
			fmt.Println(err)
			c.Response().WriteHeader(500)
			return fmt.Errorf("リクエストボディの読み込みに失敗しました: %w", err)
		}

		// 署名の検証
		// 参考: https://developers.line.biz/ja/reference/messaging-api/#signature-validation
		// x-line-signatureヘッダーから署名を取得
		decoded, err := base64.StdEncoding.DecodeString(req.Header.Get("x-line-signature"))
		if err != nil {
			fmt.Println(err)
			c.Response().WriteHeader(500)
			return fmt.Errorf("x-line-signatureをエンコードできません: %s", err.Error())
		}

		// lineシークレットをキーにしてHMAC-SHA256にしたものと署名が正しいか検証
		hash := hmac.New(sha256.New, []byte(cfg.Line.ChannelSecret))
		_, err = hash.Write(body)
		if err != nil {
			c.Response().WriteHeader(500)
			return fmt.Errorf("hash.Write(body)に失敗しました: %s", err.Error())
		}

		// Compare decoded signature and `hash.Sum(nil)` by using `hmac.Equal`
		if !hmac.Equal(decoded, hash.Sum(nil)) {
			c.Response().WriteHeader(400)
			return fmt.Errorf("署名キーが正しくありません: %s", err.Error())
		}

		req.Body = io.NopCloser(bytes.NewBuffer(body)) // リクエストボディを再利用するためにリセット

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

		// 画像メッセージの場合、先にメッセージを返す
		// events[0].Message.(*linebot.ImageMessage).Message != nil
		if events[0].Type == linebot.EventTypeMessage {
			switch events[0].Message.(type) {
			case *linebot.ImageMessage:
				_, err = bot.ReplyMessage(events[0].ReplyToken, linebot.NewTextMessage("確認中...")).Do()
			}
		}

		param := &entity.LineWebHook{
			Bot:                bot,
			Events:             events,
			ChannelSecret:      cfg.Line.ChannelSecret,
			ChannelAccessToken: cfg.Line.ChannelAccessToken,
			Request:            req,
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
