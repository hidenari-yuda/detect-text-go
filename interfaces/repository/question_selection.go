package repository

import (
	"fmt"
	"time"

	"github.com/hidenari-yuda/paychan-server/domain/entity"
	"github.com/hidenari-yuda/paychan-server/domain/utility"
	"github.com/hidenari-yuda/paychan-server/interfaces"
	"github.com/hidenari-yuda/paychan-server/usecase"
)

type QuestionSelectionRepositoryImpl struct {
	Name     string
	executer interfaces.SQLExecuter
}

// Create(param *entity.QuestionSelection) error
// GetById(id int) (*entity.QuestionSelection, error)
// GetListByUserId(userId int) ([]*entity.QuestionSelection, error)
// GetListByLineUserId(lineUserId string) ([]*entity.QuestionSelection, error)

func NewQuestionSelectionRepositoryImpl(ex interfaces.SQLExecuter) usecase.QuestionSelectionRepository {
	return &QuestionSelectionRepositoryImpl{
		Name:     "QuestionSelectionRepository",
		executer: ex,
	}
}

func (r *QuestionSelectionRepositoryImpl) Create(param *entity.QuestionSelection) error {
	_, err := r.executer.Exec(
		r.Name+"Create",
		`INSERT INTO question_selections (
			uuid,
			question_id,
			selection,
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
		param.QuestionId,
		param.Selection,
		time.Now(),
		time.Now(),
	)

	if err != nil {
		return err
	}

	return nil
}

// update
// func (r *QuestionSelectionRepositoryImpl) Update(param *entity.QuestionSelection) error {
// 	_, err := r.executer.Exec(
// 		"Update",
// 		`UPDATE question_selections SET
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

func (r *QuestionSelectionRepositoryImpl) GetListByQuestionId(questionId int) ([]*entity.QuestionSelection, error) {
	var (
		QuestionSelectionList []*entity.QuestionSelection
	)
	err := r.executer.Select(
		r.Name+"GetListByQuestionId",
		&QuestionSelectionList,
		"SELECT * FROM question_selections WHERE question_id = ?",
		questionId,
	)

	if err != nil {
		err = fmt.Errorf("failed to get QuestionSelection by id: %w", err)
		fmt.Println(err)
		return nil, err
	}

	return QuestionSelectionList, nil
}

func (r *QuestionSelectionRepositoryImpl) GetListByLineUserId(lineUserId string) ([]*entity.QuestionSelection, error) {
	var (
		QuestionSelectionList []*entity.QuestionSelection
	)
	err := r.executer.Select(
		r.Name+"GetListByLineUserId",
		&QuestionSelectionList, `
		SELECT * 
		FROM question_selections 
		WHERE line_user_id = ?
		ORDER BY id DESC
		`,
		lineUserId,
	)

	if err != nil {
		err = fmt.Errorf("failed to get QuestionSelection by line QuestionSelection id: %w", err)
		fmt.Println(err)
		return nil, err
	}

	return QuestionSelectionList, nil
}

//getByQuestionSelection
func (r *QuestionSelectionRepositoryImpl) GetListByQuestionType(questionType int) ([]*entity.QuestionSelection, error) {
	var (
		QuestionSelectionList []*entity.QuestionSelection
	)
	err := r.executer.Select(
		r.Name+"GetListByQuestionSelectionType",
		&QuestionSelectionList, `
		SELECT * 
		FROM question_selections 
		WHERE question_id IN (SELECT id FROM questions WHERE question_type = ?)
		ORDER BY id DESC
		`,
		questionType,
	)

	if err != nil {
		err = fmt.Errorf("failed to get QuestionSelection by QuestionSelection type: %w", err)
		fmt.Println(err)
		return nil, err
	}

	return QuestionSelectionList, nil
}
