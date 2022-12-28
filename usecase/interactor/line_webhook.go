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

		// メッセージのタイプで処理を分ける
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

				// for _, receiptPicture := range receiptPictures {
				// 	fmt.Println("取得したレシート画像:", receiptPicture)
				// }

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

				// コンテンツ取得
				content, err := param.Bot.GetMessageContent(message.ID).Do()
				if err != nil {
					return ok, fmt.Errorf("getMessageContentでエラー: %v", err)
				}
				defer content.Content.Close()

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

				// プレゼントを取得
				// presentList, err := i.presentRepository.GetByPointAndService(presentPrice)
				// if err != nil || len(presentList) == 0 {
				// cfg, err := config.New()
				// botToAdmin, err := linebot.New(
				// 	cfg.Line.ChannelSecret,
				// 	cfg.Line.ChannelAccessToken,
				// )

				// if _, err := botToAdmin.PushMessage(
				// 	cfg.Line.AdminUserId,
				// 	linebot.NewTextMessage(
				// 		fmt.Sprintf(
				// 			"プレゼントが取得できませんでした。\n\n・対象ユーザー\n お名前:%sさん\n一言:%s\nレシートの金額:%d\n支払いサービス:%s。エラー内容:%v",
				// 			user.LineName,
				// 			user.StatusMessage,
				// 			presentPrice.Point,
				// 			convertPaymentServiceToStr(presentPrice.PaymentService),
				// 			err,
				// 		),
				// 	),
				// ).Do(); err != nil {
				// 	return ok, fmt.Errorf("プレゼント取得通知のリプライエラー: %w", err)
				// }
				// 	return ok, fmt.Errorf("プレゼントの取得エラー: %w", err)
				// }

				if user.QuestionProgress > len(entity.QuestionMessageTitle)-1 {
					if _, err = param.Bot.ReplyMessage(
						event.ReplyToken,
						linebot.NewTextMessage(
							fmt.Sprintf(
								"チェックが完了したペイ！\n\n    %v円分のポイントをプレゼントするペイ！\n\nこれまで保有しているポイントは、%vポイントだペイ！\n\n今までのポイントを還元したい際は、メニューの「ポイント」ボタンからPayPayポイントとして受け取れるペイ！",
								presentPrice.Point,
								user.Point+presentPrice.Point,
								// convertPaymentServiceToStr(presentPrice.PaymentService),
								// presentList[0].Url,
							),
						),
					).Do(); err != nil {
						return ok, fmt.Errorf("ImageMessageのReplyMessageでエラー: %v", err)
					}

					// 質問が終わっている場合は、未回答の質問を取得する
					// if user.Age == 99 {
					// 	user.QuestionProgress = 0
					// } else if user.Gender == 99 {
					// 	user.QuestionProgress = 1
					// } else if user.Marriage == 99 {
					// 	user.QuestionProgress = 2
					// }
				}

				questionMessageSelectionlength := len(entity.QuestionMessageSelection[user.QuestionProgress])
				if questionMessageSelectionlength == 2 {

					if _, err = param.Bot.ReplyMessage(
						event.ReplyToken,
						linebot.NewTextMessage(
							fmt.Sprintf(
								"チェックが完了したペイ！\n\n    %v円分のポイントをプレゼントするペイ！\n\nこれまで保有しているポイントは、%vポイントだペイ！\n\n今までのポイントを還元したい際は、メニューの「ポイント」ボタンからPayPayポイントとして受け取れるペイ！",
								presentPrice.Point,
								user.Point+presentPrice.Point,
								// convertPaymentServiceToStr(presentPrice.PaymentService),
								// presentList[0].Url,
							),
						),
						linebot.NewTemplateMessage(
							"アンケート",
							linebot.NewButtonsTemplate(
								"",
								"",
								entity.QuestionMessageTitle[user.QuestionProgress],
								linebot.NewMessageAction(entity.QuestionMessageSelection[user.QuestionProgress][0], entity.QuestionMessageTitle[user.QuestionProgress]+":"+entity.QuestionMessageSelection[user.QuestionProgress][0]),
								linebot.NewMessageAction(entity.QuestionMessageSelection[user.QuestionProgress][1], entity.QuestionMessageTitle[user.QuestionProgress]+":"+entity.QuestionMessageSelection[user.QuestionProgress][1]),
							),
						),
					).Do(); err != nil {
						return ok, fmt.Errorf("ImageMessageのReplyMessageでエラー: %v", err)
					}
				} else if questionMessageSelectionlength == 3 {
					if _, err = param.Bot.ReplyMessage(
						event.ReplyToken,
						linebot.NewTextMessage(
							fmt.Sprintf(
								"チェックが完了したペイ！\n\n    %v円分のポイントをプレゼントするペイ！\n\nこれまで保有しているポイントは、%vポイントだペイ！\n\n今までのポイントを還元したい際は、メニューの「ポイント」ボタンからPayPayポイントとして受け取れるペイ！",
								presentPrice.Point,
								user.Point+presentPrice.Point,
								// convertPaymentServiceToStr(presentPrice.PaymentService),
								// presentList[0].Url,
							),
						),
						linebot.NewTemplateMessage(
							"アンケート",
							linebot.NewButtonsTemplate(
								"",
								"",
								entity.QuestionMessageTitle[user.QuestionProgress],
								linebot.NewMessageAction(entity.QuestionMessageSelection[user.QuestionProgress][0], entity.QuestionMessageTitle[user.QuestionProgress]+":"+entity.QuestionMessageSelection[user.QuestionProgress][0]),
								linebot.NewMessageAction(entity.QuestionMessageSelection[user.QuestionProgress][1], entity.QuestionMessageTitle[user.QuestionProgress]+":"+entity.QuestionMessageSelection[user.QuestionProgress][1]),
								linebot.NewMessageAction(entity.QuestionMessageSelection[user.QuestionProgress][2], entity.QuestionMessageTitle[user.QuestionProgress]+":"+entity.QuestionMessageSelection[user.QuestionProgress][2]),
							),
						),
					).Do(); err != nil {
						return ok, fmt.Errorf("ImageMessageのReplyMessageでエラー: %v", err)
					}
				} else if questionMessageSelectionlength == 4 {
					if _, err = param.Bot.ReplyMessage(
						event.ReplyToken,
						linebot.NewTextMessage(
							fmt.Sprintf(
								"チェックが完了したペイ！\n\n    %v円分のポイントをプレゼントするペイ！\n\nこれまで保有しているポイントは、%vポイントだペイ！\n\n今までのポイントを還元したい際は、メニューの「ポイント」ボタンからPayPayポイントとして受け取れるペイ！",
								presentPrice.Point,
								user.Point+presentPrice.Point,
								// convertPaymentServiceToStr(presentPrice.PaymentService),
								// presentList[0].Url,
							),
						),
						linebot.NewTemplateMessage(
							"アンケート",
							linebot.NewButtonsTemplate(
								"",
								"",
								entity.QuestionMessageTitle[user.QuestionProgress],
								linebot.NewMessageAction(entity.QuestionMessageSelection[user.QuestionProgress][0], entity.QuestionMessageTitle[user.QuestionProgress]+":"+entity.QuestionMessageSelection[user.QuestionProgress][0]),
								linebot.NewMessageAction(entity.QuestionMessageSelection[user.QuestionProgress][1], entity.QuestionMessageTitle[user.QuestionProgress]+":"+entity.QuestionMessageSelection[user.QuestionProgress][1]),
								linebot.NewMessageAction(entity.QuestionMessageSelection[user.QuestionProgress][2], entity.QuestionMessageTitle[user.QuestionProgress]+":"+entity.QuestionMessageSelection[user.QuestionProgress][2]),
								linebot.NewMessageAction(entity.QuestionMessageSelection[user.QuestionProgress][3], entity.QuestionMessageTitle[user.QuestionProgress]+":"+entity.QuestionMessageSelection[user.QuestionProgress][3]),
							),
						),
					).Do(); err != nil {
						return ok, fmt.Errorf("ImageMessageのReplyMessageでエラー: %v", err)
					}

				} else if questionMessageSelectionlength == 5 {
					if _, err = param.Bot.ReplyMessage(
						event.ReplyToken,
						linebot.NewTextMessage(
							fmt.Sprintf(
								"チェックが完了したペイ！\n\n    %v円分のポイントをプレゼントするペイ！\n\nこれまで保有しているポイントは、%vポイントだペイ！\n\n今までのポイントを還元したい際は、メニューの「ポイント」ボタンからPayPayポイントとして受け取れるペイ！",
								presentPrice.Point,
								user.Point+presentPrice.Point,
								// convertPaymentServiceToStr(presentPrice.PaymentService),
								// presentList[0].Url,
							),
						),
						linebot.NewTemplateMessage(
							"アンケート",
							linebot.NewCarouselTemplate(
								linebot.NewCarouselColumn(
									"",
									"",
									entity.QuestionMessageTitle[user.QuestionProgress],
									linebot.NewMessageAction(entity.QuestionMessageSelection[user.QuestionProgress][0], entity.QuestionMessageTitle[user.QuestionProgress]+":"+entity.QuestionMessageSelection[user.QuestionProgress][0]),
									linebot.NewMessageAction(entity.QuestionMessageSelection[user.QuestionProgress][1], entity.QuestionMessageTitle[user.QuestionProgress]+":"+entity.QuestionMessageSelection[user.QuestionProgress][1]),
									linebot.NewMessageAction(entity.QuestionMessageSelection[user.QuestionProgress][2], entity.QuestionMessageTitle[user.QuestionProgress]+":"+entity.QuestionMessageSelection[user.QuestionProgress][2]),
								),
								linebot.NewCarouselColumn(
									"",
									"",
									fmt.Sprint(entity.QuestionMessageTitle[user.QuestionProgress], "2"),
									linebot.NewMessageAction(entity.QuestionMessageSelection[user.QuestionProgress][3], entity.QuestionMessageTitle[user.QuestionProgress]+":"+entity.QuestionMessageSelection[user.QuestionProgress][3]),
									linebot.NewMessageAction(entity.QuestionMessageSelection[user.QuestionProgress][4], entity.QuestionMessageTitle[user.QuestionProgress]+":"+entity.QuestionMessageSelection[user.QuestionProgress][4]),
								),
							),
						),
					).Do(); err != nil {
						return ok, fmt.Errorf("ImageMessageのReplyMessageでエラー: %v", err)
					}
				} else if questionMessageSelectionlength == 6 {
					if _, err = param.Bot.ReplyMessage(
						event.ReplyToken,
						linebot.NewTextMessage(
							fmt.Sprintf(
								"チェックが完了したペイ！\n\n    %v円分のポイントをプレゼントするペイ！\n\nこれまで保有しているポイントは、%vポイントだペイ！\n\n今までのポイントを還元したい際は、メニューの「ポイント」ボタンからPayPayポイントとして受け取れるペイ！",
								presentPrice.Point,
								user.Point+presentPrice.Point,
								// convertPaymentServiceToStr(presentPrice.PaymentService),
								// presentList[0].Url,
							),
						),
						linebot.NewTemplateMessage(
							"アンケート",
							linebot.NewCarouselTemplate(
								linebot.NewCarouselColumn(
									"",
									"",
									entity.QuestionMessageTitle[user.QuestionProgress],
									linebot.NewMessageAction(entity.QuestionMessageSelection[user.QuestionProgress][0], entity.QuestionMessageTitle[user.QuestionProgress]+":"+entity.QuestionMessageSelection[user.QuestionProgress][0]),
									linebot.NewMessageAction(entity.QuestionMessageSelection[user.QuestionProgress][1], entity.QuestionMessageTitle[user.QuestionProgress]+":"+entity.QuestionMessageSelection[user.QuestionProgress][1]),
									linebot.NewMessageAction(entity.QuestionMessageSelection[user.QuestionProgress][2], entity.QuestionMessageTitle[user.QuestionProgress]+":"+entity.QuestionMessageSelection[user.QuestionProgress][2]),
								),
								linebot.NewCarouselColumn(
									"",
									"",
									fmt.Sprint(entity.QuestionMessageTitle[user.QuestionProgress], "2"),
									linebot.NewMessageAction(entity.QuestionMessageSelection[user.QuestionProgress][3], entity.QuestionMessageTitle[user.QuestionProgress]+":"+entity.QuestionMessageSelection[user.QuestionProgress][3]),
									linebot.NewMessageAction(entity.QuestionMessageSelection[user.QuestionProgress][4], entity.QuestionMessageTitle[user.QuestionProgress]+":"+entity.QuestionMessageSelection[user.QuestionProgress][4]),
									linebot.NewMessageAction(entity.QuestionMessageSelection[user.QuestionProgress][5], entity.QuestionMessageTitle[user.QuestionProgress]+":"+entity.QuestionMessageSelection[user.QuestionProgress][5]),
								),
							),
						),
					).Do(); err != nil {
						return ok, fmt.Errorf("ImageMessageのReplyMessageでエラー: %v", err)
					}
				} else if questionMessageSelectionlength == 7 {
					if _, err = param.Bot.ReplyMessage(
						event.ReplyToken,
						linebot.NewTextMessage(
							fmt.Sprintf(
								"チェックが完了したペイ！\n\n    %v円分のポイントをプレゼントするペイ！\n\nこれまで保有しているポイントは、%vポイントだペイ！\n\n今までのポイントを還元したい際は、メニューの「ポイント」ボタンからPayPayポイントとして受け取れるペイ！",
								presentPrice.Point,
								user.Point+presentPrice.Point,
								// convertPaymentServiceToStr(presentPrice.PaymentService),
								// presentList[0].Url,
							),
						),
						linebot.NewTemplateMessage(
							"アンケート",
							linebot.NewButtonsTemplate(
								"",
								"",
								entity.QuestionMessageTitle[user.QuestionProgress],
								linebot.NewMessageAction(entity.QuestionMessageSelection[user.QuestionProgress][0], entity.QuestionMessageTitle[user.QuestionProgress]+":"+entity.QuestionMessageSelection[user.QuestionProgress][0]),
								linebot.NewMessageAction(entity.QuestionMessageSelection[user.QuestionProgress][1], entity.QuestionMessageTitle[user.QuestionProgress]+":"+entity.QuestionMessageSelection[user.QuestionProgress][1]),
								linebot.NewMessageAction(entity.QuestionMessageSelection[user.QuestionProgress][2], entity.QuestionMessageTitle[user.QuestionProgress]+":"+entity.QuestionMessageSelection[user.QuestionProgress][2]),
								linebot.NewMessageAction(entity.QuestionMessageSelection[user.QuestionProgress][3], entity.QuestionMessageTitle[user.QuestionProgress]+":"+entity.QuestionMessageSelection[user.QuestionProgress][3]),
							),
						),
						linebot.NewTemplateMessage(
							"アンケート",
							linebot.NewButtonsTemplate(
								"",
								"",
								entity.QuestionMessageTitle[user.QuestionProgress],
								linebot.NewMessageAction(entity.QuestionMessageSelection[user.QuestionProgress][4], entity.QuestionMessageTitle[user.QuestionProgress]+":"+entity.QuestionMessageSelection[user.QuestionProgress][4]),
								linebot.NewMessageAction(entity.QuestionMessageSelection[user.QuestionProgress][5], entity.QuestionMessageTitle[user.QuestionProgress]+":"+entity.QuestionMessageSelection[user.QuestionProgress][5]),
								linebot.NewMessageAction(entity.QuestionMessageSelection[user.QuestionProgress][6], entity.QuestionMessageTitle[user.QuestionProgress]+":"+entity.QuestionMessageSelection[user.QuestionProgress][6]),
							),
						),
					).Do(); err != nil {
						return ok, fmt.Errorf("ImageMessageのReplyMessageでエラー: %v", err)
					}
				} else if questionMessageSelectionlength == 8 {
					if _, err = param.Bot.ReplyMessage(
						event.ReplyToken,
						linebot.NewTextMessage(
							fmt.Sprintf(
								"チェックが完了したペイ！\n\n    %v円分のポイントをプレゼントするペイ！\n\nこれまで保有しているポイントは、%vポイントだペイ！\n\n今までのポイントを還元したい際は、メニューの「ポイント」ボタンからPayPayポイントとして受け取れるペイ！",
								presentPrice.Point,
								user.Point+presentPrice.Point,
								// convertPaymentServiceToStr(presentPrice.PaymentService),
								// presentList[0].Url,
							),
						),
						linebot.NewTemplateMessage(
							"アンケート",
							linebot.NewButtonsTemplate(
								"",
								"",
								entity.QuestionMessageTitle[user.QuestionProgress],
								linebot.NewMessageAction(entity.QuestionMessageSelection[user.QuestionProgress][0], entity.QuestionMessageTitle[user.QuestionProgress]+":"+entity.QuestionMessageSelection[user.QuestionProgress][0]),
								linebot.NewMessageAction(entity.QuestionMessageSelection[user.QuestionProgress][1], entity.QuestionMessageTitle[user.QuestionProgress]+":"+entity.QuestionMessageSelection[user.QuestionProgress][1]),
								linebot.NewMessageAction(entity.QuestionMessageSelection[user.QuestionProgress][2], entity.QuestionMessageTitle[user.QuestionProgress]+":"+entity.QuestionMessageSelection[user.QuestionProgress][2]),
								linebot.NewMessageAction(entity.QuestionMessageSelection[user.QuestionProgress][3], entity.QuestionMessageTitle[user.QuestionProgress]+":"+entity.QuestionMessageSelection[user.QuestionProgress][3]),
							),
						),
						linebot.NewTemplateMessage(
							"アンケート",
							linebot.NewButtonsTemplate(
								"",
								"",
								entity.QuestionMessageTitle[user.QuestionProgress],
								linebot.NewMessageAction(entity.QuestionMessageSelection[user.QuestionProgress][4], entity.QuestionMessageTitle[user.QuestionProgress]+":"+entity.QuestionMessageSelection[user.QuestionProgress][4]),
								linebot.NewMessageAction(entity.QuestionMessageSelection[user.QuestionProgress][5], entity.QuestionMessageTitle[user.QuestionProgress]+":"+entity.QuestionMessageSelection[user.QuestionProgress][5]),
								linebot.NewMessageAction(entity.QuestionMessageSelection[user.QuestionProgress][6], entity.QuestionMessageTitle[user.QuestionProgress]+":"+entity.QuestionMessageSelection[user.QuestionProgress][6]),
								linebot.NewMessageAction(entity.QuestionMessageSelection[user.QuestionProgress][7], entity.QuestionMessageTitle[user.QuestionProgress]+":"+entity.QuestionMessageSelection[user.QuestionProgress][7]),
							),
						),
					).Do(); err != nil {
						return ok, fmt.Errorf("ImageMessageのReplyMessageでエラー: %v", err)
					}
				}

				// レシート情報をdbに登録
				receiptPicture.UserId = user.Id
				receiptPicture.LineUserId = user.LineUserId

				err = i.receiptPictureRepository.Create(receiptPicture)
				if err != nil {
					return ok, fmt.Errorf("レシート情報の登録エラー: %w", err)
				}

				// ユーザーのポイントを更新
				user.Point += presentPrice.Point

				// ユーザー情報をdbに更新
				err = i.userRepository.Update(user)

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
									"",
									"準備中",
									"実装までしばらく待ってペイ！",
									linebot.NewMessageAction("準備中", "準備中"),
								),
								// linebot.NewCarouselColumn(
								// 	"",
								// 	"準備中",
								// 	"実装までしばらく待ってくださいペイ！",
								// 	// linebot.NewMessageAction("キャンペーン", "キャンペーン"),
								// ),
								// linebot.NewCarouselColumn(
								// 	"",
								// 	"アンケート",
								// 	"簡単なアンケートに答えてポイントゲット！",
								// 	linebot.NewMessageAction("プロフィールアンケート", "プロフィールアンケート"),
								// ),
							),
						),
					).Do(); err != nil {
						return ok, fmt.Errorf("EventTypeMessageのReplyMessageでエラー: %v", err)
					}
					return ok, nil

				// case "プロフィールアンケート":
				// 	if _, err := param.Bot.ReplyMessage(
				// 		event.ReplyToken,
				// 		linebot.NewTemplateMessage(
				// 			"プロフィールアンケート",
				// 			linebot.NewButtonsTemplate(
				// 				"",
				// 				"アンケート",
				// 				"アンケートに答えてポイントゲット！",
				// 				linebot.NewMessageAction("購入アンケート", "購入アンケート"),
				// 				linebot.NewMessageAction("購入アンケート", "購入アンケート"),
				// 			),
				// 		),
				// 	).Do(); err != nil {
				// 		return ok, fmt.Errorf("EventTypeMessageのReplyMessageでエラー: %v", err)
				// 	}
				// 	return ok, nil

				// case "アンケート":
				// 	if _, err := param.Bot.ReplyMessage(
				// 		event.ReplyToken,
				// 		linebot.NewTemplateMessage(
				// 			"プロフィールアンケート",
				// 			linebot.NewButtonsTemplate(
				// 				"",
				// 				"アンケート",
				// 				"アンケートに答えてポイントゲット！",
				// 				linebot.NewMessageAction("プロフィールアンケート", "プロフィールアンケート"),
				// 				linebot.NewMessageAction("購入アンケート", "購入アンケート"),
				// 			),
				// 		),
				// 	).Do(); err != nil {
				// 		return ok, fmt.Errorf("EventTypeMessageのReplyMessageでエラー: %v", err)
				// 	}
				// 	return ok, nil

				case "ポイント":
					if _, err := param.Bot.ReplyMessage(
						event.ReplyToken,
						linebot.NewTextMessage(fmt.Sprintf("保有ポイントは %v ポイントだペイ！", user.Point)),
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

					// 年齢　    0:10代 1:20代 2:30代 3:40代 4:50代 5:60代以上
				case entity.QuestionMessageTitle[0] + ":" + entity.QuestionMessageSelection[0][0]:
					if _, err := param.Bot.ReplyMessage(
						event.ReplyToken,
						linebot.NewTextMessage("ご回答ありがペイ！\n"+"またレシートの画像を送信してみてペイ！"),
					).Do(); err != nil {
						return ok, fmt.Errorf("EventTypeMessageのReplyMessageでエラー: %v", err)
					}

					user.QuestionProgress = 1
					user.Age = 0

					// DBにメッセージを保存する処理
					err = i.userRepository.Update(user)
					return ok, nil

				case entity.QuestionMessageTitle[0] + ":" + entity.QuestionMessageSelection[0][1]:
					if _, err := param.Bot.ReplyMessage(
						event.ReplyToken,
						linebot.NewTextMessage("ご回答ありがペイ！\n"+"またレシートの画像を送信してみてペイ！"),
					).Do(); err != nil {
						return ok, fmt.Errorf("EventTypeMessageのReplyMessageでエラー: %v", err)
					}

					user.QuestionProgress = 1
					user.Age = 1

					// DBにメッセージを保存する処理
					err = i.userRepository.Update(user)
					return ok, nil

				case entity.QuestionMessageTitle[0] + ":" + entity.QuestionMessageSelection[0][2]:
					if _, err := param.Bot.ReplyMessage(
						event.ReplyToken,
						linebot.NewTextMessage("ご回答ありがペイ！\n"+"またレシートの画像を送信してみてペイ！"),
					).Do(); err != nil {
						return ok, fmt.Errorf("EventTypeMessageのReplyMessageでエラー: %v", err)
					}

					// DBにメッセージを保存する処理
					user.QuestionProgress = 1
					user.Age = 2

					// DBにメッセージを保存する処理
					err = i.userRepository.Update(user)

					return ok, nil

				case entity.QuestionMessageTitle[0] + ":" + entity.QuestionMessageSelection[0][3]:
					if _, err := param.Bot.ReplyMessage(
						event.ReplyToken,
						linebot.NewTextMessage("ご回答ありがペイ！\n"+"またレシートの画像を送信してみてペイ！"),
					).Do(); err != nil {
						return ok, fmt.Errorf("EventTypeMessageのReplyMessageでエラー: %v", err)
					}

					user.QuestionProgress = 1
					user.Age = 3

					// DBにメッセージを保存する処理
					err = i.userRepository.Update(user)
					return ok, nil

				case entity.QuestionMessageTitle[0] + ":" + entity.QuestionMessageSelection[0][4]:
					if _, err := param.Bot.ReplyMessage(
						event.ReplyToken,
						linebot.NewTextMessage("ご回答ありがペイ！\n"+"またレシートの画像を送信してみてペイ！"),
					).Do(); err != nil {
						return ok, fmt.Errorf("EventTypeMessageのReplyMessageでエラー: %v", err)
					}

					user.QuestionProgress = 1
					user.Age = 4

					// DBにメッセージを保存する処理
					err = i.userRepository.Update(user)

					return ok, nil
				case entity.QuestionMessageTitle[0] + ":" + entity.QuestionMessageSelection[0][5]:
					if _, err := param.Bot.ReplyMessage(
						event.ReplyToken,
						linebot.NewTextMessage("ご回答ありがペイ！\n"+"またレシートの画像を送信してみてペイ！"),
					).Do(); err != nil {
						return ok, fmt.Errorf("EventTypeMessageのReplyMessageでエラー: %v", err)
					}

					user.QuestionProgress = 1
					user.Age = 5

					// DBにメッセージを保存する処理
					err = i.userRepository.Update(user)

					return ok, nil

					// 性別	0：男性　1：女性　2：その他
				case entity.QuestionMessageTitle[1] + ":" + entity.QuestionMessageSelection[1][0]:
					if _, err := param.Bot.ReplyMessage(
						event.ReplyToken,
						linebot.NewTextMessage("ご回答ありがペイ！\n"+"またレシートの画像を送信してみてペイ！"),
					).Do(); err != nil {
						return ok, fmt.Errorf("EventTypeMessageのReplyMessageでエラー: %v", err)
					}

					user.QuestionProgress = 2
					user.Gender = 0

					// DBにメッセージを保存する処理
					err = i.userRepository.Update(user)

					return ok, nil

				case entity.QuestionMessageTitle[1] + ":" + entity.QuestionMessageSelection[1][1]:

					if _, err := param.Bot.ReplyMessage(
						event.ReplyToken,
						linebot.NewTextMessage("ご回答ありがペイ！\n"+"またレシートの画像を送信してみてペイ！"),
					).Do(); err != nil {
						return ok, fmt.Errorf("EventTypeMessageのReplyMessageでエラー: %v", err)
					}

					user.QuestionProgress = 2
					user.Gender = 1

					// DBにメッセージを保存する処理
					err = i.userRepository.Update(user)
					return ok, nil

				case entity.QuestionMessageTitle[1] + ":" + entity.QuestionMessageSelection[1][2]:
					if _, err := param.Bot.ReplyMessage(
						event.ReplyToken,
						linebot.NewTextMessage("ご回答ありがペイ！\n"+"またレシートの画像を送信してみてペイ！"),
					).Do(); err != nil {
						return ok, fmt.Errorf("EventTypeMessageのReplyMessageでエラー: %v", err)
					}

					user.QuestionProgress = 2
					user.Gender = 2

					// DBにメッセージを保存する処理
					err = i.userRepository.Update(user)
					return ok, nil

					// 0: 未婚　1: 既婚
				case entity.QuestionMessageTitle[2] + ":" + entity.QuestionMessageSelection[2][0]:
					if _, err := param.Bot.ReplyMessage(
						event.ReplyToken,
						linebot.NewTextMessage("ご回答ありがペイ！\n"+"またレシートの画像を送信してみてペイ！"),
					).Do(); err != nil {
						return ok, fmt.Errorf("EventTypeMessageのReplyMessageでエラー: %v", err)
					}

					user.QuestionProgress = 2
					user.Marriage = 0

					// DBにメッセージを保存する処理
					err = i.userRepository.Update(user)
					return ok, nil

				case entity.QuestionMessageTitle[2] + ":" + entity.QuestionMessageSelection[2][1]:
					if _, err := param.Bot.ReplyMessage(
						event.ReplyToken,
						linebot.NewTextMessage("ご回答ありがペイ！\n"+"またレシートの画像を送信してみてペイ！"),
					).Do(); err != nil {
						return ok, fmt.Errorf("EventTypeMessageのReplyMessageでエラー: %v", err)
					}

					user.QuestionProgress = 2
					user.Marriage = 1

					// DBにメッセージを保存する処理
					err = i.userRepository.Update(user)
					return ok, nil

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

// func PushQuestionMessage(user *entity.User, presentPrice *entity.Present param *) error {
// 	questionMessageSelectionlength := len(entity.QuestionMessageSelection[user.QuestionProgress])
// 	if questionMessageSelectionlength == 2 {

// 		if _, err := param.Bot.ReplyMessage(
// 			event.ReplyToken,
// 			linebot.NewTextMessage(
// 				fmt.Sprintf(
// 					"チェックが完了したペイ！\n\n    %v円分のポイントをプレゼントするペイ！\n\nこれまで保有しているポイントは、%vポイントだペイ！\n\n今までのポイントを還元したい際は、メニューの「ポイント」ボタンからPayPayポイントとして受け取れるペイ！",
// 					presentPrice.Point,
// 					user.Point+presentPrice.Point,
// 					// convertPaymentServiceToStr(presentPrice.PaymentService),
// 					// presentList[0].Url,
// 				),
// 			),
// 			linebot.NewTemplateMessage(
// 				"アンケート",
// 				linebot.NewButtonsTemplate(
// 					"",
// 					"",
// 					entity.QuestionMessageTitle[user.QuestionProgress],
// 					linebot.NewMessageAction(entity.QuestionMessageSelection[user.QuestionProgress][0], entity.QuestionMessageTitle[user.QuestionProgress]+":"+entity.QuestionMessageSelection[user.QuestionProgress][0]),
// 					linebot.NewMessageAction(entity.QuestionMessageSelection[user.QuestionProgress][1], entity.QuestionMessageTitle[user.QuestionProgress]+":"+entity.QuestionMessageSelection[user.QuestionProgress][1]),
// 				),
// 			),
// 		).Do(); err != nil {
// 			return ok, fmt.Errorf("ImageMessageのReplyMessageでエラー: %v", err)
// 		}
// 	} else if questionMessageSelectionlength == 3 {
// 		if _, err = param.Bot.ReplyMessage(
// 			event.ReplyToken,
// 			linebot.NewTextMessage(
// 				fmt.Sprintf(
// 					"チェックが完了したペイ！\n\n    %v円分のポイントをプレゼントするペイ！\n\nこれまで保有しているポイントは、%vポイントだペイ！\n\n今までのポイントを還元したい際は、メニューの「ポイント」ボタンからPayPayポイントとして受け取れるペイ！",
// 					presentPrice.Point,
// 					user.Point+presentPrice.Point,
// 					// convertPaymentServiceToStr(presentPrice.PaymentService),
// 					// presentList[0].Url,
// 				),
// 			),
// 			linebot.NewTemplateMessage(
// 				"アンケート",
// 				linebot.NewButtonsTemplate(
// 					"",
// 					"",
// 					entity.QuestionMessageTitle[user.QuestionProgress],
// 					linebot.NewMessageAction(entity.QuestionMessageSelection[user.QuestionProgress][0], entity.QuestionMessageTitle[user.QuestionProgress]+":"+entity.QuestionMessageSelection[user.QuestionProgress][0]),
// 					linebot.NewMessageAction(entity.QuestionMessageSelection[user.QuestionProgress][1], entity.QuestionMessageTitle[user.QuestionProgress]+":"+entity.QuestionMessageSelection[user.QuestionProgress][1]),
// 					linebot.NewMessageAction(entity.QuestionMessageSelection[user.QuestionProgress][2], entity.QuestionMessageTitle[user.QuestionProgress]+":"+entity.QuestionMessageSelection[user.QuestionProgress][2]),
// 				),
// 			),
// 		).Do(); err != nil {
// 			return ok, fmt.Errorf("ImageMessageのReplyMessageでエラー: %v", err)
// 		}
// 	} else if questionMessageSelectionlength == 4 {
// 		if _, err = param.Bot.ReplyMessage(
// 			event.ReplyToken,
// 			linebot.NewTextMessage(
// 				fmt.Sprintf(
// 					"チェックが完了したペイ！\n\n    %v円分のポイントをプレゼントするペイ！\n\nこれまで保有しているポイントは、%vポイントだペイ！\n\n今までのポイントを還元したい際は、メニューの「ポイント」ボタンからPayPayポイントとして受け取れるペイ！",
// 					presentPrice.Point,
// 					user.Point+presentPrice.Point,
// 					// convertPaymentServiceToStr(presentPrice.PaymentService),
// 					// presentList[0].Url,
// 				),
// 			),
// 			linebot.NewTemplateMessage(
// 				"アンケート",
// 				linebot.NewButtonsTemplate(
// 					"",
// 					"",
// 					entity.QuestionMessageTitle[user.QuestionProgress],
// 					linebot.NewMessageAction(entity.QuestionMessageSelection[user.QuestionProgress][0], entity.QuestionMessageTitle[user.QuestionProgress]+":"+entity.QuestionMessageSelection[user.QuestionProgress][0]),
// 					linebot.NewMessageAction(entity.QuestionMessageSelection[user.QuestionProgress][1], entity.QuestionMessageTitle[user.QuestionProgress]+":"+entity.QuestionMessageSelection[user.QuestionProgress][1]),
// 					linebot.NewMessageAction(entity.QuestionMessageSelection[user.QuestionProgress][2], entity.QuestionMessageTitle[user.QuestionProgress]+":"+entity.QuestionMessageSelection[user.QuestionProgress][2]),
// 					linebot.NewMessageAction(entity.QuestionMessageSelection[user.QuestionProgress][3], entity.QuestionMessageTitle[user.QuestionProgress]+":"+entity.QuestionMessageSelection[user.QuestionProgress][3]),
// 				),
// 			),
// 		).Do(); err != nil {
// 			return ok, fmt.Errorf("ImageMessageのReplyMessageでエラー: %v", err)
// 		}

// 	} else if questionMessageSelectionlength == 5 {
// 		if _, err = param.Bot.ReplyMessage(
// 			event.ReplyToken,
// 			linebot.NewTextMessage(
// 				fmt.Sprintf(
// 					"チェックが完了したペイ！\n\n    %v円分のポイントをプレゼントするペイ！\n\nこれまで保有しているポイントは、%vポイントだペイ！\n\n今までのポイントを還元したい際は、メニューの「ポイント」ボタンからPayPayポイントとして受け取れるペイ！",
// 					presentPrice.Point,
// 					user.Point+presentPrice.Point,
// 					// convertPaymentServiceToStr(presentPrice.PaymentService),
// 					// presentList[0].Url,
// 				),
// 			),
// 			linebot.NewTemplateMessage(
// 				"アンケート",
// 				linebot.NewButtonsTemplate(
// 					"",
// 					"",
// 					entity.QuestionMessageTitle[user.QuestionProgress],
// 					linebot.NewMessageAction(entity.QuestionMessageSelection[user.QuestionProgress][0], entity.QuestionMessageTitle[user.QuestionProgress]+":"+entity.QuestionMessageSelection[user.QuestionProgress][0]),
// 					linebot.NewMessageAction(entity.QuestionMessageSelection[user.QuestionProgress][1], entity.QuestionMessageTitle[user.QuestionProgress]+":"+entity.QuestionMessageSelection[user.QuestionProgress][1]),
// 					linebot.NewMessageAction(entity.QuestionMessageSelection[user.QuestionProgress][2], entity.QuestionMessageTitle[user.QuestionProgress]+":"+entity.QuestionMessageSelection[user.QuestionProgress][2]),
// 					linebot.NewMessageAction(entity.QuestionMessageSelection[user.QuestionProgress][3], entity.QuestionMessageTitle[user.QuestionProgress]+":"+entity.QuestionMessageSelection[user.QuestionProgress][3]),
// 				),
// 			),
// 			linebot.NewTemplateMessage(
// 				"アンケート",
// 				linebot.NewButtonsTemplate(
// 					"",
// 					"",
// 					entity.QuestionMessageTitle[user.QuestionProgress],
// 					linebot.NewMessageAction(entity.QuestionMessageSelection[user.QuestionProgress][4], entity.QuestionMessageTitle[user.QuestionProgress]+":"+entity.QuestionMessageSelection[user.QuestionProgress][4]),
// 				),
// 			),
// 		).Do(); err != nil {
// 			return ok, fmt.Errorf("ImageMessageのReplyMessageでエラー: %v", err)
// 		}
// 	} else if questionMessageSelectionlength == 6 {
// 		if _, err = param.Bot.ReplyMessage(
// 			event.ReplyToken,
// 			linebot.NewTextMessage(
// 				fmt.Sprintf(
// 					"チェックが完了したペイ！\n\n    %v円分のポイントをプレゼントするペイ！\n\nこれまで保有しているポイントは、%vポイントだペイ！\n\n今までのポイントを還元したい際は、メニューの「ポイント」ボタンからPayPayポイントとして受け取れるペイ！",
// 					presentPrice.Point,
// 					user.Point+presentPrice.Point,
// 					// convertPaymentServiceToStr(presentPrice.PaymentService),
// 					// presentList[0].Url,
// 				),
// 			),
// 			linebot.NewTemplateMessage(
// 				"アンケート",
// 				linebot.NewButtonsTemplate(
// 					"",
// 					"",
// 					entity.QuestionMessageTitle[user.QuestionProgress],
// 					linebot.NewMessageAction(entity.QuestionMessageSelection[user.QuestionProgress][0], entity.QuestionMessageTitle[user.QuestionProgress]+":"+entity.QuestionMessageSelection[user.QuestionProgress][0]),
// 					linebot.NewMessageAction(entity.QuestionMessageSelection[user.QuestionProgress][1], entity.QuestionMessageTitle[user.QuestionProgress]+":"+entity.QuestionMessageSelection[user.QuestionProgress][1]),
// 					linebot.NewMessageAction(entity.QuestionMessageSelection[user.QuestionProgress][2], entity.QuestionMessageTitle[user.QuestionProgress]+":"+entity.QuestionMessageSelection[user.QuestionProgress][2]),
// 					linebot.NewMessageAction(entity.QuestionMessageSelection[user.QuestionProgress][3], entity.QuestionMessageTitle[user.QuestionProgress]+":"+entity.QuestionMessageSelection[user.QuestionProgress][3]),
// 				),
// 			),
// 			linebot.NewTemplateMessage(
// 				"アンケート",
// 				linebot.NewButtonsTemplate(
// 					"",
// 					"",
// 					entity.QuestionMessageTitle[user.QuestionProgress],
// 					linebot.NewMessageAction(entity.QuestionMessageSelection[user.QuestionProgress][4], entity.QuestionMessageTitle[user.QuestionProgress]+":"+entity.QuestionMessageSelection[user.QuestionProgress][4]),
// 					linebot.NewMessageAction(entity.QuestionMessageSelection[user.QuestionProgress][5], entity.QuestionMessageTitle[user.QuestionProgress]+":"+entity.QuestionMessageSelection[user.QuestionProgress][5]),
// 				),
// 			),
// 		).Do(); err != nil {
// 			return ok, fmt.Errorf("ImageMessageのReplyMessageでエラー: %v", err)
// 		}
// 	} else if questionMessageSelectionlength == 7 {
// 		if _, err = param.Bot.ReplyMessage(
// 			event.ReplyToken,
// 			linebot.NewTextMessage(
// 				fmt.Sprintf(
// 					"チェックが完了したペイ！\n\n    %v円分のポイントをプレゼントするペイ！\n\nこれまで保有しているポイントは、%vポイントだペイ！\n\n今までのポイントを還元したい際は、メニューの「ポイント」ボタンからPayPayポイントとして受け取れるペイ！",
// 					presentPrice.Point,
// 					user.Point+presentPrice.Point,
// 					// convertPaymentServiceToStr(presentPrice.PaymentService),
// 					// presentList[0].Url,
// 				),
// 			),
// 			linebot.NewTemplateMessage(
// 				"アンケート",
// 				linebot.NewButtonsTemplate(
// 					"",
// 					"",
// 					entity.QuestionMessageTitle[user.QuestionProgress],
// 					linebot.NewMessageAction(entity.QuestionMessageSelection[user.QuestionProgress][0], entity.QuestionMessageTitle[user.QuestionProgress]+":"+entity.QuestionMessageSelection[user.QuestionProgress][0]),
// 					linebot.NewMessageAction(entity.QuestionMessageSelection[user.QuestionProgress][1], entity.QuestionMessageTitle[user.QuestionProgress]+":"+entity.QuestionMessageSelection[user.QuestionProgress][1]),
// 					linebot.NewMessageAction(entity.QuestionMessageSelection[user.QuestionProgress][2], entity.QuestionMessageTitle[user.QuestionProgress]+":"+entity.QuestionMessageSelection[user.QuestionProgress][2]),
// 					linebot.NewMessageAction(entity.QuestionMessageSelection[user.QuestionProgress][3], entity.QuestionMessageTitle[user.QuestionProgress]+":"+entity.QuestionMessageSelection[user.QuestionProgress][3]),
// 				),
// 			),
// 			linebot.NewTemplateMessage(
// 				"アンケート",
// 				linebot.NewButtonsTemplate(
// 					"",
// 					"",
// 					entity.QuestionMessageTitle[user.QuestionProgress],
// 					linebot.NewMessageAction(entity.QuestionMessageSelection[user.QuestionProgress][4], entity.QuestionMessageTitle[user.QuestionProgress]+":"+entity.QuestionMessageSelection[user.QuestionProgress][4]),
// 					linebot.NewMessageAction(entity.QuestionMessageSelection[user.QuestionProgress][5], entity.QuestionMessageTitle[user.QuestionProgress]+":"+entity.QuestionMessageSelection[user.QuestionProgress][5]),
// 					linebot.NewMessageAction(entity.QuestionMessageSelection[user.QuestionProgress][6], entity.QuestionMessageTitle[user.QuestionProgress]+":"+entity.QuestionMessageSelection[user.QuestionProgress][6]),
// 				),
// 			),
// 		).Do(); err != nil {
// 			return ok, fmt.Errorf("ImageMessageのReplyMessageでエラー: %v", err)
// 		}
// 	} else if questionMessageSelectionlength == 8 {
// 		if _, err = param.Bot.ReplyMessage(
// 			event.ReplyToken,
// 			linebot.NewTextMessage(
// 				fmt.Sprintf(
// 					"チェックが完了したペイ！\n\n    %v円分のポイントをプレゼントするペイ！\n\nこれまで保有しているポイントは、%vポイントだペイ！\n\n今までのポイントを還元したい際は、メニューの「ポイント」ボタンからPayPayポイントとして受け取れるペイ！",
// 					presentPrice.Point,
// 					user.Point+presentPrice.Point,
// 					// convertPaymentServiceToStr(presentPrice.PaymentService),
// 					// presentList[0].Url,
// 				),
// 			),
// 			linebot.NewTemplateMessage(
// 				"アンケート",
// 				linebot.NewButtonsTemplate(
// 					"",
// 					"",
// 					entity.QuestionMessageTitle[user.QuestionProgress],
// 					linebot.NewMessageAction(entity.QuestionMessageSelection[user.QuestionProgress][0], entity.QuestionMessageTitle[user.QuestionProgress]+":"+entity.QuestionMessageSelection[user.QuestionProgress][0]),
// 					linebot.NewMessageAction(entity.QuestionMessageSelection[user.QuestionProgress][1], entity.QuestionMessageTitle[user.QuestionProgress]+":"+entity.QuestionMessageSelection[user.QuestionProgress][1]),
// 					linebot.NewMessageAction(entity.QuestionMessageSelection[user.QuestionProgress][2], entity.QuestionMessageTitle[user.QuestionProgress]+":"+entity.QuestionMessageSelection[user.QuestionProgress][2]),
// 					linebot.NewMessageAction(entity.QuestionMessageSelection[user.QuestionProgress][3], entity.QuestionMessageTitle[user.QuestionProgress]+":"+entity.QuestionMessageSelection[user.QuestionProgress][3]),
// 				),
// 			),
// 			linebot.NewTemplateMessage(
// 				"アンケート",
// 				linebot.NewButtonsTemplate(
// 					"",
// 					"",
// 					entity.QuestionMessageTitle[user.QuestionProgress],
// 					linebot.NewMessageAction(entity.QuestionMessageSelection[user.QuestionProgress][4], entity.QuestionMessageTitle[user.QuestionProgress]+":"+entity.QuestionMessageSelection[user.QuestionProgress][4]),
// 					linebot.NewMessageAction(entity.QuestionMessageSelection[user.QuestionProgress][5], entity.QuestionMessageTitle[user.QuestionProgress]+":"+entity.QuestionMessageSelection[user.QuestionProgress][5]),
// 					linebot.NewMessageAction(entity.QuestionMessageSelection[user.QuestionProgress][6], entity.QuestionMessageTitle[user.QuestionProgress]+":"+entity.QuestionMessageSelection[user.QuestionProgress][6]),
// 					linebot.NewMessageAction(entity.QuestionMessageSelection[user.QuestionProgress][7], entity.QuestionMessageTitle[user.QuestionProgress]+":"+entity.QuestionMessageSelection[user.QuestionProgress][7]),
// 				),
// 			),
// 		).Do(); err != nil {
// 			return ok, fmt.Errorf("ImageMessageのReplyMessageでエラー: %v", err)
// 		}
// 	}
// }

// func convertPaymentServiceToStr(paymentService int) string {
// 	// switch paymentService {
// 	// case entity.PayPay:
// 	// 	return "PayPay"
// 	// case entity.LinePay:
// 	// 	return "LINE Pay"
// 	// case entity.MercariPay:
// 	// 	return "メルペイ"
// 	// case entity.Cash:
// 	// 	return "現金"
// 	// case entity.AmazonPay:
// 	// 	return "Amazon Pay"
// 	// case entity.RakutenPay:
// 	// 	return "楽天ペイ"
// 	// }
// 	return "PayPay"
// }
