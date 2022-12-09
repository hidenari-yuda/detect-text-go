package presenter

import "github.com/hidenari-yuda/paychan/domain/entity/responses"

func NewOkJSONPresenter(resp responses.OK) Presenter {
	return NewJSONPresenter(200, resp)
}
