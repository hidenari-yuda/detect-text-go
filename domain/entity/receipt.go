package entity

import "time"

type Receipt struct {
	Id               uint      `db:"id" json:"id"`
	Uuid             string    `db:"uuid" json:"uuid"`
	ReceiptPictureId uint      `db:"receipt_picture_id" json:"receipt_picture_id"`
	StoreName        string    `db:"store_name" json:"store_name"`
	TotalPrice       uint      `db:"total_price" json:"total_price"`
	PurchasedAt      time.Time `db:"purchased_at" json:"purchased_at"`
	CreatedAt        time.Time `db:"created_at" json:"created_at"`
	UpdatedAt        time.Time `db:"updated_at" json:"updated_at"`

	// 関連テーブル
	ParchasedItems []ParchasedItem `db:"parchased_items" json:"parchased_items"`
}

func NewReceipt() *Receipt {
	return &Receipt{}
}
