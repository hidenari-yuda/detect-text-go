package repository

import (
	"fmt"
	"time"

	"github.com/hidenari-yuda/paychan-server/domain/entity"
	"github.com/hidenari-yuda/paychan-server/domain/utility"
	"github.com/hidenari-yuda/paychan-server/interfaces"
	"github.com/hidenari-yuda/paychan-server/usecase"
)

type ReceiptRepositoryImpl struct {
	Name     string
	executer interfaces.SQLExecuter
}

// Create(param *entity.Receipt) error
// GetById(id uint) (*entity.Receipt, error)
// GetListByUserId(userId uint) ([]*entity.Receipt, error)
// GetListByLineUserId(lineUserId string) ([]*entity.Receipt, error)

func NewReceiptRepositoryImpl(ex interfaces.SQLExecuter) usecase.ReceiptRepository {
	return &ReceiptRepositoryImpl{
		Name:     "ReceiptRepository",
		executer: ex,
	}
}

func (r *ReceiptRepositoryImpl) Create(param *entity.Receipt) error {
	_, err := r.executer.Exec(
		"SignUp",
		`INSERT INTO Receipts (
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

func (r *ReceiptRepositoryImpl) GetById(id uint) (*entity.Receipt, error) {
	var (
		Receipt entity.Receipt
	)
	err := r.executer.Get(
		"GetByFirebaseId",
		&Receipt,
		"SELECT * FROM Receipts WHERE id = ?",
		id,
	)

	if err != nil {
		err = fmt.Errorf("failed to get Receipt by firebase id: %w", err)
		fmt.Println(err)
		return nil, err
	}

	return &Receipt, nil
}

func (r *ReceiptRepositoryImpl) GetListByLineUserId(lineUserId string) ([]*entity.Receipt, error) {
	var (
		ReceiptList []*entity.Receipt
	)
	err := r.executer.Select(
		"GetByLineReceiptId",
		&ReceiptList, `
		SELECT * 
		FROM receipts 
		WHERE user_id = (
			SELECT id
			FROM users
			WHERE line_user_id = ?
		)`,
		lineUserId,
	)

	if err != nil {
		err = fmt.Errorf("failed to get Receipt by line Receipt id: %w", err)
		fmt.Println(err)
		return nil, err
	}

	return ReceiptList, nil
}

// 今日登録されたレシートのリストを取得する
func (r *ReceiptRepositoryImpl) GetListByToday(lineUserId string) ([]*entity.Receipt, error) {
	var (
		ReceiptList []*entity.Receipt
	)
	err := r.executer.Select(
		"GetByLineReceiptId",
		&ReceiptList, `
		SELECT * 
		FROM receipts 
		WHERE user_id = (
			SELECT id
			FROM users
			WHERE line_user_id = ?
		)
		AND DATE(created_at) = CURDATE()`,
		lineUserId,
	)

	if err != nil {
		err = fmt.Errorf("failed to get Receipt by line Receipt id: %w", err)
		fmt.Println(err)
		return nil, err
	}

	return ReceiptList, nil
}
