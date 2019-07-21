package structs

type GamesJSON struct {
	Games []GamePair `json:"games"`
}

type GamePair struct {
	Id          string `json:"id"`
	Description string `json:"description"`
}

type UsersJSON struct {
	Ids []string `json:"ids"`
}

type UserNotification struct {
	Games  GamesJSON `json:"games"`
	UserID string    `json:"userId"`
}
