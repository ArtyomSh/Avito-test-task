package main

import (
	//_ "AvitoTesting/cmd/main/docs"
	"AvitoTesting/internal/config"
	"AvitoTesting/internal/handlers"
	"AvitoTesting/pkg/client/postgres"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	cfg := config.GetConfig()
	log.Println("config received", cfg)
	DB := postgres.Init(cfg)
	h := handlers.New(DB)

	fmt.Println("router created")
	router := mux.NewRouter()

	router.HandleFunc("/segments/{id}", h.GetAllSegments).Methods(http.MethodGet)
	router.HandleFunc("/segment", h.AddSegment).Methods(http.MethodPost)
	router.HandleFunc("/user/{id}", h.AddUser).Methods(http.MethodPost)
	router.HandleFunc("/segment", h.DeleteSegment).Methods(http.MethodDelete)

	//router.PathPrefix("/swagger").Handler(httpSwagger.WrapHandler)
	log.Println("API is running!")

	port := cfg.Listen.Port
	http.ListenAndServe(fmt.Sprintf(":%s", port), router)
}
