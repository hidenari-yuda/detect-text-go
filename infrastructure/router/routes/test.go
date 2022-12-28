package routes

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	vision "cloud.google.com/go/vision/apiv1"
	"github.com/hidenari-yuda/paychan-server/domain/config"
	"github.com/hidenari-yuda/paychan-server/domain/entity"
	"github.com/hidenari-yuda/paychan-server/infrastructure/database"
	"github.com/hidenari-yuda/paychan-server/infrastructure/di"
	"github.com/hidenari-yuda/paychan-server/usecase"
	"github.com/hidenari-yuda/paychan-server/usecase/interactor"
	"github.com/labstack/echo/v4"
	"github.com/line/line-bot-sdk-go/v7/linebot"
	"google.golang.org/api/option"
)

type TestRoutes struct {
}

func (r *TestRoutes) DetectTest(db *database.DB, firebase usecase.Firebase) func(c echo.Context) error {
	return func(c echo.Context) error {
		_ = di.InitializeUserHandler(db, firebase)

		filePath := fmt.Sprint(".public/snapshot/test.png") // ファイル名をユニークにする
		fileTest, err := os.Open(filePath)
		// byteData, err := io.ReadAll(fileTest)
		// if err != nil {
		// 	fmt.Println("ファイルが読み込めません:", err)
		// 	return err
		// }

		// err = os.WriteFile(filePath, byteData, 0644)
		// if err != nil {
		// 	fmt.Println("ファイルが書き出せません:", err)
		// 	return err
		// }

		// io.Writer
		w := io.Writer(os.Stdout)

		ctx := context.Background()
		// cfg, err := config.New()
		// fmt.Println("vision.NewImageAnnotatorClient(ctx)に入ります", cfg.Google.ApplicationCredentials)

		// client, err := vision.NewImageAnnotatorClient(ctx)
		client, err := vision.NewImageAnnotatorClient(ctx, option.WithCredentialsFile(".conf/google-application-credentials-prd.json"))
		if err != nil {
			fmt.Println("vision.NewImageAnnotatorClient(ctx)に失敗しました:", err)
			return err
		}

		file, err := os.Open(filePath)
		if err != nil {
			fmt.Println("os.Open(filePath)に失敗しました:", err)
			return err
		}
		defer file.Close()

		image, err := vision.NewImageFromReader(fileTest)
		if err != nil {
			fmt.Println("vision.NewImageFromReader(fileTest)に失敗しました:", err)
			return err
		}
		annotations, err := client.DetectTexts(ctx, image, nil, 10)
		if err != nil {
			fmt.Println("client.DetectTexts(ctx, image, nil, 10)に失敗しました:", err)
			return err
		}

		if len(annotations) == 0 {
			fmt.Fprintln(w, "No text found.")
			return err
		}

		fmt.Println("Texts:", annotations[0].Description)

		// s

		// renderJSON(c, annotations)
		return nil
	}
}

func (r *TestRoutes) CheckReceiptTestRoutes(db *database.DB, firebase usecase.Firebase) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			fileName string = c.Param("fileName")
		)
		_ = di.InitializeUserHandler(db, firebase)
		fmt.Println(fileName)

		filePath := fmt.Sprintf(".public/snapshot/%s", fileName) // ファイル名をユニークにする

		file, err := os.Open(filePath)
		if err != nil {
			fmt.Println(err)
			return err
		}
		defer file.Close()

		var receiptPictureList []*entity.ReceiptPicture = []*entity.ReceiptPicture{}

		receiptPicture, present, err := CheckReceiptTest(file, receiptPictureList)
		if err != nil {
			fmt.Println(err)
			return err
		}

		fmt.Println(receiptPicture, present)
		return nil
	}
}

