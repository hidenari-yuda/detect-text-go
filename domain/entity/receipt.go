package entity

import "time"

type Receipt struct {
	Id          uint      `db:"id" json:"id"`
	Uuid        string    `db:"uuid" json:"uuid"`
	UserId      uint      `db:"user_id" json:"user_id"`
	Url         string    `db:"url" json:"url"`
	StoreName   string    `db:"store_name" json:"store_name"`
	Price       uint      `db:"price" json:"price"`
	GiftPrice   uint      `db:"gift_price" json:"gift_price"`
	PurchasedAt time.Time `db:"purchased_at" json:"purchased_at"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`

	// 関連テーブル
	ParchasedItems []ParchasedItem `db:"parchased_items" json:"parchased_items"`
}

func NewReceipt() *Receipt {
	return &Receipt{}
}
