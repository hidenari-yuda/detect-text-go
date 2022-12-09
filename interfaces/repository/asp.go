package repository

import (
	"fmt"
	"time"

	"github.com/hidenari-yuda/paychan/domain/entity"
	"github.com/hidenari-yuda/paychan/domain/utility"
	"github.com/hidenari-yuda/paychan/interfaces"
	"github.com/hidenari-yuda/paychan/usecase"
)

type AspRepositoryImpl struct {
	Name     string
	executer interfaces.SQLExecuter
}

// Create(param *entity.Asp) error
// GetById(id uint) (*entity.Asp, error)
// GetListByUserId(userId uint) ([]*entity.Asp, error)
// GetListByLineUserId(lineUserId string) ([]*entity.Asp, error)

func NewAspRepositoryImpl(ex interfaces.SQLExecuter) usecase.AspRepository {
	return &AspRepositoryImpl{
		Name:     "AspRepository",
		executer: ex,
	}
}

func (r *AspRepositoryImpl) Create(param *entity.Asp) error {
	_, err := r.executer.Exec(
		"SignUp",
		`INSERT INTO Asps (
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

func (r *AspRepositoryImpl) GetById(id uint) (*entity.Asp, error) {
	var (
		Asp entity.Asp
	)
	err := r.executer.Get(
		"GetByFirebaseId",
		&Asp,
		"SELECT * FROM Asps WHERE id = ?",
		id,
	)

	if err != nil {
		err = fmt.Errorf("failed to get Asp by firebase id: %w", err)
		fmt.Println(err)
		return nil, err
	}

	return &Asp, nil
}

func (r *AspRepositoryImpl) GetListByLineUserId(lineUserId string) ([]*entity.Asp, error) {
	var (
		AspList []*entity.Asp
	)
	err := r.executer.Select(
		"GetByLineAspId",
		&AspList, `
		SELECT * 
		FROM Asps 
		WHERE user_id = (
			SELECT id
			FROM users
			WHERE line_user_id = ?
		)`,
		lineUserId,
	)

	if err != nil {
		err = fmt.Errorf("failed to get Asp by line Asp id: %w", err)
		fmt.Println(err)
		return nil, err
	}

	return AspList, nil
}
