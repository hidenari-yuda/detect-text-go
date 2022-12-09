package presenter

import "github.com/hidenari-yuda/paychan-server/domain/entity/responses"

func NewOkJSONPresenter(resp responses.OK) Presenter {
	return NewJSONPresenter(200, resp)
}
