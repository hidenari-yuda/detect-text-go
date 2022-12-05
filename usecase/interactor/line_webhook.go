package interactor

import (
	"fmt"

	"github.com/hidenari-yuda/umerun-resume/domain/entity"
	"github.com/line/line-bot-sdk-go/v7/linebot"
)

type GetLineWebHookInput struct {
	Param *entity.LineWebHookParam
}

type GetLineWebHookOutput struct {
	Ok bool
}

func (i *UserInteractorImpl) GetLineWebHook(input GetLineWebHookInput) (output GetLineWebHookOutput, err error) {

	for _, event := range input.Param.Events {
		// 既存のユーザーかどうか確認
		user, err := i.userRepository.GetByLineUserId(event.Source.UserID)
		if err != nil {
			res, err := input.Param.Bot.GetProfile(event.Source.UserID).Do()
			if err != nil {
				return output, err
			}
			// 既存のユーザーではない場合は新規登録
			err = i.userRepository.SignUp(&entity.SignUpParam{
				LineUserId:    res.UserID,
				LineName:      res.DisplayName,
				PictureUrl:    res.PictureURL,
				StatusMessage: res.StatusMessage,
				Language:      res.Language,
			})
			if err != nil {
				return output, err
			}
		}

		// 今日登録されたレシートのリストを取得する
		receiptPictures, err := i.receiptPictureRepository.GetListByToday(event.Source.UserID)
		if err != nil {
			return output, err
		}

		if event.Type == linebot.EventTypeMessage {
			switch message := event.Message.(type) {

			// 画像データの場合
			case *linebot.ImageMessage:
				fmt.Println("image message:", message)

				// コンテンツ取得
				bot, err := linebot.New("LINE_SECRET", "LINE_ACCESS_TOKEN")
				if err != nil {
					return output, fmt.Errorf("botクライアントの作成でエラー: %v", err)
				}
				content, err := bot.GetMessageContent(message.ID).Do()
				if err != nil {
					return output, fmt.Errorf("getMessageContentでエラー: %v", err)
				}

				defer content.Content.Close()
				fmt.Println("content.Content", content.Content)

				// レシートか判定
				info, err := checkReceipt(content.Content, receiptPictures)
				if err != nil {

					_, err = bot.ReplyMessage(
						event.ReplyToken,
						linebot.NewTextMessage(fmt.Sprint("レシートが認識できませんでした。\nもう1度やり直してトン")),
					).Do()

					if err != nil {
						return output, fmt.Errorf("ImageMessageのReplyMessageでエラー: %v", err)
					}

					return output, fmt.Errorf("checkReceiptでエラー: %v", err)
				}

				// 10件以上の場合は、10件以上ある旨を通知する
				if len(receiptPictures) > 10 {
					if _, err = input.Param.Bot.ReplyMessage(
						event.ReplyToken,
						linebot.NewTextMessage(
							fmt.Sprintf(
								"本日上限である10件のレシートが登録されています。\n\n%sさん、また明日来てトン！！",
								user.LineName,
							),
						)).Do(); err != nil {
						return output, nil
					}
				}

				fmt.Println("info", info)

				_, err = bot.ReplyMessage(
					event.ReplyToken,
					linebot.NewTextMessage(fmt.Sprintf(
						"レシートを受け取りました\n\n    %v円をプレゼントします！\n\n  まだ未登録の方は、下記のリンクから登録してください！\n\n    https://www.google.com",
						"値段",
					)),
				).Do()
				if err != nil {
					fmt.Println(err)
					fmt.Println("ImageMessageのReplyMessageでエラー")
				}

				// ギフトをdbに保存する

				// message.OriginalContentURL,
				// message.PreviewImageURL,
			// テキストメッセージの場合
			case *linebot.TextMessage:
				fmt.Println("text message:", message)

				if _, err := input.Param.Bot.ReplyMessage(
					event.ReplyToken,
					linebot.NewTextMessage("テキストメッセージありがトン！\n"+"レシートの画像を送信してみてトン！"),
				).Do(); err != nil {
					return output, fmt.Errorf("EventTypeMessageのReplyMessageでエラー: %v", err)
				}

				// DBにメッセージを保存する処理
				err = i.LineMessageRepository.Create(&entity.LineMessage{
					LineUserId: event.Source.UserID,
				})

			// スタンプの場合
			case *linebot.StickerMessage:
				fmt.Println("stamp message:", message)

				if _, err := input.Param.Bot.ReplyMessage(
					event.ReplyToken,
					linebot.NewTextMessage("スタンプありがトン！\n"+"レシートの画像を送信してみてトン！"),
				).Do(); err != nil {
					return output, fmt.Errorf("EventTypeMessageのReplyMessageでエラー: %v", err)
				}

			// 動画データの場合
			case *linebot.VideoMessage:
				fmt.Println("message", message)

				if _, err := input.Param.Bot.ReplyMessage(
					event.ReplyToken,
					linebot.NewTextMessage("動画ありがトン！\n"+"レシートの画像を送信してみてトン！"),
				).Do(); err != nil {
					return output, fmt.Errorf("EventTypeMessageのReplyMessageでエラー: %v", err)
				}

			// 音声データの場合
			case *linebot.AudioMessage:
				fmt.Println("message", message)

				if _, err := input.Param.Bot.ReplyMessage(
					event.ReplyToken,
					linebot.NewTextMessage("音声メッセージありがトン！\n"+"レシートの画像を送信してみてトン！"),
				).Do(); err != nil {
					return output, fmt.Errorf("EventTypeMessageのReplyMessageでエラー: %v", err)
				}

			}
		}
	}

	output.Ok = true

	return output, nil

}

// コンテンツ取得
// bot, err := linebot.New(<channel secret>, <channel token>)
// if err != nil {
// 	...
// }
// content, err := bot.GetMessageContent(<messageID>).Do()
// if err != nil {
// 	...
// }
// defer content.Content.Close()

// プロフィール情報取得
// bot, err := linebot.New(<channel secret>, <channel token>)
// if err != nil {
// 	...
// }
// res, err := bot.GetProfile(<userId>).Do();
// if err != nil {
// 	...
// }
// println(res.DisplayName)
// println(res.PictureURL)
// println(res.StatusMessage)
