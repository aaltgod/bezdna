package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/aaltgod/bezdna/internal/config"
	"github.com/aaltgod/bezdna/internal/database"
	httpHandler "github.com/aaltgod/bezdna/internal/handler/http"
	wsHandler "github.com/aaltgod/bezdna/internal/handler/ws"
	"github.com/aaltgod/bezdna/internal/repository/db"
	"github.com/aaltgod/bezdna/internal/service"
	"github.com/aaltgod/bezdna/internal/sniffer"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	log "github.com/sirupsen/logrus"
)

func main() {
	ctx := context.Background()

	log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(log.InfoLevel)

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
	if err := sn.Run(ctx); err != nil {
		log.Fatalln("couldn't run sniffer", err)
	}

	service := service.New(sn, db.New(dbAdapter))
	httpHandler := httpHandler.New(service)
	wsHandler := wsHandler.New(service)

	router := chi.NewRouter()

	router.Use(cors.AllowAll().Handler)

	router.Mount("/api", func() chi.Router {
		router.Post("/create-service", httpHandler.UpsertService)
		router.Get("/get-services", httpHandler.GetServices)

		router.Get("/get-streams-by-service", httpHandler.GetStreamsByService)

		router.Mount("/ws", func() chi.Router {
			router.HandleFunc("/get-streams-by-service", wsHandler.GetStreamsByService)
			router.HandleFunc("/get-streams", wsHandler.GetStreams)

			return router
		}())

		return router
	}())

	log.Infof("START SERVER on HOST %s and PORT `%d`", serverConfig.Host, serverConfig.Port)

	log.Fatal(http.ListenAndServe(
		fmt.Sprintf("%s:%d", serverConfig.Host, serverConfig.Port),
		router,
	))
}
