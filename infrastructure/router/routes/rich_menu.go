package routes

import (
	"fmt"
	"net/http"

	"github.com/hidenari-yuda/paychan-server/domain/config"
	"github.com/hidenari-yuda/paychan-server/infrastructure/database"
	"github.com/hidenari-yuda/paychan-server/usecase"
	"github.com/labstack/echo/v4"
	"github.com/line/line-bot-sdk-go/v7/linebot"
)

type RichMenuRoutes struct{}

// Create(db *database.DB, firebase usecase.Firebase) func(c echo.Context) error
// CreateAlias(db *database.DB, firebase usecase.Firebase) func(c echo.Context) error
// GetAll(db *database.DB, firebase usecase.Firebase) func(c echo.Context) error

func (r *RichMenuRoutes) Create(db *database.DB, firebase usecase.Firebase) func(c echo.Context) error {
	return func(c echo.Context) error {

		richMemu := linebot.RichMenu{
			Size:        linebot.RichMenuSize{Width: 2500, Height: 1686},
			Selected:    true,
			Name:        "richmenu",
			ChatBarText: "Tap to open",

			Areas: []linebot.AreaDetail{
				{
					Bounds: linebot.RichMenuBounds{X: 0, Y: 0, Width: 2500, Height: 1686},
					Action: linebot.RichMenuAction{
						Type: linebot.RichMenuActionTypePostback,
						Data: "action=campaign",
					},
				},
				{
					Bounds: linebot.RichMenuBounds{X: 0, Y: 0, Width: 2500, Height: 1686},
					Action: linebot.RichMenuAction{
						Type: linebot.RichMenuActionTypePostback,
						Data: "action=campaign",
					},
				},
			},
		}

		cfg, err := config.New()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}

		bot, err := linebot.New(cfg.Line.ChannelSecret, cfg.Line.ChannelAccessToken)
		if err != nil {
			fmt.Println(err)
			return err
		}

		res, err := bot.CreateRichMenu(richMemu).Do()
		if err != nil {
			fmt.Println(err)
			return err
		}
		println(res.RichMenuID)

		c.JSON(http.StatusOK, res)
		return nil
	}
}

func (r *RichMenuRoutes) CreateAlias(db *database.DB, firebase usecase.Firebase) func(c echo.Context) error {
	return func(c echo.Context) error {
		// 		curl -v -X POST https://api.line.me/v2/bot/richmenu/alias \
		// -H 'Authorization: Bearer {channel access token}' \
		// -H 'Content-Type: application/json' \
		// -d \
		// '{
		//     "richMenuAliasId": "richmenu-alias-a",
		//     "richMenuId": "richmenu-19682466851b21e2d7c0ed482ee0930f"
		// }'

		var (
			richMenuAliasId = c.Param("richMenuAliasId") // richmenu-alias-a
			richMenuId      = c.Param("richMenuId")      // richmenu-19682466851b21e2d7c0ed482ee0930f
		)

		if richMenuId == "" || richMenuAliasId == "" {
			return c.JSON(http.StatusBadRequest, "richMenuId and richMenuAliasId are required")
		}

		cfg, err := config.New()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}

		bot, err := linebot.New(cfg.Line.ChannelSecret, cfg.Line.ChannelAccessToken)
		if err != nil {
			fmt.Println(err)
			return err
		}

		res, err := bot.CreateRichMenuAlias(richMenuAliasId, richMenuId).Do()
		if err != nil {
			fmt.Println(err)
			return err
		}

		// h := di.InitializeRichMenuHandler(db, firebase)
		// presenter, err := h.CreateAlias(alias)
		// if err != nil {
		// 	return c.JSON(http.StatusInternalServerError, err)
		// }

		c.JSON(http.StatusOK, res)
		return nil
	}
}

// リッチメニューの配列を取得する
func (r *RichMenuRoutes) GetAll(db *database.DB, firebase usecase.Firebase) func(c echo.Context) error {
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

		res, err := bot.GetRichMenuList().Do()
		if err != nil {
			fmt.Println(err)
			return err
		}
		for _, richMenu := range res {
			println(richMenu.RichMenuID)
		}

		// レスポンスの例
		// {
		// 	"richmenus": [
		// 		{
		// 			"richMenuId": "{richMenuId}",
		// 			"size": {
		// 				"width": 2500,
		// 				"height": 1686
		// 			},
		// 			"selected": false,
		// 			"areas": [
		// 				{
		// 					"bounds": {
		// 						"x": 0,
		// 						"y": 0,
		// 						"width": 2500,
		// 						"height": 1686
		// 					},
		// 					"action": {
		// 						"type": "postback",
		// 						"data": "action=buy&itemid=123"
		// 					}
		// 				}
		// 			]
		// 		}
		// 	]
		// }

		// h := di.InitializeRichMenuHandler(db, firebase)
		// presenter, err := h.GetAll()
		// if err != nil {
		// 	return c.JSON(http.StatusInternalServerError, err)
		// }

		c.JSON(http.StatusOK, res)
		return nil
	}
}

// カメラロールアクション
// クイックリプライボタンにのみ設定できるアクションです。このアクションが関連づけられたボタンがタップされると、LINEのカメラロール画面が開きます。

// type
// String 必須
// cameraRoll

// label
// String 必須
// アクションのラベル
// 最大文字数：20

// // カメラロールアクションをかえす
// // ActionTypeCameraRoll
// func (r *UserRoutes) GetCameraRollAction() *linebot.CameraRollAction {
// 	return linebot.NewCameraRollAction("カメラロール")
// 	// return linebot.NewCameraRollAction("カメラロール")

// }

// // クイックリプライの作成
// // クイックリプライは、メッセージを送信したユーザーに対して、メッセージを送信する際に表示するボタンを設定できます。
// func (r *UserRoutes) GetQuickReply() *linebot.QuickReplyButton {
// 	// クイックリプライのアクションを作成
// 	actions := []linebot.QuickReplyButton{
// 		ImageURL: "https://example.com/sushi.png",
// 		Action:   r.GetCameraRollAction(),
// 		// r.GetCameraRollAction(),
// 		// r.GetLocationAction(),
// 		// r.GetPostbackAction(),
// 	}

// 	// クイックリプライの作成
// 	linebot.NewQuickReplyItems(actions...)
// }
