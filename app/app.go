package app

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/divyag/services/logger"
	"github.com/divyag/services/service"
	"github.com/jmoiron/sqlx"

	"github.com/gorilla/mux"
)

// Sanity check before starting application
func envCheck() {

	envVariables := []string{
		"DB_HOST",
		"DB_PORT",
		"DB_USERNAME",
		"DB_PASSWORD",
		"DB_NAME",
		"SERVER_ADDRESS",
		"SERVER_PORT",
	}

	for _, v := range envVariables {
		if os.Getenv(v) == "" {
			log.Fatalf("Application cannot be initiated, env variable %s missing!", v)
		}
	}
}

// Start application and create routes
func Start() {

	envCheck()
	dbClient := getDbClient()

	//Create routes for endpoints
	router := mux.NewRouter()

	//wiring
	serviceRepositoryDb := service.NewServiceRepositoryDb(dbClient)
	serviceVersionRepositoryDb := service.NewServiceVersionRepositoryDb(dbClient)

	//initialize serviceHandler
	serviceHandler := ServiceHandler{service.NewServicePackageService(serviceRepositoryDb)}
	serviceVersionsHandler := ServiceVersionHandler{service.NewServiceVersionService(serviceVersionRepositoryDb)}

	//TODO: app conetxt.Conetxt across
	router.HandleFunc("/services", serviceHandler.getAllServices).
		Methods(http.MethodGet)
	router.HandleFunc("/services/{service_id}", serviceHandler.getService).
		Methods(http.MethodGet)
	router.HandleFunc("/services", serviceHandler.createService).
		Methods(http.MethodPost)
	router.HandleFunc("/services/{service_id}", serviceHandler.updateService).
		Methods(http.MethodPut)
	router.HandleFunc("/services/{service_id}", serviceHandler.deleteService).
		Methods(http.MethodDelete)
	router.HandleFunc("/services/{service_id}/versions", serviceVersionsHandler.getServiceVersionsByServiceID).
		Methods(http.MethodGet)
	router.HandleFunc("/services/{service_id}/versions", serviceVersionsHandler.createServiceVersion).
		Methods(http.MethodPost)
	router.HandleFunc("/services/{service_id}/versions", serviceVersionsHandler.deleteServiceVersion).
		Methods(http.MethodDelete)

	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%s", os.Getenv("SERVER_ADDRESS"), os.Getenv("SERVER_PORT")), router))

}

func WriteResponse(w http.ResponseWriter, code int, data interface{}) {
	w.WriteHeader(code)

	if code != http.StatusNoContent {
		w.Header().Add("Content-Type", "applciation/json")
		if err := json.NewEncoder(w).Encode(data); err != nil {
			logger.Error("Unable to Write Data to response!")
		}
	}
}

func getDbClient() *sqlx.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"))

	dbClient, err := sqlx.Open("postgres", psqlInfo)
	if err != nil {
		panic("unable to establish connection")
	}

	dbClient.SetConnMaxLifetime(time.Minute * 3)
	dbClient.SetMaxOpenConns(10)
	dbClient.SetMaxIdleConns(10)
	return dbClient
}
