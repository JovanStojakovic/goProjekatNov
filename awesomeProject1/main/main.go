package main

import (
	"context"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	router := mux.NewRouter()
	router.StrictSlash(true)

	store, err := New()
	if err != nil {
		log.Fatal(err)
	}

	server := configStore{
		store: store,
	}

	router.HandleFunc("/config/", server.addConfigHandler).Methods("POST")
	router.HandleFunc("/config/{id}/", server.addConfigVersion).Methods("POST")
	router.HandleFunc("/configs/{id}", server.getConfigVersionsHandler).Methods("GET")
	router.HandleFunc("/group/", server.addConfigGroupHandler).Methods("POST")
	router.HandleFunc("/configs/", server.getAllConfigsHandler).Methods("GET")
	router.HandleFunc("/group/{id}/", server.addConfigGroupVersion).Methods("POST")
	router.HandleFunc("/groups/", server.getAllGroupsHandler).Methods("GET")
	router.HandleFunc("/group/{id}/", server.getConfigGroupVersions).Methods("GET")
	router.HandleFunc("/config/{id}/{version}/", server.deleteConfigHandler).Methods("DELETE")
	router.HandleFunc("/group/{id}/{version}/", server.deleteConfigGroupHandler).Methods("DELETE")
	router.HandleFunc("/group/{id}/{version}/", server.getConfigGroupHandler).Methods("GET")
	router.HandleFunc("/config/{id}/{version}/", server.getConfigHandler).Methods("GET")
	router.HandleFunc("/group/{id}/{version}/", server.UpdateConfigGroupHandler).Methods("PUT")
	router.HandleFunc("/group/{id}/{version}/{labels}/", server.getConfigByLabels).Methods("GET")

	// start server
	srv := &http.Server{Addr: "0.0.0.0:8080", Handler: router}
	go func() {
		log.Println("server starting")
		if err := srv.ListenAndServe(); err != nil {
			if err != http.ErrServerClosed {
				log.Fatal(err)
			}
		}
	}()

	<-quit

	log.Println("service shutting down ...")

	// gracefully stop server
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}
	log.Println("server stopped")
}
