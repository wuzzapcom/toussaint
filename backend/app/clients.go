package app

import (
	"errors"
	"fmt"
)

type ClientType int
const (
	Telegram ClientType = iota
)

//GetClientType takes source string and returns parsed type and error.
//if failed, then returns error and -1 as ClientType
func GetClientType(str string) (ClientType, error) {
	switch str {
	case "telegram":
		return Telegram, nil
	default:
		return -1, errors.New(fmt.Sprintf("not found client type %s", str))
	}
}

type RequestType int
const (
	All = iota
	Sale
)

func GetRequestType(str string) (RequestType, error) {
	switch str {
	case "all":
		return All, nil
	case "sale":
		return Sale, nil
	default:
		return -1, errors.New(fmt.Sprintf("not found request type %s", str))
	}
}

type StorableClient struct {
	Subscriptions []string `json:"subs"`
}

type Client interface {
	Type() ClientType
	Subscriptions() []string
	ID() string
	Storable() StorableClient
	AddSubscription(subscription string)
}

func NewTelegramClient(id string) Client {
	return &telegramClient{
		subscriptions: make([]string, 0),
		id: id,
	}
}

type telegramClient struct {
	subscriptions []string
	id string
}

func (tc *telegramClient) Type() ClientType {
	return Telegram
}

func (tc *telegramClient) Subscriptions() []string {
	return tc.subscriptions
}

func (tc *telegramClient) ID() string {
	return tc.id
}

func (tc *telegramClient) Storable() StorableClient {
	return StorableClient{
		Subscriptions: tc.Subscriptions(),
	}
}

func (tc *telegramClient) AddSubscription(subscription string) {
	tc.subscriptions = append(tc.subscriptions, subscription)
}
