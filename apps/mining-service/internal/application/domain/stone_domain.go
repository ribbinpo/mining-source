package domain

import (
	"time"
)

type StoneDomain struct {
	ID        string          `json:"id"`
	DomainURL string          `json:"domain_url"`
	Paths     []string        `json:"paths"`
	Status    StoneStatusEnum `json:"status"`
	Version   int             `json:"version"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}

func NewStoneDomain(s *StoneDomain) *StoneDomain {
	return &StoneDomain{
		ID:        s.ID,
		DomainURL: s.DomainURL,
		Paths:     s.Paths,
		Status:    s.Status,
		Version:   s.Version,
		CreatedAt: s.CreatedAt,
		UpdatedAt: s.UpdatedAt,
	}
}

func CreateStoneDomain(domainURL string) *StoneDomain {
	return &StoneDomain{
		DomainURL: domainURL,
		Paths:     []string{},
		Status:    StoneStatusEnumPending,
		Version:   1,
	}
}

func (s *StoneDomain) AddPaths(paths []string) {
	s.Paths = append(s.Paths, paths...)
	s.UpdatedAt = time.Now()
	s.Status = StoneStatusEnumCompleted
}

func (s *StoneDomain) Retry() {
	s.Status = StoneStatusEnumPending
	s.Version++
	s.UpdatedAt = time.Now()
}
