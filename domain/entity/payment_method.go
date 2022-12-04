package entity

import "time"

type PaymentMethod struct {
	Id        uint      `db:"id" json:"id"`
	UserId    uint      `db:"user_id" json:"user_id"`
	Method    uint      `db:"method" json:"method"` // 0: cash, 1: credit card 2:mobile (default: 0)
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}
