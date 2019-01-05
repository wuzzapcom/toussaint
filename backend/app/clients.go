package app

type ClientType int
const (
	Telegram ClientType = iota
)

type StorableClient struct {
	Subscriptions []string `json:"subs"`
}

type Client interface {
	Type() ClientType
	Subscriptions() []string
	ID() string
	Storable() StorableClient
}

type TelegramClient struct {
	subscriptions []string
	id string
}

func (tc *TelegramClient) Type() ClientType {
	return Telegram
}

func (tc *TelegramClient) Subscriptions() []string {
	return tc.subscriptions
}

func (tc *TelegramClient) ID() string {
	return tc.id
}

func (tc *TelegramClient) Storable() StorableClient {
	return StorableClient{
		Subscriptions: tc.Subscriptions(),
	}
}
