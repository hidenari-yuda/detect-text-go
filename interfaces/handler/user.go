package handler

import (
	"github.com/hidenari-yuda/paychan-server/domain/entity"
	"github.com/hidenari-yuda/paychan-server/domain/entity/responses"
	"github.com/hidenari-yuda/paychan-server/interfaces/presenter"
	"github.com/hidenari-yuda/paychan-server/usecase/interactor"
)

type UserHandler interface {
	// Gest API
	SignUp(param *entity.SignUpParam) (presenter.Presenter, error)
	SignIn(param *entity.SignInParam) (presenter.Presenter, error)
	GetByFirebaseToken(token string) (presenter.Presenter, error)
	GetByLineUserId(lineUserId string) (presenter.Presenter, error)

	// Line API
	GetLineWebHook(param *entity.LineWebHook) (presenter.Presenter, error)
}

type UserHandlerImpl struct {
	UserInteractor interactor.UserInteractor
}

func NewUserHandlerImpl(ui interactor.UserInteractor) UserHandler {
	return &UserHandlerImpl{
		UserInteractor: ui,
	}
}

func (h *UserHandlerImpl) SignUp(param *entity.SignUpParam) (presenter.Presenter, error) {

	ok, err := h.UserInteractor.SignUp(param)
	if err != nil {
		// c.JSON(c, presenter.NewErrorJsonPresenter(err))
		return nil, err
	}

	return presenter.NewOkJSONPresenter(responses.NewOK(ok)), nil

}

func (h *UserHandlerImpl) SignIn(param *entity.SignInParam) (presenter.Presenter, error) {
	user, err := h.UserInteractor.SignIn(param)
	if err != nil {
		// c.JSON(c, presenter.NewErrorJsonPresenter(err))
		return nil, err
	}

	return presenter.NewUserJSONPresenter(responses.NewUser(user)), nil

}

func (h *UserHandlerImpl) GetByFirebaseToken(token string) (presenter.Presenter, error) {
	user, err := h.UserInteractor.GetByFirebaseToken(token)

	if err != nil {
		// c.JSON(c, presenter.NewErrorJsonPresenter(err))
		return nil, err
	}

	return presenter.NewUserJSONPresenter(responses.NewUser(user)), nil
}

func (h *UserHandlerImpl) GetByLineUserId(lineUserId string) (presenter.Presenter, error) {
	user, err := h.UserInteractor.GetByLineUserId(lineUserId)

	if err != nil {
		// c.JSON(c, presenter.NewErrorJsonPresenter(err))
		return nil, err
	}

	return presenter.NewUserJSONPresenter(responses.NewUser(user)), nil
}

func (h *UserHandlerImpl) GetLineWebHook(param *entity.LineWebHook) (presenter.Presenter, error) {
	ok, err := h.UserInteractor.GetLineWebHook(param)
	if err != nil {
		return nil, err
	}

	return presenter.NewOkJSONPresenter(responses.NewOK(ok)), nil
}
