package repository

import (
	"fmt"
	"time"

	"github.com/hidenari-yuda/paychan-server/domain/entity"
	"github.com/hidenari-yuda/paychan-server/domain/utility"
	"github.com/hidenari-yuda/paychan-server/interfaces"
	"github.com/hidenari-yuda/paychan-server/usecase"
)

type LineMessageRepositoryImpl struct {
	Name     string
	executer interfaces.SQLExecuter
}

// Create(param *entity.LineMessage) error
// GetById(id uint) (*entity.LineMessage, error)
// GetListByUserId(userId uint) ([]*entity.LineMessage, error)
// GetListByLineUserId(lineUserId string) ([]*entity.LineMessage, error)

func NewLineMessageRepositoryImpl(ex interfaces.SQLExecuter) usecase.LineMessageRepository {
	return &LineMessageRepositoryImpl{
		Name:     "LineMessageRepository",
		executer: ex,
	}
}

func (r *LineMessageRepositoryImpl) Create(param *entity.LineMessage) error {
	_, err := r.executer.Exec(
		"SignUp",
		`INSERT INTO payment_methods (
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

func (r *LineMessageRepositoryImpl) GetById(id uint) (*entity.LineMessage, error) {
	var (
		LineMessage entity.LineMessage
	)
	err := r.executer.Get(
		"GetByFirebaseId",
		&LineMessage,
		"SELECT * FROM payment_methods WHERE id = ?",
		id,
	)

	if err != nil {
		err = fmt.Errorf("failed to get LineMessage by firebase id: %w", err)
		fmt.Println(err)
		return nil, err
	}

	return &LineMessage, nil
}

func (r *LineMessageRepositoryImpl) GetListByLineUserId(lineUserId string) ([]*entity.LineMessage, error) {
	var (
		LineMessageList []*entity.LineMessage
	)
	err := r.executer.Select(
		"GetByLineLineMessageId",
		&LineMessageList, `
		SELECT * 
		FROM payment_methods 
		WHERE user_id = (
			SELECT id
			FROM users
			WHERE line_user_id = ?
		)`,
		lineUserId,
	)

	if err != nil {
		err = fmt.Errorf("failed to get LineMessage by line LineMessage id: %w", err)
		fmt.Println(err)
		return nil, err
	}

	return LineMessageList, nil
}
