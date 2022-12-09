package handler

import (
	"github.com/hidenari-yuda/paychan/domain/entity"
	"github.com/hidenari-yuda/paychan/domain/entity/responses"
	"github.com/hidenari-yuda/paychan/interfaces/presenter"

	"github.com/hidenari-yuda/paychan/usecase/interactor"
)

type PresentHandler interface {
	// Gest API
	Create(param *entity.Present) (presenter.Presenter, error)
}

type PresentHandlerImpl struct {
	PresentInteractor interactor.PresentInteractor
}

func NewPresentHandlerImpl(ui interactor.PresentInteractor) PresentHandler {
	return &PresentHandlerImpl{
		PresentInteractor: ui,
	}
}

func (h *PresentHandlerImpl) Create(param *entity.Present) (presenter.Presenter, error) {

	present, err := h.PresentInteractor.Create(param)
	if err != nil {
		return nil, err
	}

	return presenter.NewPresentJSONPresenter(responses.NewPresent(present)), nil

}
