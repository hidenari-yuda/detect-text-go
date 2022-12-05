package repository

import (
	"fmt"
	"time"

	"github.com/hidenari-yuda/umerun-resume/domain/entity"
	"github.com/hidenari-yuda/umerun-resume/domain/utility"
	"github.com/hidenari-yuda/umerun-resume/interfaces"
	"github.com/hidenari-yuda/umerun-resume/usecase"
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
		"SignUp",
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

func (r *PresentRepositoryImpl) GetById(id uint) (*entity.Present, error) {
	var (
		present entity.Present
	)
	err := r.executer.Get(
		"GetByFirebaseId",
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

func (r *PresentRepositoryImpl) GetByPrice(price uint) ([]*entity.Present, error) {
	var (
		presentList []*entity.Present
	)
	err := r.executer.Select(
		"GetByPrice",
		&presentList, `
		SELECT * 
		FROM Presents 
		WHERE price >= ? && price <= ?
		ORDER BY price ASC
		`,
		price, price+20,
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

func (r *PresentRepositoryImpl) GetListByReceiptId(receiptId uint) ([]*entity.Present, error) {
	var (
		presentList []*entity.Present
	)
	err := r.executer.Select(
		"GetByLinePresentId",
		&presentList, `
		SELECT * 
		FROM Presents 
		WHERE user_id = (
			SELECT id
			FROM users
			WHERE line_user_id = ?
		)`,
		receiptId,
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
		"GetByLinePresentId",
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
