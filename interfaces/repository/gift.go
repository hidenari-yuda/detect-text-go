package repository

import (
	"fmt"
	"time"

	"github.com/hidenari-yuda/umerun-resume/domain/entity"
	"github.com/hidenari-yuda/umerun-resume/domain/utility"
	"github.com/hidenari-yuda/umerun-resume/interfaces"
	"github.com/hidenari-yuda/umerun-resume/usecase"
)

type GiftRepositoryImpl struct {
	Name     string
	executer interfaces.SQLExecuter
}

// Create(param *entity.Gift) error
// GetById(id uint) (*entity.Gift, error)
// GetListByUserId(userId uint) ([]*entity.Gift, error)
// GetListByLineUserId(lineUserId string) ([]*entity.Gift, error)

func NewGiftRepositoryImpl(ex interfaces.SQLExecuter) usecase.GiftRepository {
	return &GiftRepositoryImpl{
		Name:     "GiftRepository",
		executer: ex,
	}
}

func (r *GiftRepositoryImpl) Create(param *entity.Gift) error {
	_, err := r.executer.Exec(
		"SignUp",
		`INSERT INTO Gifts (
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

func (r *GiftRepositoryImpl) GetById(id uint) (*entity.Gift, error) {
	var (
		Gift entity.Gift
	)
	err := r.executer.Get(
		"GetByFirebaseId",
		&Gift,
		"SELECT * FROM Gifts WHERE id = ?",
		id,
	)

	if err != nil {
		err = fmt.Errorf("failed to get Gift by firebase id: %w", err)
		fmt.Println(err)
		return nil, err
	}

	return &Gift, nil
}

func (r *GiftRepositoryImpl) GetListByReceiptId(receiptId uint) ([]*entity.Gift, error) {
	var (
		GiftList []*entity.Gift
	)
	err := r.executer.Select(
		"GetByLineGiftId",
		&GiftList, `
		SELECT * 
		FROM Gifts 
		WHERE user_id = (
			SELECT id
			FROM users
			WHERE line_user_id = ?
		)`,
		receiptId,
	)

	if err != nil {
		err = fmt.Errorf("failed to get Gift by line Gift id: %w", err)
		fmt.Println(err)
		return nil, err
	}

	return GiftList, nil
}

func (r *GiftRepositoryImpl) GetListByLineUserId(lineUserId string) ([]*entity.Gift, error) {
	var (
		GiftList []*entity.Gift
	)
	err := r.executer.Select(
		"GetByLineGiftId",
		&GiftList, `
		SELECT * 
		FROM Gifts 
		WHERE user_id = (
			SELECT id
			FROM users
			WHERE line_user_id = ?
		)`,
		lineUserId,
	)

	if err != nil {
		err = fmt.Errorf("failed to get Gift by line Gift id: %w", err)
		fmt.Println(err)
		return nil, err
	}

	return GiftList, nil
}
