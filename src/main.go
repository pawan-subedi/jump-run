package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

var actionMap = make(map[string]Stat)

type Stat struct{
	Sum int
	Count int
}
type ActionRequest struct {
	Action string   `json:"action" validate:"required"`
	Time int `json:"time" validate:"required"`
}

type ActionResponse struct {
	Action string   `json:"action"`
	Avg int `json:"avg"`
}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/stats", getStats).Methods("GET")
	router.HandleFunc("/action", addAction).Methods("POST")
	http.ListenAndServe(":6656", router)
}


func getStats(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	for k,v := range actionMap {
		json.NewEncoder(w).Encode(ActionResponse{
			Action: k,
			Avg:    v.Sum/v.Count,
		})
	}
}

func addAction(w http.ResponseWriter, r *http.Request) {
	var nilEntry = Stat{}
	var body ActionRequest
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil || len(body.Action) == 0 {
		http.Error(w, "Invalid Input", http.StatusBadRequest)
	}
	if actionMap[body.Action] == nilEntry{
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