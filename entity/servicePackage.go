package entity

import (
	"time"
)

type ServicePackage struct {
	Id          string   `db:"service_id"`
	Name        string   `db:"name"`
	Description string   `db:"description"`
	CreatedAt   string   `db:"created_at"`
	UpdatedAt   string   `db:"updated_at"`
	Versions    []string `db:"versions"`
}

func (s ServicePackage) ToDto() ServiceResponseDto {

	return ServiceResponseDto(s)

}

func ToService(d ServiceResponseDto) *ServicePackage {

	return &ServicePackage{
		Id:          d.Id,
		Name:        d.Name,
		Description: d.Description,
		CreatedAt:   d.CreatedAt,
		UpdatedAt:   d.UpdatedAt,
		Versions:    d.Versions,
	}
}

// Convert ServiceRequestDto to ServicePackage
func NewServicePackage(d ServiceRequestDto) *ServicePackage {
	return &ServicePackage{
		Name:        d.Name,
		Description: d.Description,
		CreatedAt:   time.Now().Format(time.RFC3339),
		UpdatedAt:   time.Now().Format(time.RFC3339),
		Versions:    []string{"s1", "s2"},
	}
}
