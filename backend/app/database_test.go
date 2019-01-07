package app

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func removeDB() {
	err := os.Remove(databaseName)
	if err != nil {
		panic(err)
	}
}

func TestDB(t *testing.T) {

	defer removeDB()

	game := Game{
		Id: "1",
		Name: "test",
		Price: 123,
		IsPlus: true,
	}
	err := database.AddGame(game)
	assert.Nil(t, err)

	client := &telegramClient{
		subscriptions: []string{},
		id: "1",
	}

	err = database.AddUser(client)
	assert.Nil(t, err)

	err = database.AddGameToUser(game.Id, client.id, Telegram)
	assert.Nil(t, err)

	err = database.DeleteGameFromUser(game.Id, client.id, Telegram)
	assert.Nil(t, err)

	res, err := database.GetGame(game.Id)
	assert.Nil(t, err)
	assert.Equal(t, game.Name, res.Name)

	err = database.DeleteGame(game.Id)
	assert.Nil(t, err)
	assert.Equal(t, game.Name, res.Name)
}
