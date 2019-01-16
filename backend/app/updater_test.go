package app

import (
	"github.com/stretchr/testify/assert"
	"sync/atomic"
	"testing"
	"time"
)

func TestUpdater(t *testing.T) {

	defer removeDB(t)

	btw, err := time.ParseDuration("1s")
	assert.Nil(t, err)

	upd := time.Now().Add(2 * time.Second)
	assert.Nil(t, err)

	var done int32

	updater := Updater{
		betweenUpdates: btw,
		updateTime:     upd,
		updaterFunc: func() {
			atomic.AddInt32(&done, 1)
		},
		stop:     make(chan bool, 0),
		finished: make(chan bool),
	}

	go updater.Start()

	time.Sleep(time.Second * 3)
	updater.Stop()
	assert.Equal(t, int32(1), done)
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
