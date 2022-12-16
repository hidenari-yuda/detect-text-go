package entity

import "time"

type ParchasedItem struct {
	Id        int       `db:"id" json:"id"`
	Uuid      string    `db:"uuid" json:"uuid"`
	ReceiptId int       `db:"receipt_id" json:"receipt_id"`
	Name      string    `db:"name" json:"name"`     // 商品名
	Price     int       `db:"price" json:"price"`   // 価格
	Number    int       `db:"number" json:"number"` // 個数
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}
