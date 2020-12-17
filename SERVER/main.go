package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var (
	tam float64 = 0
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

type structProcesos struct {
	Pid           int     `json:"PID,omitempty"`
	Nombre        string  `json:"Nombre,omitempty"`
	Usuario       string  `json:"Usuario,omitempty"`
	Estado        string  `json:"Estado,omitempty"`
	PorcentajeRam float64 `json:"PorcentajeRam,omitempty"`
	Ppid          int     `json:"PPID,omitempty"`
}

type ListaProcesos struct {
	ListaProcesos []structProcesos `json:"lista_procesos"`
}

type Message struct {
	Name string
	Body string
	Time int64
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Homies")
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

		for {
			if err := conn.WriteMessage(messageType, b); err != nil {
				log.Println(err)
			}
			time.Sleep(2 * time.Second)
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

func main() {

	router := mux.NewRouter()
	router.HandleFunc("/", homePage)
	router.HandleFunc("/ws", wsEndPoint)
	router.HandleFunc("/cpu", cpu).Methods("GET", "OPTIONS")
	log.Fatal(http.ListenAndServe(":80", router))

}

func cpu(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	file, _ := ioutil.ReadFile("/proc/cpu_201503893")
	data := ListaProcesos{}
	_ = json.Unmarshal([]byte(file), &data) //SERIALIZZAR LA DATA

	for i := 0; i < len(data.ListaProcesos); i++ {
		var temporal float64 = data.ListaProcesos[i].PorcentajeRam
		data.ListaProcesos[i].PorcentajeRam = (temporal / tam) * 100
	}
	json.NewEncoder(w).Encode(data.ListaProcesos)
}
