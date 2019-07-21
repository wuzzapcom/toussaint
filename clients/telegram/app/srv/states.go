package srv

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"toussaint/backend/structs"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

//APIEndpoint should be set in main
var APIEndpoint string

type State int

const (
	NO_STATE State = iota
	SEARCH_GAME_WAIT_GAME
)

//HandleMessage returns message to user, should context be cached and error.
func (state State) HandleMessage(message *tgbotapi.Message, payload interface{}) (string, bool, State, interface{}, error) {
	switch state {
	case NO_STATE:
		return handleNoState(message, payload)
	case SEARCH_GAME_WAIT_GAME:
		return handleSearchGameWaitGame(message, payload)
	default:
		return unimplemented_msg_ru, false, NO_STATE, nil, nil
	}
}

//returns message, shouldCache, next state, payload and error
func handleNoState(message *tgbotapi.Message, payload interface{}) (string, bool, State, interface{}, error) {
	// /search accepts searching game as parameter in the same string
	// example: /search Witcher 3
	splitted := strings.SplitN(message.Text, " ", 2)
	command := splitted[0]
	var msg string
	if len(splitted) != 1 {
		msg = splitted[1]
	}
	switch command {
	case "/start":
		resp, err := performRequest("PUT", fmt.Sprintf("/users?client-id=%d&client-type=telegram", message.From.ID), nil)
		if err != nil {
			return register_fail_msg_ru, false, NO_STATE, nil, err
		}
		if resp.StatusCode != http.StatusCreated {
			return register_fail_msg_ru, false, NO_STATE, nil, fmt.Errorf("got status code %d", resp.StatusCode)
		}
		return register_ok_msg_ru, false, NO_STATE, nil, nil
	case "/list":
		resp, err := performRequest("GET", fmt.Sprintf("/list?client-id=%d&client-type=telegram&request-type=all", message.From.ID), nil)
		if err != nil {
			return get_list_fail_msg_ru, false, NO_STATE, nil, err
		}
		if resp.StatusCode != http.StatusOK {
			return get_list_fail_msg_ru, false, NO_STATE, nil, fmt.Errorf("got status code %d", resp.StatusCode)
		}
		listMsg, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return get_list_fail_msg_ru, false, NO_STATE, nil, err
		}
		return string(listMsg), false, NO_STATE, nil, nil
	case "/sale":
		resp, err := performRequest("GET", fmt.Sprintf("/list?client-id=%d&client-type=telegram&request-type=sale", message.From.ID), nil)
		if err != nil {
			return get_list_fail_msg_ru, false, NO_STATE, nil, err
		}
		if resp.StatusCode != http.StatusOK {
			return get_list_fail_msg_ru, false, NO_STATE, nil, fmt.Errorf("got status code %d", resp.StatusCode)
		}
		listMsg, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return get_list_fail_msg_ru, false, NO_STATE, nil, err
		}
		return string(listMsg), false, NO_STATE, nil, nil
	case "/search":
		if len(msg) == 0 {
			return search_expected_game_name_msg_ru, false, NO_STATE, nil, nil
		}
		resp, err := performRequest("GET", fmt.Sprintf("/search?name=%s", url.QueryEscape(msg)), nil)
		if err != nil {
			return get_game_fail_msg_ru, false, NO_STATE, nil, err
		}
		gamesBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return get_game_fail_msg_ru, false, NO_STATE, nil, err
		}
		var games structs.GamesJSON
		err = json.Unmarshal(gamesBytes, &games)
		if err != nil {
			return get_game_fail_msg_ru, false, NO_STATE, nil, err
		}
		ids, descs := FormatGamesListMessage(games)
		return descs, true, SEARCH_GAME_WAIT_GAME, ids, nil
	case "/delete":
		if len(msg) == 0 {
			return delete_expected_game_name_msg_ru, false, NO_STATE, nil, nil
		}
		resp, err := performRequest("GET", fmt.Sprintf("/games?name=%s", url.QueryEscape(msg)), nil)
		if err != nil {
			return get_game_fail_msg_ru, false, NO_STATE, nil, err
		}
		id, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return get_game_fail_msg_ru, false, NO_STATE, nil, err
		}

		resp, err = performRequest(
			"DELETE",
			fmt.Sprintf("/notifications?client-id=%d&client-type=telegram&game-id=%s", message.From.ID, url.QueryEscape(string(id))),
			nil,
		)

		if resp.StatusCode != http.StatusOK {
			return delete_notify_fail_msg_ru, false, NO_STATE, nil, fmt.Errorf("got status code %d", resp.StatusCode)
		}
		return delete_notify_ok_msg_ru, false, NO_STATE, nil, nil
	case "/help":
		return help_msg_ru, false, NO_STATE, nil, nil
	case "/debug":
		response := "testing notifications..."
		resp, err := performRequest("GET", fmt.Sprintf("/notifications/trigger?client-type=telegram&client-id=%d", message.From.ID), nil)
		if err != nil {
			response = fmt.Sprintf("error during testing request: %+v", err)
		}
		if resp.StatusCode != 200 {
			response = "status code != 200"
		}
		return response, false, NO_STATE, nil, nil
	default:
		return unimplemented_msg_ru, false, NO_STATE, nil, nil
	}
}

func handleSearchGameWaitGame(message *tgbotapi.Message, payload interface{}) (string, bool, State, interface{}, error) {
	ids, ok := payload.([]string)
	if !ok {
		return get_game_number_fail_msg_ru, false, NO_STATE, nil, nil
	}
	num, err := strconv.Atoi(message.Text)
	if err != nil {
		return get_game_number_fail_msg_ru, false, NO_STATE, nil, err
	}

	if num > len(ids)+1 || num <= 0 {
		return get_game_number_fail_msg_ru, false, NO_STATE, nil, fmt.Errorf("id great than array size")
	}

	if num == len(ids)+1 {
		return cancelled_ok_msg_ru, false, NO_STATE, nil, nil
	}

	fmt.Printf("wtf id is %d", num)
	resp, err := performRequest(
		"PUT",
		fmt.Sprintf("/notifications?client-id=%d&client-type=telegram&game-id=%s", message.From.ID, ids[num-1]),
		nil,
	)

	if resp.StatusCode != http.StatusCreated {
		return get_game_number_fail_msg_ru, false, NO_STATE, nil, fmt.Errorf("got status code %d", resp.StatusCode)
	}
	return notify_ok_msg_ru, false, NO_STATE, nil, nil
}

func performRequest(method string, url string, body io.Reader) (*http.Response, error) {
	cl := http.Client{}
	req, err := http.NewRequest(method, fmt.Sprintf("%s%s", APIEndpoint, url), body)
	if err != nil {
		return nil, err
	}

	return cl.Do(req)
}

func FormatGamesListMessage(games structs.GamesJSON) ([]string, string) {
	ids := make([]string, len(games.Games))
	var descs string
	for i, game := range games.Games {
		ids[i] = game.Id
		descs += fmt.Sprintf("%d) %s\n", i+1, game.Description)
	}
	descs += fmt.Sprintf("%d) %s", len(games.Games)+1, cancel_search_msg_ru)
	return ids, descs
}
