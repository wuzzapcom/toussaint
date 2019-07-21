package app

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"toussaint/backend/structs"

	"github.com/gorilla/mux"
)

var notifier = NewNotifier()

var database = NewDatabase()

func handleSearch(w http.ResponseWriter, r *http.Request) {
	name, httpErr := parseQueryParameter(r.URL.Query(), "name")
	if httpErr != nil {
		log.Printf("[ERR] GET /search: %+v", httpErr.Message)
		w.WriteHeader(httpErr.Code)
		return
	}

	searched, err := SearchByName(name)
	if err != nil {
		log.Printf("[ERR] GET /search: %+v", err)
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
		log.Printf("[ERR] GET /search: %+v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(marshalled)
}

func handleGetGames(w http.ResponseWriter, r *http.Request) {
	name, httpErr := parseQueryParameter(r.URL.Query(), "name")
	if httpErr != nil {
		log.Printf("[ERR] PUT /games: %+v", httpErr.Message)
		w.WriteHeader(httpErr.Code)
		return
	}

	id, err := database.GetIDByGameName(name)
	if err != nil {
		log.Printf("[ERR] GET /games: %+v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if id == "" {
		log.Printf("[ERR] GET /games: game not found")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Write([]byte(id))
}

func handlePutUsers(w http.ResponseWriter, r *http.Request) {
	clientID, httpErr := parseQueryParameter(r.URL.Query(), "client-id")
	if httpErr != nil {
		log.Printf("[ERR] PUT /users: %+v", httpErr.Message)
		w.WriteHeader(httpErr.Code)
		return
	}

	clientType, httpErr := parseClientType(r.URL.Query())
	if httpErr != nil {
		log.Printf("[ERR] PUT /users: %+v", httpErr.Message)
		w.WriteHeader(httpErr.Code)
		return
	}

	var client Client
	switch clientType {
	case Telegram:
		client = NewTelegramClient(clientID)
	default:
		log.Printf("[ERR] PUT /register: unhandled client type %+v", clientType)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err := database.AddUser(client)
	if err != nil {
		log.Printf("[ERR] PUT /register AddUser: %+v", err)
		w.WriteHeader(http.StatusConflict)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func handlePutNotifications(w http.ResponseWriter, r *http.Request) {

	clientID, httpErr := parseQueryParameter(r.URL.Query(), "client-id")
	if httpErr != nil {
		log.Printf("[ERR] PUT /notifications: %+v", httpErr.Message)
		w.WriteHeader(httpErr.Code)
		return
	}

	clientType, httpErr := parseClientType(r.URL.Query())
	if httpErr != nil {
		log.Printf("[ERR] PUT /notifications: %+v", httpErr.Message)
		w.WriteHeader(httpErr.Code)
		return
	}

	gameID, httpErr := parseQueryParameter(r.URL.Query(), "game-id")
	if httpErr != nil {
		log.Printf("[ERR] PUT /notifications: %+v", httpErr.Message)
		w.WriteHeader(httpErr.Code)
		return
	}

	game, err := SearchByID(gameID)
	if err != nil {
		if err != nil {
			log.Printf("[ERR] PUT /notifications GetClientType: %+v", err)
			w.WriteHeader(http.StatusNotAcceptable)
			return
		}
	}

	err = database.AddGame(game)
	if err != nil {
		log.Printf("[ERR] PUT /notifications AddGame: %+v", err)
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	err = database.AddGameToUser(game.Id, clientID, clientType)
	if err != nil {
		log.Printf("[ERR] PUT /notifications AddGameToUser: %+v", err)
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func handleDeleteNotifications(w http.ResponseWriter, r *http.Request) {

	gameID, httpErr := parseQueryParameter(r.URL.Query(), "game-id")
	if httpErr != nil {
		log.Printf("[ERR] DELETE /notifications: %+v", httpErr.Message)
		w.WriteHeader(httpErr.Code)
		return
	}

	clientID, httpErr := parseQueryParameter(r.URL.Query(), "client-id")
	if httpErr != nil {
		log.Printf("[ERR] DELETE /notifications: %+v", httpErr.Message)
		w.WriteHeader(httpErr.Code)
		return
	}

	clientType, httpErr := parseClientType(r.URL.Query())
	if httpErr != nil {
		log.Printf("[ERR] DELETE /notifications: %+v", httpErr.Message)
		w.WriteHeader(httpErr.Code)
		return
	}

	err := database.DeleteGameFromUser(gameID, clientID, clientType)
	if err != nil {
		log.Printf("[ERR] DELETE /notify DeleteGameFromUser: %+v", err)
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func handleGetList(w http.ResponseWriter, r *http.Request) {

	clientID, httpErr := parseQueryParameter(r.URL.Query(), "client-id")
	if httpErr != nil {
		log.Printf("[ERR] GET /list %+v", httpErr.Message)
		w.WriteHeader(httpErr.Code)
		return
	}

	clientType, httpErr := parseClientType(r.URL.Query())
	if httpErr != nil {
		log.Printf("[ERR] GET /list: %+v", httpErr.Message)
		w.WriteHeader(httpErr.Code)
		return
	}

	requestType, httpErr := parseRequestType(r.URL.Query())
	if httpErr != nil {
		log.Printf("[ERR] GET /list: %+v", httpErr.Message)
		w.WriteHeader(httpErr.Code)
	}

	games, err := database.GetGamesForUser(clientID, clientType, requestType)
	if err != nil {
		log.Printf("[ERR] GET /list GetGamesForUser: %+v", err)
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	msg := GenerateMessage(games, requestType == Sale)

	w.Write([]byte(msg))
}

func handleGetUsers(w http.ResponseWriter, r *http.Request) {
	clientType, httpErr := parseClientType(r.URL.Query())
	if httpErr != nil {
		log.Printf("[ERR] GET /users: %+v", httpErr.Message)
		w.WriteHeader(httpErr.Code)
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

func handleNotificationsTrigger(w http.ResponseWriter, r *http.Request) {
	clientType, httpErr := parseClientType(r.URL.Query())
	if httpErr != nil {
		log.Printf("[ERR] GET /notifications/trigger: %+v", httpErr.Message)
		w.WriteHeader(httpErr.Code)
		return
	}
	clientID, httpErr := parseQueryParameter(r.URL.Query(), "client-id")
	if httpErr != nil {
		log.Printf("[ERR] GET /notifications/trigger: %+v", httpErr.Message)
		w.WriteHeader(httpErr.Code)
		return
	}
	notifier.NotifyUser(clientType, structs.UserNotification{
		UserID: clientID,
		Games: structs.GamesJSON{
			Games: []structs.GamePair{
				structs.GamePair{
					Id:          "test_game",
					Description: "test_game",
				},
			},
		},
	})
	w.WriteHeader(http.StatusOK)
}

func SetupRestAPI(host string, port int, debug bool) *http.Server {
	router := mux.NewRouter()

	router.HandleFunc("/v1/search", handleSearch).Methods("GET")
	router.HandleFunc("/v1/games", handleGetGames).Methods("GET")
	router.HandleFunc("/v1/notifications", handlePutNotifications).Methods("PUT")
	router.HandleFunc("/v1/notifications", handleDeleteNotifications).Methods("DELETE")
	router.HandleFunc("/v1/notifications", notifier.ServeHTTP).Methods("GET")
	router.HandleFunc("/v1/list", handleGetList).Methods("GET")
	router.HandleFunc("/v1/users", handleGetUsers).Methods("GET")
	router.HandleFunc("/v1/users", handlePutUsers).Methods("PUT")

	if debug {
		log.Printf("setting GET /v1/notifications/trigger ...")
		router.HandleFunc("/v1/notifications/trigger", handleNotificationsTrigger).Methods("GET")
	}

	srv := &http.Server{
		Handler: router,
		Addr:    fmt.Sprintf("%s:%d", host, port),
	}

	return srv
}
