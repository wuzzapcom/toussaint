package app

import (
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestUpdater(t *testing.T) {

	defer removeDB(t)

	btw, err := time.ParseDuration("1s")
	assert.Nil(t, err)

	var done int32

	updater := Updater{
		betweenUpdates: btw,
		updateTime:     time.Now().Add(2 * time.Second),
		updaterFunc: func() {
			atomic.AddInt32(&done, 1)
		},
		stop:     make(chan bool, 0),
		finished: make(chan bool),
	}

	go updater.Start()

	time.Sleep(time.Second * 5)
	updater.Stop()
	assert.True(t, done > int32(0))
}

func TestUpdaterFunc(t *testing.T) {

	defer removeDB(t)

	game := Game{Id: "EP9000-CUSA11995_00-MARVELSSPIDERMAN"}

	err := database.AddGame(game)
	assert.Nil(t, err)

	update()

	loaded, err := database.GetGame(game.Id)
	assert.Nil(t, err)

	assert.NotEqual(t, game.Name, loaded.Name)
}

func TestUpdate(t *testing.T) {

	defer removeDB(t)

	// Prepare environment
	game := Game{
		Name:  "Game",
		Id:    "EP9000-CUSA11995_00-MARVELSSPIDERMAN",
		Price: 400,
	}
	user := NewTelegramClient("1")

	// Prepare database state
	assert.Nil(t, database.AddGame(game))
	assert.Nil(t, database.AddUser(user))
	assert.Nil(t, database.AddGameToUser(game.Id, user.ID(), user.Type()))

	// Mock SearchByID since it is used in update() to trigger sale notification
	SearchByID = func(id string) (Game, error) {
		return Game{
			Id:        game.Id,
			Name:      game.Name,
			Price:     game.Price,
			SalePrice: 200,
		}, nil
	}

	// run http server with Notifier
	srv := httptest.NewServer(notifier)

	finished := make(chan bool)

	// run http client that will establish connection with server
	// ant wait until any data will be received
	var cl = func() {
		resp, err := http.Get(srv.URL)
		assert.Nil(t, err)
		// data array MUST be initialized!
		var data = make([]byte, 1000)
		n, err := resp.Body.Read(data)
		assert.Nil(t, err)
		// wait until someone sends some data here
		for n == 0 {
			t.Logf("data is empty, try again...")
			n, err = resp.Body.Read(data)
			assert.Nil(t, err)
		}
		t.Logf("HTTP client received: %s", string(data))
		finished <- true
	}

	// start client in goroutine
	go cl()

	// toggle update process
	update()

	// wait until everything is OK
	<-finished

	// return old state for future tests
	SearchByID = searchByID

}
