package interactor

import (
	"fmt"

	"github.com/hidenari-yuda/umerun-resume/domain/config"
	"github.com/hidenari-yuda/umerun-resume/domain/entity"
	"github.com/line/line-bot-sdk-go/v7/linebot"
)

type GetLineWebHookInput struct {
	Param *entity.LineWebHook
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
				return output, fmt.Errorf("lineプロフィールの取得エラー: %w", err)
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
				return output, fmt.Errorf("新規登録エラー: %w", err)
			}
		}

		// 今日登録されたレシートのリストを取得する
		receiptPictures, err := i.receiptPictureRepository.GetListByToday(event.Source.UserID)
		if err != nil {
			return output, fmt.Errorf("今日登録されたレシートのリストの取得エラー: %w", err)
		}

		if event.Type == linebot.EventTypeMessage {
			switch message := event.Message.(type) {

			// 画像データの場合
			case *linebot.ImageMessage:
				fmt.Println("image message:", message)

				// コンテンツ取得
				content, err := input.Param.Bot.GetMessageContent(message.ID).Do()
				if err != nil {
					return output, fmt.Errorf("getMessageContentでエラー: %v", err)
				}
				defer content.Content.Close()

				fmt.Println("content.Content", content.Content)

				// レシートか判定

				receiptPicture, presenet, err := checkReceipt(content.Content, receiptPictures)
				if err != nil {

					if _, err = input.Param.Bot.ReplyMessage(
						event.ReplyToken,
						linebot.NewTextMessage(fmt.Sprint("レシートが認識できませんでした。\nもう1度やり直してトン")),
					).Do(); err != nil {
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

				// プレゼントを取得
				present, err := i.presentRepository.GetByPrice(presenet.Price)
				if err != nil || present == nil {
					cfg, err := config.New()
					botToAdmin, err := linebot.New(
						cfg.Line.ChannelSecret,
						cfg.Line.ChannelAccessToken,
					)
					if _, eff := botToAdmin.PushMessage(
						cfg.Line.AdminUserId,
						linebot.NewTextMessage(
							fmt.Sprintf(
								"プレゼントが取得できませんでした。\n\n・対象ユーザー\n お名前:%sさん\n一言:%s\nレシートの金額:%d\n支払いサービス:%s。",
								user.LineName,
								user.StatusMessage,
								presenet.Price,
								convertPaymentServiceToStr(presenet.PaymentService),
							),
						),
					).Do(); eff != nil {
						return output, fmt.Errorf("プレゼント取得エラー: %w", err)
					}
					return output, fmt.Errorf("プレゼントの取得エラー: %w", err)
				}

				if _, err = input.Param.Bot.ReplyMessage(
					event.ReplyToken,
					linebot.NewTextMessage(fmt.Sprintf(
						"レシートを受け取りました\n\n    %v円をプレゼントします！\n\n  まだ未登録の方は、下記のリンクから登録してください！\n\n    https://www.google.com",
						"値段",
					)),
				).Do(); err != nil {
					return output, fmt.Errorf("ImageMessageのReplyMessageでエラー: %v", err)
				}

				// レシート情報をdbに登録
				err = i.receiptPictureRepository.Create(receiptPicture)
				if err != nil {
					return output, fmt.Errorf("レシート情報の登録エラー: %w", err)
				}

				// ギフトをdbに保存する
				err = i.presentRepository.Create(&entity.Present{})

				/********** 画像メッセージ以外の場合 **********/
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
				err = i.lineMessageRepository.Create(&entity.LineMessage{
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

func convertPaymentServiceToStr(paymentService entity.PaymentService) string {
	// switch paymentService {
	// case entity.paypay:
	// 	return "楽天"
	// case entity.PaymentServicePayPay:
	// 	return "PayPay"
	// case entity.PaymentServiceLinePay:
	// 	return "LINE Pay"
	// case entity.PaymentServiceYahoo:
	// 	return "Yahoo!ショッピング"
	// case entity.PaymentServiceAmazon:
	// 	return "Amazon"
	// case entity.PaymentServiceOther:
	// 	return "その他"
	// default:
	// 	return "不明"
	// }
	return "不明"
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
