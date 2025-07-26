package database

import (
	"context"
	"errors"

	"github.com/ribbinpo/mining-service/internal/application/domain"
	"github.com/ribbinpo/mining-service/internal/application/port"
	model "github.com/ribbinpo/mining-service/internal/frameworks/database/model"
	"gorm.io/gorm"
)

type stoneRepository struct {
	db *gorm.DB
}

func NewStoneRepository(db *gorm.DB) port.StoneRepository {
	return &stoneRepository{db}
}

func (r *stoneRepository) GetStoneByID(id string) (*domain.StoneDomain, error) {
	ctx := context.Background()
	var stone model.StoneModel
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&stone).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return model.StoneModelToDomain(&stone), nil
}

func (r *stoneRepository) CreateStone(stone *domain.StoneDomain) error {
	stoneModel := model.StoneDomainToModel(stone)
	ctx := context.Background()
	if err := r.db.WithContext(ctx).Create(stoneModel).Error; err != nil {
		return err
	}
	return nil
}

func (r *stoneRepository) GetAllStones(options *port.GetAllStonesOptions) ([]*domain.StoneDomain, error) {
	ctx := context.Background()
	query := r.db.WithContext(ctx).Model(&model.StoneModel{})

	// Apply status filter if specified
	if options.StatusFilter != "" {
		query = query.Where("status = ?", string(options.StatusFilter))
	}

	// Apply sorting if specified
	if options.SortBy != "" {
		order := options.SortBy
		if options.SortOrder != "" {
			order += " " + options.SortOrder
		}
		query = query.Order(order)
	}

	// Apply pagination
	if options.PageSize > 0 {
		query = query.Limit(options.PageSize).Offset(options.Page * options.PageSize)
	}

	var stones []model.StoneModel
	if err := query.Find(&stones).Error; err != nil {
		return nil, err
	}

	domains := make([]*domain.StoneDomain, len(stones))
	for i, stone := range stones {
		domains[i] = model.StoneModelToDomain(&stone)
	}
	return domains, nil
}

func (r *stoneRepository) GetStoneByURL(url string) (*domain.StoneDomain, error) {
	ctx := context.Background()
	var stone model.StoneModel
	err := r.db.WithContext(ctx).Where("domain = ?", url).First(&stone).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return model.StoneModelToDomain(&stone), nil
}

func (r *stoneRepository) UpdateStone(stone *domain.StoneDomain) error {
	ctx := context.Background()
	stoneModel := model.StoneDomainToModel(stone)

	// Use Updates to update only non-zero fields
	result := r.db.WithContext(ctx).Model(&model.StoneModel{}).Where("id = ?", stone.ID).Updates(stoneModel)
	if result.Error != nil {
		return result.Error
	}

	// Check if any rows were affected
	if result.RowsAffected == 0 {
		return errors.New("stone not found")
	}

	return nil
}

func (r *stoneRepository) DeleteStone(id string) error {
	ctx := context.Background()

	result := r.db.WithContext(ctx).Where("id = ?", id).Delete(&model.StoneModel{})
	if result.Error != nil {
		return result.Error
	}

	// Check if any rows were affected
	if result.RowsAffected == 0 {
		return errors.New("stone not found")
	}

	return nil
}
