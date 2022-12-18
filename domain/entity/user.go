package entity

import (
	"time"
)

type User struct {
	Id            int       `db:"id" json:"id"`
	Uuid          string    `db:"uuid" json:"uuid"`
	FirebaseId    string    `db:"firebase_id" json:"firebase_id"`
	LineUserId    string    `db:"line_user_id" json:"line_user_id"`
	LineName      string    `db:"line_name" json:"line_name"`
	PictureUrl    string    `db:"picture_url" json:"picture_url"`
	StatusMessage string    `db:"status_message" json:"status_message"`
	Language      string    `db:"language" json:"language"`
	Point         int       `db:"point" json:"point"`
	Prefecture    int       `db:"prefecture" json:"prefecture"`
	Age           int       `db:"age" json:"age"`                     // 0: 10代, 1: 20代, 2: 30代, 3: 40代, 4: 50代, 5: 60代, 6: 70代以上
	Gender        int       `db:"gender" json:"gender"`               // 0: 男性, 1: 女性 2: その他
	Occupation    int       `db:"occupation" json:"occupation"`       // 0: 学生, 1: 会社員, 2: 自営業, 3: 公務員, 4: その他
	Married       int       `db:"married" json:"married"`             // 0: 未婚, 1: 既婚
	AnnualIncome  int       `db:"annual_income" json:"annual_income"` // 0: 100万円台, 1: 200万円台, 2: 300万円台, 3: 400万円台, 4: 500万円台, 5: 600万円台, 6: 700万円台, 7: 800万円台, 8: 900万円台, 9: 1000万円以上
	Name          string    `db:"name" json:"name"`
	Email         string    `db:"email" json:"email"`
	Password      string    `db:"password" json:"password"`
	CreatedAt     time.Time `db:"created_at" json:"created_at"`
	UpdatedAt     time.Time `db:"updated_at" json:"updated_at"`

	// 関連テーブル
	ReceiptPictures []ReceiptPicture `db:"-" json:"receipt_pictures"`
	PaymentMethods  []PaymentMethod  `db:"-" json:"payment_methods"`
	Presents        []Present        `db:"-" json:"presents"`
	LineMessages    []LineMessage    `db:"-" json:"line_messages"`

	// dbにはないが、返却用に追加
	TotalPrice    int    `db:"-" json:"total_price"`
	TotalResheets int    `db:"-" json:"total_resheets"`
	LineOneWord   string `db:"-" json:"line_one_word"`
	Active        bool   `db:"-" json:"active"`
}

func NewUser() *User {
	return &User{}
}

type SignUpParam struct {
	FirebaseId    string `db:"firebase_id" json:"firebase_id"`
	Name          string `db:"name" json:"name"`
	Email         string `db:"email" json:"email"`
	Password      string `db:"password" json:"password"`
	LineUserId    string `db:"line_user_id" json:"line_user_id"`
	LineName      string `db:"line_name" json:"line_name"`
	PictureUrl    string `db:"picture_url" json:"picture_url"`
	StatusMessage string `db:"status_message" json:"status_message"`
	Language      string `db:"language" json:"language"`
}

type SignInParam struct {
	Token string `json:"token" validate:"required"`
}
