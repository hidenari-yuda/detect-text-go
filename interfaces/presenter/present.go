package presenter

import "github.com/hidenari-yuda/paychan-server/domain/entity/responses"

func NewPresentJSONPresenter(resp responses.Present) Presenter {
	return NewJSONPresenter(200, resp)
}

func NewPresentListJSONPresenter(resp responses.PresentList) Presenter {
	return NewJSONPresenter(200, resp)
}
