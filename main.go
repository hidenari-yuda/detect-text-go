package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/hidenari-yuda/paychan/domain/config"
	"github.com/hidenari-yuda/paychan/domain/entity"
	"github.com/hidenari-yuda/paychan/domain/utility"
	"github.com/hidenari-yuda/paychan/infrastructure/batch"
	"github.com/hidenari-yuda/paychan/infrastructure/database"
	"github.com/hidenari-yuda/paychan/infrastructure/driver"
	infrastructure "github.com/hidenari-yuda/paychan/infrastructure/router"

	"github.com/hidenari-yuda/paychan/usecase"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/slack-go/slack"
)

func init() {
	time.Local = utility.Tokyo

	goEnv := os.Getenv("GO_ENV")
	if goEnv == "" {
		goEnv = "local"
	}

	if err := godotenv.Load(fmt.Sprintf(".conf/%s.env", goEnv)); err != nil {
		panic("Failed to load .env file")
	}
	fmt.Println("load .env.local:", os.Getenv("GO_ENV"))
	fmt.Println("load .env.local:", os.Getenv("APP_ENV"))
}

func main() {
	cfg, err := config.New()
	if err != nil {
		panic(err)
	}

	switch cfg.App.Service {
	case "api":
		// 一旦 apiコンテナを立ち上げる時にマイグレーションする
		db := database.NewDB(cfg.DB, true)
		err := db.MigrateUp(".migrations")
		if err != nil {
			fmt.Println(err)
		}
		// cache := driver.NewRedisCacheImpl(cfg.Redis)
		if cfg.App.Env == "local" {
			firebase := driver.NewFirebaseImpl(cfg.Firebase)
			fmt.Println("getTestUserToken:", uuid.New().String())
			getTestUserToken(firebase, uuid.New().String())
		}
		r := infrastructure.NewRouter(cfg)

		// // エラーハンドラー（dev or prdのみSlack通知）
		if cfg.App.Env != "local" {
			r.Engine.HTTPErrorHandler = customHTTPErrorHandler
		}

		// ルーティング
		r.SetUp().Start()

	case "batch":
		batch.NewBatch(cfg).Start()
	}
}

func getTestUserToken(fb usecase.Firebase, uuid string) {
	customToken, _ := fb.GetCustomToken(uuid)
	idToken, err := fb.GetIDToken(customToken)
	if err != nil {
		panic(err)
	}
	fmt.Println("test token is :", idToken)
}

func customHTTPErrorHandler(err error, c echo.Context) {
	var (
		cfg, _        = config.New()
		code, message = entity.ErrorInfo(err)
		statusCode    = strconv.Itoa(code)
		path          = c.Path()
		method        = c.Request().Method
		errText       = err.Error()
	)

	fmt.Println(err)

	te := "*開発環境 Error*\n" +
		">>>status: " + message + "（" + statusCode + "）" + "\n" +
		"method: " + method + "\n" +
		"uri: " + path + "\n" +
		"error: `" + errText + "` \n"

	// アクセストークンを使用してクライアントを生成する
	// https://api.slack.com/apps からトークン取得
	// 参考: https://risaki-masa.com/how-to-get-api-token-in-slack/
	tkn := cfg.Slack.AccessToken
	chanelID := cfg.Slack.ChanelID
	s := slack.New(tkn)

	// MsgOptionText() の第二引数に true を設定すると特殊文字をエスケープする
	_, _, err = s.PostMessage(chanelID, slack.MsgOptionBlocks(
		&slack.SectionBlock{
			Type: slack.MBTSection,
			Text: &slack.TextBlockObject{
				Type: "mrkdwn",
				Text: te,
			},
		},
	))
	if err != nil {
		fmt.Println(err)
	}

	c.Logger().Error(err)
}
