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
	Update(param *entity.Present) (presenter.Presenter, error)
	GetById(id uint) (presenter.Presenter, error)
	GetByLineUserId(lineUserId string) (presenter.Presenter, error)
	GetAll() (presenter.Presenter, error)
	DeleteByExpired() (presenter.Presenter, error)
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
func (h *PresentHandlerImpl) Update(param *entity.Present) (presenter.Presenter, error) {

	present, err := h.PresentInteractor.Update(param)
	if err != nil {
		return nil, err
	}

	return presenter.NewPresentJSONPresenter(responses.NewPresent(present)), nil
}

func (h *PresentHandlerImpl) GetById(id uint) (presenter.Presenter, error) {
	present, err := h.PresentInteractor.GetById(id)

	if err != nil {
		return nil, err
	}

	return presenter.NewPresentJSONPresenter(responses.NewPresent(present)), nil
}

func (h *PresentHandlerImpl) GetByLineUserId(lineUserId string) (presenter.Presenter, error) {
	presentList, err := h.PresentInteractor.GetByLineUserId(lineUserId)

	if err != nil {
		return nil, err
	}

	return presenter.NewPresentListJSONPresenter(responses.NewPresentList(presentList)), nil
}

func (h *PresentHandlerImpl) GetAll() (presenter.Presenter, error) {
	presentList, err := h.PresentInteractor.GetAll()

	if err != nil {
		return nil, err
	}

	return presenter.NewPresentListJSONPresenter(responses.NewPresentList(presentList)), nil
}

func (h *PresentHandlerImpl) DeleteByExpired() (presenter.Presenter, error) {

	ok, err := h.PresentInteractor.DeleteByExpired()
	if err != nil {
		return nil, err
	}

	return presenter.NewOkJSONPresenter(responses.NewOK(ok)), nil
}
