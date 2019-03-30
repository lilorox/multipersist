package main

import (
	"fmt"
	"log"
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

func (s SearchResults) CSV() string {
	csv := "size;maxPersistence;numbersCount;nbMostPersistent;searchTime\n"
	for i := 0; i < len(s.results); i++ {
		sr := s.results[i]
		csv += fmt.Sprintf(
			"%d;%d;%d;%d;%f\n",
			sr.size,
			sr.maxPersistence,
			sr.numbersCount,
			len(sr.mostPersistent),
			float64(sr.searchTime)/float64(time.Second),
		)
	}
	return csv
}

func (s SearchResults) Print() {
	fmt.Println("Size       maxPersistence   numbersCount   nbMostPersistent   searchTime")
	for i := 0; i < len(s.results); i++ {
		sr := s.results[i]
		fmt.Printf(
			"%-10d %-16d %-14d %-18d %.4fms\n",
			sr.size,
			sr.maxPersistence,
			sr.numbersCount,
			len(sr.mostPersistent),
			float64(sr.searchTime)/float64(time.Millisecond),
		)
	}
}

func search(size int, maxSize int) SearchResults {
	log.Printf("Starting searching with %d digits\n", size)
	n := NewNumber(size)
	r := SearchResults{
		results: make([]SearchResult, 0),
	}
	sr := NewSearchResult(size)
	start := time.Now()

	for {
		steps := n.Persistence()
		if steps > sr.maxPersistence {
			sr.maxPersistence = steps
			sr.mostPersistent = []string{n.String()}
		} else if steps == sr.maxPersistence {
			sr.mostPersistent = append(sr.mostPersistent, n.String())
		}
		sr.numbersCount++

		if !n.Increment() {
			sr.searchTime = time.Since(start)
			r.results = append(r.results, sr)
			if size >= maxSize {
				break
			}
			log.Printf(
				"Max persistence for %d digits: %d (%.2fms)\n",
				size,
				sr.maxPersistence,
				float64(sr.searchTime)/float64(time.Millisecond),
			)
			size++
			n.Resize(size)
			sr = NewSearchResult(size)
			start = time.Now()
		}
	}
	log.Println("Maximum number of digits to look for attained.")

	return r
}
