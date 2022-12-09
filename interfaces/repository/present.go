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
// GetById(id uint) (*entity.Present, error)
// GetListByUserId(userId uint) ([]*entity.Present, error)
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
		`INSERT INTO Presents (
			uuid,
			firebase_id,
			name, 
			email, 
			password,
			created_at,
			updated_at
			) VALUES (
				?,
				?,
				?, 
				?,
				?,
				?,
				?
		)`,
		utility.CreateUUID(),
		"",
		"ゲスト",
		// param.Email,
		// param.Password,
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
		`UPDATE Presents
		SET
			user_id = ?,
			receipt_picture_id = ?,
			price = ?,
			payment_service = ?,
			updated_at = ?
		WHERE
			id = ?`,
		present.UserId,
		present.ReceiptPictureId,
		present.Price,
		present.PaymentService,
		time.Now(),
		present.Id,
	)

	if err != nil {
		return err
	}

	return nil
}

func (r *PresentRepositoryImpl) GetById(id uint) (*entity.Present, error) {
	var (
		present entity.Present
	)
	err := r.executer.Get(
		r.Name+"GetById",
		&present,
		"SELECT * FROM Presents WHERE id = ?",
		id,
	)

	if err != nil {
		err = fmt.Errorf("failed to get Present by firebase id: %w", err)
		fmt.Println(err)
		return nil, err
	}

	return &present, nil
}

func (r *PresentRepositoryImpl) GetByPriceAndService(present *entity.Present) ([]*entity.Present, error) {
	var (
		presentList []*entity.Present
	)
	err := r.executer.Select(
		r.Name+"GetByPriceAndService",
		&presentList, `
		SELECT * 
		FROM Presents 
		WHERE price >= ? 
		AND price <= ?
		AND payment_service = ?
		ORDER BY price ASC
		`,
		present.Price, present.Price+20, present.PaymentService,
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

func (r *PresentRepositoryImpl) GetListByReceiptPictureId(receiptPictureId uint) ([]*entity.Present, error) {
	var (
		presentList []*entity.Present
	)
	err := r.executer.Select(
		r.Name+"GetByReceiptPictureId",
		&presentList, `
		SELECT * 
		FROM Presents 
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
		FROM Presents 
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
		FROM Presents`,
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
		`DELETE FROM Presents
		WHERE expirary < ?`,
		time.Now(),
	)

	if err != nil {
		return err
	}

	return nil
}
