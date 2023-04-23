package main

import (
	"net/http"

	"github.com/aaltgod/bezdna/internal/config"
	"github.com/aaltgod/bezdna/internal/database"
	"github.com/aaltgod/bezdna/internal/repository/db"
	"github.com/aaltgod/bezdna/internal/service"
	"github.com/aaltgod/bezdna/internal/sniffer"
	"github.com/go-chi/chi/v5"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(log.DebugLevel)

	snifferConfig, err := config.ProvideSnifferConfig()
	if err != nil {
		log.Fatal("couldn't provide sniffer config", err)
	}

	sn := sniffer.New(snifferConfig.Interface)
	if err := sn.Run(); err != nil {
		log.Fatal("couldn't run sniffer", err)
	}

	dbConfig, err := config.ProvideDBConfig()
	if err != nil {
		log.Fatal("couldn't provide db config", err)
	}

	dbAdapter, err := database.New(dbConfig)
	if err != nil {
		log.Fatal("couldn't create db adapter", err)
	}

	service := service.New(sn, db.New(dbAdapter))

	router := chi.NewRouter()

	router.Get("/services", service.GetServices)
	router.Post("/service", service.AddService)

	http.ListenAndServe(":2137", router)
}
