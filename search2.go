package main

import (
	"fmt"
	"math/big"
	"sync"
	"time"
)

type result struct {
	size             int           // Number of digits of the integer
	maxPersistence   int           // Maximum persistence found across all integer of the same size
	numbersCount     int           // Total count of number that have been computed
	nbMostPersistent int           // Count of the numbers that reach the maximum persistence
	searchTime       time.Duration // Duration of the search
	mutex            sync.Mutex
}

func (r *result) updateSteps(steps int) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if steps > r.maxPersistence {
		r.maxPersistence = steps
		r.nbMostPersistent++
	} else if steps == r.maxPersistence {
		r.nbMostPersistent = 1
	}
	r.numbersCount++
}

type Results []*result

func NewResults() *Results {
	var r Results
	return &r
}

func (r *Results) Print() {
	fmt.Printf("Print %+v\n", r)
	fmt.Println("┍    Size    ┯  maxPers.  ┯  mostPers. ┯ TotalCount  ┯    Time    ┑")
	for i := 0; i < len(*r); i++ {
		res := (*r)[i]
		fmt.Printf(
			"│ %-11d│ %-11d│ %-11d│ %-11d│ %-5.4fms│\n",
			res.size,
			res.maxPersistence,
			res.nbMostPersistent,
			res.numbersCount,
			float64(res.searchTime)/float64(time.Millisecond),
		)
	}
	fmt.Println("┕━━━━━━━━━━━━┷━━━━━━━━━━━━┷━━━━━━━━━━━━┷━━━━━━━━━━━━━┷━━━━━━━━━━━━┙")
}

type Workload struct {
	product *big.Int
	result  *result
}

func generateNumbers(size int, maxSize int) (chan *Workload, chan *Results) {
	work := make(chan *Workload, 1)
	end := make(chan *Results, 1)

	go func() {
		var results Results
		defer func() {
			end <- &results
		}()
		var m *big.Int
		n := NewNumber(size)
		r := &result{size: size}
		start := time.Now()

		for {
			m = n.Product()
			work <- &Workload{
				product: m,
				result:  r,
			}

			m = n.Increment()
			// End of the current search now we move to the next size or end here
			if m == nil {
				r.searchTime = time.Since(start)
				results.Append(r)
				log.Printf(
					"Maximum persistence for size %d: %d (%.4fms)\n",
					r.maxPersistence,
					r.size,
					float64(r.searchTime)/float64(time.Millisecond),
				)
				fmt.Printf("New len(results)=%d\n", len(results))
				if size >= maxSize {
					break
				}

				size++
				n.Resize(size)

				r = &result{size: size}
				start = time.Now()
			}
		}
	}()

	return work, end
}

func search2(work chan *Workload, end chan *Results) {
	for {
		select {
		case results := <-end:
			results.Print()
			return
		case w := <-work:
			w.result.updateSteps(persistence(w.product, 1))
		}
	}
}
