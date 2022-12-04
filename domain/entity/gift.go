package entity

import "time"

type Gift struct {
	Id        uint      `db:"id" json:"id"`
	Uuid      string    `db:"uuid" json:"uuid"`
	UserId    uint      `db:"user_id" json:"user_id"`
	ReceiptId uint      `db:"receipt_id" json:"receipt_id"`
	Price     uint      `db:"price" json:"price"`
	Url       string    `db:"url" json:"url"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}
