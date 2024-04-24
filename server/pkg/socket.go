package pkg

type MessageType int32

const (
	CONNECT             MessageType = 0
	UPDATE_PLAYER_STATE MessageType = 1
)

type ConnectGameRequest struct {
	GameID   string `json:"game_id"`
	Nickname string `json:"nickname"`
}

type ConnectGameResponse struct {
	GameID   string `json:"game_id"`
	Player   Player `json:"player"`
	Opponent Player `json:"opponent"`
}

type SocketMessage struct {
	Type MessageType `json:"type"`
	Data any         `json:"data"`
}

type SocketMessageConnectGameRequest struct {
	Data ConnectGameRequest `json:"connect_game_request"`
}

type SocketMessageConnectGameResponse struct {
	Data ConnectGameResponse `json:"connect_game_response"`
}

type SocketMessageUpdatePlayerState struct {
	Data PlayerState `json:"player_state"`
}
