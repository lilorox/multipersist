package main

import (
	"fmt"
	"log"
	"os"
	"time"
)

// Result holds the result of the search for the maximum persistence of an
// integer of a certain size.
type Result struct {
	size             int           // Number of digits of the integer
	maxPersistence   int           // Maximum persistence found across all integer of the same size
	numbersCount     int           // Total count of number that have been computed
	nbMostPersistent int           // Count of numbers reaching the maximum persistence
	searchTime       time.Duration // Duration of the search
}

// Results gathers a collection of results of a run
type Results []*Result

func NewResults() *Results {
	var r Results
	return &r
}

func (r *Results) ExportToCSV(filename string) {
	f, err := os.Create(filename)
	if err != nil {
		log.Fatalf("Cannot create file %s: %s", filename, err)
	}
	defer f.Close()

	csv := "size;maxPersistence;numbersCount;nbMostPersistent;searchTime\n"
	for i := 0; i < len(*r); i++ {
		res := (*r)[i]
		csv += fmt.Sprintf(
			"%d;%d;%d;%d;%f\n",
			res.size,
			res.maxPersistence,
			res.numbersCount,
			res.nbMostPersistent,
			float64(res.searchTime)/float64(time.Second),
		)
	}
	_, err = f.WriteString(csv)
	if err != nil {
		log.Fatalf("Cannot write to file %s: %s", filename, err)
	}
	f.Sync()
}

func (r *Results) Print() {
	fmt.Println("┍    Size    ┯  maxPers.  ┯  mostPers. ┯ TotalCount ┯   Time (ms)  ┑")
	for i := 0; i < len(*r); i++ {
		res := (*r)[i]
		fmt.Printf(
			"│ %-11d│ %-11d│ %-11d│ %-11d│ %-13.4f│\n",
			res.size,
			res.maxPersistence,
			res.nbMostPersistent,
			res.numbersCount,
			float64(res.searchTime)/float64(time.Millisecond),
		)
	}
	fmt.Println("┕━━━━━━━━━━━━┷━━━━━━━━━━━━┷━━━━━━━━━━━━┷━━━━━━━━━━━━┷━━━━━━━━━━━━━━┙")
}
