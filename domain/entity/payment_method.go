package entity

import "time"

type PaymentMethod struct {
	Id             uint           `db:"id" json:"id"`
	UserId         uint           `db:"user_id" json:"user_id"`
	PaymentService PaymentService `db:"payment_service" json:"payment_service"` // 0: amazon, 1: rakuten, 2: yahoo, 3: other (default: 3)
	CreatedAt      time.Time      `db:"created_at" json:"created_at"`
	UpdatedAt      time.Time      `db:"updated_at" json:"updated_at"`
}
