package app

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func removeDB(t *testing.T) {
	err := os.Remove(databaseName)
	if err != nil {
		t.Log(err)
	}
}

func TestDB(t *testing.T) {

	defer removeDB(t)

	game := Game{
		Id:     "EP9000-CUSA11995_00-0000000000MSMDLX",
		Name:   "test",
		Price:  123,
		IsPlus: true,
	}

	client := &telegramClient{
		subscriptions: []string{},
		id:            "1",
	}

	err := database.AddUser(client)
	assert.Nil(t, err)

	err = database.AddGameToUser(game.Id, client.id, Telegram)
	assert.Nil(t, err)

	err = database.DeleteGameFromUser(game.Id, client.id, Telegram)
	assert.Nil(t, err)
}

func TestGetUsers(t *testing.T) {

	defer removeDB(t)

	users, err := database.GetUsers(Telegram)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(users))

	err = database.AddUser(&telegramClient{
		subscriptions: []string{},
		id:            "1",
	})
	assert.Nil(t, err)

	err = database.AddUser(&telegramClient{
		subscriptions: []string{},
		id:            "2",
	})
	assert.Nil(t, err)

	users, err = database.GetUsers(Telegram)
	assert.Nil(t, err)
	assert.Equal(t, 2, len(users))
	assert.Equal(t, "1", users[0])
	assert.Equal(t, "2", users[1])
}
