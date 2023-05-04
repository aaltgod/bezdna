package main

import (
	"fmt"
	"net/http"

	"github.com/aaltgod/bezdna/internal/config"
	"github.com/aaltgod/bezdna/internal/database"
	"github.com/aaltgod/bezdna/internal/handler"
	"github.com/aaltgod/bezdna/internal/repository/db"
	"github.com/aaltgod/bezdna/internal/service"
	"github.com/aaltgod/bezdna/internal/sniffer"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(log.DebugLevel)

	serverConfig, err := config.ProvideServerConfig()
	if err != nil {
		log.Fatalln("couldn't provide server config", err)
	}

	dbConfig, err := config.ProvideDBConfig()
	if err != nil {
		log.Fatalln("couldn't provide db config", err)
	}

	dbAdapter, err := database.New(dbConfig)
	if err != nil {
		log.Fatalln("couldn't create db adapter", err)
	}

	snifferConfig, err := config.ProvideSnifferConfig()
	if err != nil {
		log.Fatalln("couldn't provide sniffer config", err)
	}

	sn := sniffer.New(snifferConfig.Interface, db.New(dbAdapter))
	if err := sn.Run(); err != nil {
		log.Fatalln("couldn't run sniffer", err)
	}

	handler := handler.New(service.New(sn, db.New(dbAdapter)))

	router := chi.NewRouter()

	router.Use(cors.AllowAll().Handler)

	router.Mount("/api", func() chi.Router {
		r := chi.NewRouter()

		r.Post("/service", handler.AddService)
		r.Get("/services", handler.GetServices)

		r.Get("/streams-by-service", handler.GetStreamsByService)

		r.HandleFunc("/ws", handler.WSGetStreams)

		return r
	}())

	log.Infof("START SERVER on PORT `%d`", serverConfig.Port)

	http.ListenAndServe(
		fmt.Sprintf("%s:%d", serverConfig.Host, serverConfig.Port),
		router,
	)
}
