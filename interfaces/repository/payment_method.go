package repository

import (
	"fmt"
	"time"

	"github.com/hidenari-yuda/paychan-server/domain/entity"
	"github.com/hidenari-yuda/paychan-server/domain/utility"
	"github.com/hidenari-yuda/paychan-server/interfaces"
	"github.com/hidenari-yuda/paychan-server/usecase"
)

type PaymentMethodRepositoryImpl struct {
	Name     string
	executer interfaces.SQLExecuter
}

// Create(param *entity.PaymentMethod) error
// GetById(id int) (*entity.PaymentMethod, error)
// GetListByUserId(userId int) ([]*entity.PaymentMethod, error)
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
			user_id,
			payment_service,
			created_at,
			updated_at
			) VALUES (
				?,
				?,
				?, 
				?,
				?
		)`,
		utility.CreateUUID(),
		param.UserId,
		param.PaymentService,
		time.Now(),
		time.Now(),
	)

	if err != nil {
		return err
	}

	return nil
}

// update
func (r *PaymentMethodRepositoryImpl) Update(param *entity.PaymentMethod) error {
	_, err := r.executer.Exec(
		"Update",
		`UPDATE payment_methods SET
			payment_service = ?,
			updated_at = ?
			WHERE id = ?`,
		param.PaymentService,
		time.Now(),
		param.Id,
	)

	if err != nil {
		return err
	}

	return nil
}

func (r *PaymentMethodRepositoryImpl) GetById(id int) (*entity.PaymentMethod, error) {
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
