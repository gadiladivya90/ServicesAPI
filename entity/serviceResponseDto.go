package entity

import (
	"math"
)

type ServiceResponseDto struct {
	Id          string   `json:"service_id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	CreatedAt   string   `json:"created_at"`
	UpdatedAt   string   `json:"updated_at"`
	Versions    []string `json:"versions"`
}

type PaginationResponseDto struct {
	Data      []ServiceResponseDto
	Count     int
	Page      int
	PageCount int
	Total     int
}

// Paginate return paginated result for GetMany Request
func Paginate(d []ServiceResponseDto, total int, filter FilterParams) *PaginationResponseDto {
	var page, pageCount int
	if filter.Limit != 0 && filter.Offset != 0 {
		page = int(math.Floor(float64(filter.Offset)/float64(filter.Limit)) + 1)
	} else {
		page = 1
	}

	if filter.Limit != 0 && total != 0 {
		pageCount = int(math.Ceil(float64(total) / float64(filter.Limit)))
	} else {
		pageCount = 1
	}

	return &PaginationResponseDto{
		Data:      d,
		Count:     len(d),
		Total:     total,
		Page:      page,
		PageCount: pageCount,
	}
}
