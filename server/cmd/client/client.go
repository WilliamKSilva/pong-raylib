package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
)

// I need to learn how to import packages from the same module
type Player struct {
	ID       string `json:"id"`
	Nickname string `json:"nickname"`
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
	X float64 `json:"x"`
	Y float64 `json:"y"`
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

var addr = flag.String("addr", "localhost:3000", "http service address")
var gameId = flag.String("gameId", "", "first player to connect")

func main() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: *addr, Path: "/ws"}
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	defer c.Close()

	if err != nil {
		log.Fatal("error:", err)
	}

	// which player this websocket connection refers to
	var connectGameData ConnectMessageResponse
	done := make(chan struct{})
	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Printf("recv: %s", message)
			var socketMessage SocketMessage
			err = json.Unmarshal(message, &socketMessage)

			if err != nil {
				log.Println("unmarshal:", err)
				return
			}

			if socketMessage.Type == CONNECT {
				connectGameData = socketMessage.Data.(ConnectMessageResponse)
			}
		}
	}()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return
		case _ = <-ticker.C:
			// Player not connected yet
			if connectGameData.GameID == "" {
				// The client is testing the first Player
				if *gameId == "" {
					socketMessage := SocketMessage{
						Type: CONNECT,
						Data: ConnectGameData{
							GameID:   "",
							Nickname: "player1",
						},
					}

					socketMessageData, err := json.Marshal(socketMessage)

					if err != nil {
						log.Fatal("marshal:", err)
						return
					}

					log.Println(socketMessageData)

					c.WriteMessage(websocket.BinaryMessage, socketMessageData)
				} else {
					socketMessage := SocketMessage{
						Type: CONNECT,
						Data: ConnectGameData{
							GameID:   *gameId,
							Nickname: "player2",
						},
					}

					socketMessageData, err := json.Marshal(socketMessage)

					if err != nil {
						log.Fatal("marshal:", err)
						return
					}

					c.WriteMessage(websocket.BinaryMessage, socketMessageData)
				}
			}

			// If player is already connected the next messages will be Player State update
			socketMessage := SocketMessage{
				Type: UPDATE_PLAYER_STATE,
				Data: PlayerState{
					GameID:   connectGameData.GameID,
					PlayerID: connectGameData.Player.ID,
					Position: Position{
						X: 20.0,
						Y: 20.0,
					},
				},
			}

			socketMessageData, err := json.Marshal(socketMessage)

			if err != nil {
				log.Fatal("marshal:", err)
				return
			}

			c.WriteMessage(websocket.BinaryMessage, socketMessageData)

		case <-interrupt:
			log.Println("interrupt")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}
