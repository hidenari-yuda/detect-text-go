package interactor

import (
	"github.com/hidenari-yuda/paychan-server/domain/entity"
	"github.com/hidenari-yuda/paychan-server/usecase"
)

type PresentInteractor interface {
	// Gest API
	Create(param *entity.Present) (present *entity.Present, err error)
	Update(param *entity.Present) (present *entity.Present, err error)
	GetById(id int) (present *entity.Present, err error)
	GetByLineUserId(LineUserId string) (presentList []*entity.Present, err error)
	GetAll() (presentList []*entity.Present, err error)
	DeleteByExpired() (ok bool, err error)
}

type PresentInteractorImpl struct {
	firebase                 usecase.Firebase
	userRepository           usecase.UserRepository
	receiptPictureRepository usecase.ReceiptPictureRepository
	receiptRepository        usecase.ReceiptRepository
	parchasedItemRepository  usecase.ParchasedItemRepository
	paymentMethodRepository  usecase.PaymentMethodRepository
	presentRepository        usecase.PresentRepository
	lineMessageRepository    usecase.LineMessageRepository
	aspRepository            usecase.AspRepository
}

func NewPresentInteractorImpl(
	fb usecase.Firebase,
	uR usecase.UserRepository,
	rpR usecase.ReceiptPictureRepository,
	rR usecase.ReceiptRepository,
	piR usecase.ParchasedItemRepository,
	pmR usecase.PaymentMethodRepository,
	pR usecase.PresentRepository,
	lmR usecase.LineMessageRepository,
	aR usecase.AspRepository,
) PresentInteractor {
	return &PresentInteractorImpl{
		firebase:                 fb,
		userRepository:           uR,
		receiptPictureRepository: rpR,
		receiptRepository:        rR,
		parchasedItemRepository:  piR,
		paymentMethodRepository:  pmR,
		presentRepository:        pR,
		lineMessageRepository:    lmR,
		aspRepository:            aR,
	}
}

func (i *PresentInteractorImpl) Create(param *entity.Present) (present *entity.Present, err error) {
	// if param.Expirary == nil {
	// 	param.Expirary = utility.GetNow().AddDate(0, 0, 1)
	// }

	present = param
	// ユーザー登録
	err = i.presentRepository.Create(present)
	if err != nil {
		return nil, err
	}

	return present, nil
}

func (i *PresentInteractorImpl) Update(param *entity.Present) (present *entity.Present, err error) {
	present = param
	// ユーザー登録
	err = i.presentRepository.Update(present)
	if err != nil {
		return nil, err
	}

	return param, nil
}

func (i *PresentInteractorImpl) GetById(id int) (present *entity.Present, err error) {
	// ユーザー登録
	present, err = i.presentRepository.GetById(id)
	if err != nil {
		return nil, err
	}

	return present, nil
}

func (i *PresentInteractorImpl) GetByLineUserId(LineUserId string) (presentList []*entity.Present, err error) {
	// ユーザー登録
	presentList, err = i.presentRepository.GetListByLineUserId(LineUserId)
	if err != nil {
		return nil, err
	}

	return presentList, nil
}

func (i *PresentInteractorImpl) GetAll() (presentList []*entity.Present, err error) {
	// ユーザー登録
	presentList, err = i.presentRepository.GetAll()
	if err != nil {
		return nil, err
	}

	return presentList, nil
}

func (i *PresentInteractorImpl) DeleteByExpired() (ok bool, err error) {
	// ユーザー登録
	err = i.presentRepository.DeleteByExpired()
	if err != nil {
		return ok, err
	}

	ok = true

	return ok, nil
}
