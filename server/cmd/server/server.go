package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"

	. "github.com/WilliamKSilva/pong/server/pkg"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "localhost:3000", "http service address")

var upgrader = websocket.Upgrader{}

var games []Game

func connect_game(connectGame ConnectGameRequest, c *websocket.Conn) {
	var game Game
	if connectGame.GameID == "" {
		game.ID = uuid.NewString()
		game.PlayerOne = Player{
			ID:         uuid.NewString(),
			Nickname:   connectGame.Nickname,
			SocketConn: c,
		}
		log.Println("game created:", game.ID)
		games = append(games, game)

		player_and_opponent := get_player_and_opponent(game.PlayerOne.ID, game)

		var connectMessageResponse ConnectGameResponse
		if player_and_opponent.Opponent == nil {
			connectMessageResponse = ConnectGameResponse{
				GameID:   game.ID,
				Player:   *player_and_opponent.Player,
				Opponent: Player{},
			}
		} else {
			connectMessageResponse = ConnectGameResponse{
				GameID:   game.ID,
				Player:   *player_and_opponent.Player,
				Opponent: *player_and_opponent.Opponent,
			}
		}

		socketMessage := SocketMessage{
			Type: CONNECT,
			Data: connectMessageResponse,
		}

		socketMessageData, err := json.Marshal(socketMessage)

		if err != nil {
			log.Println("marshal:", err)
			return
		}

		err = c.WriteMessage(websocket.BinaryMessage, socketMessageData)

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
					ID:         uuid.NewString(),
					Nickname:   connectGame.Nickname,
					SocketConn: c,
				}

				game = games[i]
				player_and_opponent := get_player_and_opponent(game.PlayerTwo.ID, game)
				var connectMessageResponse ConnectGameResponse
				if player_and_opponent.Opponent == nil {
					connectMessageResponse = ConnectGameResponse{
						GameID:   game.ID,
						Player:   *player_and_opponent.Player,
						Opponent: Player{},
					}
				}

				connectMessageResponse = ConnectGameResponse{
					GameID:   game.ID,
					Player:   *player_and_opponent.Player,
					Opponent: *player_and_opponent.Opponent,
				}

				socketMessage := SocketMessage{
					Type: CONNECT,
					Data: connectMessageResponse,
				}

				socketMessageData, err := json.Marshal(socketMessage)

				if err != nil {
					log.Println("marshal:", err)
					return
				}

				// Send ConnectMessageResponse to player
				c.WriteMessage(websocket.BinaryMessage, socketMessageData)

				// Send ConnectMessageResponse to opponent with
				// new connection data
				connectMessageResponse = ConnectGameResponse{
					GameID:   game.ID,
					Player:   *player_and_opponent.Opponent,
					Opponent: *player_and_opponent.Player,
				}

				socketMessage.Data = connectMessageResponse

				player_and_opponent.Opponent.SocketConn.WriteJSON(socketMessage)
				break
			}
		}
	}
}

func update_player_state(playerState PlayerState) {
	for i := 0; i < len(games); i++ {
		if games[i].ID == playerState.GameID {
			player_and_opponent := get_player_and_opponent(playerState.PlayerID, games[i])
			if player_and_opponent.Opponent != nil {
				player_and_opponent.Opponent.SocketConn.WriteJSON(playerState)
			}
			break
		}
	}
}

func get_player_and_opponent(playerID string, game Game) PlayerAndOpponent {
	if game.PlayerOne.ID == playerID {
		if game.PlayerTwo.Nickname != "" {
			return PlayerAndOpponent{
				Player:   &game.PlayerOne,
				Opponent: &game.PlayerTwo,
			}
		}

		return PlayerAndOpponent{
			Player:   &game.PlayerOne,
			Opponent: nil,
		}
	}

	if game.PlayerTwo.ID == playerID {
		return PlayerAndOpponent{
			Player:   &game.PlayerTwo,
			Opponent: &game.PlayerOne,
		}
	}

	return PlayerAndOpponent{Player: nil, Opponent: nil}
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
			var connectGameData SocketMessageConnectGameRequest
			// Waste of resources, but it works for now
			err := json.Unmarshal(message, &connectGameData)

			if err != nil {
				log.Fatal("unmarshal:", err)
			}

			connect_game(connectGameData.Data, c)
		}

		// Packet sent on websocket conn is related to player state
		if socketMessage.Type == UPDATE_PLAYER_STATE {
			var playerState SocketMessageUpdatePlayerState
			log.Println(playerState.Data.PlayerID)

			// Waste of resources, but it works for now
			err := json.Unmarshal(message, &playerState)

			if err != nil {
				log.Fatal("unmarshal:", err)
			}

			update_player_state(playerState.Data)
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
