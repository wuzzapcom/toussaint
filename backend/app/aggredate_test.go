package app

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestMessageWithSales(t *testing.T) {

	parsedTime, err := time.Parse(timeFormat, "01.01.2019")
	assert.Nil(t, err)
	var games = make([]Game, 3)
	games[0] = Game{
		Name:      "Game 1",
		Price:     10,
		IsPlus:    true,
		SalePrice: 5,
		SaleEnd:   parsedTime,
	}
	games[1] = Game{
		Name:  "Game 2",
		Price: 20,
	}
	games[2] = Game{
		Name:      "Game 3",
		Price:     40,
		SalePrice: 20,
		SaleEnd:   parsedTime,
	}

	assert.Equal(t, `Ваши игры:
Game 1 стоила 10, цена со скидкой -- 5 рублей. Только для подписчиков PS Plus. Акция продлится до 01.01.2019
Game 2 стоит 20
Game 3 стоила 40, цена со скидкой -- 20 рублей. Акция продлится до 01.01.2019
`, GenerateMessage(games, false))

	assert.Equal(t, `Ваши игры:
Game 1 стоила 10, цена со скидкой -- 5 рублей. Только для подписчиков PS Plus. Акция продлится до 01.01.2019
Game 3 стоила 40, цена со скидкой -- 20 рублей. Акция продлится до 01.01.2019
`, GenerateMessage(games, true))

	assert.Equal(t, noGame, GenerateMessage([]Game{}, false))
	assert.Equal(t, noSaleGame, GenerateMessage([]Game{}, true))
}

func TestDescribeGames(t *testing.T) {
	parsedTime, err := time.Parse(timeFormat, "01.01.2019")
	assert.Nil(t, err)
	var games = make([]Game, 3)

	games[0] = Game{
		Id:        "1",
		Name:      "Game 1",
		Price:     10,
		IsPlus:    true,
		SalePrice: 5,
		SaleEnd:   parsedTime,
	}
	games[1] = Game{
		Id:    "2",
		Name:  "Game 2",
		Price: 20,
	}
	games[2] = Game{
		Id:        "3",
		Name:      "Game 3",
		Price:     40,
		SalePrice: 20,
		SaleEnd:   parsedTime,
	}

	ids, msgs := DescribeGames(games)
	assert.Equal(t, len(games), len(msgs))

	assert.Equal(t,
		"1",
		ids[0])
	assert.Equal(t,
		"Game 1, цена по скидке: 5 рублей. Полная стоимость: 10. Акция продлится до 01.01.2019. Только для подписчиков PS Plus.",
		msgs[0])

	assert.Equal(t,
		"2",
		ids[1])
	assert.Equal(t,
		"Game 2, 20 рублей",
		msgs[1])

	assert.Equal(t,
		"3",
		ids[2])
	assert.Equal(t,
		"Game 3, цена по скидке: 20 рублей. Полная стоимость: 40. Акция продлится до 01.01.2019.",
		msgs[2])
}
