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
