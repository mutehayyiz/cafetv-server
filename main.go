package main

import (
	"crypto/tls"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/mutehayyiz/cafetv-server/config"
	"github.com/mutehayyiz/cafetv-server/storage"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
)

const Version = "0.1"

func intro() {
	logrus.Info("cafetv ", Version)
	logrus.Info("Copyright Ahmet Şenharputlu 2020")
	logrus.Info("https://cafetv.com")
}

func main() {

	intro()

	config.Global.Load("config.json")
	storage.Connect()

	router := GenerateRouter()

	addr := fmt.Sprintf(":%d", config.Global.Port)
	certFile := config.Global.ServerOptions.CertFile
	keyFile := config.Global.ServerOptions.KeyFile

	if config.Global.ServerOptions.EnableTLS {
		srv := &http.Server{
			Addr:         addr,
			Handler:      router,
			TLSConfig:    &config.Global.ServerOptions.TLSConfig,
			TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0),
		}
		log.Fatal(srv.ListenAndServeTLS(certFile, keyFile))
	} else {
		log.Fatal(http.ListenAndServe(addr, router))
	}
}

func GenerateRouter() *mux.Router {
	r := mux.NewRouter()
	// Endpoints called by product owners

	r.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("./public"))))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))


	r.HandleFunc("/api/media", AddMedia).Methods(http.MethodPost)


	// media

	r.HandleFunc("/api/media", GetAllMedias).Methods(http.MethodGet)
	r.HandleFunc("/api/media/{id}", GetMediaByID).Methods(http.MethodGet)
	r.HandleFunc("/api/category/{category}", GetMediaByCategory).Methods(http.MethodGet)

	r.HandleFunc("/api/category", GetCategories).Methods(http.MethodGet)

	r.HandleFunc("/api/media/{id}", Update).Methods(http.MethodPut)

	r.HandleFunc("/api/media/{id}/delete", DeleteMediaByID).Methods(http.MethodDelete)
	r.HandleFunc("/api/media/delete", DeleteAllMedias).Methods(http.MethodDelete)


	// admins


	adminRouter := r.PathPrefix("/admin").Subrouter()
	//adminRouter.Use(AuthenticationMiddleware)
	adminRouter.HandleFunc("/admin/media", GetAllMedias).Methods(http.MethodGet)
	adminRouter.HandleFunc("/admin/media", AddMedia).Methods(http.MethodPost)
	adminRouter.HandleFunc("/admin/media/{id}", GetMediaByID).Methods(http.MethodGet)
	adminRouter.HandleFunc("/admin/media/{id}/delete", DeleteMediaByID).Methods(http.MethodDelete)

	r.HandleFunc("/api/ping", Ping).Methods(http.MethodPost)

	return r
}


