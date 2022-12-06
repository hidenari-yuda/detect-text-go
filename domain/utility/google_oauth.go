package utility

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
)

// GoogleOAuthのリダイレクトURL作成用Configを生成
func NewGoogleAuthConf() *oauth2.Config {

	err := godotenv.Load(".env.local")
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	credentialsJSON := os.Getenv("GOOGLE_OAUTH_CREDENTIALS_JSON")

	// 第2引数にサービス権限（カレンダーの全権限を設定）
	config, err := google.ConfigFromJSON([]byte(credentialsJSON), calendar.CalendarScope)
	if err != nil {
		log.Fatalf("config生成エラー: %v", err)
	}
	return config

}
