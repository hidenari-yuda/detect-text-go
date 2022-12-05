package handler

import (
	"fmt"

	"github.com/hidenari-yuda/detect-text/domain/entity"
	"github.com/hidenari-yuda/detect-text/domain/entity/responses"
	"github.com/hidenari-yuda/detect-text/interfaces/presenter"

	"github.com/hidenari-yuda/detect-text/usecase/interactor"
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

	output, err := h.UserInteractor.SignUp(interactor.SignUpInput{
		Param: param,
	})
	fmt.Println(output, err)

	if err != nil {
		// c.JSON(c, presenter.NewErrorJsonPresenter(err))
		return nil, err
	}

	return presenter.NewOkJSONPresenter(responses.NewOK(output.Ok)), nil

}

func (h *UserHandlerImpl) SignIn(param *entity.SignInParam) (presenter.Presenter, error) {
	output, err := h.UserInteractor.SignIn(interactor.SignInInput{
		Param: param,
	})
	if err != nil {
		// c.JSON(c, presenter.NewErrorJsonPresenter(err))
		return nil, err
	}

	return presenter.NewUserJSONPresenter(responses.NewUser(output.User)), nil

}

func (h *UserHandlerImpl) GetByFirebaseToken(token string) (presenter.Presenter, error) {
	output, err := h.UserInteractor.GetByFirebaseToken(interactor.GetByFirebaseTokenInput{
		Token: token,
	})

	if err != nil {
		// c.JSON(c, presenter.NewErrorJsonPresenter(err))
		return nil, err
	}

	return presenter.NewUserJSONPresenter(responses.NewUser(output.User)), nil
}

func (h *UserHandlerImpl) GetByLineUserId(lineUserId string) (presenter.Presenter, error) {
	output, err := h.UserInteractor.GetByLineUserId(interactor.GetByLineUserIdInput{
		LineUserId: lineUserId,
	})

	if err != nil {
		// c.JSON(c, presenter.NewErrorJsonPresenter(err))
		return nil, err
	}

	return presenter.NewUserJSONPresenter(responses.NewUser(output.User)), nil
}

func (h *UserHandlerImpl) GetLineWebHook(param *entity.LineWebHook) (presenter.Presenter, error) {
	output, err := h.UserInteractor.GetLineWebHook(interactor.GetLineWebHookInput{
		Param: param,
	})
	if err != nil {
		return nil, err
	}

	return presenter.NewOkJSONPresenter(responses.NewOK(output.Ok)), nil
}
