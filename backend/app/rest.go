package app

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
	"toussaint/backend/structs"
)

var database = NewDatabase()

func handleGetGame(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	searched, err := SearchByName(name)
	if err != nil {
		log.Printf("[ERR] GET /game: %+v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	ids, msgs := DescribeGames(searched)

	games := structs.GamesJSON{
		Games: make([]structs.GamePair, len(ids)),
	}

	for i := 0; i < len(ids); i++ {
		games.Games[i].Id = ids[i]
		games.Games[i].Description = msgs[i]
	}

	marshalled, err := json.Marshal(games)
	if err != nil {
		log.Printf("[ERR] GET /game: %+v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(marshalled)
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

	gameId := r.URL.Query().Get("game-id")
	if gameId == "" {
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	game, err := SearchByID(gameId)
	if err != nil {
		if err != nil {
			log.Printf("[ERR] PUT /notify GetClientType: %+v", err)
			w.WriteHeader(http.StatusNotAcceptable)
			return
		}
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

	msg := GenerateMessage(games, requestType == Sale)

	w.Write([]byte(msg))
}

func handleGetUsers(w http.ResponseWriter, r *http.Request) {
	clientTypeStr := r.URL.Query().Get("client-type")
	if clientTypeStr == "" {
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	clientType, err := GetClientType(clientTypeStr)
	if err != nil {
		log.Printf("[ERR] GET /users GetClientType: %+v", err)
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	users, err := database.GetUsers(clientType)
	if err != nil {
		log.Printf("[ERR] GET /users database.GetUsers: %+v", err)
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	usersJSON := structs.UsersJSON{
		Ids: users,
	}

	marshalled, err := json.Marshal(usersJSON)
	if err != nil {
		log.Printf("[ERR] GET /users json.Marshal: %+v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(marshalled)
}

func SetupRestApi(host string, port int) *http.Server {
	router := mux.NewRouter()

	router.HandleFunc("/v1/game", handleGetGame).Methods("GET")
	router.HandleFunc("/v1/register", handlePutRegister).Methods("PUT")
	router.HandleFunc("/v1/notify", handlePutNotify).Methods("PUT")
	router.HandleFunc("/v1/notify", handlePutNotify).Methods("DELETE")
	router.HandleFunc("/v1/list", handleGetList).Methods("GET")
	router.HandleFunc("/v1/users", handleGetUsers).Methods("GET")

	srv := &http.Server{
		Handler: router,
		Addr:    fmt.Sprintf("%s:%d", host, port),
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	return srv
}
