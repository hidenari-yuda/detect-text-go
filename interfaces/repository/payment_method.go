package repository

import (
	"fmt"
	"time"

	"github.com/hidenari-yuda/detect-text/domain/entity"
	"github.com/hidenari-yuda/detect-text/domain/utility"
	"github.com/hidenari-yuda/detect-text/interfaces"
	"github.com/hidenari-yuda/detect-text/usecase"
)

type PaymentMethodRepositoryImpl struct {
	Name     string
	executer interfaces.SQLExecuter
}

// Create(param *entity.PaymentMethod) error
// GetById(id uint) (*entity.PaymentMethod, error)
// GetListByUserId(userId uint) ([]*entity.PaymentMethod, error)
// GetListByLineUserId(lineUserId string) ([]*entity.PaymentMethod, error)

func NewPaymentMethodRepositoryImpl(ex interfaces.SQLExecuter) usecase.PaymentMethodRepository {
	return &PaymentMethodRepositoryImpl{
		Name:     "PaymentMethodRepository",
		executer: ex,
	}
}

func (r *PaymentMethodRepositoryImpl) Create(param *entity.PaymentMethod) error {
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

func (r *PaymentMethodRepositoryImpl) GetById(id uint) (*entity.PaymentMethod, error) {
	var (
		PaymentMethod entity.PaymentMethod
	)
	err := r.executer.Get(
		"GetByFirebaseId",
		&PaymentMethod,
		"SELECT * FROM payment_methods WHERE id = ?",
		id,
	)

	if err != nil {
		err = fmt.Errorf("failed to get PaymentMethod by firebase id: %w", err)
		fmt.Println(err)
		return nil, err
	}

	return &PaymentMethod, nil
}

func (r *PaymentMethodRepositoryImpl) GetListByLineUserId(lineUserId string) ([]*entity.PaymentMethod, error) {
	var (
		PaymentMethodList []*entity.PaymentMethod
	)
	err := r.executer.Select(
		"GetByLinePaymentMethodId",
		&PaymentMethodList, `
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
		err = fmt.Errorf("failed to get PaymentMethod by line PaymentMethod id: %w", err)
		fmt.Println(err)
		return nil, err
	}

	return PaymentMethodList, nil
}
