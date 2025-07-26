package database

import "github.com/ribbinpo/mining-service/internal/application/domain"

func StoneDomainToModel(domain *domain.StoneDomain) *StoneModel {
	paths := make([]PathModel, len(domain.Paths))
	for i, path := range domain.Paths {
		paths[i] = PathModel{
			Name: path,
		}
	}
	return &StoneModel{
		ID:        domain.ID,
		Domain:    domain.DomainURL,
		Version:   domain.Version,
		Status:    string(domain.Status),
		Paths:     paths,
		CreatedAt: domain.CreatedAt,
		UpdatedAt: domain.UpdatedAt,
	}
}

func StoneModelToDomain(model *StoneModel) *domain.StoneDomain {
	paths := make([]string, len(model.Paths))
	for i, path := range model.Paths {
		paths[i] = path.Name
	}
	return &domain.StoneDomain{
		ID:        model.ID,
		DomainURL: model.Domain,
		Version:   model.Version,
		Status:    domain.StoneStatusEnum(model.Status),
		Paths:     paths,
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
	}
}
