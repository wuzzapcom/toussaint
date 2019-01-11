package app

import (
	"fmt"
	"strings"
)

const timeFormat = "02.01.2006"

const gameMessagePreface = "Ваши игры:"
const saleMessageFormat = "%s стоила %d, цена со скидкой -- %d рублей. %sАкция продлится до %s"
const saleMessagePlusAppendix = "Только для подписчиков PS Plus. "
const gameMessageFormat = "%s стоит %d"

const noGame = "У вас пока нет игр"
const noSaleGame = "У вас пока нет игр со скидкой"

const describeGameFormat = "%s, %d рублей"
const describeSaleGameFormat = "%s, цена по скидке: %d рублей. Полная стоимость: %d. Акция продлится до %s."
const describeSalePlusGameFormat = describeSaleGameFormat + " Только для подписчиков PS Plus."

func GenerateMessage(games []Game, filterForSales bool) string {

	if len(games) == 0 {
		if filterForSales {
			return noSaleGame
		}
		return noGame
	}

	builder := strings.Builder{}
	builder.WriteString(gameMessagePreface + "\n")

	for _, game := range games {
		if game.SalePrice != 0 {
			var isPlusAppendix string
			if game.IsPlus {
				isPlusAppendix = saleMessagePlusAppendix
			}
			builder.WriteString(fmt.Sprintf(
				saleMessageFormat,
				game.Name,
				game.Price,
				game.SalePrice,
				isPlusAppendix,
				game.SaleEnd.Format(timeFormat),
			) + "\n")
		} else {
			if !filterForSales {
				builder.WriteString(fmt.Sprintf(
					gameMessageFormat,
					game.Name,
					game.Price,
				) + "\n")
			}
		}
	}
	return builder.String()
}

func DescribeGames(games []Game) ([]string, []string) {
	ids := make([]string, 0)
	descriptions := make([]string, 0)
	for _, game := range games {
		ids = append(ids, game.Id)
		if game.SalePrice != 0 {
			if game.IsPlus {
				descriptions = append(descriptions, fmt.Sprintf(
					describeSalePlusGameFormat,
					game.Name,
					game.SalePrice,
					game.Price,
					game.SaleEnd.Format(timeFormat),
				))
				continue
			}
			descriptions = append(descriptions, fmt.Sprintf(
				describeSaleGameFormat,
				game.Name,
				game.SalePrice,
				game.Price,
				game.SaleEnd.Format(timeFormat),
			))
			continue
		}
		descriptions = append(descriptions, fmt.Sprintf(
			describeGameFormat,
			game.Name,
			game.Price,
		))
	}
	return ids, descriptions
}
