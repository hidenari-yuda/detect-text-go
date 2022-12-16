package entity

import "time"

type Asp struct {
	Id        int       `db:"id" json:"id"`
	Uuid      string    `db:"uuid" json:"uuid"`
	UserId    int       `db:"user_id" json:"user_id"`
	Service   int       `db:"service" json:"service"` // 0: amazon, 1: rakuten 2:yahoo (default: 0)
	Url       string    `db:"url" json:"url"`
	Price     int       `db:"price" json:"price"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}
