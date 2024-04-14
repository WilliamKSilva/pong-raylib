package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

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
	X string `json:"x"`
	Y string `json:"y"`
}

type PlayerState struct {
	PlayerID string   `json:"player_id"`
	Position Position `json:"position"`
}

var addr = flag.String("addr", "localhost:3000", "http service address")

var upgrader = websocket.Upgrader{}

var games []Game

func connect_game(connectGame ConnectGameData) {
	if connectGame.GameID == "" {
		var game Game
		game.ID = uuid.NewString()
		game.PlayerOne = Player{
			ID:       uuid.NewString(),
			Nickname: connectGame.Nickname,
		}
		log.Println("game created:", game.ID)
		games = append(games, game)
	} else {
		// Game already exists, player connecting is PlayerTwo
		for i := 0; i < len(games); i++ {
			if games[i].ID == connectGame.GameID {
				log.Println("game updated:", games[i].ID)
				games[i].PlayerTwo = Player{
					ID:       uuid.NewString(),
					Nickname: connectGame.Nickname,
				}
				break
			}
		}
	}
}

func handle_socket_conn(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}

		var connectGameData ConnectGameData
		err = json.Unmarshal(message, &connectGameData)

		// Packet sent on websocket conn is related to game connection
		if err == nil {
			connect_game(connectGameData)
		}

		log.Printf("recv: %s", message)

		err = c.WriteMessage(mt, message)
		if err != nil {
			log.Println("write:", err)
			break
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
