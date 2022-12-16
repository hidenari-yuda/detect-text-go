package utility

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"fmt"

	"github.com/hidenari-yuda/paychan-server/domain/config"
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
			// token := b.cfg.App
		}
	}
	if p == "" || u == "" {
		return false, nil
	}

	mac := hmac.New(sha256.New, []byte(b.cfg.App.BasicSecret))
	mac.Write([]byte(p))
	expected := hex.EncodeToString(mac.Sum(nil))

	fmt.Println("username: ", username)
	fmt.Println("password:", password)
	fmt.Println("expected:", expected)

	if subtle.ConstantTimeCompare([]byte(username), []byte(u)) == 1 &&
		subtle.ConstantTimeCompare([]byte(password), []byte(expected)) == 1 {
		return true, nil
	}

	return false, nil
}
