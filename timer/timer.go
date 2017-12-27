package timer

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"log"
)

//Request JSON represent the timer JSON reuests
type Request struct {
	ID          string `json:"id"`
	Destination string `json:"destination"`
	Interval    uint   `json:"interval"`
	Message     string `json:"message"`
	Enabled     bool   `json:"enabled"`
}

//Bind the data to json request
func (s *Request) Bind(r *http.Request) error {
	return nil
}

//SingleTimer is the timer Object
type SingleTimer struct {
	ID          string
	Destination string
	Message     string
	Running     bool
	Timing      uint `json:"interval"`
	interval    time.Duration
	Enabled     bool
	ticker      *time.Ticker
	quit        chan struct{}
}

var timers map[string]*SingleTimer
var userTimers map[string][]*SingleTimer

const resolution time.Duration = time.Second

func init() {
	timers = make(map[string]*SingleTimer)
}

//GetAllTimers Get All timers
func GetAllTimers() *map[string]*SingleTimer {
	return &timers
}

//Health Check
func Health() {
	log.Println("200")
}

//Create a new timer
func Create(s Request) (*SingleTimer, error) {
	if timers[s.ID] != nil {
		return nil, fmt.Errorf("timer id %s already exists", s.ID)
	}
	t := SingleTimer{
		ID:          s.ID,
		Destination: s.Destination,
		Message:     s.Message,
		Enabled:     s.Enabled,
	}
	log.Println("new timer", t.Message)
	t.interval = time.Duration(s.Interval) * resolution
	t.Timing = s.Interval
	t.quit = make(chan struct{})
	timers[s.ID] = &t
	return &t, nil
}

//Get - get timer by ID
func Get(id string) (*SingleTimer, error) {
	t := timers[id]
	if t == nil {
		return nil, fmt.Errorf("timer id %q does not exists", id)
	}
	return t, nil
}

// Update the timer parameters
func (s *SingleTimer) Update(r Request) {
	if r.Destination != "" {
		s.Destination = r.Destination
	}
	if r.Destination != "" {
		s.Message = r.Message
	}
	s.Enabled = r.Enabled
	if r.Interval > 0 {
		s.SetInterval(time.Duration(r.Interval) * resolution)
	}
}

//Stop the Timer
func (s *SingleTimer) Stop() {
	if !s.Running {
		return
	}
	s.Running = false

	s.quit <- struct{}{}
}

//Delete the timer
func (s *SingleTimer) Delete() {
	delete(timers, s.ID)
}

//SetInterval - set the timer interval
func (s *SingleTimer) SetInterval(interval time.Duration) {
	s.interval = interval
	s.Timing = uint(interval / resolution)
	if s.ticker != nil {
		s.Stop()
		s.Start()
	}
}

//Start - Make the timer work
func (s *SingleTimer) Start() error {
	if s.Running {
		return nil
	}

	if s.Destination == "" {
		return errors.New("cannot start timer with no destination")
	}
	s.Running = true
	go func() {
		s.ticker = time.NewTicker(s.interval)

		defer s.ticker.Stop()
		for {
			select {
			case <-s.ticker.C:
				if s.Enabled {
					log.Println(s.Message)

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
	return nil
}
