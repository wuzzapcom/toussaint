package app

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetGame(t *testing.T) {

	defer removeDB(t)

	srv := httptest.NewServer(http.HandlerFunc(handleGetGames))

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
	srv := httptest.NewServer(http.HandlerFunc(handlePutUsers))

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

	srv := httptest.NewServer(http.HandlerFunc(handlePutNotifications))

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

	srv := httptest.NewServer(http.HandlerFunc(handleDeleteNotifications))

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
