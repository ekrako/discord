package main

import (
	"log"
	"time"

	"github.com/ekrako/discord/timer"
)

func main() {
	timer.Health()
	t := timer.Create("ID", "https://requestb.in/14tpopr1", "Hello World 1", true, 1)
	t = timer.Create("ID2", "https://requestb.in/14tpopr1", "Hello World 2", true, 1)
	t.SetInterval(500 * time.Millisecond)
	log.Println("Change 0.5  time.Second")
	t.Run()
	time.Sleep(5 * time.Second)
	t.Message = "Hi There 2"
	log.Println("Changed message to timer 2 to hi there")
	t = timer.Get("ID3")
	t.Run()
	time.Sleep(5 * time.Second)
	t.SetInterval(time.Second)
	log.Println("Change interval Second")
	time.Sleep(5 * time.Second)
	log.Println("stopping ticker...")
	t.Stop()
	// log.Fatal(http.ListenAndServe(":8080", nil))
}

// ID:          "s.ID",
// Destination: "https://requestb.in/14tpopr1",
// Message:     "s.Message",
// Enabled:     true,
