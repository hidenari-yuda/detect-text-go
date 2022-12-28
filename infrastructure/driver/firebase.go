package driver

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"cloud.google.com/go/firestore"
	"cloud.google.com/go/storage"
	gcStorage "cloud.google.com/go/storage"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"github.com/google/uuid"
	"github.com/hidenari-yuda/paychan-server/domain/config"
	"github.com/hidenari-yuda/paychan-server/domain/entity"
	"github.com/hidenari-yuda/paychan-server/usecase"
	"github.com/pkg/errors"
	gcContext "golang.org/x/net/context"
	"google.golang.org/api/option"
)

type FirebaseImpl struct {
	auth       *auth.Client
	firestore  *firestore.Client
	storage    *storage.Client
	client     *gcStorage.Client
	bucketName string
	webAPIKey  string
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

func (d *FirebaseImpl) UploadImage(content io.ReadCloser, messageId string) (string, error) {
	ctx := gcContext.Background()
	bkt := d.client.Bucket(d.bucketName)
	failureText := "画像の取得に失敗しました"

	_, err := bkt.Attrs(ctx)
	if err != nil {
		// TODO: Handle error.
		return failureText, err
	}

	fileName := fmt.Sprintf("ReceiptPicture/%s.jpeg", messageId)
	obj := bkt.Object(fileName)

	//create an id
	id := uuid.New()

	w := obj.NewWriter(ctx)

	//Set the attribute （ダウンロードトークン付与）
	w.ObjectAttrs.Metadata = map[string]string{"firebaseStorageDownloadTokens": id.String()}
	w.ObjectAttrs.ACL = []gcStorage.ACLRule{
		{
			Entity: gcStorage.AllUsers,
			Role:   gcStorage.RoleReader,
		},
	}

	data, _ := ioutil.ReadAll(content)
	w.Write(data)

	// Close, just like writing a file.
	if err := w.Close(); err != nil {
		// TODO: Handle error.
		return failureText, err
	}

	// urlを生成する処理
	// 形式: https://firebasestorage.googleapis.com/v0/b/プロジェクト名.appspot.com/o/ファイルパス?alt=xxxx&token=xxxxx
	fileURL := fmt.Sprintf(
		"https://firebasestorage.googleapis.com/v0/b/%s/o/%s?alt=media&token=%s",
		d.bucketName, url.QueryEscape(fileName), id,
	)

	return fileURL, nil
}
