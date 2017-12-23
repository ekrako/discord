package timer

import (
	"net/http"
	"strings"
	"time"

	"log"
)

type TimerSettings struct {
	ID          string `json:"_id"`
	Destination string `json:"destination"`
	Interval    uint   `json:"interval"`
	Message     string `json:"message"`
	Enabled     bool   `json:"enabled"`
}

type SingleTimer struct {
	ID          string
	Destination string
	Message     string
	interval    time.Duration
	Enabled     bool
	ticker      *time.Ticker
	quit        chan struct{}
	deleted     bool
}

var timers map[string]SingleTimer

func init() {
	timers = make(map[string]SingleTimer)
}

//Health Check
func Health() {
	log.Println("200")
}

//Create a new timer
func Create(id string, destination string, messege string, enabled bool, interval uint) *SingleTimer {
	// func Create() *SingleTimer {
	t := SingleTimer{

		ID:          id,
		Destination: destination,
		Message:     messege,
		Enabled:     enabled,
	}
	log.Println("new timer", t.Message)
	t.interval = time.Duration(interval) * time.Second
	t.quit = make(chan struct{})
	timers[id] = t
	return &t
}

//Get - get timer by ID
func Get(id string) *SingleTimer {
	t := timers[id]
	return &t
}

//Stop the Timer
func (s *SingleTimer) Stop() {
	if s.deleted {
		return
	}
	s.deleted = true

	close(s.quit)
}

//SetInterval - set the timer interval
func (s *SingleTimer) SetInterval(interval time.Duration) {
	s.interval = interval
	if s.ticker != nil {
		s.ticker.Stop()
		s.ticker = time.NewTicker(s.interval)
	}
}

//Run - Make the timer work
func (s *SingleTimer) Run() {

	go func() {
		s.ticker = time.NewTicker(s.interval)
		defer s.ticker.Stop()
		for {
			select {
			case <-s.ticker.C:
				if s.Enabled {
					log.Println(s.Message)
					log.Println(s.interval)

					resp, err := http.Post(s.Destination, "application/json", strings.NewReader(s.Message))
					if err != nil {
						log.Printf("faild to send message %s to %s", s.Message, s.Destination)
					}
					defer resp.Body.Close()
				}
			case <-s.quit:
				return
			}
		}
	}()
}
