package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
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

	defer removeDB()

	srv := httptest.NewServer(http.HandlerFunc(handleGetGame))

	resp, err := http.Get(fmt.Sprintf("%s%s", srv.URL, "/games?name=battlefield"))
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var games []Game
	body, err := ioutil.ReadAll(resp.Body)
	assert.Nil(t, err)

	err = json.Unmarshal(body, &games)
	assert.Nil(t, err)

	t.Log(games)
}

func TestGetGame3(t *testing.T) {

	defer removeDB()

	srv := httptest.NewServer(http.HandlerFunc(handleGetGame))

	resp, err := http.Get(fmt.Sprintf("%s%s", srv.URL, "/games?name="))
	assert.Nil(t, err)
	assert.Equal(t, http.StatusNotAcceptable, resp.StatusCode)

	body, err := ioutil.ReadAll(resp.Body)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(body))
}

func TestPutRegister1(t *testing.T) {

	defer removeDB()
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

	defer removeDB()

	game := Game{
		Id: "some-id",
	}
	putNotify(t, http.StatusNotAcceptable, game)
	registerUser(t, http.StatusCreated)
	putNotify(t, http.StatusCreated, game)
}

func putNotify(t *testing.T, expectedStatus int, game Game) {

	srv := httptest.NewServer(http.HandlerFunc(handlePutNotify))

	data, err := json.Marshal(game)
	assert.Nil(t, err)

	client := http.Client{}
	req, err := http.NewRequest("PUT", fmt.Sprintf("%s%s", srv.URL, "/notify?client-id=1&client-type=telegram"), bytes.NewReader(data))
	assert.Nil(t, err)

	resp, err := client.Do(req)
	assert.Nil(t, err)

	assert.Equal(t, expectedStatus, resp.StatusCode)
}

func TestDeleteNotify1(t *testing.T) {

	defer removeDB()

	srv := httptest.NewServer(http.HandlerFunc(handleDeleteNotify))

	game := Game{
		Id: "some-id",
	}

	data, err := json.Marshal(game)
	assert.Nil(t, err)

	client := http.Client{}
	req, err := http.NewRequest("PUT", fmt.Sprintf("%s%s", srv.URL, "/notify?game-id=1&client-id=1&client-type=telegram"), bytes.NewReader(data))
	assert.Nil(t, err)

	resp, err := client.Do(req)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusNotAcceptable, resp.StatusCode)

	putNotify(t, http.StatusCreated, game)

	resp, err = client.Do(req)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestGetList(t *testing.T) {

	defer removeDB()

	registerUser(t, http.StatusCreated)

	srv := httptest.NewServer(http.HandlerFunc(handleGetList))

	resp, err := http.Get(fmt.Sprintf("%s%s", srv.URL, "/list?client-id=1&client-type=telegram&request-type=all"))
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	putNotify(t, http.StatusCreated, Game{
		Id: "1",
		SalePrice: 10,
	})

	putNotify(t, http.StatusCreated, Game{
		Id: "2",
	})

	resp, err = http.Get(fmt.Sprintf("%s%s", srv.URL, "/list?client-id=1&client-type=telegram&request-type=all"))
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	data, err := ioutil.ReadAll(resp.Body)
	assert.Nil(t, err)

	var games []Game

	err = json.Unmarshal(data, &games)
	assert.Nil(t, err)

	assert.Equal(t, 2, len(games))

	resp, err = http.Get(fmt.Sprintf("%s%s", srv.URL, "/list?client-id=1&client-type=telegram&request-type=sale"))
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	data, err = ioutil.ReadAll(resp.Body)
	assert.Nil(t, err)

	games = nil

	err = json.Unmarshal(data, &games)
	assert.Nil(t, err)

	assert.Equal(t, 1, len(games))
}
