package entity

import "time"

type Present struct {
	Id               int       `db:"id" json:"id"`
	Uuid             string    `db:"uuid" json:"uuid"`
	UserId           int       `db:"user_id" json:"user_id"`
	LineUserId       string    `db:"line_user_id" json:"line_user_id"`
	ReceiptPictureId int       `db:"receipt_picture_id" json:"receipt_picture_id"`
	PaymentService   int       `db:"payment_service" json:"payment_service"` // 0: amazon, 1: rakuten 2:yahoo (default: 0)
	Point            int       `db:"point" json:"point"`
	Expirary         time.Time `db:"expirary" json:"expirary"`
	Used             bool      `db:"used" json:"used"`
	Url              string    `db:"url" json:"url"`
	CreatedAt        time.Time `db:"created_at" json:"created_at"`
	UpdatedAt        time.Time `db:"updated_at" json:"updated_at"`
}

func NewPresent() *Present {
	return &Present{}
}
