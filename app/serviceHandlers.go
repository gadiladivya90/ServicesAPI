package app

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/divyag/services/dto"
	"github.com/divyag/services/logger"
	"github.com/divyag/services/service"

	"github.com/gorilla/mux"
)

type ServiceHandler struct {
	service service.ServicePackageService
}

func (sh *ServiceHandler) greet(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprint(w, "Hello World!")
}

func (sh *ServiceHandler) getAllServices(w http.ResponseWriter, r *http.Request) {
	fitlerParams, err := getFilterParams(r)
	if err != nil {
		writeResponse(w, http.StatusBadRequest, err.Error())
	}

	services, appErr := sh.service.GetAllServices(fitlerParams)
	if appErr != nil {
		writeResponse(w, appErr.Code, appErr.GetMessage())
	} else {
		writeResponse(w, http.StatusOK, services)
	}
}

func (sh *ServiceHandler) getService(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	service, err := sh.service.GetServiceByID(vars["service_id"])
	if err != nil {
		writeResponse(w, err.Code, err.GetMessage())
	} else {
		writeResponse(w, http.StatusOK, service)
	}
}

func (sh *ServiceHandler) createService(w http.ResponseWriter, r *http.Request) {
	var request dto.ServiceRequestDto
	decodeErr := json.NewDecoder(r.Body).Decode(&request)
	if decodeErr != nil {
		writeResponse(w, http.StatusBadRequest, nil)
		return
	}
	c, err := sh.service.CreateService(request)
	if err != nil {
		writeResponse(w, err.Code, err.GetMessage())
	} else {
		writeResponse(w, http.StatusCreated, c)
	}
}

func (sh *ServiceHandler) updateService(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	_, err := sh.service.GetServiceByID(vars["service_id"])
	if err != nil {
		writeResponse(w, err.Code, err.GetMessage())
		return
	}

	var request dto.ServiceRequestDto
	decodeErr := json.NewDecoder(r.Body).Decode(&request)
	if decodeErr != nil {
		writeResponse(w, http.StatusBadRequest, nil)
		return
	}

	c, err := sh.service.UpdateService(vars["service_id"], request)
	if err != nil {
		writeResponse(w, err.Code, err.GetMessage())
	} else {
		writeResponse(w, http.StatusAccepted, c)
	}

}

func (sh *ServiceHandler) deleteService(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	_, err := sh.service.GetServiceByID(vars["service_id"])
	if err != nil {
		writeResponse(w, err.Code, err.GetMessage())
	}

	err = sh.service.DeleteService(vars["service_id"])
	if err != nil {
		writeResponse(w, err.Code, err.GetMessage())
	} else {
		writeResponse(w, http.StatusNoContent, nil)
	}

}

func writeResponse(w http.ResponseWriter, code int, data interface{}) {
	w.WriteHeader(code)

	if code != http.StatusNoContent {
		w.Header().Add("Content-Type", "applciation/json")
		if err := json.NewEncoder(w).Encode(data); err != nil {
			logger.Error("Unable to Write Data to response!")
		}
	}
}

func getFilterParams(r *http.Request) (dto.FilterParams, error) {
	var err error
	var filters dto.FilterParams
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
