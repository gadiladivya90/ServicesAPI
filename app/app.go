package app

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/divyag/services/schema"
	"github.com/divyag/services/service"

	"github.com/gorilla/mux"
)

func envCheck() {
	if os.Getenv("DB_HOST") == "" ||
		os.Getenv("DB_PORT") == "" ||
		os.Getenv("DB_USERNAME") == "" ||
		os.Getenv("DB_PASSWORD") == "" ||
		os.Getenv("DB_NAME") == "" {
		log.Fatal("Database env variables missing!")
	}

	if os.Getenv("SERVER_ADDRESS") == "" || os.Getenv("SERVER_PORT") == "" {
		log.Fatal("Application server env variables missing!")

	}
}

func Start() {

	envCheck()
	handlers := ServiceHandler{service: service.NewServicePackageService(schema.NewServiceRepositoryDb())}

	router := mux.NewRouter()
	router.HandleFunc("/greet", handlers.greet).Methods(http.MethodGet)
	router.HandleFunc("/services", handlers.getAllServices).Methods(http.MethodGet)
	router.HandleFunc("/services/{service_id}", handlers.getService).Methods(http.MethodGet)
	router.HandleFunc("/services", handlers.createService).Methods(http.MethodPost)
	router.HandleFunc("/services/{service_id}", handlers.updateService).Methods(http.MethodPut)
	router.HandleFunc("/services/{service_id}", handlers.deleteService).Methods(http.MethodDelete)

	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%s", os.Getenv("SERVER_ADDRESS"), os.Getenv("SERVER_PORT")), router))
}
