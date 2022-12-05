package interactor

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	vision "cloud.google.com/go/vision/apiv1"
	"github.com/hidenari-yuda/umerun-resume/domain/entity"
)

func checkReceipt(content io.ReadCloser, receiptPictureList []*entity.ReceiptPicture) (*entity.ReceiptPicture, *entity.Present, error) {
	var (
		receiptPicture *entity.ReceiptPicture
		present        *entity.Present
		err            error
	)

	byteData, err := io.ReadAll(content)
	if err != nil {
		return receiptPicture, present, err
	}

	filePath := fmt.Sprint("./public/snapshot/", time.Now(), "-receipt.jpg") // ファイル名をユニークにする

	os.WriteFile(filePath, byteData, 0644)

	// io.Writer
	w := io.Writer(os.Stdout)

	ctx := context.Background()

	client, err := vision.NewImageAnnotatorClient(ctx)
	if err != nil {
		return receiptPicture, present, err
	}

	file, err := os.Open(filePath)
	if err != nil {
		return receiptPicture, present, err
	}
	defer file.Close()

	image, err := vision.NewImageFromReader(file)
	if err != nil {
		return receiptPicture, present, err
	}
	annotations, err := client.DetectTexts(ctx, image, nil, 10)
	if err != nil {
		return receiptPicture, present, err
	}

	if len(annotations) == 0 {
		fmt.Fprintln(w, "No text found.")
		return receiptPicture, present, err
	}

	receiptPicture.DetectedText = annotations[0].Description

	for _, v := range receiptPictureList {
		if receiptPicture.DetectedText == v.DetectedText {
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

	if strMap["支払い履歴"] || strMap["過去1ヶ月"] {
		receiptPicture.Service = 0
	} else if strMap["取引履歴"] || strMap["PayPay"] {
		receiptPicture.Service = 1
	} else if strMap["出品"] || strMap["毎月のご利用状況"] {
		receiptPicture.Service = 2
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

	len := len(annotations)
	present = &entity.Present{
		Price: uint(len),
	}

	fmt.Println("len is:", present)

	// dbに保存
	return receiptPicture, present, nil
}

// func detectTextFromReceipt(content io.ReadCloser) (uint, error) {
// 	return receiptPicture, present, nil
// }

func ContainsList(s string, list ...string) (strMap map[string]bool) {
	for _, v := range list {
		if strings.Contains(s, v) {
			strMap[v] = true
		}
	}
	return strMap
}
