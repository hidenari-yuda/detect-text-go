package entity

import (
	"time"
)

type User struct {
	Id            uint      `db:"id" json:"id"`
	Uuid          string    `db:"uuid" json:"uuid"`
	FirebaseId    string    `db:"firebase_id" json:"firebase_id"`
	LineUserId    string    `db:"line_user_id" json:"line_user_id"`
	LineName      string    `db:"-" json:"line_name"`
	PictureUrl    string    `json:"picture_url"`
	StatusMessage string    `json:"status_message"`
	Language      string    `json:"language"`
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
	TotalResheets uint   `db:"-" json:"total_resheets"`
	LineOneWord   string `db:"-" json:"line_one_word"`
	Active        bool   `db:"-" json:"active"`
}

func NewUser() *User {
	return &User{}
}

type SignUpParam struct {
	Email         string `db:"email" json:"email"`
	Password      string `db:"password" json:"password"`
	LineUserId    string `db:"line_user_id" json:"line_user_id"`
	LineName      string `db:"-" json:"line_name"`
	PictureUrl    string `json:"picture_url"`
	StatusMessage string `json:"status_message"`
	Language      string `json:"language"`
}

type SignInParam struct {
	Token string `json:"token" validate:"required"`
}
