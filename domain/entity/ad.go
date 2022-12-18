package entity

import "time"

type Ad struct {
	Id          int       `db:"id" json:"id"`
	Uuid        string    `db:"uuid" json:"uuid"`
	Service     int       `db:"service" json:"service"` // 0: amazon, 1: rakuten 2:yahoo (default: 0)
	Url         string    `db:"url" json:"url"`
	ImageUrl    string    `db:"image_url" json:"image_url"`
	Price       int       `db:"price" json:"price"`
	Title       string    `db:"title" json:"title"`
	Description string    `db:"description" json:"description"`
	Impression  int       `db:"impression" json:"impression"`
	Click       int       `db:"click" json:"click"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}
