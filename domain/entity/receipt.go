package entity

import "time"

type Receipt struct {
	Id            uint      `db:"id" json:"id"`
	Uuid          string    `db:"uuid" json:"uuid"`
	UserId        uint      `db:"user_id" json:"user_id"`
	Url           string    `db:"url" json:"url"`
	DetectedText  string    `db:"detected_text" json:"detected_text"`
	Service       uint      `db:"service" json:"service"` // 0: amazon, 1: rakuten, 2: yahoo, 3: other (default: 3)
	StoreName     string    `db:"store_name" json:"store_name"`
	PaymentMethod uint      `db:"payment_method" json:"payment_method"` // 0: cash, 1: credit card 2:mobile (default: 0)
	TotalPrice    uint      `db:"total_price" json:"total_price"`
	GiftPrice     uint      `db:"gift_price" json:"gift_price"`
	PurchasedAt   time.Time `db:"purchased_at" json:"purchased_at"`
	CreatedAt     time.Time `db:"created_at" json:"created_at"`
	UpdatedAt     time.Time `db:"updated_at" json:"updated_at"`

	// 関連テーブル
	ParchasedItems []ParchasedItem `db:"parchased_items" json:"parchased_items"`
}
