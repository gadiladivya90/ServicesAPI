package service

import (
	"github.com/divyag/services/entity"
	"github.com/divyag/services/errs"
)

type ServiceVersionService interface {
	GetServiceVersionsByServiceID(string) ([]entity.ServiceVersionResponseDto, *errs.AppErr)
	CreateServiceVersion(entity.ServiceVersionRequestDto) (entity.ServiceVersionResponseDto, *errs.AppErr)
	DeleteServiceVersion(string) *errs.AppErr
}

type DefaultServiceVersionService struct {
	repo ServiceVersionRepositoryDb
}

func (s DefaultServiceVersionService) GetServiceVersionsByServiceID(id string) ([]entity.ServiceVersionResponseDto, *errs.AppErr) {

	return s.repo.FindServiceVersionsByID(id)

}

func (s DefaultServiceVersionService) CreateServiceVersion(d entity.ServiceVersionRequestDto) (entity.ServiceVersionResponseDto, *errs.AppErr) {
	err := entity.ServiceVersionRequestDto.Validate(d)
	if err != nil {
		return entity.ServiceVersionResponseDto{}, err
	}

	service := entity.NewServiceVersion(d)
	newService, err := s.repo.SaveServiceVersion(service)
	if err != nil {
		return entity.ServiceVersionResponseDto{}, err
	}

	return newService.ToSVDto(), nil
}

func (s DefaultServiceVersionService) DeleteServiceVersion(id string) *errs.AppErr {

	err := s.repo.DeleteServiceVersion(id)
	if err != nil {
		return err
	}

	return nil
}

func NewServiceVersionService(repo ServiceVersionRepositoryDb) DefaultServiceVersionService {
	return DefaultServiceVersionService{repo: repo}
}
