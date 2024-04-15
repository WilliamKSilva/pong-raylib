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

var addr = flag.String("addr", "localhost:3000", "http service address")

func main() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: *addr, Path: "/ws"}
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	defer c.Close()

	if err != nil {
		log.Fatal("error:", err)
	}

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
		}
	}()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	joinedGame := false
	for {
		select {
		case <-done:
			return
		case t := <-ticker.C:
			if !joinedGame {
				connectGameData := ConnectGameData{
					GameID:   "",
					Nickname: "player_one",
				}
				jsonData, err := json.Marshal(connectGameData)

				if err != nil {
					log.Fatal("error:", err)
				}

				err = c.WriteMessage(websocket.BinaryMessage, jsonData)

				if err != nil {
					log.Fatal("write:", err)
					return
				}
			}

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
