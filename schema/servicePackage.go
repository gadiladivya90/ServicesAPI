package schema

import (
	"time"

	"github.com/divyag/services/dto"
)

type ServicePackage struct {
	Id          string `db:"service_id"`
	Name        string `db:"name"`
	Description string `db:"description"`
	CreatedAt   string `db:"created_at"`
	UpdatedAt   string `db:"updated_at"`
}

func (s ServicePackage) ToDto() dto.ServiceResponseDto {

	return dto.ServiceResponseDto{
		Id:          s.Id,
		Description: s.Description,
		Name:        s.Name,
		CreatedAt:   s.CreatedAt,
		UpdatedAt:   s.UpdatedAt,
	}
}

func ToService(d dto.ServiceResponseDto) *ServicePackage {

	return &ServicePackage{
		Id:          d.Id,
		Name:        d.Name,
		Description: d.Description,
		CreatedAt:   d.CreatedAt,
		UpdatedAt:   d.UpdatedAt,
	}
}

func NewServicePackage(d dto.ServiceRequestDto) *ServicePackage {
	return &ServicePackage{
		Name:        d.Name,
		Description: d.Description,
		CreatedAt:   time.Now().Format(time.RFC3339),
		UpdatedAt:   time.Now().Format(time.RFC3339),
	}
}
