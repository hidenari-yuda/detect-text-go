package entity

import "github.com/line/line-bot-sdk-go/v7/linebot"

type LineWebHookParam struct {
	Bot    *linebot.Client
	Events []*linebot.Event
}
