package ws

type Msg struct {
	*Ping         `json:"ping"`
	*ServerUpdate `json:"server_update"`
	*ClientUpdate `json:"client_update"`
}

type Ping struct {
}

type ServerUpdate struct {
	Status  string                 `json:"status"`
	Players map[string]*PlayerData `json:"players"`
}

type ClientUpdate struct {
	Status string     `json:"status"`
	Player PlayerData `json:"player"`
}

// todo: find some way to unify this with the player entity
type PlayerData struct {
	UUID      string  `json:"uuid"`
	Connected bool    `json:"connected"`
	Name      string  `json:"name"`
	X         float64 `json:"x"`
	Y         float64 `json:"y"`
	DX        float64 `json:"dx"`
	DY        float64 `json:"dy"`
	Dir       int     `json:"dir"`

	ClientUpdated bool
}
