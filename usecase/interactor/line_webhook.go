package interactor

import (
	"fmt"

	"github.com/hidenari-yuda/paychan-server/domain/config"
	"github.com/hidenari-yuda/paychan-server/domain/entity"
	"github.com/line/line-bot-sdk-go/v7/linebot"
)

func (i *UserInteractorImpl) GetLineWebHook(param *entity.LineWebHook) (ok bool, err error) {

	for _, event := range param.Events {

		// 既存のユーザーかどうか確認
		user, err := i.userRepository.GetByLineUserId(event.Source.UserID)
		if err != nil {

			// 既存のユーザーではない場合は新規登録
			res, err := param.Bot.GetProfile(event.Source.UserID).Do()
			if err != nil {
				return ok, fmt.Errorf("lineプロフィールの取得エラー: %w", err)
			}
			err = i.userRepository.SignUp(&entity.SignUpParam{
				LineUserId:    res.UserID,
				LineName:      res.DisplayName,
				PictureUrl:    res.PictureURL,
				StatusMessage: res.StatusMessage,
				Language:      res.Language,
			})
			if err != nil {
				return ok, fmt.Errorf("新規登録エラー: %w", err)
			}
		}

		if event.Type == linebot.EventTypeMessage {
			switch message := event.Message.(type) {

			// 画像データの場合
			case *linebot.ImageMessage:
				fmt.Println("image message:", message)

				// 今日登録されたレシートのリストを取得する
				receiptPictures, err := i.receiptPictureRepository.GetListByToday(event.Source.UserID)
				if err != nil {
					return ok, fmt.Errorf("今日登録されたレシートのリストの取得エラー: %w", err)
				}

				// コンテンツ取得
				content, err := param.Bot.GetMessageContent(message.ID).Do()
				if err != nil {
					return ok, fmt.Errorf("getMessageContentでエラー: %v", err)
				}
				defer content.Content.Close()

				fmt.Println("content.Content", content.Content)

				// レシートか判定
				receiptPicture, presentPrice, err := CheckReceipt(content.Content, receiptPictures)
				if err != nil {

					if _, err = param.Bot.ReplyMessage(
						event.ReplyToken,
						linebot.NewTextMessage(fmt.Sprint("レシートが認識できませんでした。\nもう1度やり直してペイ")),
					).Do(); err != nil {
						return ok, fmt.Errorf("ImageMessageのReplyMessageでエラー: %v", err)
					}

					return ok, fmt.Errorf("checkReceiptでエラー: %v", err)
				}

				// 10件以上の場合は、10件以上ある旨を通知する
				if len(receiptPictures) > 10 {
					if _, err = param.Bot.ReplyMessage(
						event.ReplyToken,
						linebot.NewTextMessage(
							fmt.Sprintf(
								"本日上限である10件のレシートが登録されています。\n\n%sさん、また明日来てペイ！！",
								user.LineName,
							),
						)).Do(); err != nil {

						return ok, nil
					}
				}

				// プレゼントを取得
				// presentList, err := i.presentRepository.GetByPointAndService(presentPrice)
				// if err != nil || len(presentList) == 0 {
				cfg, err := config.New()
				botToAdmin, err := linebot.New(
					cfg.Line.ChannelSecret,
					cfg.Line.ChannelAccessToken,
				)

				if _, err := botToAdmin.PushMessage(
					cfg.Line.AdminUserId,
					linebot.NewTextMessage(
						fmt.Sprintf(
							"プレゼントが取得できませんでした。\n\n・対象ユーザー\n お名前:%sさん\n一言:%s\nレシートの金額:%d\n支払いサービス:%s。エラー内容:%v",
							user.LineName,
							user.StatusMessage,
							presentPrice.Point,
							convertPaymentServiceToStr(presentPrice.PaymentService),
							err,
						),
					),
				).Do(); err != nil {
					return ok, fmt.Errorf("プレゼント取得通知のリプライエラー: %w", err)
				}
				// 	return ok, fmt.Errorf("プレゼントの取得エラー: %w", err)
				// }

				if _, err = param.Bot.ReplyMessage(
					event.ReplyToken,
					linebot.NewTextMessage(fmt.Sprintf(
						"チェックが完了したペイ！\n\n    %v円分のプレゼントを%vで送るペイ！\n\n  今までのポイントを還元したい際は、メニューの「ポイント」ボタンから受け取ってレシ！\n\n",
						presentPrice.Point,
						convertPaymentServiceToStr(presentPrice.PaymentService),
						// presentList[0].Url,
					)),
				).Do(); err != nil {
					return ok, fmt.Errorf("ImageMessageのReplyMessageでエラー: %v", err)
				}

				// レシート情報をdbに登録
				err = i.receiptPictureRepository.Create(receiptPicture)
				if err != nil {
					return ok, fmt.Errorf("レシート情報の登録エラー: %w", err)
				}

				fmt.Println("image function is ok!!")

				/********** 画像メッセージ以外の場合 **********/
			// テキストメッセージの場合
			case *linebot.TextMessage:
				fmt.Println("text message:", message)

				switch message.Text {

				// case "アップ":
				// 	if _, err := param.Bot.ReplyMessage(
				// 		event.ReplyToken,
				// 		linebot.NewTextMessage("PayPay又はLINE Payの購入履歴のスクリーンショットを選んでだぺイ！"),
				// 	).Do(); err != nil {
				// 		return ok, fmt.Errorf("TextMessageのReplyMessageでエラー: %v", err)
				// 	}

				// 	time.Sleep(1 * time.Second)
				// 	linebot.NewCameraRollAction("アップロード").QuickReplyAction()

				// if _, err := param.Bot.Reply
				// 	event.ReplyToken,
				// 	linebot.NewCameraRollAction("アップロード"),
				// ).Do(); err != nil {
				// 	return ok, fmt.Errorf("TextMessageのReplyMessageでエラー: %v", err)
				// }

				case "キャンペーン":
					if _, err := param.Bot.ReplyMessage(
						event.ReplyToken,
						linebot.NewTemplateMessage(
							"開催中のキャンペーン",
							linebot.NewCarouselTemplate(
								linebot.NewCarouselColumn(
									"https://paypay-qr.s3-ap-northeast-1.amazonaws.com/qr/qr_2021_06_01_01_00_00_0000000000000000000000000000000000000000000000000000000000000000.png",
									"キャンペーン1",
									"キャンペーン",
									linebot.NewMessageAction("キャンペーン", "キャンペーン"),
								),
								linebot.NewCarouselColumn(
									"https://paypay-qr.s3-ap-northeast-1.amazonaws.com/qr/qr_2021_06_01_01_00_00_0000000000000000000000000000000000000000000000000000000000000000.png",
									"キャンペーン2",
									"キャンペーン",
									linebot.NewMessageAction("キャンペーン", "キャンペーン"),
								),
								linebot.NewCarouselColumn(
									"https://paypay-qr.s3-ap-northeast-1.amazonaws.com/qr/qr_2021_06_01_01_00_00_0000000000000000000000000000000000000000000000000000000000000000.png",
									"キャンペーン3",
									"キャンペーン",
									linebot.NewMessageAction("キャンペーン", "キャンペーン"),
								),
							),
						),
					).Do(); err != nil {
						return ok, fmt.Errorf("EventTypeMessageのReplyMessageでエラー: %v", err)
					}
					return ok, nil

				case "アンケート":
					if _, err := param.Bot.ReplyMessage(
						event.ReplyToken,
						linebot.NewTemplateMessage(

							"アンケート",
							linebot.NewCarouselTemplate(
								linebot.NewCarouselColumn(
									"",
									"アンケート",
									"アンケートを回答するペイ！",
									linebot.NewURIAction("アンケートを回答するペイ！", "line://app/1653824439-5jQXjz5A"),
								),
							),
						),
					).Do(); err != nil {
						return ok, fmt.Errorf("EventTypeMessageのReplyMessageでエラー: %v", err)
					}
					return ok, nil

				case "ポイント":
					if _, err := param.Bot.ReplyMessage(
						event.ReplyToken,
						linebot.NewTextMessage(fmt.Sprintf("保有ポイントは %v ポイントです。", user.Point)),
					).Do(); err != nil {
						return ok, fmt.Errorf("EventTypeMessageのReplyMessageでエラー: %v", err)
					}

					if _, err := param.Bot.ReplyMessage(
						event.ReplyToken,
						linebot.NewTemplateMessage(
							"下のボタンから還元する方法を選べるペイ！",
							linebot.NewCarouselTemplate(
								linebot.NewCarouselColumn(
									"",
									"PayPayポイントに還元",
									"還元",
									linebot.NewMessageAction("PayPayポイントに還元", "PayPayポイントに還元"),
								),
							),
						),
					).Do(); err != nil {
						return ok, fmt.Errorf("EventTypeMessageのReplyMessageでエラー: %v", err)
					}
					return ok, nil

				case "PayPayポイントに還元":
					presentList, err := i.presentRepository.GetByPointAndService(&entity.Present{
						Point:          user.Point,
						PaymentService: 0,
					})
					if err != nil || presentList == nil || len(presentList) == 0 {
						// エラー処理
						cfg, err := config.New()
						botToAdmin, err := linebot.New(
							cfg.Line.ChannelSecret,
							cfg.Line.ChannelAccessToken,
						)

						if _, eff := botToAdmin.PushMessage(
							cfg.Line.AdminUserId,
							linebot.NewTextMessage(
								fmt.Sprintf(
									"プレゼントが取得できませんでした。\n\n・対象ユーザー\n お名前:%sさん\n一言:%s\nレシートの金額:%d\n支払いサービス:%s。エラー内容:%v",
									user.LineName,
									user.StatusMessage,
									user.Point,
									"PayPay",
									err,
								),
							),
						).Do(); eff != nil {
							return ok, fmt.Errorf("プレゼント取得通知のリプライエラー: %w", err)
						}
					}

					if _, err := param.Bot.ReplyMessage(
						event.ReplyToken,
						linebot.NewTextMessage("PayPayポイントはこちらから確認できるペイ！\nhttps://paypay.ne.jp/point/"),
					).Do(); err != nil {
						return ok, fmt.Errorf("EventTypeMessageのReplyMessageでエラー: %v", err)
					}

					for _, present := range presentList {
						// ギフトをdbに保存する
						err = i.presentRepository.Update(&entity.Present{
							Id:     present.Id,
							UserId: user.Id,
							// ReceiptPictureId: receiptPicture.Id,
							Point:          present.Point,
							PaymentService: present.PaymentService,
							Url:            present.Url,
						})
					}
					return ok, nil

					// case "LINEPayポイントに還元":

					// 	presentList, err := i.presentRepository.GetByPointAndService(&entity.Present{
					// 		Point:          user.Point,
					// 		PaymentService: 1,
					// 	},
					// 	)
					// 	if err != nil || len(presentList) == 0 {
					// 		cfg, err := config.New()
					// 		botToAdmin, err := linebot.New(
					// 			cfg.Line.ChannelSecret,
					// 			cfg.Line.ChannelAccessToken,
					// 		)

					// 		if _, eff := botToAdmin.PushMessage(
					// 			cfg.Line.AdminUserId,
					// 			linebot.NewTextMessage(
					// 				fmt.Sprintf(
					// 					"プレゼントが取得できませんでした。\n\n・対象ユーザー\n お名前:%sさん\n一言:%s\nレシートの金額:%d\n支払いサービス:%s。エラー内容:%v",
					// 					user.LineName,
					// 					user.StatusMessage,
					// 					user.Point,
					// 					"PayPay",
					// 					err,
					// 				),
					// 			),
					// 		).Do(); eff != nil {
					// 			return ok, fmt.Errorf("プレゼント取得通知のリプライエラー: %w", err)
					// 		}
					// 	}

					// 	if _, err := param.Bot.ReplyMessage(
					// 		event.ReplyToken,
					// 		linebot.NewTextMessage("LINEポイントはこちらから確認できるペイ！\nhttps://point.line.me/"),
					// 	).Do(); err != nil {
					// 		return ok, fmt.Errorf("EventTypeMessageのReplyMessageでエラー: %v", err)
					// 	}

					// 	// ギフトをdbに保存する
					// 	err = i.presentRepository.Update(&entity.Present{
					// 		Id:     presentList[0].Id,
					// 		UserId: user.Id,
					// 		// ReceiptPictureId: receiptPicture.Id,
					// 		Point:          presentList[0].Point,
					// 		PaymentService: presentList[0].PaymentService,
					// 		Url:            presentList[0].Url,
					// 	})
					// 	return ok, nil

				}

				if _, err := param.Bot.ReplyMessage(
					event.ReplyToken,
					linebot.NewTextMessage("テキストメッセージありがペイ！\n"+"レシートの画像を送信してみてペイ！"),
				).Do(); err != nil {
					return ok, fmt.Errorf("EventTypeMessageのReplyMessageでエラー: %v", err)
				}

				// // DBにメッセージを保存する処理
				// err = i.lineMessageRepository.Create(&entity.LineMessage{
				// 	LineUserId: event.Source.UserID,
				// })

			// スタンプの場合
			case *linebot.StickerMessage:
				fmt.Println("stamp message:", message)

				if _, err := param.Bot.ReplyMessage(
					event.ReplyToken,
					linebot.NewTextMessage("スタンプありがペイ！\n"+"レシートの画像を送信してみてペイ！"),
				).Do(); err != nil {
					return ok, fmt.Errorf("EventTypeMessageのReplyMessageでエラー: %v", err)
				}

			// 動画データの場合
			case *linebot.VideoMessage:
				fmt.Println("message", message)

				if _, err := param.Bot.ReplyMessage(
					event.ReplyToken,
					linebot.NewTextMessage("動画ありがペイ！\n"+"レシートの画像を送信してみてペイ！"),
				).Do(); err != nil {
					return ok, fmt.Errorf("EventTypeMessageのReplyMessageでエラー: %v", err)
				}

			// 音声データの場合
			case *linebot.AudioMessage:
				fmt.Println("message", message)

				if _, err := param.Bot.ReplyMessage(
					event.ReplyToken,
					linebot.NewTextMessage("音声メッセージありがペイ！\n"+"レシートの画像を送信してみてペイ！"),
				).Do(); err != nil {
					return ok, fmt.Errorf("EventTypeMessageのReplyMessageでエラー: %v", err)
				}

			}
		}
	}

	ok = true

	return ok, nil

}

func convertPaymentServiceToStr(paymentService int) string {
	// switch paymentService {
	// case entity.PayPay:
	// 	return "PayPay"
	// case entity.LinePay:
	// 	return "LINE Pay"
	// case entity.MercariPay:
	// 	return "メルペイ"
	// case entity.Cash:
	// 	return "現金"
	// case entity.AmazonPay:
	// 	return "Amazon Pay"
	// case entity.RakutenPay:
	// 	return "楽天ペイ"
	// }
	return "PayPay"
}
