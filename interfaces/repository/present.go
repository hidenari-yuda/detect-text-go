package repository

import (
	"fmt"
	"time"

	"github.com/hidenari-yuda/paychan-server/domain/entity"
	"github.com/hidenari-yuda/paychan-server/domain/utility"
	"github.com/hidenari-yuda/paychan-server/interfaces"
	"github.com/hidenari-yuda/paychan-server/usecase"
)

type PresentRepositoryImpl struct {
	Name     string
	executer interfaces.SQLExecuter
}

// Create(param *entity.Present) error
// GetById(id int) (*entity.Present, error)
// GetListByUserId(userId int) ([]*entity.Present, error)
// GetListByLineUserId(lineUserId string) ([]*entity.Present, error)

func NewPresentRepositoryImpl(ex interfaces.SQLExecuter) usecase.PresentRepository {
	return &PresentRepositoryImpl{
		Name:     "PresentRepository",
		executer: ex,
	}
}

func (r *PresentRepositoryImpl) Create(param *entity.Present) error {
	_, err := r.executer.Exec(
		r.Name+"Create",
		`INSERT INTO presents (
			uuid,
			user_id,
			line_user_id,
			receipt_picture_id,
			payment_service,
			point,
			url,
			expirary,
			used,
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
		param.ReceiptPictureId,
		param.PaymentService,
		param.Point,
		param.Url,
		param.Expirary,
		param.Used,
		time.Now(),
		time.Now(),
	)

	if err != nil {
		return err
	}

	return nil
}

// Update
func (r *PresentRepositoryImpl) Update(present *entity.Present) error {
	_, err := r.executer.Exec(
		r.Name+"Update",
		`UPDATE presents
		SET
			user_id = ?,
			line_user_id = ?,
			receipt_picture_id = ?,
			payment_service = ?,
			point = ?,
			url = ?,
			expirary = ?,
			used = ?,
			updated_at = ?
		WHERE
			id = ?`,
		present.UserId,
		present.LineUserId,
		present.ReceiptPictureId,
		present.PaymentService,
		present.Point,
		present.Url,
		present.Expirary,
		present.Used,
		time.Now(),
		present.Id,
	)

	if err != nil {
		return err
	}

	return nil
}

func (r *PresentRepositoryImpl) GetById(id int) (*entity.Present, error) {
	var (
		present entity.Present
	)
	err := r.executer.Get(
		r.Name+"GetById",
		&present,
		"SELECT * FROM presents WHERE id = ?",
		id,
	)

	if err != nil {
		err = fmt.Errorf("failed to get Present by firebase id: %w", err)
		fmt.Println(err)
		return nil, err
	}

	return &present, nil
}

func (r *PresentRepositoryImpl) GetByPointAndService(present *entity.Present) ([]*entity.Present, error) {
	var (
		presentList []*entity.Present
	)
	err := r.executer.Select(
		r.Name+"GetByPointAndService",
		&presentList, `
		SELECT * 
		FROM presents 
		WHERE Point >= ? 
		AND point <= ?
		AND payment_service = ?
		AND used = 0
		ORDER BY point ASC
		`,
		present.Point, present.Point+20, present.PaymentService,
	)

	if err != nil {
		err = fmt.Errorf("failed to get Present by firebase id: %w", err)
		fmt.Println(err)
		return nil, err
	}

	// if len(presentList) == 0 {
	// 	// 管理lineに通知
	// 	return nil, nil
	// }

	// present = *presentList[0]

	return presentList, nil
}

func (r *PresentRepositoryImpl) GetListByReceiptPictureId(receiptPictureId int) ([]*entity.Present, error) {
	var (
		presentList []*entity.Present
	)
	err := r.executer.Select(
		r.Name+"GetByReceiptPictureId",
		&presentList, `
		SELECT * 
		FROM presents 
		WHERE receipt_picture_id = ?`,
		receiptPictureId,
	)

	if err != nil {
		err = fmt.Errorf("failed to get Present by line Present id: %w", err)
		fmt.Println(err)
		return nil, err
	}

	return presentList, nil
}

func (r *PresentRepositoryImpl) GetListByLineUserId(lineUserId string) ([]*entity.Present, error) {
	var (
		presentList []*entity.Present
	)
	err := r.executer.Select(
		r.Name+"GetByLinePresentId",
		&presentList, `
		SELECT * 
		FROM presents 
		WHERE user_id = (
			SELECT id
			FROM users
			WHERE line_user_id = ?
		)`,
		lineUserId,
	)

	if err != nil {
		err = fmt.Errorf("failed to get Present by line Present id: %w", err)
		fmt.Println(err)
		return nil, err
	}

	return presentList, nil
}

// GetAll
func (r *PresentRepositoryImpl) GetAll() ([]*entity.Present, error) {
	var (
		presentList []*entity.Present
	)
	err := r.executer.Select(
		r.Name+"GetAll",
		&presentList, `
		SELECT * 
		FROM presents`,
	)

	if err != nil {
		err = fmt.Errorf("failed to get Present by line Present id: %w", err)
		fmt.Println(err)
		return nil, err
	}

	return presentList, nil
}

// DeleteByExpired
func (r *PresentRepositoryImpl) DeleteByExpired() error {
	_, err := r.executer.Exec(
		r.Name+"DeleteByExpired",
		`DELETE FROM presents
		WHERE expirary < ?`,
		time.Now(),
	)

	if err != nil {
		return err
	}

	return nil
}
