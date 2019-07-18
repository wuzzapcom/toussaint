package app

import (
	"log"
	"time"
	"toussaint/backend/structs"
)

const TimeFormat = "15:04"

func NewUpdater(betweenUpdates, updateTime string) (*Updater, error) {
	btw, err := time.ParseDuration(betweenUpdates)
	if err != nil {
		return nil, err
	}

	upd, err := time.Parse(TimeFormat, updateTime)
	if err != nil {
		return nil, err
	}

	return &Updater{
		betweenUpdates: btw,
		updateTime:     upd,
		updaterFunc:    update,
		stop:           make(chan bool, 0),
		finished:       make(chan bool),
	}, nil
}

type Updater struct {
	betweenUpdates time.Duration
	updateTime     time.Time

	updaterFunc func()

	stop     chan bool
	finished chan bool
}

func (updater *Updater) Start() {
	sleepTime := updater.updateTime.Sub(time.Now())
	time.Sleep(sleepTime)

	ticker := time.NewTicker(updater.betweenUpdates)
	for {
		select {
		case <-ticker.C:
			updater.updaterFunc()
		case <-updater.stop:
			ticker.Stop()
			updater.finished <- true
		}
	}
}

func (updater *Updater) Stop() {
	updater.stop <- true
	<-updater.finished
}

func update() {
	// TODO: this code is bad in terms of performance:
	// 1/ store users and games in memory is bad
	// 2/ work with database is not optimal at all
	games, err := database.GetGames()
	if err != nil {
		log.Printf("[ERR] Update database.GetGame: %+v", err)
		return
	}

	var updatedGames = make([]Game, 0)
	var saleStartedGameIds = make([]string, 0)

	for i := range games {
		newGame, err := SearchByID(games[i].Id)
		if err != nil {
			log.Printf("[ERR] Update SearchByID: %+v", err)
			return
		}
		if newGame != games[i] {
			updatedGames = append(updatedGames, newGame)
		}
		if newGame.SalePrice != 0 && games[i].SalePrice == 0 {
			saleStartedGameIds = append(saleStartedGameIds, newGame.Id)
		}
	}

	err = database.AddGames(updatedGames)
	if err != nil {
		log.Printf("[ERR] Update AddGames: %+v", err)
		return
	}

	//TODO: implement multiple client types
	users, err := database.GetUsers(Telegram)
	if err != nil {
		log.Printf("[ERR] Update GetUsers: %+v", err)
		return
	}

	for _, user := range users {
		games, err := database.GetGamesForUser(user, Telegram, All)
		if err != nil {
			log.Printf("[ERR] Update GetGamesForUser: %+v", err)
			return
		}

		var sales = make([]Game, 0)
		for _, saleGame := range saleStartedGameIds {
			for _, userGame := range games {
				if saleGame == userGame.Id {
					sales = append(sales, userGame)
				}
			}
		}
		ids, descs := DescribeGames(sales)
		notif := structs.UserNotification{
			Games: structs.GamesJSON{
				Games: make([]structs.GamePair, len(ids)),
			},
			UserID: user,
		}
		for i := range ids {
			notif.Games.Games[i] = structs.GamePair{
				Id:          ids[i],
				Description: descs[i],
			}
		}
		err = notifier.NotifyUser(Telegram, notif)
		if err != nil {
			log.Printf("[ERR] Update failed to notify user: %+v", err)
		}
	}
}
