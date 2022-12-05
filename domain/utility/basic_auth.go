package utility

import (
	"crypto/subtle"

	"github.com/hidenari-yuda/detect-text/domain/config"
	"github.com/labstack/echo/v4"
)

type BasicAuth struct {
	cfg config.Config
}

func NewBasicAuth(cfg config.Config) *BasicAuth {
	return &BasicAuth{
		cfg: cfg,
	}
}

func (b *BasicAuth) BasicAuthValidator(username, password string, c echo.Context) (bool, error) {
	var (
		u, p string
	)

	for i, name := range b.cfg.App.BasicUsers {
		if username == name {
			u = b.cfg.App.BasicUsers[i]
			p = b.cfg.App.BasicPasswords[i]
		}
	}
	if p == "" || u == "" {
		return false, nil
	}

	if subtle.ConstantTimeCompare([]byte(username), []byte(u)) == 1 &&
		subtle.ConstantTimeCompare([]byte(password), []byte(p)) == 1 {
		return true, nil
	}

	return false, nil
}
