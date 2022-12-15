package handler

import (
	"github.com/hidenari-yuda/paychan-server/domain/entity/responses"
	"github.com/hidenari-yuda/paychan-server/interfaces/presenter"
	"github.com/hidenari-yuda/paychan-server/usecase/interactor"
)

type AdminHandler interface {
	Authorize(username, password string) (presenter.Presenter, error)
}

type AdminHandlerImpl struct {
	adminInteractor interactor.AdminInteractor
}

func NewAdminHandlerImpl(aI interactor.AdminInteractor) AdminHandler {
	return &AdminHandlerImpl{adminInteractor: aI}
}

func (h *AdminHandlerImpl) Authorize(username, password string) (presenter.Presenter, error) {
	var (
		input = interactor.AdminAuthorizeInput{Username: username, Password: password}
	)

	output, err := h.adminInteractor.Authorize(input)
	if err != nil {
		return nil, err
	}

	return presenter.NewOkJSONPresenter(responses.NewOK(output.OK)), nil
}
