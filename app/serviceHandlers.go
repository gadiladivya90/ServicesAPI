package app

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/divyag/services/entity"
	"github.com/divyag/services/logger"
	"github.com/divyag/services/service"

	"github.com/gorilla/mux"
)

type ServiceHandler struct {
	service service.ServicePackageService
}

/*
getAllServices retrieves List of Services and statusCode
Handle Bad request errors or UnExpectedErrors and write to response
*/
func (sh *ServiceHandler) getAllServices(w http.ResponseWriter, r *http.Request) {
	fitlerParams, err := getFilterParams(r)
	if err != nil {
		logger.Error(fmt.Sprintf("Error during Parsing filter params err: %+v\n", err.Error()))
		WriteResponse(w, http.StatusBadRequest, err.Error())
	}

	services, appErr := sh.service.GetAllServices(fitlerParams)
	if appErr != nil {
		WriteResponse(w, appErr.Code, appErr.GetMessage())
	} else {
		WriteResponse(w, http.StatusOK, services)
	}
}

/*
getService retrieves Services of given service_id
Handle errors and write result to response
*/
func (sh *ServiceHandler) getService(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	service, err := sh.service.GetServiceByID(vars["service_id"])
	if err != nil {
		WriteResponse(w, err.Code, err.GetMessage())
	} else {
		WriteResponse(w, http.StatusOK, service)
	}
}

/*
createService handles retrieving payload,
Handle errors during creation of service
*/
func (sh *ServiceHandler) createService(w http.ResponseWriter, r *http.Request) {
	var request entity.ServiceRequestDto
	//get the request payload
	decodeErr := json.NewDecoder(r.Body).Decode(&request)
	if decodeErr != nil {
		WriteResponse(w, http.StatusBadRequest, nil)
		return
	}
	//
	c, err := sh.service.CreateService(request)
	if err != nil {
		WriteResponse(w, err.Code, err.GetMessage())
	} else {
		WriteResponse(w, http.StatusCreated, c)
	}
}

/*
updateService handles retrieving payload,
Handle errors during updating service
*/
func (sh *ServiceHandler) updateService(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	existingService, err := sh.service.GetServiceByID(vars["service_id"])
	if err != nil {
		WriteResponse(w, err.Code, err.GetMessage())
		return
	}

	var request entity.ServiceRequestDto
	decodeErr := json.NewDecoder(r.Body).Decode(&request)
	if decodeErr != nil {
		WriteResponse(w, http.StatusBadRequest, nil)
		return
	}

	err = entity.ServiceRequestDto.Validate(request)
	if err != nil {
		WriteResponse(w, http.StatusBadRequest, nil)
		return
	}

	// create ServicePackge and updating Id and createdAt from existingService details
	service := entity.NewServicePackage(request)
	service.Id = existingService.Id
	service.CreatedAt = existingService.CreatedAt

	c, err := sh.service.UpdateService(service)
	if err != nil {
		WriteResponse(w, err.Code, err.GetMessage())
	} else {
		WriteResponse(w, http.StatusAccepted, c)
	}
}

/*
deleteService deletes service of a given service_id
*/
func (sh *ServiceHandler) deleteService(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	//verify if the service exists
	_, err := sh.service.GetServiceByID(vars["service_id"])
	if err != nil {
		WriteResponse(w, err.Code, err.GetMessage())
	}

	// delete service
	err = sh.service.DeleteService(vars["service_id"])
	if err != nil {
		WriteResponse(w, err.Code, err.GetMessage())
	} else {
		WriteResponse(w, http.StatusNoContent, nil)
	}

}

func getFilterParams(r *http.Request) (entity.FilterParams, error) {
	var err error
	var filters entity.FilterParams
	if r.URL.Query().Has("filter") {
		filters.Filter = r.URL.Query().Get("filter")
	}
	if r.URL.Query().Has("limit") {
		filters.Limit, err = strconv.ParseUint(r.URL.Query().Get("limit"), 10, 32)
		if err != nil {
			return filters, err
		}
	}

	if r.URL.Query().Has("offset") {
		filters.Offset, err = strconv.ParseUint(r.URL.Query().Get("offset"), 10, 32)
		if err != nil {
			return filters, err
		}
	}

	if r.URL.Query().Has("sort") {
		filters.Sort = r.URL.Query().Get("sort")
	} else {
		filters.Sort = "name"
	}

	return filters, nil
}
