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

func serverFile(w http.ResponseWriter, req * http.Request) {
	path:=req.URL.Path[1:]

	http.ServeFile(w, req, "./public/"+path)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, fmt.Sprintf("./%s/index.html", "./public"))
}

func GenerateRouter() *mux.Router {
	r := mux.NewRouter()
	// Endpoints called by product owners

	r.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("./public"))))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))


	r.HandleFunc("/api/media", AddMedia).Methods(http.MethodPost)

	r.HandleFunc("/api/media", GetAllMedias).Methods(http.MethodGet)
	r.HandleFunc("/api/media/{id}", GetMediaByID).Methods(http.MethodGet)
	r.HandleFunc("/api/media/category/{category}", GetMediaByCategory).Methods(http.MethodGet)

//	r.HandleFunc("/api/media/{id}", Update).Methods(http.MethodPut)

	r.HandleFunc("/api/media/{id}/delete", DeleteMedia).Methods(http.MethodDelete)
	r.HandleFunc("/api/media/delete", DeleteAllMedias).Methods(http.MethodDelete)

	adminRouter := r.PathPrefix("/admin").Subrouter()
	//adminRouter.Use(AuthenticationMiddleware)
	adminRouter.HandleFunc("/media", GetAllMedias).Methods(http.MethodGet)
	adminRouter.HandleFunc("/media", AddMedia).Methods(http.MethodPost)
	adminRouter.HandleFunc("/media/{id}", GetMediaByID).Methods(http.MethodGet)
	adminRouter.HandleFunc("/media/{id}/delete", DeleteMedia).Methods(http.MethodDelete)

	r.HandleFunc("/license/ping", Ping).Methods(http.MethodPost)

	return r
}


