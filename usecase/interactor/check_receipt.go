package interactor

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	vision "cloud.google.com/go/vision/apiv1"
	"github.com/hidenari-yuda/paychan-server/domain/entity"
)

func CheckReceipt(
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
	// w := io.Writer(os.Stdout)

	ctx := context.Background()

	client, err := vision.NewImageAnnotatorClient(ctx)
	if err != nil {
		fmt.Println(err)
		return receiptPicture, present, err
	}

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

	annotations, err := client.DetectTexts(ctx, image, nil, 10)
	if err != nil {
		fmt.Println(err)
		return receiptPicture, present, err
	}

	if len(annotations) == 0 {
		fmt.Println(err)
		return receiptPicture, present, err
	}
	// receiptPicture.DetectedText = annotations[0].Description

	for _, v := range receiptPictureList {
		if receiptPicture.DetectedText == v.DetectedText {
			fmt.Println(err)
			return receiptPicture, present, err
		}
	}

	// "ホーム", "お知らせ", "出品", "支払い", "マイページ", // mercari
	// "取引履歴", "PayPay", // paypay
	// "支払い履歴", "LINE", // linePay
	// amazonPay
	// rakutenPay
	// d払い
	// auPay
	// famiPay

	strMap := ContainsList(annotations[0].Description,
		"合計", "円", "¥",
		"毎月のご利用状況", "出品", // mercari
		"取引履歴", "PayPay", // paypay
		"支払い履歴", "過去1ヶ月", // linePay
		// "取引履歴", "売上", "出品", "購入", "お知らせ", "マイページ", // yahoo
		// "ホーム", "お知らせ", "出品", "購入", "マイページ", // rakuma
		// "ホーム", "お知らせ", "出品", "購入", "マイページ", // yahoo
		// "お買い上げ", "合計金額", "お支払い", "お届け先", "ご注文内容", // rakuten suica  pasmo nanaco waon// manaca // piapo// my number // edy // nimoca
	)

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

	receiptPicture = &entity.ReceiptPicture{
		DetectedText: annotations[0].Description,
	}

	if strMap["取引履歴"] || strMap["PayPay"] {
		receiptPicture.Service = 0
		receiptPicture.PaymentService = 0
	} else if strMap["支払い履歴"] || strMap["過去1ヶ月"] {
		receiptPicture.Service = 1
		receiptPicture.PaymentService = 1
	} else if strMap["出品"] || strMap["毎月のご利用状況"] {
		receiptPicture.Service = 2
		receiptPicture.PaymentService = 2
	}

	// strings.Contains(annotations[0].Description, "合計", "円", "¥")

	// for i, _ := range annotations {
	// 	fmt.Println("number is:", i, "\n", "detected text is:", annotations[i].Description)

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
	// }

	fmt.Println("detected text is:", annotations[0].Description)
	fmt.Println("len is:", len)

	// dbに保存
	return receiptPicture, present, nil
}

func ContainsList(s string, list ...string) map[string]bool {
	var (
		strMap = make(map[string]bool)
	)
	for _, v := range list {
		if strings.Contains(s, v) {
			strMap[v] = true
		}
	}
	return strMap
}
