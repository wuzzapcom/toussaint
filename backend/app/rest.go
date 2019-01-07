package app

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

var database = NewDatabase()

func handleGetGame(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	searched, err := Search(name)
	if err != nil {
		log.Printf("[ERR] GET /game: %+v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	games, err := json.Marshal(searched)
	if err != nil {
		log.Printf("[ERR] GET /game: %+v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(games)
}

func handlePutRegister(w http.ResponseWriter, r *http.Request) {
	clientId := r.URL.Query().Get("client-id")
	if clientId == "" {
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	clientTypeStr := r.URL.Query().Get("client-type")
	if clientTypeStr == "" {
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	clientType, err := GetClientType(clientTypeStr)
	if err != nil {
		log.Printf("[ERR] PUT /register GetClientType: %+v", err)
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	var client Client
	switch clientType {
	case Telegram:
		client = NewTelegramClient(clientId)
	default:
		log.Printf("[ERR] PUT /register: unhandled client type %+v", clientType)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = database.AddUser(client)
	if err != nil {
		log.Printf("[ERR] PUT /register AddUser: %+v", err)
		w.WriteHeader(http.StatusConflict)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func handlePutNotify(w http.ResponseWriter, r *http.Request) {

	clientId := r.URL.Query().Get("client-id")
	if clientId == "" {
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	clientTypeStr := r.URL.Query().Get("client-type")
	if clientTypeStr == "" {
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	clientType, err := GetClientType(clientTypeStr)
	if err != nil {
		log.Printf("[ERR] PUT /notify GetClientType: %+v", err)
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("[ERR] PUT /notify ioutil.ReadAll: %+v", err)
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	var game Game
	err = json.Unmarshal(body, &game)
	if err != nil {
		log.Printf("[ERR] PUT /notify json.Unmarshal: %+v", err)
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	err = database.AddGame(game)
	if err != nil {
		log.Printf("[ERR] PUT /notify AddGame: %+v", err)
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	err = database.AddGameToUser(game.Id, clientId, clientType)
	if err != nil {
		log.Printf("[ERR] PUT /notify AddGameToUser: %+v", err)
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	w.WriteHeader(http.StatusCreated)
}


func handleDeleteNotify(w http.ResponseWriter, r *http.Request) {

	gameId := r.URL.Query().Get("game-id")
	if gameId == "" {
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	clientId := r.URL.Query().Get("client-id")
	if clientId == "" {
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	clientTypeStr := r.URL.Query().Get("client-type")
	if clientTypeStr == "" {
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	clientType, err := GetClientType(clientTypeStr)
	if err != nil {
		log.Printf("[ERR] DELETE /notify GetClientType: %+v", err)
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	err = database.DeleteGameFromUser(gameId, clientId, clientType)
	if err != nil {
		log.Printf("[ERR] DELETE /notify DeleteGameFromUser: %+v", err)
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func handleGetList(w http.ResponseWriter, r *http.Request) {

	clientId := r.URL.Query().Get("client-id")
	if clientId == "" {
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	clientTypeStr := r.URL.Query().Get("client-type")
	if clientTypeStr == "" {
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	clientType, err := GetClientType(clientTypeStr)
	if err != nil {
		log.Printf("[ERR] GET /list GetClientType: %+v", err)
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	requestTypeStr := r.URL.Query().Get("request-type")
	if requestTypeStr == "" {
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	requestType, err := GetRequestType(requestTypeStr)
	if err != nil {
		log.Printf("[ERR] GET /list GetRequestType: %+v", err)
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	games, err := database.GetGamesForUser(clientId, clientType, requestType)
	if err != nil {
		log.Printf("[ERR] GET /list GetGamesForUser: %+v", err)
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	data, err := json.Marshal(games)
	if err != nil {
		log.Printf("[ERR] GET /list json.Marshal: %+v", err)
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	w.Write(data)

}


func SetupRestApi(host string, port int) *http.Server {
	router := mux.NewRouter()

	router.HandleFunc("/game", handleGetGame).Methods("GET")
	router.HandleFunc("/register", handlePutRegister).Methods("PUT")
	router.HandleFunc("/notify", handlePutNotify).Methods("PUT")
	router.HandleFunc("/notify", handlePutNotify).Methods("DELETE")
	router.HandleFunc("/list", handleGetList).Methods("GET")

	srv := &http.Server{
		Handler:      router,
		Addr:         fmt.Sprintf("%s:%d", host, port),
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	return srv
}