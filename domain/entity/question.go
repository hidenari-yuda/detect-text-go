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
	0: "ご年齢を教えてペイ！",   // 0: 10代, 1: 20代, 2: 30代, 3: 40代, 4: 50代, 5: 60代, 6: 70代以上
	2: "性別を教えてペイ！",    // 0: 男性, 1: 女性
	3: "ご職業を教えてペイ！",   // 0: 会社員, 1: 公務員, 2: 自営業, 3: 学生, 4: 主婦, 5: その他
	4: "ご家族構成を教えてペイ！", // 0: 一人暮らし, 1: 夫婦, 2: 子供あり, 3: その他
	5: "年収を教えてペイ！",    // 0: 100万円台, 1: 200万円台, 2: 300万円台, 3: 400万円台, 4: 500万円台, 5: 600万円台, 6: 700万円台, 7: 800万円台, 8: 900万円台, 9: 1000万円以上
}

var QuestionMessageSelection [][]string = [][]string{
	[]string{"10代", "20代", "30代", "40代", "50代", "60代", "70代以上"},
	[]string{"男性", "女性", "その他"},
	[]string{"会社員", "公務員", "自営業", "学生", "主婦", "その他"},
	[]string{"一人暮らし", "夫婦", "子供あり", "その他"},
	[]string{"100万円台", "200万円台", "300万円台", "400万円台", "500万円台", "600万円台", "700万円台", "800万円台", "900万円台", "1000万円以上"},
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
