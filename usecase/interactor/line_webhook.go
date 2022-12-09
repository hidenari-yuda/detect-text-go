package interactor

import (
	"fmt"

	"github.com/hidenari-yuda/paychan/domain/config"
	"github.com/hidenari-yuda/paychan/domain/entity"
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

			// 既存のユーザーではない場合は新規登録
			res, err := input.Param.Bot.GetProfile(event.Source.UserID).Do()
			if err != nil {
				return output, fmt.Errorf("lineプロフィールの取得エラー: %w", err)
			}
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
				receiptPicture, presentPrice, err := checkReceipt(content.Content, receiptPictures)
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
				presentList, err := i.presentRepository.GetByPriceAndService(presentPrice)
				if err != nil || len(presentList) == 0 {
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
								presentPrice.Price,
								convertPaymentServiceToStr(presentPrice.PaymentService),
								err,
							),
						),
					).Do(); eff != nil {
						return output, fmt.Errorf("プレゼント取得通知のリプライエラー: %w", err)
					}
					return output, fmt.Errorf("プレゼントの取得エラー: %w", err)
				}

				if _, err = input.Param.Bot.ReplyMessage(
					event.ReplyToken,
					linebot.NewTextMessage(fmt.Sprintf(
						"チェックが完了したトン！\n\n    %v円分のプレゼントを%vで送るトン！\n\n  以下のリンクから受け取ってね！\n\n    %v",
						presentPrice.Price,
						convertPaymentServiceToStr(presentPrice.PaymentService),
						presentList[0].Url,
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
				err = i.presentRepository.Update(&entity.Present{
					Id:               presentList[0].Id,
					UserId:           user.Id,
					ReceiptPictureId: receiptPicture.Id,
					Price:            presentPrice.Price,
					PaymentService:   presentList[0].PaymentService,
					Url:              presentList[0].Url,
				})

				fmt.Println("image function is ok!!")

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
	switch paymentService {
	case entity.PayPay:
		return "PayPay"
	case entity.LinePay:
		return "LINE Pay"
	case entity.MercariPay:
		return "メルペイ"
	case entity.Cash:
		return "現金"
	case entity.AmazonPay:
		return "Amazon Pay"
	case entity.RakutenPay:
		return "楽天ペイ"
	}
	return "PayPay"
}
