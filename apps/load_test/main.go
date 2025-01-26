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
	OrgID          string `json:"org_id"`
	HousingType    string `json:"housing_type"`
	PrivacyType    string `json:"privacy_type"`
	TripType       string `json:"trip_type"`
	Status         string `json:"status"`
	VolunteerLimit int    `json:"volunteer_limit"`
	Name           string `json:"name"`
	Description    string `json:"description"`
}

// Prepare mock data for trips
var trips = []Trip{
	{"load_test", "house", "complete", "domestic", "listed", 2, "cat watcher", "Need someone to watch over the neighbor's cat while they're away."},
	{"load_test", "apartment", "complete", "domestic", "listed", 4, "library organizer", "We need help organizing a local library. Volunteers needed."},
	{"load_test", "apartment", "complete", "domestic", "listed", 1, "concert assistant", "Assist with setup and management of the concert."},
	{"load_test", "hostel", "shared", "international", "draft", 2, "pope sitter", "We need someone to watch the pope while he goes about his day. It's very important."},
	{"load_test", "apartment", "complete", "domestic", "draft", 3, "event assistant", "Help assist with the upcoming event, need an extra hand."},
	{"load_test", "hostel", "shared", "domestic", "draft", 1, "food distributor", "Someone is needed to distribute food to the local shelter."},
	{"load_test", "house", "complete", "domestic", "draft", 2, "child care assistant", "Help take care of kids at a daycare center."},
	{"load_test", "house", "complete", "domestic", "draft", 3, "elder care", "Provide care and assistance for elderly individuals."},
	{"load_test", "hotel", "complete", "domestic", "draft", 2, "sports event helper", "Help organize and manage a local sports event."},
	{"load_test", "hostel", "shared", "domestic", "draft", 2, "festival coordinator", "We need help coordinating the festival, volunteers appreciated."},
}

const apiURL = "http://localhost:8080/v1/trips?org_id=load_test"

// Function to send a single request
func sendTripRequest(trip Trip, wg *sync.WaitGroup) {
	defer wg.Done()

	jsonData := fmt.Sprintf(`{
		"org_id": "%s",
		"housing_type": "%s",
		"privacy_type": "%s",
		"trip_type": "%s",
		"status": "%s",
		"volunteer_limit": %d,
		"name": "%s",
		"description": "%s"
	}`, trip.OrgID, trip.HousingType, trip.PrivacyType, trip.TripType, trip.Status, trip.VolunteerLimit, trip.Name, trip.Description)

	resp, err := http.Post(apiURL, "application/json", bytes.NewBuffer([]byte(jsonData)))
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()

	fmt.Printf("Sent request for: %s | Status Code: %d\n", trip.Name, resp.StatusCode)
}

func main() {
	var (
		wg      sync.WaitGroup
		limiter = time.Tick(40 * time.Millisecond)
	)

	fmt.Println("Starting load test...")
	time.Sleep(1 * time.Second)
	for i := 0; i < 100000; i++ {
		wg.Add(1)
		trip := trips[i%len(trips)]
		go func(trip Trip) {
			<-limiter
			sendTripRequest(trip, &wg)
		}(trip)
	}
	wg.Wait()
	fmt.Println("All requests sent!")
}
