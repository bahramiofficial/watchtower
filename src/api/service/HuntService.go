package service

import (
	"context"
	"time"

	"github.com/bahramiofficial/watchtower/src/api/model"
	"github.com/bahramiofficial/watchtower/src/database"
	"gorm.io/gorm"
)

type HuntService struct {
	database *gorm.DB
}

func NewHuntService() *HuntService {
	return &HuntService{
		database: database.GetDb(),
	}
}

func (s *HuntService) Create(ctx context.Context, req *model.CreateHuntRequest) (*model.HuntResponse, error) {
	Hunt := model.Hunt{Name: req.Name}
	Hunt.CreatedAt = time.Now().UTC()

	tx := s.database.WithContext(ctx).Begin()
	err := tx.Create(&Hunt).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	tx.Commit()
	dd := &model.HuntResponse{Id: Hunt.Id, Name: Hunt.Name}
	return dd, nil

}

func (s *HuntService) Update(ctx context.Context, id int, req *model.UpdateHuntRequest) (*model.HuntResponse, error) {
	UpdateMap := map[string]interface{}{
		"Name": req.Name,
	}

	tx := s.database.WithContext(ctx).Begin()
	err := tx.Model(&model.Hunt{}).Where("id = ?", id).Updates(UpdateMap).Error

	if err != nil {
		tx.Rollback()
		return nil, err
	}
	hunt := &model.Hunt{}

	err = tx.Model(&model.Hunt{}).Where("id = ?", id).First(&hunt).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	tx.Commit()

	dd := &model.HuntResponse{Id: hunt.Id, Name: hunt.Name}
	return dd, nil
}
