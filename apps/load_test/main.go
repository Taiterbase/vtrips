package main

import (
	"bytes"
	"fmt"
	"net/http"
	"sync"
	"time"
)

// Define the data structure for a trip request.
type Trip struct {
	ClientID       string `json:"client_id"`
	Status         string `json:"status"`
	VolunteerLimit int    `json:"volunteer_limit"`
	Name           string `json:"name"`
	Description    string `json:"description"`
}

// Prepare mock data for trips
var trips = []Trip{
	{"load", "draft", 2, "pope sitter", "We need someone to watch the pope while he goes about his day. It's very important."},
	{"load", "draft", 3, "event assistant", "Help assist with the upcoming event, need an extra hand."},
	{"load", "draft", 1, "food distributor", "Someone is needed to distribute food to the local shelter."},
	{"load", "draft", 2, "cat watcher", "Need someone to watch over the neighbor's cat while they're away."},
	{"load", "draft", 4, "library organizer", "We need help organizing a local library. Volunteers needed."},
	{"load", "draft", 1, "concert assistant", "Assist with setup and management of the concert."},
	{"load", "draft", 2, "child care assistant", "Help take care of kids at a daycare center."},
	{"load", "draft", 3, "elder care", "Provide care and assistance for elderly individuals."},
	{"load", "draft", 2, "sports event helper", "Help organize and manage a local sports event."},
	{"load", "draft", 2, "festival coordinator", "We need help coordinating the festival, volunteers appreciated."},
}

const apiURL = "http://localhost:8080/v1/trips?client_id=load"

// Function to send a single request
func sendTripRequest(trip Trip, wg *sync.WaitGroup) {
	defer wg.Done()

	// Convert trip data to JSON
	jsonData := fmt.Sprintf(`{
		"client_id": "%s",
		"status": "%s",
		"volunteer_limit": %d,
		"name": "%s",
		"description": "%s"
	}`, trip.ClientID, trip.Status, trip.VolunteerLimit, trip.Name, trip.Description)

	// Send the request
	resp, err := http.Post(apiURL, "application/json", bytes.NewBuffer([]byte(jsonData)))
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()

	// Print status of the request
	fmt.Printf("Sent request for: %s | Status Code: %d\n", trip.Name, resp.StatusCode)
}

func main() {
	fmt.Println("Starting load test...")
	time.Sleep(1 * time.Second)

	var wg sync.WaitGroup

	// Set up a rate limiter to avoid overloading the API
	limiter := time.Tick(1 * time.Millisecond)

	// Spin off 1000 requests concurrently
	for i := 0; i < 1000; i++ {
		wg.Add(1)

		// Choose a random trip (or you could cycle through them)
		trip := trips[i%len(trips)]

		go func(trip Trip) {
			<-limiter // Rate-limiting to avoid spamming requests too fast
			sendTripRequest(trip, &wg)
		}(trip)
	}

	// Wait for all goroutines to finish
	wg.Wait()

	fmt.Println("All requests sent!")
}
