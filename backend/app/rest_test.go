package app

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"toussaint/backend/structs"

	"github.com/stretchr/testify/assert"
)

//func TestGetGame1(t *testing.T) {
//	srv := SetupRestApi("", 9999)
//	srv.ListenAndServe()
//	//go srv.ListenAndServe()
//	//defer srv.Shutdown(nil)
//	//
//	//resp, err := http.Get("http://localhost:9999/games?name=battlefield")
//	//assert.Nil(t, err)
//	//assert.Equal(t, http.StatusOK, resp.StatusCode)
//	//
//	//var games []Game
//	//body, err := ioutil.ReadAll(resp.Body)
//	//assert.Nil(t, err)
//	//
//	//err = json.Unmarshal(body, games)
//	//assert.Nil(t, err)
//	//
//	//t.Log(games)
//}

func TestGetGame2(t *testing.T) {

	defer removeDB(t)

	srv := httptest.NewServer(http.HandlerFunc(handleGetGame))

	resp, err := http.Get(fmt.Sprintf("%s%s", srv.URL, "/games?name=battlefield"))
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var games structs.GamesJSON
	body, err := ioutil.ReadAll(resp.Body)
	assert.Nil(t, err)

	err = json.Unmarshal(body, &games)
	assert.Nil(t, err)

	t.Log(games)
}

func TestGetGame3(t *testing.T) {

	defer removeDB(t)

	srv := httptest.NewServer(http.HandlerFunc(handleGetGame))

	resp, err := http.Get(fmt.Sprintf("%s%s", srv.URL, "/games?name="))
	assert.Nil(t, err)
	assert.Equal(t, http.StatusNotAcceptable, resp.StatusCode)

	body, err := ioutil.ReadAll(resp.Body)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(body))
}

func TestPutRegister1(t *testing.T) {

	defer removeDB(t)
	registerUser(t, http.StatusCreated)
	registerUser(t, http.StatusConflict)
}

func registerUser(t *testing.T, expectedCode int) {
	srv := httptest.NewServer(http.HandlerFunc(handlePutRegister))

	client := http.Client{}
	req, err := http.NewRequest("PUT", fmt.Sprintf("%s%s", srv.URL, "/register?client-id=1&client-type=telegram"), nil)
	assert.Nil(t, err)

	resp, err := client.Do(req)
	assert.Nil(t, err)

	assert.Equal(t, expectedCode, resp.StatusCode)

}

func TestPutNotify1(t *testing.T) {

	defer removeDB(t)

	gameId := "EP9000-CUSA11995_00-MARVELSSPIDERMAN"
	putNotify(t, http.StatusNotAcceptable, gameId)
	registerUser(t, http.StatusCreated)
	putNotify(t, http.StatusCreated, gameId)
}

func putNotify(t *testing.T, expectedStatus int, gameId string) {

	srv := httptest.NewServer(http.HandlerFunc(handlePutNotify))

	client := http.Client{}
	req, err := http.NewRequest("PUT",
		fmt.Sprintf("%s%s",
			srv.URL,
			fmt.Sprintf("/notify?client-id=1&client-type=telegram&game-id=%s", gameId),
		), nil)
	assert.Nil(t, err)

	resp, err := client.Do(req)
	assert.Nil(t, err)

	assert.Equal(t, expectedStatus, resp.StatusCode)
}

func TestDeleteNotify1(t *testing.T) {

	defer removeDB(t)

	srv := httptest.NewServer(http.HandlerFunc(handleDeleteNotify))

	gameId := "EP9000-CUSA11995_00-MARVELSSPIDERMAN"

	client := http.Client{}
	req, err := http.NewRequest("DELETE",
		fmt.Sprintf("%s%s",
			srv.URL,
			fmt.Sprintf("/notify?client-id=1&client-type=telegram&game-id=%s", gameId),
		), nil)
	assert.Nil(t, err)

	resp, err := client.Do(req)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusNotAcceptable, resp.StatusCode)

	registerUser(t, http.StatusCreated)
	putNotify(t, http.StatusCreated, gameId)

	resp, err = client.Do(req)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestGetList(t *testing.T) {

	defer removeDB(t)

	registerUser(t, http.StatusCreated)

	srv := httptest.NewServer(http.HandlerFunc(handleGetList))

	resp, err := http.Get(fmt.Sprintf("%s%s", srv.URL, "/list?client-id=1&client-type=telegram&request-type=all"))
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	putNotify(t, http.StatusCreated, "EP9000-CUSA11995_00-MARVELSSPIDERMAN")

	putNotify(t, http.StatusCreated, "EP9000-CUSA11995_00-0000000000MSMDLX")

	resp, err = http.Get(fmt.Sprintf("%s%s", srv.URL, "/list?client-id=1&client-type=telegram&request-type=all"))
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	data, err := ioutil.ReadAll(resp.Body)
	assert.Nil(t, err)

	t.Log(string(data))

	resp, err = http.Get(fmt.Sprintf("%s%s", srv.URL, "/list?client-id=1&client-type=telegram&request-type=sale"))
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	data, err = ioutil.ReadAll(resp.Body)
	assert.Nil(t, err)

	t.Log(string(data))
}

func getNotify(t *testing.T, expectedStatus int) structs.GamesJSON {
	srv := httptest.NewServer(http.HandlerFunc(handleGetNotify))

	resp, err := http.Get(fmt.Sprintf("%s%s", srv.URL, "/notify?client-id=1&client-type=telegram"))
	assert.Nil(t, err)
	assert.Equal(t, expectedStatus, resp.StatusCode)

	data, err := ioutil.ReadAll(resp.Body)
	assert.Nil(t, err)

	var games structs.GamesJSON
	err = json.Unmarshal(data, &games)
	assert.Nil(t, err)

	return games
}

func makeGameActive(t *testing.T, gameID string) {
	err := database.AddNotifications([]string{gameID})
	assert.Nil(t, err)
}

func TestGetNotifty(t *testing.T) {
	defer removeDB(t)

	registerUser(t, http.StatusCreated)

	g := getNotify(t, http.StatusOK)
	assert.Equal(t, 0, len(g.Games))

	gameID := "EP9000-CUSA11995_00-MARVELSSPIDERMAN"

	putNotify(t, http.StatusCreated, gameID)
	makeGameActive(t, gameID)

	g = getNotify(t, http.StatusOK)
	assert.Equal(t, 1, len(g.Games))
	t.Log(g)
}
