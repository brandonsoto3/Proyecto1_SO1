package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"os/exec"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/shirou/gopsutil/cpu"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

//STRUCTS DEL WEB SERVICE
type structRam struct {
	Memoria_Total     int `json:"Total_de_memoria_RAM_del_servidor,omitempty"`
	Memoria_en_uso    int `json:"Total_de_memoria_RAM_consumida,omitempty"`
	Porcentaje_en_uso int `json:"Porcentaje_de_consumo_de_RAM,omitempty"`
}

type StructListaRam struct {
	StructListaRam []structRam `json:"lista_ram"`
}

type structProcesos struct {
	Pid           int     `json:"PID,omitempty"`
	Nombre        string  `json:"Nombre,omitempty"`
	Usuario       string  `json:"Usuario,omitempty"`
	Estado        string  `json:"Estado,omitempty"`
	PorcentajeRam float64 `json:"PorcentajeRam,omitempty"`
	Ppid          int     `json:"PPID,omitempty"`
}

type StructListaProcesos struct {
	StructListaProcesos []structProcesos `json:"lista_procesos"`
}

type structKill struct {
	Pid string `json:"pid,omitempty"`
}

type Porcentaje_CPU struct {
	valor float64
}

//VARIABLES
var (
	tamanio float64 = 0
)

func reader3(conn *websocket.Conn) {

	for {
		messageType, p, err := conn.ReadMessage()

		if err != nil {
			log.Println(err)
			return
		}
		//MENSAJE RECIBIDO DESDE EL CLIENTE
		log.Println(string(p))
		valor, _ := cpu.Percent(0, false)
		percentaje := valor[0]
		m := Porcentaje_CPU{math.Ceil(percentaje*100) / 100}

		b, err := json.Marshal(m)
		for {
			if err := conn.WriteMessage(messageType, b); err != nil {
				log.Println(err)
			}
			time.Sleep(2 * time.Second)
		}
	}

}

func reader2(conn *websocket.Conn) {

	for {
		messageType, p, err := conn.ReadMessage()

		if err != nil {
			log.Println(err)
			return
		}
		//MENSAJE RECIBIDO DESDE EL CLIENTE
		log.Println(string(p))

		file, err := os.Open("/proc/cpu_201503893")
		if err != nil {
			log.Fatal(err)
		}
		defer func() {
			if err = file.Close(); err != nil {
				log.Fatal(err)
			}
		}()

		b, err := ioutil.ReadAll(file)

		for {
			if err := conn.WriteMessage(messageType, b); err != nil {
				log.Println(err)
			}
			time.Sleep(2 * time.Second)
		}

	}

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

		file, err := os.Open("/proc/memo_201503893")
		if err != nil {
			log.Fatal(err)
		}
		defer func() {
			if err = file.Close(); err != nil {
				log.Fatal(err)
			}
		}()

		b, err := ioutil.ReadAll(file)

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

func wsEndPoint2(w http.ResponseWriter, r *http.Request) {

	ws, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Println(err)
	}

	log.Println("Conexion establecida")
	reader2(ws)

}

func wsEndPoint3(w http.ResponseWriter, r *http.Request) {

	ws, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Println(err)
	}

	log.Println("Conexion establecida")
	reader3(ws)

}

func porcentaje(w http.ResponseWriter, r *http.Request) {

	valor, _ := cpu.Percent(0, false)
	por := valor[0]

	fmt.Println(math.Ceil(por*100) / 100) //TRUE SI QUEREMOS VALOR POR CPU

}

func main() {
	//Obtenemos el tamanio
	leerInicio()
	//Inicio el codigo del servidor
	router := mux.NewRouter()

	router.HandleFunc("/", inicio)
	router.HandleFunc("/procesos", enviarProcesos).Methods("GET", "OPTIONS")
	router.HandleFunc("/ram", informacionRAM).Methods("GET", "OPTIONS")
	router.HandleFunc("/kill/{id}", matarProceso).Methods("POST", "OPTIONS")
	router.HandleFunc("/ws", wsEndPoint)
	router.HandleFunc("/ws2", wsEndPoint2)
	router.HandleFunc("/ws3", wsEndPoint3)
	router.HandleFunc("/porcentaje", porcentaje)

	fmt.Println("El servidor se ha iniciado en el puerto 80")
	log.Fatal(http.ListenAndServe(":80", router))
}

//FUNCIONES DEL SERVIDOR-------------------------------------------------------------------------------------------------------------------------------
//Esta funcion solo se agrego para que se vea bonito el servidor cuando inicia :)
func inicio(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	fmt.Fprintf(w, "BRANDON SOTO")
}

//Esta funcion va a devolver el json con la info de los procesos
func enviarProcesos(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	file, _ := ioutil.ReadFile("/proc/cpu_201503893")
	data := StructListaProcesos{}
	_ = json.Unmarshal([]byte(file), &data)
	for i := 0; i < len(data.StructListaProcesos); i++ {
		var temporal float64 = data.StructListaProcesos[i].PorcentajeRam
		data.StructListaProcesos[i].PorcentajeRam = (temporal / tamanio) * 100
	}
	json.NewEncoder(w).Encode(data.StructListaProcesos)
}

//Esta funcion va a devolver el json con la informacion de la pagina de la RAM
func informacionRAM(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	file, a := ioutil.ReadFile("/proc/memo_201503893")
	data := StructListaRam{}
	a = json.Unmarshal([]byte(file), &data)

	fmt.Println(a)
	json.NewEncoder(w).Encode(data.StructListaRam[0])

}

//Esta funcion va a matar el proceso especificado
func matarProceso(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	params := mux.Vars(req)
	var valor structKill
	valor.Pid = params["id"]

	cmd := exec.Command("sudo", "kill", "-9", valor.Pid)
	_, err := cmd.Output()
	//_, err := exec.Command("sh", "-c", "sudo -9 kill "+valor.Pid).Output()
	if err != nil {
		fmt.Printf("Error matando el proceso: %v", err)
	}
	json.NewEncoder(w).Encode(valor)
}

//FUNCIONES ADICIONALES-------------------------------------------------------------------------------------------------------------------------------
func leerInicio() {
	file, _ := ioutil.ReadFile("/proc/memo_201503893")
	data := StructListaRam{}
	_ = json.Unmarshal([]byte(file), &data)
	tamanio = float64(data.StructListaRam[0].Memoria_Total)
}
