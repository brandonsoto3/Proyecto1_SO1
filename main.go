package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

type task struct {
	ID      int    `json:ID`
	Name    string `json:Name`
	Content string `json:Content`
}

type Message struct {
	Name string
	Body string
	Time int64
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Home")
}

func reader(conn *websocket.Conn) {

	for {
		messageType, p, err := conn.ReadMessage()

		if err != nil {
			log.Println(err)
			return
		}
		//MENSAJE RECIBIDO DESDE EL CLIENTE
		log.Println(string(p))
		m := Message{"AÑON", "Hello", 1294706395881547000}

		b, err := json.Marshal(m)

		for i := 0; i < 10; i++ {
			if err := conn.WriteMessage(messageType, b); err != nil {
				time.Sleep(10000 * time.Millisecond)
				log.Println(err)
			}
		}

	}

}

func wsEndPoint(w http.ResponseWriter, r *http.Request) {

	ws, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Println(err)
	}

	log.Println("Conexion establecida")
	reader(ws)

}

func setupRoutes() {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/ws", wsEndPoint)
}

func main() {
	setupRoutes()
	log.Fatal(http.ListenAndServe(":3000", nil))

}
