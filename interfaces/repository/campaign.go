package repository

import (
	"fmt"
	"time"

	"github.com/hidenari-yuda/paychan-server/domain/entity"
	"github.com/hidenari-yuda/paychan-server/domain/utility"
	"github.com/hidenari-yuda/paychan-server/interfaces"
	"github.com/hidenari-yuda/paychan-server/usecase"
)

type CampaignRepositoryImpl struct {
	Name     string
	executer interfaces.SQLExecuter
}

// Create(param *entity.Campaign) error
// GetById(id int) (*entity.Campaign, error)
// GetListByUserId(userId int) ([]*entity.Campaign, error)
// GetListByLineUserId(lineUserId string) ([]*entity.Campaign, error)

func NewCampaignRepositoryImpl(ex interfaces.SQLExecuter) usecase.CampaignRepository {
	return &CampaignRepositoryImpl{
		Name:     "CampaignRepository",
		executer: ex,
	}
}

func (r *CampaignRepositoryImpl) Create(param *entity.Campaign) error {
	_, err := r.executer.Exec(
		r.Name+"Create",
		`INSERT INTO campaigns (
			uuid,
			service,
			url,
			picture_url,
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
		param.PictureUrl,
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
func (r *CampaignRepositoryImpl) Update(param *entity.Campaign) error {
	_, err := r.executer.Exec(
		r.Name+"Update",
		`UPDATE campaigns SET
			service = ?,
			url = ?,
			picture_url = ?,
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
		param.PictureUrl,
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
func (r *CampaignRepositoryImpl) UpdateImpression(param *entity.Campaign) error {
	_, err := r.executer.Exec(
		r.Name+"UpdateImpression",
		`UPDATE campaigns SET
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
func (r *CampaignRepositoryImpl) UpdateClick(param *entity.Campaign) error {
	_, err := r.executer.Exec(
		r.Name+"UpdateClick",
		`UPDATE campaigns SET
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

func (r *CampaignRepositoryImpl) GetById(id int) (*entity.Campaign, error) {
	var (
		Campaign entity.Campaign
	)
	err := r.executer.Get(
		r.Name+"GetById",
		&Campaign,
		"SELECT * FROM campaigns WHERE id = ?",
		id,
	)

	if err != nil {
		err = fmt.Errorf("failed to get Campaign by firebase id: %w", err)
		fmt.Println(err)
		return nil, err
	}

	return &Campaign, nil
}

func (r *CampaignRepositoryImpl) GetListByLineUserId(lineUserId string) ([]*entity.Campaign, error) {
	var (
		CampaignList []*entity.Campaign
	)
	err := r.executer.Select(
		r.Name+"GetListByLineUserId",
		&CampaignList, `
		SELECT * 
		FROM campaigns 
		WHERE user_id = (
			SELECT id
			FROM users
			WHERE line_user_id = ?
		)`,
		lineUserId,
	)

	if err != nil {
		err = fmt.Errorf("failed to get Campaign by line Campaign id: %w", err)
		fmt.Println(err)
		return nil, err
	}

	return CampaignList, nil
}
