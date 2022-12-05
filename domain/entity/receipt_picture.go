package entity

import "time"

type ReceiptPicture struct {
	Id             uint           `db:"id" json:"id"`
	Uuid           string         `db:"uuid" json:"uuid"`
	UserId         uint           `db:"user_id" json:"user_id"`
	LineUserId     string         `db:"line_user_id" json:"line_user_id"`
	Url            string         `db:"url" json:"url"`
	DetectedText   string         `db:"detected_text" json:"detected_text"`
	Service        Service        `db:"service" json:"service"`                 // 0: amazon, 1: rakuten, 2: yahoo, 3: other (default: 3)
	PaymentService PaymentService `db:"payment_service" json:"payment_service"` // 0: amazon, 1: rakuten, 2: yahoo, 3: other (default: 3)
	TotalPrice     uint           `db:"total_price" json:"total_price"`
	CreatedAt      time.Time      `db:"created_at" json:"created_at"`
	UpdatedAt      time.Time      `db:"updated_at" json:"updated_at"`

	// 関連テーブル
	Receips []Receipt `db:"-" json:"receipts"`
}

func NewReceiptPicture() *ReceiptPicture {
	return &ReceiptPicture{}
}

type Service uint

const (
	paypay Service = iota
	line
	mercari
	amazon
	rakuten
)

type PaymentService uint

const (
	PayPay PaymentService = iota
	LinePay
	MercariPay
	AmazonPay
	RakutenPay
	Docomo
	AuPay
	FamiPay
	Cash
	NullPaymentService = 99
)
