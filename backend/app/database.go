package app

import (
	"encoding/json"
	"errors"
	"github.com/boltdb/bolt"
	"log"
)

/*
	документация: https://github.com/boltdb/bolt#opening-a-database

	продумать плоскую key-value структуру базы

	bucket games:
		все закешированные игры
		политика обновления кэша - если есть скидка, то не обновлять до SaleEnd даты, иначе - обновлять раз в N времени

	bucket users:
		пользователи, ключ - автоинкрементный инт, значение - json, содержащий тип записи(пока один - телеграм), массив id
		(игр, на которые подписан пользователь) и вложенный JSON информации, уникальной для типа
		в случае телеграма это будет одно поле ID.

		ИСПОЛЬЗОВАТЬ ВМЕСТО JSON ВЛОЖЕННЫЕ БАКЕТЫ??

	описать в комментариях структуру базы!

 */

type Database struct {
	db *bolt.DB
}

var databaseName = "toussaint.db"

const (
	gamesBucketName = "games"
	usersTelegramBucketName = "telegram"
)

func NewDatabase() *Database {
	database := &Database{}
	database.open()
	err := database.init()
	if err != nil {
		log.Fatal(err)
	}

	return database
}


func (database *Database) open() {
	db, err := bolt.Open(databaseName, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	database.db = db
}

func (database *Database) init() error {
	return database.db.Update(func(tx *bolt.Tx) error {
		buckets := []string{
			gamesBucketName,
			usersTelegramBucketName,
		}

		for _, bucketName := range buckets {
			_, err := tx.CreateBucketIfNotExists([]byte(bucketName))
			if err != nil {
				return err
			}
		}

		return nil
	})
}

func (database *Database) AddGame(game Game) error {
	return database.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(gamesBucketName))
		if bucket == nil {
			return errors.New("bucket for games was not found")
		}

		buf, err := json.Marshal(game)
		if err != nil {
			return err
		}

		err = bucket.Put([]byte(game.Id), buf)
		if err != nil {
			return err
		}

		return nil
	})
}

func (database *Database) GetGame(id string) (game Game, err error) {
	err = database.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(gamesBucketName))
		if bucket == nil {
			return errors.New("bucket for games was not found")
		}

		g := bucket.Get([]byte(id))
		if g == nil {
			return errors.New("not found")
		}

		e := json.Unmarshal(g, &game)
		if e != nil {
			return e
		}

		return nil
	})

	return game, err
}

func (database *Database) AddUser(client Client) error {
	bucketName := selectBucketByClient(client.Type())
	return database.db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(bucketName))
		if err != nil {
			return err
		}

		buf, err := json.Marshal(client.Storable())
		if err != nil {
			return err
		}

		err = bucket.Put([]byte(client.ID()), buf)
		if err != nil {
			return err
		}

		return nil
	})
}

func (database *Database) AddGameToUser(gameID string, clientID string, clientType ClientType) error {
	bucketName := selectBucketByClient(clientType)
	return database.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))
		if bucket == nil {
			return errors.New("bucket for users was not found")
		}
		v := bucket.Get([]byte(clientID))

		var storable StorableClient
		err := json.Unmarshal(v, &storable)
		if err != nil {
			return err
		}
		storable.Subscriptions = append(storable.Subscriptions, gameID)

		buf, err := json.Marshal(storable)
		if err != nil {
			return err
		}

		err = bucket.Put([]byte(clientID), buf)
		if err != nil {
			return err
		}

		return nil
	})
}

func (database *Database) Close() {
	err := database.db.Close()
	if err != nil {
		log.Println(err)
	}
}

func selectBucketByClient(tp ClientType) (res string) {
	switch tp {
	case Telegram:
		res = usersTelegramBucketName
	default:
		log.Fatal("Unknown client type")
	}

	return
}