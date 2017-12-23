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

var timers []SingleTimer

func Health() {
	log.Println("200")
}

func Create() *SingleTimer {
	t := SingleTimer{

		ID:          "s.ID",
		Destination: "https://requestb.in/14tpopr1",
		Message:     "s.Message",
		Enabled:     true,
	}
	log.Println("new timer", t.Message)
	t.interval = time.Second
	t.quit = make(chan struct{})
	return &t
}
func (s *SingleTimer) Stop() {
	if s.deleted {
		return
	}
	s.deleted = true

	close(s.quit)
}
func (s *SingleTimer) SetInterval(interval time.Duration) {
	s.interval = interval
	s.ticker.Stop()
	s.ticker = time.NewTicker(s.interval)
}

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
