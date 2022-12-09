package interactor

import (
	"github.com/hidenari-yuda/paychan/domain/entity"
	"github.com/hidenari-yuda/paychan/usecase"
)

type PresentInteractor interface {
	// Gest API
	Create(param *entity.Present) (*entity.Present, error)
}

type PresentInteractorImpl struct {
	firebase                 usecase.Firebase
	PresentRepository        usecase.PresentRepository
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
	uR usecase.PresentRepository,
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
		PresentRepository:        uR,
		receiptPictureRepository: rpR,
		receiptRepository:        rR,
		parchasedItemRepository:  piR,
		paymentMethodRepository:  pmR,
		presentRepository:        pR,
		lineMessageRepository:    lmR,
		aspRepository:            aR,
	}
}

func (i *PresentInteractorImpl) Create(param *entity.Present) (*entity.Present, error) {
	// ユーザー登録
	err := i.PresentRepository.Create(param)
	if err != nil {
		return nil, err
	}

	return param, nil
}
