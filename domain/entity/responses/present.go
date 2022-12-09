package responses

import "github.com/hidenari-yuda/paychan/domain/entity"

type Present struct {
	Present *entity.Present `json:"present"`
}

func NewPresent(present *entity.Present) Present {
	return Present{
		Present: present,
	}
}

type PresentList struct {
	PresentList []*entity.Present `json:"present_list"`
}

func NewPresentList(presents []*entity.Present) PresentList {
	return PresentList{
		PresentList: presents,
	}
}
