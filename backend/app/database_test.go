package app

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDB(t *testing.T) {
	databaseName = "toussaint_test.db"
	database := NewDatabase()
	game := Game{
		Id: "1",
		Name: "test",
		Price: 123,
		IsPlus: true,
	}
	err := database.AddGame(game)
	assert.Nil(t, err)

	client := &TelegramClient{
		subscriptions: []string{},
		id: "1",
	}

	err = database.AddUser(client)
	assert.Nil(t, err)

	err = database.AddGameToUser(game.Id, client.id, Telegram)
	assert.Nil(t, err)

	res, err := database.GetGame(game.Id)
	assert.Nil(t, err)
	assert.Equal(t, game.Name, res.Name)
}
