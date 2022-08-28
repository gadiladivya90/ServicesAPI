package service

import (
	"github.com/divyag/services/dto"
	"github.com/divyag/services/errs"
	"github.com/divyag/services/schema"
)

type ServicePackageService interface {
	GetAllServices(dto.FilterParams) (*dto.PaginationResponseDto, *errs.AppErr)
	GetServiceByID(string) (dto.ServiceResponseDto, *errs.AppErr)
	CreateService(dto.ServiceRequestDto) (dto.ServiceResponseDto, *errs.AppErr)
	UpdateService(string, dto.ServiceRequestDto) (dto.ServiceResponseDto, *errs.AppErr)
	DeleteService(string) *errs.AppErr
}

type DefaultServicePackageService struct {
	repo schema.ServiceRepositoryDb
}

func (s DefaultServicePackageService) GetAllServices(filters dto.FilterParams) (*dto.PaginationResponseDto, *errs.AppErr) {
	r, err := s.repo.FindAll(filters)
	if err != nil {
		return nil, err
	}
	return r, nil

}

func (s DefaultServicePackageService) GetServiceByID(id string) (dto.ServiceResponseDto, *errs.AppErr) {
	service, err := s.repo.FindServiceByID(id)
	if err != nil {
		return dto.ServiceResponseDto{}, err
	}

	return service.ToDto(), nil
}

func (s DefaultServicePackageService) CreateService(d dto.ServiceRequestDto) (dto.ServiceResponseDto, *errs.AppErr) {
	err := dto.ServiceRequestDto.Validate(d)
	if err != nil {
		return dto.ServiceResponseDto{}, err
	}

	service := schema.NewServicePackage(d)
	newService, err := s.repo.SaveService(service)
	if err != nil {
		return dto.ServiceResponseDto{}, err
	}

	return newService.ToDto(), nil
}

func (s DefaultServicePackageService) UpdateService(id string, d dto.ServiceRequestDto) (dto.ServiceResponseDto, *errs.AppErr) {
	err := dto.ServiceRequestDto.Validate(d)
	if err != nil {
		return dto.ServiceResponseDto{}, err
	}

	service := schema.NewServicePackage(d)
	service.Id = id
	udpatedService, err := s.repo.UpdateService(id, service)
	if err != nil {
		return dto.ServiceResponseDto{}, err
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

func NewServicePackageService(repo schema.ServiceRepositoryDb) DefaultServicePackageService {
	return DefaultServicePackageService{repo: repo}
}
