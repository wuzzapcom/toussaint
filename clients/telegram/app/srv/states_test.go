package srv

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"toussaint/backend/structs"

	"github.com/stretchr/testify/assert"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func TestHandleNoStateStart1(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(http.StatusCreated)
	}))

	APIEndpoint = srv.URL

	msg := &tgbotapi.Message{
		Text: "/start",
		From: &tgbotapi.User{
			ID: 1,
		},
	}

	str, shouldCache, state, p, err := handleNoState(msg, nil)
	assert.Nil(t, err)
	assert.Equal(t, register_ok_msg_ru, str)
	assert.False(t, shouldCache)
	assert.Equal(t, NO_STATE, state)
	assert.Nil(t, p)
}

func TestHandleNoStateStart2(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(http.StatusInternalServerError)
	}))

	APIEndpoint = srv.URL

	msg := &tgbotapi.Message{
		Text: "/start",
		From: &tgbotapi.User{
			ID: 1,
		},
	}

	str, shouldCache, state, p, err := handleNoState(msg, nil)
	assert.NotNil(t, err)
	assert.Equal(t, register_fail_msg_ru, str)
	assert.False(t, shouldCache)
	assert.Equal(t, NO_STATE, state)
	assert.Nil(t, p)
}

func TestHandleNoStateList(t *testing.T) {
	testMsg := "OK"

	srv := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.Write([]byte(testMsg))
	}))

	APIEndpoint = srv.URL

	msg := &tgbotapi.Message{
		Text: "/list",
		From: &tgbotapi.User{
			ID: 1,
		},
	}

	str, shouldCache, state, p, err := handleNoState(msg, nil)
	assert.Nil(t, err)
	assert.Equal(t, testMsg, str)
	assert.False(t, shouldCache)
	assert.Equal(t, NO_STATE, state)
	assert.Nil(t, p)
}

func TestHandleNoStateSearch(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		games := structs.GamesJSON{}
		games.Games = make([]structs.GamePair, 1)
		games.Games[0].Id = "ID"
		games.Games[0].Description = "Description"
		b, err := json.Marshal(games)
		assert.Nil(t, err)

		rw.Write(b)
	}))

	APIEndpoint = srv.URL

	msg := &tgbotapi.Message{
		Text: "/search Apex Legends",
		From: &tgbotapi.User{
			ID: 1,
		},
	}

	str, shouldCache, state, payload, err := handleNoState(msg, nil)
	assert.Nil(t, err)
	assert.Equal(t, "1) Description\n2) Отменить", str)
	assert.True(t, shouldCache)
	assert.Equal(t, SEARCH_GAME_WAIT_GAME, state)
	_, ok := payload.([]string)
	assert.True(t, ok)
}

func TestHandleSearchGameWaitGame(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(http.StatusCreated)
	}))

	APIEndpoint = srv.URL

	msg := &tgbotapi.Message{
		Text: "1",
		From: &tgbotapi.User{
			ID: 1,
		},
	}

	payload := make([]string, 1)
	payload[0] = "ID"
	str, shouldCache, state, _, err := handleSearchGameWaitGame(msg, payload)
	assert.Nil(t, err)
	assert.Equal(t, notify_ok_msg_ru, str)
	assert.False(t, shouldCache)
	assert.Equal(t, NO_STATE, state)
}
