package repository

import (
	"fmt"
	"time"

	"github.com/hidenari-yuda/paychan-server/domain/entity"
	"github.com/hidenari-yuda/paychan-server/domain/utility"
	"github.com/hidenari-yuda/paychan-server/interfaces"
	"github.com/hidenari-yuda/paychan-server/usecase"
)

type AdRepositoryImpl struct {
	Name     string
	executer interfaces.SQLExecuter
}

// Create(param *entity.Ad) error
// GetById(id int) (*entity.Ad, error)
// GetListByUserId(userId int) ([]*entity.Ad, error)
// GetListByLineUserId(lineUserId string) ([]*entity.Ad, error)

func NewAdRepositoryImpl(ex interfaces.SQLExecuter) usecase.AdRepository {
	return &AdRepositoryImpl{
		Name:     "AdRepository",
		executer: ex,
	}
}

func (r *AdRepositoryImpl) Create(param *entity.Ad) error {
	_, err := r.executer.Exec(
		r.Name+"Create",
		`INSERT INTO ads (
			uuid,
			service,
			url,
			image_url,
			price,
			title,
			description,
			impression,
			click,
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
				?
		)`,
		utility.CreateUUID(),
		param.Service,
		param.Url,
		param.ImageUrl,
		param.Price,
		param.Title,
		param.Description,
		0,
		0,
		time.Now(),
		time.Now(),
	)

	if err != nil {
		return err
	}

	return nil
}

// update
func (r *AdRepositoryImpl) Update(param *entity.Ad) error {
	_, err := r.executer.Exec(
		r.Name+"Update",
		`UPDATE ads SET
			service = ?,
			url = ?,
			image_url = ?,
			price = ?,
			title = ?,
			description = ?,
			impression = ?,
			click = ?,
			updated_at = ?
			WHERE id = ?
		`,
		param.Service,
		param.Url,
		param.ImageUrl,
		param.Price,
		param.Title,
		param.Description,
		param.Impression,
		param.Click,
		time.Now(),
		param.Id,
	)

	if err != nil {
		return err
	}

	return nil
}

// updateImpression
func (r *AdRepositoryImpl) UpdateImpression(param *entity.Ad) error {
	_, err := r.executer.Exec(
		r.Name+"UpdateImpression",
		`UPDATE ads SET
			impression = ?,
			updated_at = ?
			WHERE id = ?
		`,
		param.Impression,
		time.Now(),
		param.Id,
	)

	if err != nil {
		return err
	}

	return nil
}

// updateClick
func (r *AdRepositoryImpl) UpdateClick(param *entity.Ad) error {
	_, err := r.executer.Exec(
		r.Name+"UpdateClick",
		`UPDATE ads SET
			click = ?,
			updated_at = ?
			WHERE id = ?
		`,
		param.Click,
		time.Now(),
		param.Id,
	)

	if err != nil {
		return err
	}

	return nil
}

func (r *AdRepositoryImpl) GetById(id int) (*entity.Ad, error) {
	var (
		Ad entity.Ad
	)
	err := r.executer.Get(
		r.Name+"GetById",
		&Ad,
		"SELECT * FROM ads WHERE id = ?",
		id,
	)

	if err != nil {
		err = fmt.Errorf("failed to get Ad by firebase id: %w", err)
		fmt.Println(err)
		return nil, err
	}

	return &Ad, nil
}

func (r *AdRepositoryImpl) GetListByLineUserId(lineUserId string) ([]*entity.Ad, error) {
	var (
		AdList []*entity.Ad
	)
	err := r.executer.Select(
		r.Name+"GetListByLineUserId",
		&AdList, `
		SELECT * 
		FROM ads 
		WHERE user_id = (
			SELECT id
			FROM users
			WHERE line_user_id = ?
		)`,
		lineUserId,
	)

	if err != nil {
		err = fmt.Errorf("failed to get Ad by line Ad id: %w", err)
		fmt.Println(err)
		return nil, err
	}

	return AdList, nil
}
