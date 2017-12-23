package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

type Timer struct {
	ID        string   `json:"id,omitempty"`
	Firstname string   `json:"firstname,omitempty"`
	Lastname  string   `json:"lastname,omitempty"`
	Address   *Address `json:"address,omitempty"`
}

type Address struct {
	City  string `json:"city,omitempty"`
	State string `json:"state,omitempty"`
}

var timers map[string]Timer
var timersIndex int

func GetTimerEndpoint(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("content-type", "application/json")
	id := chi.URLParam(req, "id")
	json.NewEncoder(w).Encode(timers[id])
}

func GetTimersEndpoint(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(timers)
}

func CreateTimerEndpoint(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("content-type", "application/json")
	var timer Timer
	_ = json.NewDecoder(req.Body).Decode(&timer)
	timer.ID = strconv.Itoa(timersIndex)
	timers[timer.ID] = timer
	timersIndex++
	json.NewEncoder(w).Encode(timer)
}

func DeleteTimerEndpoint(w http.ResponseWriter, req *http.Request) {
	id := chi.URLParam(req, "id")
	delete(timers, id)
}

func main() {
	router := chi.NewRouter()
	timers = make(map[string]Timer)
	timers["1"] = Timer{ID: "1", Firstname: "Nic", Lastname: "Raboy", Address: &Address{City: "Dublin", State: "CA"}}
	timers["2"] = Timer{ID: "2", Firstname: "Maria", Lastname: "Raboy"}
	timersIndex = 3
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Get("/timers", GetTimersEndpoint)
	router.Get("/timer/{id}", GetTimerEndpoint)
	router.Delete("/timer/{id}", DeleteTimerEndpoint)
	router.Post("/timers/", CreateTimerEndpoint)
	log.Fatal(http.ListenAndServe(":8080", router))
}
