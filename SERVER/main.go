package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/shirou/gopsutil/internal/common"
)

var ClocksPerSec = float64(100)

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

type Message struct {
	Name string
	Body string
	Time int64
}

//VARIABLES
var (
	tamanio float64 = 0
)

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
	router.HandleFunc("/porcentaje", prueba)

	fmt.Println("El servidor se ha iniciado en el puerto 80")
	log.Fatal(http.ListenAndServe(":80", router))
}

func prueba(w http.ResponseWriter, r *http.Request) {
	fmt.Println(Percent(0, true))
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

type TimesStat struct {
	CPU       string  `json:"cpu"`
	User      float64 `json:"user"`
	System    float64 `json:"system"`
	Idle      float64 `json:"idle"`
	Nice      float64 `json:"nice"`
	Iowait    float64 `json:"iowait"`
	Irq       float64 `json:"irq"`
	Softirq   float64 `json:"softirq"`
	Steal     float64 `json:"steal"`
	Guest     float64 `json:"guest"`
	GuestNice float64 `json:"guestNice"`
}

type InfoStat struct {
	CPU        int32    `json:"cpu"`
	VendorID   string   `json:"vendorId"`
	Family     string   `json:"family"`
	Model      string   `json:"model"`
	Stepping   int32    `json:"stepping"`
	PhysicalID string   `json:"physicalId"`
	CoreID     string   `json:"coreId"`
	Cores      int32    `json:"cores"`
	ModelName  string   `json:"modelName"`
	Mhz        float64  `json:"mhz"`
	CacheSize  int32    `json:"cacheSize"`
	Flags      []string `json:"flags"`
	Microcode  string   `json:"microcode"`
}

type lastPercent struct {
	sync.Mutex
	lastCPUTimes    []TimesStat
	lastPerCPUTimes []TimesStat
}

var lastCPUPercent lastPercent
var invoke common.Invoker = common.Invoke{}

func init() {
	lastCPUPercent.Lock()
	lastCPUPercent.lastCPUTimes, _ = Times(false)
	lastCPUPercent.lastPerCPUTimes, _ = Times(true)
	lastCPUPercent.Unlock()
}

// Counts returns the number of physical or logical cores in the system
func Counts(logical bool) (int, error) {
	return CountsWithContext(context.Background(), logical)
}

func CountsWithContext(ctx context.Context, logical bool) (int, error) {
	if logical {
		ret := 0
		// https://github.com/giampaolo/psutil/blob/d01a9eaa35a8aadf6c519839e987a49d8be2d891/psutil/_pslinux.py#L599
		procCpuinfo := common.HostProc("cpuinfo")
		lines, err := common.ReadLines(procCpuinfo)
		if err == nil {
			for _, line := range lines {
				line = strings.ToLower(line)
				if strings.HasPrefix(line, "processor") {
					ret++
				}
			}
		}
		if ret == 0 {
			procStat := common.HostProc("stat")
			lines, err = common.ReadLines(procStat)
			if err != nil {
				return 0, err
			}
			for _, line := range lines {
				if len(line) >= 4 && strings.HasPrefix(line, "cpu") && '0' <= line[3] && line[3] <= '9' { // `^cpu\d` regexp matching
					ret++
				}
			}
		}
		return ret, nil
	}
	// physical cores
	// https://github.com/giampaolo/psutil/blob/8415355c8badc9c94418b19bdf26e622f06f0cce/psutil/_pslinux.py#L615-L628
	var threadSiblingsLists = make(map[string]bool)
	// These 2 files are the same but */core_cpus_list is newer while */thread_siblings_list is deprecated and may disappear in the future.
	// https://www.kernel.org/doc/Documentation/admin-guide/cputopology.rst
	// https://github.com/giampaolo/psutil/pull/1727#issuecomment-707624964
	// https://lkml.org/lkml/2019/2/26/41
	for _, glob := range []string{"devices/system/cpu/cpu[0-9]*/topology/core_cpus_list", "devices/system/cpu/cpu[0-9]*/topology/thread_siblings_list"} {
		if files, err := filepath.Glob(common.HostSys(glob)); err == nil {
			for _, file := range files {
				lines, err := common.ReadLines(file)
				if err != nil || len(lines) != 1 {
					continue
				}
				threadSiblingsLists[lines[0]] = true
			}
			ret := len(threadSiblingsLists)
			if ret != 0 {
				return ret, nil
			}
		}
	}
	// https://github.com/giampaolo/psutil/blob/122174a10b75c9beebe15f6c07dcf3afbe3b120d/psutil/_pslinux.py#L631-L652
	filename := common.HostProc("cpuinfo")
	lines, err := common.ReadLines(filename)
	if err != nil {
		return 0, err
	}
	mapping := make(map[int]int)
	currentInfo := make(map[string]int)
	for _, line := range lines {
		line = strings.ToLower(strings.TrimSpace(line))
		if line == "" {
			// new section
			id, okID := currentInfo["physical id"]
			cores, okCores := currentInfo["cpu cores"]
			if okID && okCores {
				mapping[id] = cores
			}
			currentInfo = make(map[string]int)
			continue
		}
		fields := strings.Split(line, ":")
		if len(fields) < 2 {
			continue
		}
		fields[0] = strings.TrimSpace(fields[0])
		if fields[0] == "physical id" || fields[0] == "cpu cores" {
			val, err := strconv.Atoi(strings.TrimSpace(fields[1]))
			if err != nil {
				continue
			}
			currentInfo[fields[0]] = val
		}
	}
	ret := 0
	for _, v := range mapping {
		ret += v
	}
	return ret, nil
}

func Times(percpu bool) ([]TimesStat, error) {
	return TimesWithContext(context.Background(), percpu)
}

func TimesWithContext(ctx context.Context, percpu bool) ([]TimesStat, error) {
	filename := common.HostProc("stat")
	var lines = []string{}
	if percpu {
		statlines, err := common.ReadLines(filename)
		if err != nil || len(statlines) < 2 {
			return []TimesStat{}, nil
		}
		for _, line := range statlines[1:] {
			if !strings.HasPrefix(line, "cpu") {
				break
			}
			lines = append(lines, line)
		}
	} else {
		lines, _ = common.ReadLinesOffsetN(filename, 0, 1)
	}

	ret := make([]TimesStat, 0, len(lines))

	for _, line := range lines {
		ct, err := parseStatLine(line)
		if err != nil {
			continue
		}
		ret = append(ret, *ct)

	}
	return ret, nil
}

func parseStatLine(line string) (*TimesStat, error) {
	fields := strings.Fields(line)

	if len(fields) == 0 {
		return nil, errors.New("stat does not contain cpu info")
	}

	if strings.HasPrefix(fields[0], "cpu") == false {
		return nil, errors.New("not contain cpu")
	}

	cpu := fields[0]
	if cpu == "cpu" {
		cpu = "cpu-total"
	}
	user, err := strconv.ParseFloat(fields[1], 64)
	if err != nil {
		return nil, err
	}
	nice, err := strconv.ParseFloat(fields[2], 64)
	if err != nil {
		return nil, err
	}
	system, err := strconv.ParseFloat(fields[3], 64)
	if err != nil {
		return nil, err
	}
	idle, err := strconv.ParseFloat(fields[4], 64)
	if err != nil {
		return nil, err
	}
	iowait, err := strconv.ParseFloat(fields[5], 64)
	if err != nil {
		return nil, err
	}
	irq, err := strconv.ParseFloat(fields[6], 64)
	if err != nil {
		return nil, err
	}
	softirq, err := strconv.ParseFloat(fields[7], 64)
	if err != nil {
		return nil, err
	}

	ct := &TimesStat{
		CPU:     cpu,
		User:    user / ClocksPerSec,
		Nice:    nice / ClocksPerSec,
		System:  system / ClocksPerSec,
		Idle:    idle / ClocksPerSec,
		Iowait:  iowait / ClocksPerSec,
		Irq:     irq / ClocksPerSec,
		Softirq: softirq / ClocksPerSec,
	}
	if len(fields) > 8 { // Linux >= 2.6.11
		steal, err := strconv.ParseFloat(fields[8], 64)
		if err != nil {
			return nil, err
		}
		ct.Steal = steal / ClocksPerSec
	}
	if len(fields) > 9 { // Linux >= 2.6.24
		guest, err := strconv.ParseFloat(fields[9], 64)
		if err != nil {
			return nil, err
		}
		ct.Guest = guest / ClocksPerSec
	}
	if len(fields) > 10 { // Linux >= 3.2.0
		guestNice, err := strconv.ParseFloat(fields[10], 64)
		if err != nil {
			return nil, err
		}
		ct.GuestNice = guestNice / ClocksPerSec
	}

	return ct, nil
}

func (c TimesStat) String() string {
	v := []string{
		`"cpu":"` + c.CPU + `"`,
		`"user":` + strconv.FormatFloat(c.User, 'f', 1, 64),
		`"system":` + strconv.FormatFloat(c.System, 'f', 1, 64),
		`"idle":` + strconv.FormatFloat(c.Idle, 'f', 1, 64),
		`"nice":` + strconv.FormatFloat(c.Nice, 'f', 1, 64),
		`"iowait":` + strconv.FormatFloat(c.Iowait, 'f', 1, 64),
		`"irq":` + strconv.FormatFloat(c.Irq, 'f', 1, 64),
		`"softirq":` + strconv.FormatFloat(c.Softirq, 'f', 1, 64),
		`"steal":` + strconv.FormatFloat(c.Steal, 'f', 1, 64),
		`"guest":` + strconv.FormatFloat(c.Guest, 'f', 1, 64),
		`"guestNice":` + strconv.FormatFloat(c.GuestNice, 'f', 1, 64),
	}

	return `{` + strings.Join(v, ",") + `}`
}

// Total returns the total number of seconds in a CPUTimesStat
func (c TimesStat) Total() float64 {
	total := c.User + c.System + c.Nice + c.Iowait + c.Irq + c.Softirq +
		c.Steal + c.Idle
	return total
}

func (c InfoStat) String() string {
	s, _ := json.Marshal(c)
	return string(s)
}

func getAllBusy(t TimesStat) (float64, float64) {
	busy := t.User + t.System + t.Nice + t.Iowait + t.Irq +
		t.Softirq + t.Steal
	return busy + t.Idle, busy
}

func calculateBusy(t1, t2 TimesStat) float64 {
	t1All, t1Busy := getAllBusy(t1)
	t2All, t2Busy := getAllBusy(t2)

	if t2Busy <= t1Busy {
		return 0
	}
	if t2All <= t1All {
		return 100
	}
	return math.Min(100, math.Max(0, (t2Busy-t1Busy)/(t2All-t1All)*100))
}

func calculateAllBusy(t1, t2 []TimesStat) ([]float64, error) {
	// Make sure the CPU measurements have the same length.
	if len(t1) != len(t2) {
		return nil, fmt.Errorf(
			"received two CPU counts: %d != %d",
			len(t1), len(t2),
		)
	}

	ret := make([]float64, len(t1))
	for i, t := range t2 {
		ret[i] = calculateBusy(t1[i], t)
	}
	return ret, nil
}

// Percent calculates the percentage of cpu used either per CPU or combined.
// If an interval of 0 is given it will compare the current cpu times against the last call.
// Returns one value per cpu, or a single value if percpu is set to false.
func Percent(interval time.Duration, percpu bool) ([]float64, error) {
	return PercentWithContext(context.Background(), interval, percpu)
}

func PercentWithContext(ctx context.Context, interval time.Duration, percpu bool) ([]float64, error) {
	if interval <= 0 {
		return percentUsedFromLastCall(percpu)
	}

	// Get CPU usage at the start of the interval.
	cpuTimes1, err := Times(percpu)
	if err != nil {
		return nil, err
	}

	if err := common.Sleep(ctx, interval); err != nil {
		return nil, err
	}

	// And at the end of the interval.
	cpuTimes2, err := Times(percpu)
	if err != nil {
		return nil, err
	}

	return calculateAllBusy(cpuTimes1, cpuTimes2)
}

func percentUsedFromLastCall(percpu bool) ([]float64, error) {
	cpuTimes, err := Times(percpu)
	if err != nil {
		return nil, err
	}
	lastCPUPercent.Lock()
	defer lastCPUPercent.Unlock()
	var lastTimes []TimesStat
	if percpu {
		lastTimes = lastCPUPercent.lastPerCPUTimes
		lastCPUPercent.lastPerCPUTimes = cpuTimes
	} else {
		lastTimes = lastCPUPercent.lastCPUTimes
		lastCPUPercent.lastCPUTimes = cpuTimes
	}

	if lastTimes == nil {
		return nil, fmt.Errorf("error getting times for cpu percent. lastTimes was nil")
	}
	return calculateAllBusy(lastTimes, cpuTimes)
}
