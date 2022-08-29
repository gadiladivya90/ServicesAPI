package entity

import (
	"regexp"

	"github.com/divyag/services/errs"
)

type ServiceRequestDto struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (s ServiceRequestDto) Validate() *errs.AppErr {

	//verify if name and description texts only has alphanumberic characters
	if s.Description == "" || !regexp.MustCompile(`^[a-zA-Z0-9]*$`).MatchString(s.Description) {
		return errs.BadRequestError("Description cannot be empty and should only contain alphanumber!")
	}

	if s.Name == "" || !regexp.MustCompile(`^[a-zA-Z0-9]*$`).MatchString(s.Name) {
		return errs.BadRequestError("Name cannot be empty and should only contain alphanumber!")
	}

	return nil
}

type FilterParams struct {
	Offset uint64
	Limit  uint64
	Filter string
	Sort   string
}

type InputFilter struct {
	Key       string
	Value     string
	Condition string
}
