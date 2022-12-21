package repository

import (
	"fmt"
	"time"

	"github.com/hidenari-yuda/paychan-server/domain/entity"
	"github.com/hidenari-yuda/paychan-server/domain/utility"
	"github.com/hidenari-yuda/paychan-server/interfaces"
	"github.com/hidenari-yuda/paychan-server/usecase"
)

type UserRepositoryImpl struct {
	Name     string
	executer interfaces.SQLExecuter
}

func NewUserRepositoryImpl(ex interfaces.SQLExecuter) usecase.UserRepository {
	return &UserRepositoryImpl{
		Name:     "UserRepository",
		executer: ex,
	}
}

func (r *UserRepositoryImpl) SignUp(param *entity.SignUpParam) error {
	_, err := r.executer.Exec(
		r.Name+"SignUp",
		`INSERT INTO users (
			uuid,
			firebase_id,
			name, 
			email, 
			password,
			line_user_id,
			line_name,
			picture_url,
			status_message,
			language,
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
				?,
				?,
				?,
				?,
				?
		)`,
		utility.CreateUUID(),
		param.FirebaseId,
		param.Name,
		param.Email,
		param.Password,
		param.LineUserId,
		param.LineName,
		param.PictureUrl,
		param.StatusMessage,
		param.Language,
		time.Now(),
		time.Now(),
	)

	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepositoryImpl) SignIn(email, password string) (user *entity.User, err error) {
	err = r.executer.Get(
		r.Name+"SignIn",
		user,
		"SELECT * FROM users WHERE email = ? AND password = ?",
		email,
		password,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepositoryImpl) GetByFirebaseId(firebaseId string) (*entity.User, error) {
	var (
		user entity.User
	)
	err := r.executer.Get(
		r.Name+"GetByFirebaseId",
		&user,
		"SELECT * FROM users WHERE firebase_id = ?",
		firebaseId,
	)

	if err != nil {
		err = fmt.Errorf("failed to get user by firebase id: %w", err)
		fmt.Println(err)
		return nil, err
	}

	return &user, nil
}

func (r *UserRepositoryImpl) GetByLineUserId(lineUserId string) (*entity.User, error) {
	var (
		user entity.User
	)
	err := r.executer.Get(
		r.Name+"GetByLineUserId",
		&user,
		"SELECT * FROM users WHERE line_user_id = ?",
		lineUserId,
	)

	if err != nil {
		err = fmt.Errorf("failed to get user by line user id: %w", err)
		fmt.Println(err)
		return nil, err
	}

	return &user, nil
}

// getAll
func (r *UserRepositoryImpl) GetAll() ([]*entity.User, error) {
	var (
		users []*entity.User
	)
	err := r.executer.Select(
		r.Name+"GetAll",
		&users,
		"SELECT * FROM users ORDER BY id DESC",
	)

	if err != nil {
		err = fmt.Errorf("failed to get all users: %w", err)
		fmt.Println(err)
		return nil, err
	}

	return users, nil
}

func (r *UserRepositoryImpl) Update(user *entity.User) error {
	_, err := r.executer.Exec(
		r.Name+"Update",
		`UPDATE users SET
			name = ?,
			email = ?,
			password = ?,
			line_user_id = ?,
			line_name = ?,
			picture_url = ?,
			status_message = ?,
			language = ?,
			prefecture = ?,
			age = ?,
			gender = ?,
			occupation = ?,
			married = ?,
			annual_income = ?,
			updated_at = ?
		WHERE line_user_id = ?`,
		user.Name,
		user.Email,
		user.Password,
		user.LineUserId,
		user.LineName,
		user.PictureUrl,
		user.StatusMessage,
		user.Language,
		user.Prefecture,
		user.Age,
		user.Gender,
		user.Occupation,
		user.Married,
		user.AnnualIncome,
		time.Now(),
		user.LineUserId,
	)

	if err != nil {
		err = fmt.Errorf("failed to update user: %w", err)
		fmt.Println(err)
		return err
	}

	return nil
}

func (r *UserRepositoryImpl) UpdateColumnStr(lineUserId, column, value string) error {
	_, err := r.executer.Exec(
		r.Name+"UpdateColumn",
		"UPDATE users SET "+column+" = ? WHERE line_user_id = ?",
		value,
		lineUserId,
	)

	if err != nil {
		err = fmt.Errorf("failed to update user column: %w", err)
		fmt.Println(err)
		return err
	}

	return nil
}

func (r *UserRepositoryImpl) UpdateColumnInt(lineUserId, column string, value int) error {
	_, err := r.executer.Exec(
		r.Name+"UpdateColumn",
		"UPDATE users SET "+column+" = ? WHERE line_user_id = ?",
		value,
		lineUserId,
	)

	if err != nil {
		err = fmt.Errorf("failed to update user column: %w", err)
		fmt.Println(err)
		return err
	}

	return nil
}
