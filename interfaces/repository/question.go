package repository

import (
	"fmt"
	"time"

	"github.com/hidenari-yuda/paychan-server/domain/entity"
	"github.com/hidenari-yuda/paychan-server/domain/utility"
	"github.com/hidenari-yuda/paychan-server/interfaces"
	"github.com/hidenari-yuda/paychan-server/usecase"
)

type QuestionRepositoryImpl struct {
	Name     string
	executer interfaces.SQLExecuter
}

// Create(param *entity.Question) error
// GetById(id int) (*entity.Question, error)
// GetListByUserId(userId int) ([]*entity.Question, error)
// GetListByLineUserId(lineUserId string) ([]*entity.Question, error)

func NewQuestionRepositoryImpl(ex interfaces.SQLExecuter) usecase.QuestionRepository {
	return &QuestionRepositoryImpl{
		Name:     "QuestionRepository",
		executer: ex,
	}
}

func (r *QuestionRepositoryImpl) Create(param *entity.Question) error {
	_, err := r.executer.Exec(
		r.Name+"Create",
		`INSERT INTO questions (
			uuid,
			user_id,
			line_user_id,
			receipt_picture_id,
			question,
			text,
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
				?
		)`,
		utility.CreateUUID(),
		param.UserId,
		param.LineUserId,
		param.ReceiptPictureId,
		param.Question,
		param.Text,
		time.Now(),
		time.Now(),
	)

	if err != nil {
		return err
	}

	return nil
}

// update
// func (r *QuestionRepositoryImpl) Update(param *entity.Question) error {
// 	_, err := r.executer.Exec(
// 		"Update",
// 		`UPDATE questions SET
// 			payment_service = ?,
// 			updated_at = ?
// 			WHERE id = ?`,
// 		param.PaymentService,
// 		time.Now(),
// 		param.Id,
// 	)

// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

func (r *QuestionRepositoryImpl) GetById(id int) (*entity.Question, error) {
	var (
		Question entity.Question
	)
	err := r.executer.Get(
		r.Name+"GetById",
		&Question,
		"SELECT * FROM questions WHERE id = ?",
		id,
	)

	if err != nil {
		err = fmt.Errorf("failed to get Question by id: %w", err)
		fmt.Println(err)
		return nil, err
	}

	return &Question, nil
}

func (r *QuestionRepositoryImpl) GetListByLineUserId(lineUserId string) ([]*entity.Question, error) {
	var (
		QuestionList []*entity.Question
	)
	err := r.executer.Select(
		r.Name+"GetListByLineUserId",
		&QuestionList, `
		SELECT * 
		FROM questions 
		WHERE line_user_id = ?
		ORDER BY id DESC
		`,
		lineUserId,
	)

	if err != nil {
		err = fmt.Errorf("failed to get Question by line Question id: %w", err)
		fmt.Println(err)
		return nil, err
	}

	return QuestionList, nil
}

//getByQuestion
func (r *QuestionRepositoryImpl) GetListByQuestionType(questionType int) ([]*entity.Question, error) {
	var (
		QuestionList []*entity.Question
	)
	err := r.executer.Select(
		r.Name+"GetListByQuestionType",
		&QuestionList, `
		SELECT * 
		FROM questions 
		WHERE question = ?
		ORDER BY id DESC
		`,
		questionType,
	)

	if err != nil {
		err = fmt.Errorf("failed to get Question by Question type: %w", err)
		fmt.Println(err)
		return nil, err
	}

	return QuestionList, nil
}
