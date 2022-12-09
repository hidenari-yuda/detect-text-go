package driver

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"github.com/hidenari-yuda/paychan-server/domain/config"
	"github.com/hidenari-yuda/paychan-server/domain/entity"
	"github.com/hidenari-yuda/paychan-server/usecase"
	"github.com/pkg/errors"
	"google.golang.org/api/option"
)

type FirebaseImpl struct {
	auth      *auth.Client
	firestore *firestore.Client
	webAPIKey string
}

func NewFirebaseImpl(fbConfig config.Firebase) usecase.Firebase {
	ctx := context.Background()
	opts := option.WithCredentialsFile(fbConfig.JSONFilePath)

	app, err := firebase.NewApp(ctx, nil, opts)
	if err != nil {
		panic(err.Error())
	}

	auth, err := app.Auth(ctx)
	if err != nil {
		panic(err.Error())
	}

	firestore, err := app.Firestore(ctx)
	if err != nil {
		panic(err.Error())
	}

	return &FirebaseImpl{
		auth:      auth,
		firestore: firestore,
		webAPIKey: fbConfig.WebAPIKey,
	}
}

func (d *FirebaseImpl) VerifyIDToken(idToken string) (string, error) {
	token, err := d.auth.VerifyIDToken(context.Background(), idToken)
	if err != nil {
		switch {
		case strings.Contains(err.Error(), "ID token has expired"):
			return "", entity.ErrFirebaseExpiredToken
		case strings.Contains(err.Error(), "failed to verify token signature"): // token文字列が不正
			return "", entity.ErrFirebaseFailedToVerify
		case strings.Contains(err.Error(), "ID token has invalid"): // tokenが不正（異なるfirebaseプロジェクトなど）
			return "", entity.ErrFirebaseInvalidToken
		case strings.Contains(err.Error(), "ID token issued at future timestamp"):
			return "", entity.ErrFirebaseFutureIssued
		default:
			fmt.Println(err)
			return "", entity.ErrServerError
		}
	}
	return token.UID, nil
}

func (d *FirebaseImpl) GetCustomToken(uid string) (string, error) {
	token, err := d.auth.CustomToken(context.Background(), uid)
	if err != nil {
		return "", entity.ErrServerError
	}
	return token, nil
}

func (d *FirebaseImpl) GetIDToken(token string) (string, error) {

	url := "https://www.googleapis.com/identitytoolkit/v3/relyingparty/verifyCustomToken?key=" + d.webAPIKey
	input := &struct {
		Token             string `json:"token"`
		ReturnSecureToken bool   `json:"returnSecureToken"`
	}{
		Token:             token,
		ReturnSecureToken: true,
	}
	params, err := json.Marshal(input)
	if err != nil {
		return "", errors.Wrap(entity.ErrServerError, err.Error())
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(params))
	if err != nil {
		return "", errors.Wrap(entity.ErrServerError, err.Error())
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", errors.Wrap(entity.ErrServerError, err.Error())
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", errors.Wrap(entity.ErrServerError, err.Error())
	}

	if resp.StatusCode != http.StatusOK {
		return "", errors.Wrap(entity.ErrServerError, "status code is wrong")
	}

	output := &struct {
		IDToken string `json:"idToken"`
	}{}

	err = json.Unmarshal(body, &output)
	if err != nil {
		return "", errors.Wrap(entity.ErrServerError, err.Error())
	}

	return output.IDToken, nil
}

func (d *FirebaseImpl) GetPhoneNumber(uid string) (string, error) {
	record, err := d.auth.GetUser(context.Background(), uid)
	if err != nil {
		return "", fmt.Errorf("%s:%w", err.Error(), entity.ErrServerError)
	}
	return record.PhoneNumber, nil
}

func (d *FirebaseImpl) Set(doc string, data map[string]interface{}) error {
	_, err := d.firestore.
		Collection("pitacoin").
		Doc(doc).
		Set(context.Background(), data, firestore.MergeAll)
	if err != nil {
		return errors.Wrap(entity.ErrServerError, err.Error())
	}
	return nil
}

// https://qiita.com/nisitanisubaru/items/3ff4e0b08b20700f408c
func (d *FirebaseImpl) CreateUser(email, password string) (string, error) {
	user := (&auth.UserToCreate{}).
		Email(email).
		EmailVerified(false).
		Password(password).
		Disabled(false)

	u, err := d.auth.CreateUser(context.Background(), user)
	if err != nil {
		fmt.Println("email: ", email)
		return "", fmt.Errorf("%s:%w", strings.Split(err.Error(), ""), entity.ErrFirebaseEmailExists)
	}

	return u.UID, err
}

// https://qiita.com/nisitanisubaru/items/3ff4e0b08b20700f408c
func (d *FirebaseImpl) UpdateEmail(email, uid string) error {
	user := (&auth.UserToUpdate{}).
		Email(email).
		EmailVerified(true)

	_, err := d.auth.UpdateUser(context.Background(), uid, user)
	if err != nil {
		return fmt.Errorf("%s:%w", strings.Split(err.Error(), ""), entity.ErrFirebaseEmailExists)
	}

	return err
}

func (d *FirebaseImpl) UpdatePassword(password, uid string) error {
	user := (&auth.UserToUpdate{}).Password(password)

	_, err := d.auth.UpdateUser(context.Background(), uid, user)
	if err != nil {
		return fmt.Errorf("%s:%w", strings.Split(err.Error(), ""), entity.ErrFirebaseEmailExists)
	}

	return err
}
