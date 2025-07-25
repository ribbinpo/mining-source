package port

import "github.com/ribbinpo/mining-service/internal/application/domain"

type StoneUsecase interface {
	ExecuteDigFromQueue() error
	ExecuteDig(targetURL string) ([]string, error)
	RegisterDig(url string, refresh bool) error
	GetStoneByURL(url string) (*domain.StoneDomain, error)
	// GetStone() (*domain.StoneDomain, error)
}
