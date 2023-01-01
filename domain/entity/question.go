package entity

import "time"

type Question struct {
	Id               int       `db:"id" json:"id"`
	Uuid             string    `db:"uuid" json:"uuid"`
	UserId           int       `db:"user_id" json:"user_id"`
	LineUserId       string    `db:"line_user_id" json:"line_user_id"`
	ReceiptPictureId int       `db:"receipt_picture_id" json:"receipt_picture_id"`
	Question         int       `db:"question" json:"question"` // 0: お店名, 1: お店の住所, 2: お店の電話番号, 3: お店の営業時間, 4: お店の定休日, 5: お店のURL, 6: お店のメールアドレス, 7: お店のその他の連絡先, 8: お店のその他の情報
	Text             string    `db:"text" json:"text"`         // お店の情報
	CreatedAt        time.Time `db:"created_at" json:"created_at"`
	UpdatedAt        time.Time `db:"updated_at" json:"updated_at"`

	// 関連
	QuestionSelections []QuestionSelection `db:"question_selections" json:"question_selections"`
}

type QuestionSelection struct {
	Id         int       `db:"id" json:"id"`
	Uuid       string    `db:"uuid" json:"uuid"`
	QuestionId int       `db:"question_id" json:"question_id"`
	Selection  int       `db:"selection" json:"selection"` // 0: お店名, 1: お店の住所, 2: お店の電話番号, 3: お店の営業時間, 4: お店の定休日, 5: お店のURL, 6: お店のメールアドレス, 7: お店のその他の連絡先, 8: お店のその他の情報
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
	UpdatedAt  time.Time `db:"updated_at" json:"updated_at"`
}

var QuestionMessageTitle map[int]string = map[int]string{
	0: "ご年齢を教えてペイ！",        // 0: 10代, 1: 20代, 2: 30代, 3: 40代, 4: 50代, 5: 60代, 6: 70代以上
	1: "ご性別を教えてペイ！",        // 0: 男性, 1: 女性
	2: "世帯年収を教えてペイ！",       // 0: 100万円台, 1: 200万円台, 2: 300万円台, 3: 400万円台, 4: 500万円台, 5: 600万円台, 6: 700万円台, 7: 800万円台, 8: 900万円台, 9: 1000万円以上
	3: "ご職業を教えてペイ！",        // 0: 会社員, 1: 公務員, 2: 自営業, 3: 学生, 4: 主婦, 5: その他
	4: "お仕事の業界について教えてペイ!",  // 0: IT, 1: 金融, 2: サービス, 3: メーカー, 4: その他
	5: "ご婚姻状況を教えてペイ！",      // 0: 未婚, 1: 既婚
	6: "現在の世帯人数について教えてペイ！", // 0: 1人, 1: 2人, 2: 3人, 3: 4人, 4: 5人以上
	7: "お子さんの有無について教えてペイ！", // 0: いる, 1: いない
}

var QuestionMessageSelection [][]string = [][]string{
	[]string{"10代", "20代", "30代", "40代", "50代", "60代", "70代以上"},
	[]string{"男性", "女性", "その他"},
	[]string{"100万円台", "200万円台", "300万円台", "400万円台", "500万円台", "600万円台", "700万円台", "800万円以上"},
	[]string{"会社員（一般）", "会社員（管理職）", "会社経営（経営者・役員）", "公務員", "派遣社員・契約社員", "学生", "主婦", "その他"},
	[]string{"IT・通信", "金融・保険", "小売・サービス", "製造・メーカー", "運輸", "電気・ガス", "不動産・建築", "その他"},
	[]string{"未婚", "既婚"},
	[]string{"1人", "2人", "3人", "4人", "5人以上"},
	[]string{"いる", "いない"},
}

// linebot.NewTemplateMessage(
// 	"アンケート",
// 	linebot.NewButtonsTemplate(
// 		"",
// 		"",
// 		entity.QuestionMessageTitle[user.QuestionProgress],
// 		linebot.NewMessageAction("はい", "はい"),
// 		linebot.NewMessageAction("いいえ", "いいえ"),
// 	),
