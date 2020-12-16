package main

import (
	socketio "github.com/googollee/go-socket.io"
	"encoding/json"
	"fmt"
	"log"
	"net/http"    
	"github.com/gorilla/mux"
)

type task struct {
	ID      int    `json:ID`
	Name    string `json:Name`
	Content string `json:Content`
}

type allTasks []task

var tasks = allTasks{
	{
		ID:      1,
		Name:    "Brandon Soto",
		Content: "Some content",
	}, {
		ID:      1,
		Name:    "Brandon Soto",
		Content: "Some content",
	},
}

func indexRoute(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "WELCOME")
}

func getAll(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(tasks)
}

func main() {

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", indexRoute)
	router.HandleFunc("/getAll", getAll)
	log.Fatal(http.ListenAndServe(":3000", router))
}
