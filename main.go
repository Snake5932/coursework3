package main

import (
	"coursework3/internal/certs"
	"crypto/tls"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strings"
)

func mainHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "/srv/coursework/web/index.html")
}

func fileServerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/") {
			http.NotFound(w, r)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func redirectToTls(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "https://shtjnk.ru"+r.RequestURI, http.StatusMovedPermanently)
}

func main() {
	r := mux.NewRouter()
	//api := r.PathPrefix("/api/").Subrouter()
	static := r.PathPrefix("/static/")
	fileServer := http.FileServer(http.Dir("/srv/coursework/web/static"))
	static.Handler(http.StripPrefix("/static", fileServerMiddleware(fileServer)))
	r.HandleFunc("/", mainHandler).Methods("GET")
	cert, _ := tls.X509KeyPair([]byte(certs.CertChain), []byte(certs.Key))
	server := http.Server{
		Addr:      "185.188.183.121:443",
		Handler:   r,
		TLSConfig: &tls.Config{Certificates: []tls.Certificate{cert}},
	}
	go func() {
		if err := http.ListenAndServe("185.188.183.121:80", http.HandlerFunc(redirectToTls)); err != nil {
			log.Fatalf("ListenAndServe error: %v", err)
		}
	}()
	err := server.ListenAndServeTLS("", "")
	if err != nil {
		log.Fatalf("ListenAndServe error: %v", err)
	}
}
