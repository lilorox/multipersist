package main

import (
	"fmt"
	"time"
)

// SearchResult holds the result of the search for the maximum persistence of an
// integer of a certain size.
type SearchResult struct {
	size           int           // Number of digits of the integer
	maxPersistence int           // Maximum persistence found across all integer of the same size
	numbersCount   int           // Total count of number that have been computed
	mostPersistent []string      // List of the numbers that reach the maximum persistence
	searchTime     time.Duration // Duration of the search
}

func NewSearchResult(size int) SearchResult {
	sr := SearchResult{
		size:           size,
		mostPersistent: make([]string, 1),
	}
	return sr
}

type SearchResults struct {
	results []SearchResult
}

func (s SearchResults) ToCSV() string {
	csv := "size;maxPersistence;numbersCount;searchTime\n"
	for i := 0; i < len(s.results); i++ {
		sr := s.results[i]
		csv += fmt.Sprintf(
			"%d;%d;%d;%f\n",
			sr.size,
			sr.maxPersistence,
			sr.numbersCount,
			float64(sr.searchTime)/float64(time.Second),
		)
	}
	return csv
}
