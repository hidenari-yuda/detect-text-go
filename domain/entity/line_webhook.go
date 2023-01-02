package entity

import (
	"net/http"

	"github.com/line/line-bot-sdk-go/v7/linebot"
)

type LineWebHook struct {
	Bot                *linebot.Client
	Events             []*linebot.Event
	ChannelSecret      string
	ChannelAccessToken string
	AdminUserId        string
	Request            *http.Request
}
