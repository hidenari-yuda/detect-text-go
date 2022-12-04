package entity

import "time"

type ParchasedItem struct {
	Name      string    `db:"name" json:"name"`     // 商品名
	Price     uint      `db:"price" json:"price"`   // 価格
	Number    uint      `db:"number" json:"number"` // 個数
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}
