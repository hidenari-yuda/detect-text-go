package repository

import (
	"fmt"
	"time"

	"github.com/hidenari-yuda/paychan-server/domain/entity"
	"github.com/hidenari-yuda/paychan-server/domain/utility"
	"github.com/hidenari-yuda/paychan-server/interfaces"
	"github.com/hidenari-yuda/paychan-server/usecase"
)

type ReceiptPictureRepositoryImpl struct {
	Name     string
	executer interfaces.SQLExecuter
}

// Create(param *entity.ReceiptPicture) error
// GetById(id uint) (*entity.ReceiptPicture, error)
// GetListByUserId(userId uint) ([]*entity.ReceiptPicture, error)
// GetListByLineUserId(lineUserId string) ([]*entity.ReceiptPicture, error)

func NewReceiptPictureRepositoryImpl(ex interfaces.SQLExecuter) usecase.ReceiptPictureRepository {
	return &ReceiptPictureRepositoryImpl{
		Name:     "ReceiptPictureRepository",
		executer: ex,
	}
}

func (r *ReceiptPictureRepositoryImpl) Create(param *entity.ReceiptPicture) error {
	_, err := r.executer.Exec(
		"SignUp",
		`INSERT INTO ReceiptPictures (
			uuid,
			user_id,
			line_user_id,
			url,
			detected_text,
			service,
			payment_service,
			point,
			total_price,
			created_at,
			updated_at
			) VALUES (
				?,
				?,
				?, 
				?,
				?,
				?,
				?,
				?,
				?,
				?,
				?
		)`,
		utility.CreateUUID(),
		param.UserId,
		param.LineUserId,
		param.Url,
		param.DetectedText,
		param.Service,
		param.PaymentService,
		param.Point,
		param.TotalPrice,
		time.Now(),
		time.Now(),
	)

	if err != nil {
		return err
	}

	return nil
}

func (r *ReceiptPictureRepositoryImpl) GetById(id uint) (*entity.ReceiptPicture, error) {
	var (
		ReceiptPicture entity.ReceiptPicture
	)
	err := r.executer.Get(
		"GetByFirebaseId",
		&ReceiptPicture,
		"SELECT * FROM ReceiptPictures WHERE id = ?",
		id,
	)

	if err != nil {
		err = fmt.Errorf("failed to get ReceiptPicture by firebase id: %w", err)
		fmt.Println(err)
		return nil, err
	}

	return &ReceiptPicture, nil
}

func (r *ReceiptPictureRepositoryImpl) GetListByLineUserId(lineUserId string) ([]*entity.ReceiptPicture, error) {
	var (
		ReceiptPictureList []*entity.ReceiptPicture
	)
	err := r.executer.Select(
		"GetByLineReceiptPictureId",
		&ReceiptPictureList, `
		SELECT * 
		FROM receiptPictures 
		WHERE user_id = (
			SELECT id
			FROM users
			WHERE line_user_id = ?
		)`,
		lineUserId,
	)

	if err != nil {
		err = fmt.Errorf("failed to get ReceiptPicture by line ReceiptPicture id: %w", err)
		fmt.Println(err)
		return nil, err
	}

	return ReceiptPictureList, nil
}

// 今日登録されたレシートのリストを取得する
func (r *ReceiptPictureRepositoryImpl) GetListByToday(lineUserId string) ([]*entity.ReceiptPicture, error) {
	var (
		ReceiptPictureList []*entity.ReceiptPicture
	)
	err := r.executer.Select(
		"GetByLineReceiptPictureId",
		&ReceiptPictureList, `
		SELECT * 
		FROM receiptPictures 
		WHERE user_id = (
			SELECT id
			FROM users
			WHERE line_user_id = ?
		)
		AND DATE(created_at) = CURDATE()`,
		lineUserId,
	)

	if err != nil {
		err = fmt.Errorf("failed to get ReceiptPicture by line ReceiptPicture id: %w", err)
		fmt.Println(err)
		return nil, err
	}

	return ReceiptPictureList, nil
}
