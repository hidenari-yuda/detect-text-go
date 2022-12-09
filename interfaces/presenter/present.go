package presenter

import "github.com/hidenari-yuda/paychan/domain/entity/responses"

func NewPresentJSONPresenter(resp responses.Present) Presenter {
	return NewJSONPresenter(200, resp)
}

func NewPresentListJSONPresenter(resp responses.PresentList) Presenter {
	return NewJSONPresenter(200, resp)
}
