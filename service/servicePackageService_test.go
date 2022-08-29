package service

import (
	"testing"
	"time"

	"github.com/ashishjuyal/banking-lib/errs"
	"github.com/divyag/services/entity"
	"github.com/golang/mock/gomock"
)

func Test_should_return_error_response_when_the_request_is_not_validated(t *testing.T) {
	// Arrange
	request := entity.ServiceRequestDto{
		Name:        "mockservice",
		Description: "mockDescription",
	}
	service := entity.NewServicePackage(request)
	// Act
	_, appError := service.NewServicePackageService(request)
	// Assert
	if appError == nil {
		t.Error("failed while testing the new account validation")
	}
}

var mockRepo *entity.ServicePackageService
var service ServicePackageService

func setup(t *testing.T) func() {
	ctrl := gomock.NewController(t)
	mockRepo = entity.NewMockAccountRepository(ctrl)
	service = NewServicePackageService(mockRepo)
	return func() {
		service = nil
		defer ctrl.Finish()
	}
}

func Test_should_return_an_error_if_the_new_service_cannot_be_created(t *testing.T) {
	// Arrange
	teardown := setup(t)
	defer teardown()

	req := entity.NewServicePackage{
		Name:        "",
		Description: "mockServiceDescription",
	}
	mockService := entity.ServicePackage{
		Name:        req.Name,
		Description: req.Description,
		CreatedAt:   time.Now().Format(time.RFC3339),
		UpdatedAt:   time.Now().Format(time.RFC3339),
		Versions:    []string{},
	}
	mockRepo.EXPECT().SaveService(req).Return(nil, errs.NewUnexpectedError("Name cannot be empty and should only contain alphanumber!"))

	// Assert
	if appError == nil {
		t.Error("Test failed while validating error for new servicePackage")
	}

}

func Test_should_return_new_account_response_when_a_new_account_is_saved_successfully(t *testing.T) {
	// Arrange
	teardown := setup(t)
	defer teardown()

	req := entity.NewServicePackage{
		Name:        "mockName",
		Description: "mockServiceDescription",
	}
	mockService := entity.ServicePackage{
		Name:        req.Name,
		Description: req.Description,
		CreatedAt:   time.Now().Format(time.RFC3339),
		UpdatedAt:   time.Now().Format(time.RFC3339),
		Versions:    []string{},
	}
	mockService.Id = "mockUUID"
	mockRepo.EXPECT().SaveService(mockService).Return(&mockService.Id, nil)

	// Assert
	if appError != nil {
		t.Error("Test failed while creating new account")
	}

}
