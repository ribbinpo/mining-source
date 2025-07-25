package port

import "github.com/ribbinpo/mining-service/internal/application/domain"

type GetAllStonesOptions struct {
	Page         int
	PageSize     int
	StatusFilter domain.StoneStatusEnum
	SortBy       string
	SortOrder    string
}

type StoneRepository interface {
	CreateStone(stone *domain.StoneDomain) error
	GetAllStones(options *GetAllStonesOptions) ([]*domain.StoneDomain, error)
	GetStoneByID(id string) (*domain.StoneDomain, error)
	GetStoneByURL(url string) (*domain.StoneDomain, error)
	UpdateStone(stone *domain.StoneDomain) error
	DeleteStone(id string) error
}
