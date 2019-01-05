package app

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"time"
)

func Exit(err error) {
	fmt.Println(err)
	os.Exit(-1)
}

func BuildURLForName(name string) (string, error) {
	fetchSize := "99999"
	u, err := url.Parse("https://store.playstation.com")
	if err != nil {
		return "", err
	}
	u.Path += fmt.Sprintf("/chihiro-api/bucket-search/Ru/ru/19/%s", name)
	params := url.Values{}
	params.Add("size", fetchSize)
	params.Add("start", "0")
	u.RawQuery = params.Encode()
	return u.String(), nil
}

type Game struct {
	Id        string    `json:"-"`
	Name      string    `json:"name"`
	Price     int       `json:"price"`
	IsPlus    bool      `json:"is_plus"`
	SalePrice int       `json:"sale_price"`
	SaleEnd   time.Time `json:"sale_end"`
}

func checkGameType(tp string) bool {
	switch tp {
	case "Полная версия":
		return true
	case "Комплект":
		return true
	default:
		return false
	}
}

//takes interface{}, converts it to []map[string]interface{} and takes first elem
func takeFirstMap(m interface{}) (map[string]interface{}, error) {
	temp1, ok := m.([]interface{})
	if !ok {
		return nil, errors.New("failed convertion to []interface{}")
	}
	if len(temp1) != 1 {
		return nil, errors.New(fmt.Sprintf("size of []interface{} is %d", len(temp1)))
	}
	res, ok := temp1[0].(map[string]interface{})
	if !ok {
		return nil, errors.New("failed convertion to map[string]interface{}")
	}
	return res, nil
}

func parseSearchAnswer(data []byte) ([]Game, error) {
	var answer map[string]*json.RawMessage
	err := json.Unmarshal(data, &answer)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(*answer["categories"], &answer)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(*answer["games"], &answer)
	if err != nil {
		return nil, err
	}

	var games []map[string]interface{}

	err = json.Unmarshal(*answer["links"], &games)
	if err != nil {
		return nil, err
	}

	resultGames := make([]Game, 0)

	for _, game := range games {
		currentGame := Game{}

		contentType, ok := game["game_contentType"].(string)
		if !ok {
			continue
		}
		if !checkGameType(contentType) {
			continue
		}

		currentGame.Id = game["id"].(string)
		currentGame.Name = game["name"].(string)

		skus, err := takeFirstMap(game["skus"])
		if err != nil {
			return nil, err
		}

		price := skus["price"].(float64)
		currentGame.Price = int(price) / 100

		rewards, err := takeFirstMap(skus["rewards"])
		if err != nil {
			//no sale
			resultGames = append(resultGames, currentGame)
			continue
		}

		salePrice := rewards["price"].(float64)
		currentGame.SalePrice = int(salePrice) / 100
		isPlus := rewards["isPlus"].(bool)
		currentGame.IsPlus = isPlus

		saleEndStr := rewards["end_date"].(string)
		saleEnd, err := time.Parse(time.RFC3339, saleEndStr)
		if err != nil {
			return nil, err
		}

		currentGame.SaleEnd = saleEnd

		resultGames = append(resultGames, currentGame)

	}
	return resultGames, nil
}

func Search(name string) ([]Game, error) {
	u, err := BuildURLForName(name)
	if err != nil {
		return nil, err
	}

	resp, err := http.Get(u)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("Got status code %s", resp.Status))
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return parseSearchAnswer(body)
}
