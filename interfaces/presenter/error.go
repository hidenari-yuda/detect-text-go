package presenter

import "github.com/hidenari-yuda/detect-text/domain/entity"

type ErrorJsonPresenter struct {
	PresenterImpl
}

func NewErrorJsonPresenter(err error) Presenter {
	code, message := entity.ErrorInfo(err)
	return &ErrorJsonPresenter{
		PresenterImpl: PresenterImpl{
			statusCode: code,
			data: map[string]interface{}{
				"error": message,
			},
		},
	}
}
