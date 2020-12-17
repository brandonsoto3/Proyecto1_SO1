/*package main

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

type strRam struct {
	Memoria_Total     int `json:"Memoria_Total,omitempty"`
	Memoria_en_uso    int `json:"Memoria_en_uso,omitempty"`
	Porcentaje_en_uso int `json:"Porcentaje_en_uso,omitempty"`
}

type ListaRam struct {
	ListaRam []strRam `json:"lista_ram"`
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

	funcion()
	router := mux.NewRouter()
	router.HandleFunc("/", homePage)
	router.HandleFunc("/ws", wsEndPoint)
	router.HandleFunc("/cpu", cpu).Methods("GET", "OPTIONS")
	router.HandleFunc("/prueba", homePage)
	log.Fatal(http.ListenAndServe(":80", router))
	fmt.Println("Servidor iniciado")

}

func funcion() {
	file, _ := ioutil.ReadFile("/proc/memo_201503893")
	data := ListaRam{}
	_ = json.Unmarshal([]byte(file), &data)
	tam = float64(data.ListaRam[0].Memoria_Total)
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
}*/

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"

	"github.com/gorilla/mux"
)

//STRUCTS DEL WEB SERVICE
type structRam struct {
	Memoria_Total     int `json:"Memoria_Total,omitempty"`
	Memoria_en_uso    int `json:"Memoria_en_uso,omitempty"`
	Porcentaje_en_uso int `json:"Porcentaje_en_uso,omitempty"`
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

//VARIABLES
var (
	tamanio float64 = 0
)

func main() {
	//Obtenemos el tamanio
	leerInicio()
	//Inicio el codigo del servidor
	router := mux.NewRouter()

	router.HandleFunc("/", inicio)
	router.HandleFunc("/procesos", enviarProcesos).Methods("GET", "OPTIONS")
	router.HandleFunc("/ram", informacionRAM).Methods("GET", "OPTIONS")
	router.HandleFunc("/kill/{id}", matarProceso).Methods("POST", "OPTIONS")

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

	//fmt.Println("Tamanio: ", len(data.StructListaProcesos))
	for i := 0; i < len(data.StructListaProcesos); i++ {
		var temporal float64 = data.StructListaProcesos[i].PorcentajeRam
		data.StructListaProcesos[i].PorcentajeRam = (temporal / tamanio) * 100
		//fmt.Println("Valor: ", data.StructListaProcesos[i])
	}

	json.NewEncoder(w).Encode(data.StructListaProcesos)
}

//Esta funcion va a devolver el json con la informacion de la pagina de la RAM
func informacionRAM(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	file, _ := ioutil.ReadFile("/proc/memo_201503893")

	data := StructListaRam{}

	_ = json.Unmarshal([]byte(file), &data)

	/*fmt.Println("Tamanio: ", len(data.StructListaRam))
	for i := 0; i < len(data.StructListaRam); i++ {
		fmt.Println("Valor: ", data.StructListaRam[i])
	}*/

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
	/*fmt.Println("Tamanio: ", tamanio)
	for i := 0; i < len(data.StructListaRam); i++ {
		fmt.Println("Valor: ", data.StructListaRam[i])
	}*/
}