func CheckReceiptTest(
	content io.ReadCloser,
	receiptPictureList []*entity.ReceiptPicture,
) (
	receiptPicture *entity.ReceiptPicture,
	present *entity.Present,
	err error) {

	byteData, err := io.ReadAll(content)
	if err != nil {
		fmt.Println(err)
		return receiptPicture, present, err
	}

	filePath := fmt.Sprint(".public/snapshot/", time.Now().Format("20060102150405"), "-receipt.jpg") // ファイル名をユニークにする

	err = os.WriteFile(filePath, byteData, 0644)
	if err != nil {
		fmt.Println("ファイルが書き出せません", err)
		return receiptPicture, present, err
	}

	// io.Writer
	w := io.Writer(os.Stdout)

	ctx := context.Background()

	client, err := vision.NewImageAnnotatorClient(ctx)
	if err != nil {
		fmt.Println(err)
		return receiptPicture, present, err
	}

	fmt.Println("vision.NewImageAnnotatorClient(ctx)に入ります")

	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println(err)
		return receiptPicture, present, err
	}
	defer file.Close()
	defer os.Remove(filePath)

	image, err := vision.NewImageFromReader(file)
	if err != nil {
		fmt.Println(err)
		return receiptPicture, present, err
	}

	fmt.Println("vision.NewImageFromReader(file)に入ります")
	annotations, err := client.DetectTexts(ctx, image, nil, 10)
	if err != nil {
		fmt.Println(err)
		return receiptPicture, present, err
	}
	fmt.Println("client.DetectTexts(ctx, image, nil, 10)に入ります")

	if len(annotations) == 0 {
		fmt.Fprintln(w, "No text found.")
		fmt.Println(err)
		return receiptPicture, present, err
	}
	fmt.Println("Texts:", annotations[0].Description)
	// receiptPicture.DetectedText = annotations[0].Description

	for _, v := range receiptPictureList {
		if receiptPicture.DetectedText == v.DetectedText {
			fmt.Println(err)
			return receiptPicture, present, err
		}
	}
	fmt.Println("receiptPictureListに入ります")

	// "ホーム", "お知らせ", "出品", "支払い", "マイページ", // mercari
	// "取引履歴", "PayPay", // paypay
	// "支払い履歴", "LINE", // linePay
	// amazonPay
	// rakutenPay
	// d払い
	// auPay
	// famiPay

	strMap := interactor.ContainsList(annotations[0].Description,
		"合計", "円", "¥",
		"毎月のご利用状況", "出品", // mercari
		"取引履歴", "PayPay", // paypay
		"支払い履歴", "過去1ヶ月", // linePay
		// "取引履歴", "売上", "出品", "購入", "お知らせ", "マイページ", // yahoo
		// "ホーム", "お知らせ", "出品", "購入", "マイページ", // rakuma
		// "ホーム", "お知らせ", "出品", "購入", "マイページ", // yahoo
		// "お買い上げ", "合計金額", "お支払い", "お届け先", "ご注文内容", // rakuten
		// "ご利用金額", "ご利用明細", "ご利用日", "ご利用店舗", "ご利用店舗", // suica
		// "ご利用金額", "ご利用明細", "ご利用日", "ご利用店舗", "ご利用店舗", // pasmo
		// "ご利用金額", "ご利用明細", "ご利用日", "ご利用店舗", "ご利用店舗", // nanaco
		// "ご利用金額", "ご利用明細", "ご利用日", "ご利用店舗", "ご利用店舗", // waon
		// "ご利用金額", "ご利用明細", "ご利用日", "ご利用店舗", "ご利用店舗", // manaca
		// "ご利用金額", "ご利用明細", "ご利用日", "ご利用店舗", "ご利用店舗", // piapo
		// "ご利用金額", "ご利用明細", "ご利用日", "ご利用店舗", "ご利用店舗", // my number
		// "ご利用金額", "ご利用明細", "ご利用日", "ご利用店舗", "ご利用店舗", // edy
		// "ご利用金額", "ご利用明細", "ご利用日", "ご利用店舗", "ご利用店舗", // nimoca
	)

	receiptPicture = &entity.ReceiptPicture{
		DetectedText: annotations[0].Description,
	}

	if strMap["取引履歴"] || strMap["PayPay"] {
		receiptPicture.Service = 0
		receiptPicture.PaymentService = 0
	} else if strMap["支払い履歴"] || strMap["過去1ヶ月"] {
		receiptPicture.Service = 1
		receiptPicture.PaymentService = 0
	} else if strMap["出品"] || strMap["毎月のご利用状況"] {
		receiptPicture.Service = 2
		receiptPicture.PaymentService = 0
	}

	// strings.Contains(annotations[0].Description, "合計", "円", "¥")

	for i, _ := range annotations {
		fmt.Println("number is:", i, "\n", "detected text is:", annotations[i].Description)

		// receipt.StoreName = annotations[i].Descriptions

		// parchasedItem := entity.ParchasedItem{
		// 	Name: annotations[i].Description,
		// 	// Price: annotations[i].Description,
		// 	Price: 0,
		// }

		// if ContainsList(annotations[i].Description, "円", "¥") {
		// 	receipt.ParchasedItems = append(receipt.ParchasedItems, entity.ParchasedItem{
		// 		Name:  annotations[i-1].Description,
		// 		Price: 0,
		// 	})
		// }

		// receipt.ParchasedItems = append(receipt.ParchasedItems, parchasedItem)
	}

	var point int
	// len := annotations[0].BoundingPoly.Vertices[2].X - annotations[0].BoundingPoly.Vertices[0].X
	len := len(annotations[0].Description)
	if len < 50 {
		point = 3
	} else if len < 100 {
		point = 5
	} else if len < 150 {
		point = 7
	} else if len < 200 {
		point = 10
	}

	present = &entity.Present{
		Point: point,
	}

	fmt.Println("len is:", len)
	fmt.Println("point is:", point)
	fmt.Println("present is:", present)
	fmt.Println("receiptPicture is:", receiptPicture)

	// dbに保存
	return receiptPicture, present, nil
}

