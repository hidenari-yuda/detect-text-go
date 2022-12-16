package usecase

import "github.com/hidenari-yuda/paychan-server/domain/entity"

type UserRepository interface {
	// Gest API
	SignUp(param *entity.SignUpParam) error
	SignIn(email, password string) (user *entity.User, err error)
	GetByFirebaseId(firebaseId string) (*entity.User, error)
	GetByLineUserId(lineUserId string) (*entity.User, error)

	// admin
	GetAll() ([]*entity.User, error)
}

type PaymentMethodRepository interface {
	// Gest API
	Create(param *entity.PaymentMethod) error
	Update(param *entity.PaymentMethod) error
	// Delete(id int) error
	GetById(id int) (*entity.PaymentMethod, error)
	// GetListByUserId(userId int) ([]*entity.PaymentMethod, error)
	GetListByLineUserId(lineUserId string) ([]*entity.PaymentMethod, error)
}

type LineMessageRepository interface {
	// Gest API
	Create(param *entity.LineMessage) error
	GetById(id int) (*entity.LineMessage, error)
	// GetListByUserId(userId int) ([]*entity.LineMessage, error)
	GetListByLineUserId(lineUserId string) ([]*entity.LineMessage, error)
}

type ReceiptPictureRepository interface {
	// Gest API
	Create(param *entity.ReceiptPicture) error
	Update(param *entity.ReceiptPicture) error
	GetById(id int) (*entity.ReceiptPicture, error)
	// GetListByUserId(userId int) ([]*entity.ReceiptPicture, error)
	GetListByLineUserId(lineUserId string) ([]*entity.ReceiptPicture, error)
	GetListByToday(lineUserId string) ([]*entity.ReceiptPicture, error)
}

type ReceiptRepository interface {
	// Gest API
	Create(param *entity.Receipt) error
	Update(param *entity.Receipt) error
	GetById(id int) (*entity.Receipt, error)
	// GetListByUserId(userId int) ([]*entity.Receipt, error)
	GetListByLineUserId(lineUserId string) ([]*entity.Receipt, error)
	GetListByToday(lineUserId string) ([]*entity.Receipt, error) // 今日登録されたレシートのリストを取得する
}

type ParchasedItemRepository interface {
	// Gest API
	Create(param *entity.ParchasedItem) error
	Update(param *entity.ParchasedItem) error
	GetById(id int) (*entity.ParchasedItem, error)
	GetListByReceiptId(receiptId int) ([]*entity.ParchasedItem, error)
	GetListByLineUserId(lineUserId string) ([]*entity.ParchasedItem, error)
}

type PresentRepository interface {
	// Gest API
	Create(param *entity.Present) error
	Update(param *entity.Present) error
	GetById(id int) (*entity.Present, error)
	GetByPointAndService(present *entity.Present) ([]*entity.Present, error)
	// GetListByUserId(userId int) ([]*entity.Present, error)
	GetListByReceiptPictureId(receiptPictureId int) ([]*entity.Present, error)
	GetListByLineUserId(lineUserId string) ([]*entity.Present, error)
	GetAll() ([]*entity.Present, error)

	DeleteByExpired() error
}

type AspRepository interface {
	// Gest API
	Create(param *entity.Asp) error
	GetById(id int) (*entity.Asp, error)
	// GetListByUserId(userId int) ([]*entity.Asp, error)
	GetListByLineUserId(lineUserId string) ([]*entity.Asp, error)
}
