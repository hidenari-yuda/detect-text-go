package usecase

import "io"

type Firebase interface {
	VerifyIDToken(idToken string) (string, error)
	GetCustomToken(uid string) (string, error)
	GetIDToken(token string) (string, error)
	GetPhoneNumber(uid string) (string, error)
	Set(doc string, data map[string]interface{}) error
	CreateUser(email, password string) (string, error)
	UpdateEmail(email, uid string) error
	UpdatePassword(password, uid string) error
	UploadImage(content io.ReadCloser, messageId string) (string, error)
}

type Cache interface {
	GetBytes(key string) ([]byte, error)
	GetString(key string) (string, error)
	Set(key string, obj interface{}, ttl int) (interface{}, error)
	Do(commandName string, args ...interface{}) (interface{}, error)
	Values(reply interface{}, err error) ([]interface{}, error)
}
