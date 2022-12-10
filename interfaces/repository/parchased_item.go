package repository

import (
	"fmt"
	"time"

	"github.com/hidenari-yuda/paychan-server/domain/entity"
	"github.com/hidenari-yuda/paychan-server/domain/utility"
	"github.com/hidenari-yuda/paychan-server/interfaces"
	"github.com/hidenari-yuda/paychan-server/usecase"
)

type ParchasedItemRepositoryImpl struct {
	Name     string
	executer interfaces.SQLExecuter
}

// Create(param *entity.ParchasedItem) error
// GetById(id uint) (*entity.ParchasedItem, error)
// GetListByUserId(userId uint) ([]*entity.ParchasedItem, error)
// GetListByLineUserId(lineUserId string) ([]*entity.ParchasedItem, error)

func NewParchasedItemRepositoryImpl(ex interfaces.SQLExecuter) usecase.ParchasedItemRepository {
	return &ParchasedItemRepositoryImpl{
		Name:     "ParchasedItemRepository",
		executer: ex,
	}
}

func (r *ParchasedItemRepositoryImpl) Create(param *entity.ParchasedItem) error {
	_, err := r.executer.Exec(
		"SignUp",
		`INSERT INTO parchased_items (
			uuid,
			receipt_id,
			name,
			price,
			number,
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
		param.ReceiptId,
		param.Name,
		param.Price,
		param.Number,
		time.Now(),
		time.Now(),
	)

	if err != nil {
		return err
	}

	return nil
}

func (r *ParchasedItemRepositoryImpl) GetById(id uint) (*entity.ParchasedItem, error) {
	var (
		ParchasedItem entity.ParchasedItem
	)
	err := r.executer.Get(
		"GetByFirebaseId",
		&ParchasedItem,
		"SELECT * FROM parchased_items WHERE id = ?",
		id,
	)

	if err != nil {
		err = fmt.Errorf("failed to get ParchasedItem by firebase id: %w", err)
		fmt.Println(err)
		return nil, err
	}

	return &ParchasedItem, nil
}

func (r *ParchasedItemRepositoryImpl) GetListByReceiptId(receiptId uint) ([]*entity.ParchasedItem, error) {
	var (
		ParchasedItemList []*entity.ParchasedItem
	)
	err := r.executer.Select(
		"GetByLineParchasedItemId",
		&ParchasedItemList, `
		SELECT * 
		FROM parchased_items 
		WHERE user_id = (
			SELECT id
			FROM users
			WHERE receipt_id = ?
		)`,
		receiptId,
	)

	if err != nil {
		err = fmt.Errorf("failed to get ParchasedItem by line ParchasedItem id: %w", err)
		fmt.Println(err)
		return nil, err
	}

	return ParchasedItemList, nil
}

func (r *ParchasedItemRepositoryImpl) GetListByLineUserId(lineUserId string) ([]*entity.ParchasedItem, error) {
	var (
		ParchasedItemList []*entity.ParchasedItem
	)
	err := r.executer.Select(
		"GetByLineParchasedItemId",
		&ParchasedItemList, `
		SELECT * 
		FROM parchased_items 
		WHERE receipt_id IN (
			SELECT id
			FROM receipts
			WHERE line_user_id = ?
	)`,
		lineUserId,
	)

	if err != nil {
		err = fmt.Errorf("failed to get ParchasedItem by line ParchasedItem id: %w", err)
		fmt.Println(err)
		return nil, err
	}

	return ParchasedItemList, nil
}
