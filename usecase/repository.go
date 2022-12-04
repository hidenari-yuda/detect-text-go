package usecase

import "github.com/hidenari-yuda/umerun-resume/domain/entity"

type UserRepository interface {
	// Gest API
	SignUp(param *entity.SignUpParam) error
	GetByFirebaseId(firebaseId string) (*entity.User, error)
	GetByLineUserId(lineUserId string) (*entity.User, error)
}

type PaymentMethodRepository interface {
	// Gest API
	Create(param *entity.PaymentMethod) error
	// Delete(id uint) error
	GetById(id uint) (*entity.PaymentMethod, error)
	// GetListByUserId(userId uint) ([]*entity.PaymentMethod, error)
	GetListByLineUserId(lineUserId string) ([]*entity.PaymentMethod, error)
}

type LineMessageRepository interface {
	// Gest API
	Create(param *entity.LineMessage) error
	GetById(id uint) (*entity.LineMessage, error)
	// GetListByUserId(userId uint) ([]*entity.LineMessage, error)
	GetListByLineUserId(lineUserId string) ([]*entity.LineMessage, error)
}

type ReceiptRepository interface {
	// Gest API
	Create(param *entity.Receipt) error
	GetById(id uint) (*entity.Receipt, error)
	// GetListByUserId(userId uint) ([]*entity.Receipt, error)
	GetListByLineUserId(lineUserId string) ([]*entity.Receipt, error)
}

type ParchasedItemRepository interface {
	// Gest API
	Create(param *entity.ParchasedItem) error
	GetById(id uint) (*entity.ParchasedItem, error)
	GetListByReceiptId(receiptId uint) ([]*entity.ParchasedItem, error)
	GetListByLineUserId(lineUserId string) ([]*entity.ParchasedItem, error)
}

type GiftRepository interface {
	// Gest API
	Create(param *entity.Gift) error
	GetById(id uint) (*entity.Gift, error)
	// GetListByUserId(userId uint) ([]*entity.Gift, error)
	GetListByReceiptId(receiptId uint) ([]*entity.Gift, error)
	GetListByLineUserId(lineUserId string) ([]*entity.Gift, error)
}

type AspRepository interface {
	// Gest API
	Create(param *entity.Asp) error
	GetById(id uint) (*entity.Asp, error)
	// GetListByUserId(userId uint) ([]*entity.Asp, error)
	GetListByLineUserId(lineUserId string) ([]*entity.Asp, error)
}
