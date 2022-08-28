package dto

import "github.com/divyag/services/errs"

type ServiceRequestDto struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (c ServiceRequestDto) Validate() *errs.AppErr {

	//TODO : Update validation here
	return nil
}

type FilterParams struct {
	Offset uint64
	Limit  uint64
	Filter string
	Sort   string
}
