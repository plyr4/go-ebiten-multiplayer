package ws

type IMsg interface {
}

type Msg struct {
	*Ping         `json:"ping"`
	*ServerUpdate `json:"server_update"`
	*ClientUpdate `json:"client_update"`
}

type Ping struct {
}

type ServerUpdate struct {
	Status           string `json:"status"`
	ConnectedPlayers int    `json:"connected_players"`
}

type ClientUpdate struct {
	Status string `json:"status"`
	Foo    int    `json:"foo"`
}
