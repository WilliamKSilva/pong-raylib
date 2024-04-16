package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Player struct {
	ID          string `json:"id"`
	Nickname    string `json:"nickname"`
	socket_conn *websocket.Conn
}

type Game struct {
	ID        string `json:"id"`
	PlayerOne Player `json:"player_one"`
	PlayerTwo Player `json:"player_two"`
}

type ConnectGameData struct {
	GameID   string `json:"game_id"`
	Nickname string `json:"nickname"`
}

type Position struct {
	X string `json:"x"`
	Y string `json:"y"`
}

type PlayerState struct {
	GameID   string   `json:"game_id"`
	PlayerID string   `json:"player_id"`
	Position Position `json:"position"`
}

type ConnectMessageResponse struct {
	GameID   string `json:"game_id"`
	Player   Player `json:"player"`
	Opponent Player `json:"opponent"`
}

type MessageTypes string

const (
	CONNECT             MessageTypes = "CONNECT"
	UPDATE_PLAYER_STATE MessageTypes = "UPDATE_PLAYER_STATE"
)

type SocketMessage struct {
	Type MessageTypes `json:"type"`
	Data any          `json:"data"`
}

type SocketMessageConnectGame struct {
	Data ConnectGameData `json:"data"`
}

type PlayerAndOpponent struct {
	player   *Player
	opponent *Player
}

var addr = flag.String("addr", "localhost:3000", "http service address")

var upgrader = websocket.Upgrader{}

var games []Game

func connect_game(connectGame ConnectGameData, c *websocket.Conn) {
	var game Game
	if connectGame.GameID == "" {
		game.ID = uuid.NewString()
		game.PlayerOne = Player{
			ID:       uuid.NewString(),
			Nickname: connectGame.Nickname,
		}
		log.Println("game created:", game.ID)
		games = append(games, game)

		player_and_opponent := get_player_and_opponent(game.PlayerOne.ID, game)

		var connectMessageResponse ConnectMessageResponse
		if player_and_opponent.opponent == nil {
			connectMessageResponse = ConnectMessageResponse{
				GameID:   game.ID,
				Player:   *player_and_opponent.player,
				Opponent: Player{},
			}
		} else {
			connectMessageResponse = ConnectMessageResponse{
				GameID:   game.ID,
				Player:   *player_and_opponent.player,
				Opponent: *player_and_opponent.opponent,
			}
		}

		connectMessageResponseData, err := json.Marshal(connectMessageResponse)

		if err != nil {
			log.Println("marshal:", err)
			return
		}

		err = c.WriteMessage(websocket.BinaryMessage, connectMessageResponseData)

		if err != nil {
			log.Fatal("write to socket:", err)
			return
		}
	} else {
		// Game already exists, player connecting is PlayerTwo
		for i := 0; i < len(games); i++ {
			if games[i].ID == connectGame.GameID {
				log.Println("game updated:", games[i].ID)
				games[i].PlayerTwo = Player{
					ID:       uuid.NewString(),
					Nickname: connectGame.Nickname,
				}

				game = games[i]
				player_and_opponent := get_player_and_opponent(game.PlayerTwo.ID, game)
				var connectMessageResponse ConnectMessageResponse
				if player_and_opponent.opponent == nil {
					connectMessageResponse = ConnectMessageResponse{
						GameID:   game.ID,
						Player:   *player_and_opponent.player,
						Opponent: Player{},
					}
				}

				connectMessageResponse = ConnectMessageResponse{
					GameID:   game.ID,
					Player:   *player_and_opponent.player,
					Opponent: *player_and_opponent.opponent,
				}

				connectMessageResponseData, err := json.Marshal(connectMessageResponse)

				if err != nil {
					log.Println("marshal:", err)
					return
				}

				c.WriteMessage(websocket.BinaryMessage, connectMessageResponseData)
				break
			}
		}
	}
}

func update_player_state(playerState PlayerState) {
	for i := 0; i < len(games); i++ {
		if games[i].ID == playerState.GameID {
			player_and_opponent := get_player_and_opponent(playerState.GameID, games[i])
			if player_and_opponent.opponent != nil {
				player_and_opponent.opponent.socket_conn.WriteJSON(playerState)
			}
			break
		}
	}
}

func get_player_and_opponent(playerID string, game Game) PlayerAndOpponent {
	if game.PlayerOne.ID == playerID {
		if game.PlayerTwo.Nickname != "" {
			return PlayerAndOpponent{
				player:   &game.PlayerOne,
				opponent: &game.PlayerTwo,
			}
		}

		return PlayerAndOpponent{
			player:   &game.PlayerOne,
			opponent: nil,
		}
	}

	if game.PlayerTwo.ID == playerID {
		return PlayerAndOpponent{
			player:   &game.PlayerTwo,
			opponent: &game.PlayerOne,
		}
	}

	return PlayerAndOpponent{nil, nil}
}

func handle_socket_conn(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}

		var socketMessage SocketMessage
		err = json.Unmarshal(message, &socketMessage)

		if err != nil {
			log.Println("unmarshal:", err)
			return
		}

		// Packet sent on websocket conn is related to game connection
		if socketMessage.Type == CONNECT {
			var connectGameData SocketMessageConnectGame
			// Waste of resources, but it works for now
			err := json.Unmarshal(message, &connectGameData)

			if err != nil {
				log.Fatal("unmarshal:", err)
			}

			connect_game(connectGameData.Data, c)
		}

		// Packet sent on websocket conn is related to player state
		if socketMessage.Type == UPDATE_PLAYER_STATE {
			var playerState PlayerState
			// Waste of resources, but it works for now
			err := json.Unmarshal(message, &playerState)

			if err != nil {
				log.Fatal("unmarshal:", err)
			}
			update_player_state(socketMessage.Data.(PlayerState))
		}
	}
}

func handle_default(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("<html>Listening!</html>"))
}

func main() {
	fmt.Println("Listening on port 3000")
	http.HandleFunc("/", handle_default)
	http.HandleFunc("/ws", handle_socket_conn)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
