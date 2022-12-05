package entity

import "time"

type Present struct {
	Id               uint           `db:"id" json:"id"`
	Uuid             string         `db:"uuid" json:"uuid"`
	UserId           uint           `db:"user_id" json:"user_id"`
	ReceiptPictureId uint           `db:"receipt_picture_id" json:"receipt_picture_id"`
	PaymentService   PaymentService `db:"payment_service" json:"payment_service"` // 0: amazon, 1: rakuten 2:yahoo (default: 0)
	Price            uint           `db:"price" json:"price"`
	Url              string         `db:"url" json:"url"`
	CreatedAt        time.Time      `db:"created_at" json:"created_at"`
	UpdatedAt        time.Time      `db:"updated_at" json:"updated_at"`
}

func NewPresent() *Present {
	return &Present{}
}

type PresentParam struct {
	UserId    uint   `json:"user_id"`
	ReceiptId uint   `json:"receipt_id"`
	Service   uint   `json:"service"`
	Price     uint   `json:"price"`
	Url       string `json:"url"`
}