func (r *TestRoutes) PushMessageTest(db *database.DB, firebase usecase.Firebase) func(c echo.Context) error {
	return func(c echo.Context) error {

		cfg, err := config.New()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}

		bot, err := linebot.New(cfg.Line.ChannelSecret, cfg.Line.ChannelAccessToken)
		if err != nil {
			fmt.Println(err)
			return err
		}

		if _, err := bot.PushMessage(
			cfg.Line.AdminUserId,
			// linebot.NewButtonsTemplate
			linebot.NewTemplateMessage(
				"下のボタンから還元する方法を選べるペイ！",
				// 	linebot.NewButtonsTemplate(
				// 		"",
				// 		"",
				// 		"下のボタンから還元する方法を選べるペイ！",
				// 		linebot.NewMessageAction("PayPayポイントに還元", "PayPayポイントに還元"),
				// 		linebot.NewMessageAction("PayPayカードに還元", "PayPayカードに還元"),
				// 		linebot.NewMessageAction("PayPay残高に還元", "PayPay残高に還元"),
				// 		linebot.NewMessageAction("sss", "s"),
				// 	),
				// ),
				// linebot.NewTemplateMessage(
				// 	"下のボタンから還元する方法を選べるペイ！",
				// 	linebot.NewButtonsTemplate(
				// 		"",
				// 		"",
				// 		"下のボタンから還元する方法を選べるペイ！",
				// 		linebot.NewMessageAction("PayPayポイントに還元", "PayPayポイントに還元"),
				// 		linebot.NewMessageAction("PayPayカードに還元", "PayPayカードに還元"),
				// 		linebot.NewMessageAction("PayPay残高に還元", "PayPay残高に還元"),
				// 		linebot.NewMessageAction("sss", "s"),
				// 	),
				// ),

				linebot.NewCarouselTemplate(
					linebot.NewCarouselColumn(
						"",
						"",
						"キャンペーン",
						linebot.NewMessageAction("PayPayポイントに還元", "PayPayポイントに還元"),
						linebot.NewMessageAction("PayPayカードに還元", "PayPayカードに還元"),
						linebot.NewMessageAction("PayPay残高に還元", "PayPay残高に還元"),
					),
					linebot.NewCarouselColumn(
						"",
						"",
						"キャンペーン",
						linebot.NewMessageAction("PayPayポイントに還元", "PayPayポイントに還元"),
						linebot.NewMessageAction("PayPayカードに還元", "PayPayカードに還元"),
						linebot.NewMessageAction("PayPay残高に還元", "PayPay残高に還元"),
					),
				),
				// linebot.NewButtonsTemplate(
				// 	".public/snapshot/test.png",
				// 	"PayPayポイントに還元",
				// 	"PayPayポイントに還元",
				// 	linebot.NewMessageAction("PayPayポイントに還元", "PayPayポイントに還元"),
				// // linebot.NewMessageAction("LINEPayポイントに還元", "LINEPayポイントに還元"),
				// // linebot.NewURIAction("PayPayポイントに還元", "line://app/1653824439-5jQXjz5A"),
				// ),
			),
		).Do(); err != nil {
			return fmt.Errorf("EventTypeMessageのReplyMessageでエラー: %v", err)
		}

		return c.JSON(http.StatusOK, "ok")
	}
}
