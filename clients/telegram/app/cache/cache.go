package cache

import (
	"log"
	"sync"
	"time"
	"toussaint/clients/telegram/app/srv"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)

var delay = time.Duration(5 * time.Hour)
var frequency = time.Duration(10 * time.Hour)

//contains all current sessions. Key: ID, Value: State
//all unfinished sessions should be removed after some delay(e.g. 5 hours) in separate goroutine
//only multistage sessions should be cached(e.g. GET /notify should not be cached)
var shared *cache

func Init() {
	shared = &cache{
		cache: make(map[int]*Context),

		validator:     validator,
		stopValidator: make(chan bool, 0),
	}
	go shared.validator(shared.stopValidator)
}

type cache struct {
	cache map[int]*Context

	validator     func(chan bool)
	stopValidator chan bool
}

func HandleMessage(message *tgbotapi.Message) (string, error) {

	var ok bool
	var c *Context

	c, ok = shared.cache[message.From.ID]
	if !ok {
		c = &Context{
			state:     srv.NO_STATE,
			updatedAt: time.Now(),

			mu: sync.Mutex{},
		}
	}

	c.updatedAt = time.Now()

	answer, shouldCache, err := c.HandleMessage(message)
	if shouldCache {
		if !ok {
			shared.cache[message.From.ID] = c
		}
	} else {
		if ok {
			delete(shared.cache, message.From.ID)
		}
	}

	return answer, err
}

type Context struct {
	state     srv.State
	updatedAt time.Time
	payload   interface{}

	mu sync.Mutex
}

func (c *Context) HandleMessage(message *tgbotapi.Message) (string, bool, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	answer, shouldCache, newState, payload, err := c.state.HandleMessage(message, c.payload)

	c.state = newState
	c.payload = payload

	return answer, shouldCache, err
}

func validator(stopValidator chan bool) {
	ticker := time.Tick(frequency)
	log.Printf("[INF] Start validator")
	for {
		select {
		case <-ticker:
			log.Printf("[INF] Validation ...")
			current := time.Now()
			for key, value := range shared.cache {
				if current.Sub(value.updatedAt) > delay {
					delete(shared.cache, key)
				}
			}
		case <-stopValidator:
			log.Printf("[INF] Stop validator")
		}
	}
}
