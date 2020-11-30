package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"strconv"
)

const MAXTIME = 10000 //Max supported time
const DEFAULTPORT = "6656"

var actionMap = make(map[string]Stat) //This is a global variable that will hold the data. If we want to make it persistent, we can use a DB library instead.

//This struct is created for each action that holds the sum and count to calculate the average.
type Stat struct {
	Sum   int
	Count int
}

//This struct is the input for an action
type ActionRequest struct {
	Action string `json:"action" validate:"required"`
	Time   int    `json:"time" validate:"required"`
}

//This struct is the output for the action
type ActionResponse struct {
	Action string `json:"action"`
	Avg    int    `json:"avg"`
}

//This struct is the input for deletion
type DeleteResponse struct {
	Action string `json:"action"`
}

func main() {
	var port = os.Getenv("JRPORT")
	_, err := strconv.Atoi(port)
	if err != nil {
		port = DEFAULTPORT
	}
	port = ":" + port
	fmt.Println("Running on port " + port)
	//go http takes cares of concurrency. Each request is spawned in a sub-routine
	router := mux.NewRouter()

	router.HandleFunc("/stats", GetStats).Methods("GET")
	router.HandleFunc("/action", AddAction).Methods("POST")
	router.HandleFunc("/delete", RemoveAction).Methods("POST")
	log.Fatal(http.ListenAndServe(port, router))
}

func GetStats(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("action")

	w.Header().Set("Content-Type", "application/json")

	for k, v := range actionMap {
		if len(query) > 0 && query != k { //If a query string is provided, only send that
			continue
		}
		json.NewEncoder(w).Encode(ActionResponse{
			Action: k,
			Avg:    v.Sum / v.Count,
		})
	}
}

func AddAction(w http.ResponseWriter, r *http.Request) {
	var nilEntry = Stat{}
	var body ActionRequest
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil || len(body.Action) == 0 || body.Time < 0 || body.Time > MAXTIME { //Simple Validation. Could use a validation library for much better vaidation
		http.Error(w, "Invalid Input", http.StatusBadRequest)
		return
	}
	if actionMap[body.Action] == nilEntry {
		actionMap[body.Action] = Stat{
			Sum:   0,
			Count: 0,
		}
	}

	var stats = actionMap[body.Action]
	stats.Count = actionMap[body.Action].Count + 1
	stats.Sum = actionMap[body.Action].Sum + body.Time
	actionMap[body.Action] = stats

	json.NewEncoder(w).Encode("Success")
}

func RemoveAction(w http.ResponseWriter, r *http.Request) {
	var body DeleteResponse
	err := json.NewDecoder(r.Body).Decode(&body)
	if err == nil {
		if body.Action == "" {
			for k, _ := range actionMap {
				delete(actionMap, k)
			}
		} else {
			delete(actionMap, body.Action)
		}
	}
	json.NewEncoder(w).Encode("Success")
}
