package app

import (
	"encoding/json"
	"net/http"

	"github.com/divyag/services/entity"
	"github.com/divyag/services/service"

	"github.com/gorilla/mux"
)

type ServiceVersionHandler struct {
	service service.ServiceVersionService
}

/*
getServiceVersionsByServiceID retrieves ServiceVersions of given service_id
Handle errors and write result to response
*/
func (sv *ServiceVersionHandler) getServiceVersionsByServiceID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	service, err := sv.service.GetServiceVersionsByServiceID(vars["service_id"])
	if err != nil {
		WriteResponse(w, err.Code, err.GetMessage())
	} else {
		WriteResponse(w, http.StatusOK, service)
	}
}

/*
createServiceVersion handles retrieving payload,
Handle errors during creation of service
*/
func (sv *ServiceVersionHandler) createServiceVersion(w http.ResponseWriter, r *http.Request) {
	var request entity.ServiceVersionRequestDto
	//get the request payload
	decodeErr := json.NewDecoder(r.Body).Decode(&request)
	if decodeErr != nil {
		WriteResponse(w, http.StatusBadRequest, nil)
		return
	}
	//
	c, err := sv.service.CreateServiceVersion(request)
	if err != nil {
		WriteResponse(w, err.Code, err.GetMessage())
	} else {
		WriteResponse(w, http.StatusCreated, c)
	}
}

/*
DeleteServiceVersion deletes service of a given service_id
*/
func (sh *ServiceVersionHandler) deleteServiceVersion(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	//verify if the service exists
	_, err := sh.service.GetServiceVersionsByServiceID(vars["service_id"])
	if err != nil {
		WriteResponse(w, err.Code, err.GetMessage())
	}

	// delete service
	err = sh.service.DeleteServiceVersion(vars["service_id"])
	if err != nil {
		WriteResponse(w, err.Code, err.GetMessage())
	} else {
		WriteResponse(w, http.StatusNoContent, nil)
	}

}
