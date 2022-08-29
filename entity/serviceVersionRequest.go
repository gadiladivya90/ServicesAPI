package entity

import (
	"regexp"

	"github.com/divyag/services/errs"
)

type ServiceVersionRequestDto struct {
	Name      string `json:"name"`
	ServiceId string `json:"service_package_id"`
}

func (s ServiceVersionRequestDto) Validate() *errs.AppErr {

	if s.Name == "" || !regexp.MustCompile(`^[a-zA-Z0-9]*$`).MatchString(s.Name) {
		return errs.BadRequestError("Name cannot be empty and should only contain alphanumber!")
	}

	return nil
}
