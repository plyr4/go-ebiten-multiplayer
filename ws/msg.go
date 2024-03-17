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
	Status  string       `json:"status"`
	Players []PlayerData `json:"players"`
}

type ClientUpdate struct {
	Status string     `json:"status"`
	Player PlayerData `json:"player"`
}

// todo: find some way to unify this with the player entity
type PlayerData struct {
	UUID string  `json:"uuid"`
	X    float64 `json:"x"`
	Y    float64 `json:"y"`
}
