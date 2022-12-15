package interactor

import (
	"crypto/subtle"

	"github.com/hidenari-yuda/paychan-server/domain/config"
)

// AdminInteractor is an Interface
//
type AdminInteractor interface {
	Authorize(input AdminAuthorizeInput) (AdminAuthorizeOutput, error)
}

// AdminInteractorImpl is an implementation of AdminInteractor
type AdminInteractorImpl struct {
	basicUsers     []string
	basicPasswords []string
}

// NewAdminInteractorImpl is an initializer for AdminInteractorImpl
//
func NewAdminInteractorImpl(appConfig config.App) AdminInteractor {
	return &AdminInteractorImpl{
		basicUsers:     appConfig.BasicUsers,
		basicPasswords: appConfig.BasicPasswords,
	}
}

type AdminAuthorizeInput struct {
	Username string
	Password string
}

type AdminAuthorizeOutput struct {
	OK bool
}

func (i *AdminInteractorImpl) Authorize(input AdminAuthorizeInput) (AdminAuthorizeOutput, error) {
	var (
		output = AdminAuthorizeOutput{}
		u, p   string
	)

	for idx, name := range i.basicUsers {
		if input.Username == name {
			u = i.basicUsers[idx]
			p = i.basicPasswords[idx]
		}
	}
	if p == "" || u == "" {
		return output, nil
	}

	if subtle.ConstantTimeCompare([]byte(input.Username), []byte(u)) == 1 &&
		subtle.ConstantTimeCompare([]byte(input.Password), []byte(p)) == 1 {
		output.OK = true
	}

	return output, nil
}
