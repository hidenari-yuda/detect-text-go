package interactor

import (
	"fmt"

	"github.com/hidenari-yuda/umerun-resume/domain/entity"
	"github.com/hidenari-yuda/umerun-resume/usecase"
	"github.com/line/line-bot-sdk-go/v7/linebot"
)

type UserInteractor interface {
	// Gest API
	SignUp(input SignUpInput) (output SignUpOutput, err error)
	SignIn(input SignInInput) (output SignInOutput, err error)
	GetByFirebaseToken(input GetByFirebaseTokenInput) (output GetByFirebaseTokenOutput, err error)
	GetByLineUserId(input GetByLineUserIdInput) (output GetByLineUserIdOutput, err error)

	// line
	GetLineWebHook(input GetLineWebHookInput) (output GetLineWebHookOutput, err error)

	// resume
}

type UserInteractorImpl struct {
	firebase       usecase.Firebase
	userRepository usecase.UserRepository
}

func NewUserInteractorImpl(
	fb usecase.Firebase,
	uR usecase.UserRepository,
) UserInteractor {
	return &UserInteractorImpl{
		firebase:       fb,
		userRepository: uR,
	}
}

type SignUpInput struct {
	Param *entity.SignUpParam
}

type SignUpOutput struct {
	Ok bool
}

func (i *UserInteractorImpl) SignUp(input SignUpInput) (output SignUpOutput, err error) {
	// ユーザー登録
	err = i.userRepository.SignUp(input.Param)
	if err != nil {
		return output, err
	}

	output.Ok = true

	return output, nil
}

type SignInInput struct {
	Param *entity.SignInParam
}

type SignInOutput struct {
	User *entity.User
}

func (i *UserInteractorImpl) SignIn(input SignInInput) (output SignInOutput, err error) {

	firebaseId, err := i.firebase.VerifyIDToken(input.Param.Token)
	if err != nil {
		return output, err
	}

	fmt.Println("exported firebaseToken is:", input.Param.Token)
	fmt.Println("exported firebaseId is:", firebaseId)

	// ユーザー登録
	output.User, err = i.userRepository.GetByFirebaseId(firebaseId)
	if err != nil {
		err = fmt.Errorf("failed to get user by firebaseId: %w", err)
		return output, err
	}

	return output, nil

}

type GetByFirebaseTokenInput struct {
	Token string
}

type GetByFirebaseTokenOutput struct {
	User *entity.User
}

func (i *UserInteractorImpl) GetByFirebaseToken(input GetByFirebaseTokenInput) (output GetByFirebaseTokenOutput, err error) {

	firebaseId, err := i.firebase.VerifyIDToken(input.Token)
	if err != nil {
		return output, err
	}

	fmt.Println("exported firebaseId is:", firebaseId)

	output.User, err = i.userRepository.GetByFirebaseId(firebaseId)
	if err != nil {
		return output, err
	}

	fmt.Println("exported user is:", output.User)

	return output, nil
}

type GetByLineUserIdInput struct {
	LineUserId string
}

type GetByLineUserIdOutput struct {
	User *entity.User
}

func (i *UserInteractorImpl) GetByLineUserId(input GetByLineUserIdInput) (output GetByLineUserIdOutput, err error) {

	output.User, err = i.userRepository.GetByLineUserId(input.LineUserId)
	if err != nil {
		return output, err
	}

	fmt.Println("exported user is:", output.User)

	return output, nil
}

type GetLineWebHookInput struct {
	Param *entity.LineWebHookParam
}

type GetLineWebHookOutput struct {
	Ok bool
}

func (i *UserInteractorImpl) GetLineWebHook(input GetLineWebHookInput) (output GetLineWebHookOutput, err error) {

	for _, event := range input.Param.Events {
		fmt.Println("event", event)
		fmt.Println("送信者のID", event.Source.UserID)

		if event.Type == linebot.EventTypeMessage {
			switch message := event.Message.(type) {
			// テキストメッセージの場合
			case *linebot.TextMessage:
				fmt.Println("message", message)

				fmt.Println("TextMessage")
				fmt.Println("----------------")
				// if _, err := bot.ReplyMessage(
				// 	event.ReplyToken,
				// 	linebot.NewTextMessage(message.Text),
				// ).Do(); err != nil {
				// 	fmt.Println(err)
				// 	fmt.Println("EventTypeMessageのReplyMessageでエラー")
				// }
				// DBにメッセージを保存する処理

				// event.Message.(*linebot.TextMessage).Text,

			// スタンプの場合
			case *linebot.StickerMessage:
				fmt.Println("message", message)

				fmt.Println("StickerMessage")
				fmt.Println("----------------")
				// if _, err := bot.ReplyMessage(
				// 	event.ReplyToken,
				// 	linebot.NewStickerMessage(11111, 11111),
				// ).Do(); err != nil {
				// 	fmt.Println(err)
				// 	fmt.Println("StickerMessageのReplyMessageでエラー")
				// }

			// 画像データの場合
			case *linebot.ImageMessage:
				fmt.Println("message", message)
				fmt.Println("ImageMessage")
				fmt.Println("message.ContentProvider.Type", message.ContentProvider.Type)
				fmt.Println("----------------")

				// コンテンツ取得
				bot, err := linebot.New("LINE_SECRET", "LINE_ACCESS_TOKEN")
				if err != nil {
					fmt.Println(err)
					fmt.Println("linebot.Newでエラー")
				}
				content, err := bot.GetMessageContent(message.ID).Do()
				if err != nil {
					fmt.Println(err)
					fmt.Println("GetMessageContentでエラー")
				}
				defer content.Content.Close()
				fmt.Println("content.Content", content.Content)

				// レシートか判定
				info, err := checkReceipt(content.Content)
				if err != nil {
					fmt.Println(err)
					fmt.Println("detectReceiptでエラー")
					_, err = bot.ReplyMessage(
						event.ReplyToken,
						linebot.NewTextMessage(fmt.Sprint("画像を受け取りました")),
					).Do()
					if err != nil {
						fmt.Println(err)
						fmt.Println("ImageMessageのReplyMessageでエラー")
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

				// message.OriginalContentURL,
				// message.PreviewImageURL,
			// 動画データの場合
			case *linebot.VideoMessage:
				fmt.Println("message", message)

				fmt.Println("VideoMessage")
				fmt.Println("----------------")
				// if _, err := bot.ReplyMessage(
				// 	event.ReplyToken,
				// 	linebot.NewVideoMessage(message.OriginalContentURL, message.PreviewImageURL),
				// ).Do(); err != nil {
				// 	fmt.Println(err)
				// 	fmt.Println("VideoMessageのReplyMessageでエラー")
				// }
				// DBにメッセージを保存する処理

				// message.OriginalContentURL,
				// message.PreviewImageURL,

			// 音声データの場合
			case *linebot.AudioMessage:
				fmt.Println("message", message)

				fmt.Println("AudioMessage")
				fmt.Println("----------------")
				// if _, err := bot.ReplyMessage(
				// 	event.ReplyToken,
				// 	linebot.NewAudioMessage(message.OriginalContentURL, message.Duration),
				// ).Do(); err != nil {
				// 	fmt.Println(err)
				// 	fmt.Println("AudioMessageのReplyMessageでエラー")
				// }
				// DBにメッセージを保存する処理

				// message.OriginalContentURL,

				// null.NewInt(int64(message.Duration), true),

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
