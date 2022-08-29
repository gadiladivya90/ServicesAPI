package app

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/divyag/services/entity"
	"github.com/divyag/services/errs"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
)

var router *mux.Router
var ch ServiceHandlers

//todo:generate service mock using mockgen
var mockService *service.MockServicePackageService

func setup(t *testing.T) func() {
	ctrl := gomock.NewController(t)
	mockService = service.NewMockServicePackageService(ctrl)
	ch = ServiceHandlers{mockService}
	router = mux.NewRouter()
	router.HandleFunc("/services", ch.getAllServicces)
	return func() {
		router = nil
		defer ctrl.Finish()
	}
}

func Test_should_return_services_with_status_code_200(t *testing.T) {
	// Arrange
	teardown := setup(t)
	defer teardown()

	dummyServices := []entity.ServiceResponseDto{
		{"a81bc81b-dead-4e5d-abff-90865d1e13b1", "s1", "s1d1", "2020-1-30 13:10:53.163", "2020-12-30 13:10:53.163", []string{}},
		{"bc1bc81b-dead-4e5d-abff-90865d1e13b1", "s2", "s2d2", "2020-2-20 13:10:53.163", "2020-12-30 13:10:53.163", []string{}},
	}
	mockService.EXPECT().getAllServices("").Return(dummyServices, nil)
	request, _ := http.NewRequest(http.MethodGet, "/services", nil)

	// Act
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	// Assert
	if recorder.Code != http.StatusOK {
		t.Error("Failed while testing the status code")
	}
}

func Test_should_return_status_code_500_with_error_message(t *testing.T) {
	// Arrange
	teardown := setup(t)
	defer teardown()
	mockService.EXPECT().GetAllServices("").Return(nil, errs.InternalServerError("some database error"))
	request, _ := http.NewRequest(http.MethodGet, "/services", nil)

	// Act
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	// Assert
	if recorder.Code != http.StatusInternalServerError {
		t.Error("Failed while testing the status code")
	}
}
