package main

import (
	"elnewsAPI/api"
	"elnewsAPI/loads"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
)

const (
	Port = ":8181"
	Key = "key_here"
	Cert = "cert_here"
)

func init() {

	log.SetFlags(log.LstdFlags | log.Lshortfile)
	loads.LoadConfiguration()
}


func main() {
	var allowedOrigins []string

	if loads.ElDevelopment {
		allowedOrigins = []string{"https://elnewspaper.com", "https://www.elnewspaper.com", "http://localhost:8080", "http://127.0.0.1:8080"}
	} else {
		allowedOrigins = []string{"https://elnewspaper.com", "https://www.elnewspaper.com"}
	}

	r := mux.NewRouter()

	// Events SubRouter
	eventsRouter := r.PathPrefix("/events").Subrouter()

	eventsRouter.Methods("GET").Path("/first_posts_json").HandlerFunc(api.PostsUnionJSON)
	eventsRouter.Methods("GET").Path("/next_posts_json").HandlerFunc(api.PostsUnionJSON)
	eventsRouter.Methods("GET").Path("/saying_today_array").HandlerFunc(api.SayingTodayJSONArray)
	eventsRouter.Methods("POST").Path("/weather_today_normal").HandlerFunc(api.WeatherDongNeTodayJSON)

	// Accounts SubRouter
	accountsRouter := r.PathPrefix("/accounts").Subrouter()

	accountsRouter.Methods("POST").Path("/get_user_info").HandlerFunc(api.GetUserInfo)
	accountsRouter.Methods("GET").Path("/oauth/login_url").HandlerFunc(api.GetLoginURLJson)
	accountsRouter.Methods("GET").Path("/oauth/{my_type}/callback").HandlerFunc(api.Authenticate)
	accountsRouter.Methods("POST").Path("/logout").HandlerFunc(api.LogoutSession)

	//headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	originsOk := handlers.AllowedOrigins(allowedOrigins)
	credentialsOk := handlers.AllowCredentials()

	server := handlers.CORS(originsOk, credentialsOk)(r)

	log.Println("starting...")

	if loads.ElDevelopment {
		err := http.ListenAndServe(Port, server)
		if err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	} else {
		err := http.ListenAndServeTLS(Port, Cert, Key, server)
		if err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}


}
