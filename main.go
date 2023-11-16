package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

type Interval struct {
	Start int
	End   int
}

func handleAPI(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Includes []Interval `json:"includes"`
		Excludes []Interval `json:"excludes"`
	}

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	output := processIntervals(request.Includes, request.Excludes)

	response := struct {
		Output []Interval `json:"output"`
	}{
		Output: output,
	}

	log.Info("Received an API request")
	log.Infof("Request: %+v", request)
	log.Infof("Response: %+v", response)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func mergeIntervals(intervals []Interval) []Interval {
	// Initialize an empty slice to store the merged intervals
	merged := []Interval{}

	// Sort the intervals based on the start time
	sort.Slice(intervals, func(i, j int) bool {
		return intervals[i].Start < intervals[j].Start
	})

	// Iterate over the sorted intervals
	for _, interval := range intervals {
		// If the merged slice is empty or if the current interval does not overlap with the previous, append it to merged
		if len(merged) == 0 || merged[len(merged)-1].End < interval.Start {
			merged = append(merged, interval)
		} else {
			// Otherwise, there is an overlap, so we merge the current and previous intervals.
			merged[len(merged)-1].End = max(merged[len(merged)-1].End, interval.End)
		}
	}
	// Return the merged intervals
	return merged
}

func processIntervals(includes, excludes []Interval) []Interval {
	result := []Interval{}

	// Merge overlapping intervals in includes and excludes
	includes = mergeIntervals(includes)
	excludes = mergeIntervals(excludes)

	// Process intervals
	for _, include := range includes {
		for _, exclude := range excludes {
			// If exclude interval is within include interval
			if exclude.Start > include.End || exclude.End < include.Start {
				continue
			} else {
				// Adjust include interval based on exclude interval
				if exclude.Start <= include.Start && exclude.End >= include.End {
					// Exclude interval completely covers include interval
					include = Interval{}
					break
				} else if exclude.Start <= include.Start {
					// Exclude interval overlaps with the start of include interval
					include.Start = exclude.End + 1
				} else if exclude.End >= include.End {
					// Exclude interval overlaps with the end of include interval
					include.End = exclude.Start - 1
				} else {
					// Exclude interval is in the middle of include interval
					result = append(result, Interval{Start: include.Start, End: exclude.Start - 1})
					include.Start = exclude.End + 1
				}
			}
		}
		if include != (Interval{}) {
			result = append(result, include)
		}
	}

	return result
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func main() {
	router := mux.NewRouter()

	// Enable CORS
	headersOk := handlers.AllowedHeaders([]string{"Content-Type"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

	// API endpoint
	router.HandleFunc("/api/process", handleAPI).Methods("POST")

	// Start the server
	fmt.Println("Server is running on port 8080...")
	// Start the server with CORS
	http.ListenAndServe(":8080", handlers.CORS(originsOk, headersOk, methodsOk)(router))
}
