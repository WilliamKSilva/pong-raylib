package pkg

import "github.com/gorilla/websocket"

type Game struct {
	ID        string `json:"id"`
	PlayerOne Player `json:"player_one"`
	PlayerTwo Player `json:"player_two"`
}

type Player struct {
	ID         string          `json:"id"`
	Nickname   string          `json:"nickname"`
	SocketConn *websocket.Conn `json:"socket_conn"`
}

type PlayerState struct {
	GameID   string   `json:"game_id"`
	PlayerID string   `json:"player_id"`
	Position Position `json:"position"`
}

type Position struct {
	X int32 `json:"x"`
	Y int32 `json:"y"`
}

type PlayerAndOpponent struct {
	Player   *Player `json:"player"`
	Opponent *Player `json:"opponent"`
}
