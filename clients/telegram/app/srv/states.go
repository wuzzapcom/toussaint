package srv

import "github.com/go-telegram-bot-api/telegram-bot-api"

type State int

const (
	NO_STATE State = iota
	SEARCH_GAME_WAIT_NAME
	SEARCH_GAME_WAIT_NAME_NO_EXPANSIONS
	SEARCH_GAME_WAIT_ID
	SEARCH_GAME_WAIT_ACTION
	DELETE_GAME
	LIST_NOTIFICATIONS
)

//HandleMessage returns message to user, should context be cached and error
func (state State) HandleMessage(message *tgbotapi.Message) (string, bool, error) {
	return "test", false, nil
}

/*
   NO_STATE(),
   SEARCH_GAME_WAIT_NAME_STATE("/search"),
   SEARCH_GAME_WAIT_NAME_NO_EXPANSIONS_STATE("/search_no_expansions"),
   SEARCH_GAME_WAIT_NUMBER_STATE(),
   SEARCH_GAME_WAIT_ACTION(),
   DELETE_GAME_STATE("/delete"),
   LIST_NOTIFICATION_GAMES("/list")
*/
