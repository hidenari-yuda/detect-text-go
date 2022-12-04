package entity

import "time"

type Asp struct {
	Id        uint      `db:"id" json:"id"`
	Uuid      string    `db:"uuid" json:"uuid"`
	UserId    uint      `db:"user_id" json:"user_id"`
	Service   uint      `db:"service" json:"service"` // 0: amazon, 1: rakuten 2:yahoo (default: 0)
	Url       string    `db:"url" json:"url"`
	Price     uint      `db:"price" json:"price"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}
