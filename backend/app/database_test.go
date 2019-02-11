package app

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func removeDB(t *testing.T) {
	err := os.Remove(databaseName)
	if err != nil {
		t.Log(err)
	}
	database = NewDatabase()
}

func TestDB(t *testing.T) {

	defer removeDB(t)

	game := Game{
		Id:     "EP9000-CUSA11995_00-0000000000MSMDLX",
		Name:   "test",
		Price:  123,
		IsPlus: true,
	}

	err := database.AddGame(game)
	assert.Nil(t, err)

	client := &telegramClient{
		subscriptions: []string{},
		id:            "1",
	}

	err = database.AddUser(client)
	assert.Nil(t, err)

	err = database.AddGameToUser(game.Id, client.id, Telegram)
	assert.Nil(t, err)

	games, err := database.GetGamesForUser(client.id, Telegram, All)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(games))
	assert.Equal(t, game, games[0])

	err = database.DeleteGameFromUser(game.Id, client.id, Telegram)
	assert.Nil(t, err)

	games, err = database.GetGamesForUser(client.id, Telegram, All)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(games))
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

func TestGetGames(t *testing.T) {

	defer removeDB(t)

	game1 := Game{Id: "1"}
	err := database.AddGame(game1)
	assert.Nil(t, err)

	game2 := Game{Id: "2"}
	game3 := Game{Id: "3"}
	err = database.AddGames([]Game{game2, game3})

	games, err := database.GetGames()
	assert.Nil(t, err)
	assert.Equal(t, 3, len(games))
	assert.Equal(t, game1, games[0])
	assert.Equal(t, game2, games[1])
	assert.Equal(t, game3, games[2])
}

func TestGetIdByName(t *testing.T) {
	defer removeDB(t)

	game1 := Game{Id: "1", Name: "123"}
	err := database.AddGame(game1)
	assert.Nil(t, err)

	id, err := database.GetIDByGameName(game1.Name)
	assert.Nil(t, err)
	assert.Equal(t, game1.Id, id)
}
