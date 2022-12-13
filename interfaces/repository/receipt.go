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
			receipt_picture_id,
			store_name,
			total_price,
			purchased_at,
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
		)`,
		utility.CreateUUID(),
		param.ReceiptPictureId,
		param.StoreName,
		param.TotalPrice,
		param.PurchasedAt,
		time.Now(),
		time.Now(),
	)

	if err != nil {
		return err
	}

	return nil
}

// update
func (r *ReceiptRepositoryImpl) Update(param *entity.Receipt) error {
	_, err := r.executer.Exec(
		"Update",
		`UPDATE Receipts SET
			store_name = ?,
			total_price = ?,
			purchased_at = ?,
			updated_at = ?
			WHERE id = ?`,
		param.StoreName,
		param.TotalPrice,
		param.PurchasedAt,
		time.Now(),
		param.Id,
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
