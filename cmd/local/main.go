package main

import (
	"coursework3/internal/antenna"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func mainHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "web/index.html")
}

func makePhysicsHandler(as *antenna.AntennaSet) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		indS := vars["ind"]
		ind, err := strconv.Atoi(indS)
		if err != nil {
			log.Println(err)
			w.WriteHeader(400)
			return
		}
		json, err := as.Marshal(ind)
		if err != nil {
			log.Println(err)
			w.WriteHeader(500)
		} else {
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(200)
			w.Write(json)
		}
	}
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
	http.Redirect(w, r, "https://www.shtjnk.ru"+r.RequestURI, http.StatusMovedPermanently)
}

func wwwHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "https://www.shtjnk.ru"+r.RequestURI, http.StatusMovedPermanently)
}

func main() {
	set := antenna.MakeSet(
		antenna.NewSlotAntenna(15, 61, 81, 0.31, 0.01),
		antenna.NewSlotAntenna(11, 61, 81, 0.23, 0.01),
		antenna.NewSlotAntenna(5, 61, 81, 0.11, 0.01))
	r := mux.NewRouter()
	api := r.PathPrefix("/api/").Subrouter()

	api.HandleFunc("/get_antenna/{ind}", makePhysicsHandler(set)).Methods("GET")

	static := r.PathPrefix("/static/")

	fileServer := http.FileServer(http.Dir("web/static"))
	static.Handler(http.StripPrefix("/static", fileServerMiddleware(fileServer)))
	r.HandleFunc("/", mainHandler).Methods("GET")

	server := http.Server{
		Addr:    "localhost:8081",
		Handler: r,
	}
	err := server.ListenAndServe()

	if err != nil {
		log.Fatalf("ListenAndServe error: %v", err)
	}
}
