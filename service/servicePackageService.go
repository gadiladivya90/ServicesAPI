package service

import (
	"fmt"

	"github.com/divyag/services/entity"
	"github.com/divyag/services/errs"
)

type ServicePackageService interface {
	GetAllServices(entity.FilterParams) (*entity.PaginationResponseDto, *errs.AppErr)
	GetServiceByID(string) (entity.ServiceResponseDto, *errs.AppErr)
	CreateService(entity.ServiceRequestDto) (entity.ServiceResponseDto, *errs.AppErr)
	UpdateService(*entity.ServicePackage) (entity.ServiceResponseDto, *errs.AppErr)
	DeleteService(string) *errs.AppErr
}

type DefaultServicePackageService struct {
	repo ServiceRepositoryDb
}

func (s DefaultServicePackageService) GetAllServices(filters entity.FilterParams) (*entity.PaginationResponseDto, *errs.AppErr) {
	r, err := s.repo.FindAll(filters)
	if err != nil {
		return nil, err
	}
	return r, nil

}

func (s DefaultServicePackageService) GetServiceByID(id string) (entity.ServiceResponseDto, *errs.AppErr) {
	service, err := s.repo.FindServiceByID(id)
	if err != nil {
		return entity.ServiceResponseDto{}, err
	}

	return service.ToDto(), nil
}

func (s DefaultServicePackageService) CreateService(d entity.ServiceRequestDto) (entity.ServiceResponseDto, *errs.AppErr) {
	fmt.Println(d)
	err := entity.ServiceRequestDto.Validate(d)
	if err != nil {
		return entity.ServiceResponseDto{}, err
	}

	service := entity.NewServicePackage(d)
	newService, err := s.repo.SaveService(service)
	if err != nil {
		return entity.ServiceResponseDto{}, err
	}

	return newService.ToDto(), nil
}

func (s DefaultServicePackageService) UpdateService(service *entity.ServicePackage) (entity.ServiceResponseDto, *errs.AppErr) {

	udpatedService, err := s.repo.UpdateService(service)
	if err != nil {
		return entity.ServiceResponseDto{}, err
	}

	return udpatedService.ToDto(), nil
}

func (s DefaultServicePackageService) DeleteService(id string) *errs.AppErr {

	err := s.repo.DeleteService(id)
	if err != nil {
		return err
	}

	return nil
}

func NewServicePackageService(repo ServiceRepositoryDb) DefaultServicePackageService {
	return DefaultServicePackageService{repo: repo}
}
