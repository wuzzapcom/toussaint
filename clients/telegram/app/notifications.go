package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"toussaint/backend/structs"
	"toussaint/clients/telegram/app/srv"
)

func NewNotifier(tg *Telegram) *Notifier {
	return &Notifier{
		tg: tg,
		cl: http.Client{},
	}
}

type Notifier struct {
	shouldStop bool
	tg         *Telegram
	cl         http.Client
}

// Start runs goroutine inside
func (notifier *Notifier) Start() {
	loop := func() error {
		req, err := http.NewRequest("GET", fmt.Sprintf("%s/notifications", srv.APIEndpoint), nil)
		if err != nil {
			return err
		}
		resp, err := notifier.cl.Do(req)
		if err != nil {
			return err
		}
		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("/v1/notifications returned unexpected status code: %d", resp.StatusCode)
		}
		var data = make([]byte, 1000)

		for !notifier.shouldStop {
			n, extra, err := notifier.read(resp.Body, data)
			if err != nil {
				resp.Body.Close()
				return err
			}
			if extra != nil {
				err = notifier.handleMessage(extra)
				if err != nil {
					return err
				}
				continue
			}
			if n != 0 {
				err = notifier.handleMessage(data)
				if err != nil {
					return err
				}
			}
		}
		return nil
	}
	wrapper := func() {
		err := loop()
		if err != nil {
			panic(err)
		}
	}
	go wrapper()
}

// read reads from io.Reader
// Usually writes read bytes into data parameter, returns nil, nil
// If received message length is greater than data, total message will be read
// and returned as result value
func (notifier *Notifier) read(source io.Reader, data []byte) (size int, res []byte, err error) {
	size, err = source.Read(data)
	if err != nil {
		return 0, nil, err
	}
	if size == len(data) {
		// read full message and return as func return value
		res := make([]byte, len(data))
		copy(data, res)

		var n = size
		for n == len(data) {
			n, err = source.Read(data)
			if err != nil {
				return 0, nil, err
			}
			res = append(res, data...)
			size += n
		}
		return size, res, nil
	}
	// info stored in data
	return size, nil, nil
}

func (notifier *Notifier) handleMessage(msg []byte) error {
	log.Printf("notifier: handling message %+v", string(msg))
	var data structs.UserNotification
	msg = bytes.Trim(msg, string(byte(0)))
	err := json.Unmarshal(msg, &data)
	if err != nil {
		return err
	}
	_, descs := srv.FormatGamesListMessage(data.Games)
	id, err := strconv.ParseInt(data.UserID, 10, 64)
	if err != nil {
		return err
	}
	log.Printf("notifier: sending message to used %d: %s", id, descs)
	return notifier.tg.answer(id, fmt.Sprintf("%s\n%s", srv.New_sales_msg_ru, descs))
}

func (notifier *Notifier) Stop() {
	notifier.shouldStop = true
}
