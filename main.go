package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"math/big"
	"os"
	"runtime/pprof"
	"strconv"
	"time"
)

// Global variables
var dCache int
var dCacheLimit *big.Int
var productCache map[string]*big.Int
var powers10 []*big.Int

var cacheHits = 0
var cacheMisses = 0

// Frequently used big constants
var big0 = big.NewInt(0)
var big10 = big.NewInt(10)

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

// Represent a number as an array of digits with the added bonus of computing
// the partial products while incrementing it.
//
// What this struct looks like for the number 2778889:
//   digits    = [ 9      8      8      8      7      7      2 ]
//   pProducts = [ 129024 14336  1792   224    28     14     2 ]
//   (index)       0      1      2      3      4      5      6
//   (designation) lowest                                    highest
type Number struct {
	size      int        // Number of digits of the number
	digits    []int      // Array of digits
	pProducts []*big.Int // Partial products of the digits
}

func NewNumber(size int) *Number {
	n := Number{
		size:      size,
		digits:    make([]int, size),
		pProducts: make([]*big.Int, size),
	}

	// Starting point is 2666666....
	for i := 0; i < size-1; i++ {
		n.digits[i] = 6
	}
	n.digits[size-1] = 2

	// Array of the partial products of the digits from the highest to the
	// lowest digit
	n.pProducts[size-1] = big.NewInt(2)
	for i := size - 2; i >= 0; i-- {
		n.pProducts[i] = new(big.Int).Mul(n.pProducts[i+1], big.NewInt(int64(n.digits[i])))
	}

	return &n
}

func (n *Number) Product() *big.Int {
	return n.pProducts[0]
}

func (n *Number) Increment() bool {
	highest := n.incRecursive(0)
	if highest == -1 {
		return false
	}

	// Update all the partial products down from the highest updated digit
	for i := highest; i >= 0; i-- {
		if i == n.size-1 {
			n.pProducts[i] = big.NewInt(int64(n.digits[i]))
		} else {
			n.pProducts[i].Mul(n.pProducts[i+1], big.NewInt(int64(n.digits[i])))
		}
	}

	return true
}

func (n *Number) incRecursive(i int) int {
	highest := i

	switch n.digits[i] {
	case 4:
		// Jump the 5
		n.digits[i] += 2
	case 9:
		// If this is the last digit and it is already a 9, this is the end
		if i == n.size-1 {
			return -1
		}

		// Edge case when we hit 2999..., the next one is 3466...
		if i == n.size-3 && n.digits[i+1] == 9 && n.digits[i+2] == 2 {
			n.digits[i+2] = 3
			n.digits[i+1] = 4
			n.digits[i] = 6
			return i + 2
		}

		// Increment the next digit
		highest = n.incRecursive(i + 1)

		// If we haven't reached the highest number of this size, now that the
		// next digit has been incremented, place the current one to the same
		// value to avoid duplicate permutations
		if highest != -1 {
			n.digits[i] = n.digits[i+1]
		}
	default:
		n.digits[i]++
	}
	return highest
}

func (n *Number) Persistence() int {
	return persistRecursive(n.Product(), 1)
}

func (n *Number) String() string {
	s := ""
	for i := n.size - 1; i >= 0; i-- {
		s += strconv.Itoa(n.digits[i])
	}
	return s
}

func (n *Number) Details() string {
	return fmt.Sprintf(
		"%s (%d):\n - digits:    %v\n - pProducts: %v\n",
		n.String(), n.size, n.digits, n.pProducts,
	)
}

func persistRecursive(n *big.Int, step int) int {
	p := multiplyDigits(n)
	//fmt.Printf(" [step=%d]: p=%s\n", step, p.String()) // XXX DEBUG

	if p.Cmp(big.NewInt(10)) == -1 {
		return step + 1
	}

	return persistRecursive(p, step+1)
}

func multiplyDigits(n *big.Int) *big.Int {
	p := big.NewInt(1)
	q := new(big.Int)
	r := new(big.Int)

	for n.Cmp(dCacheLimit) >= 0 {
		q.QuoRem(n, dCacheLimit, r)
		//fmt.Printf("  n = %s || q = %s || r = %s || p = %s\n", n.String(), q.String(), r.String(), p.String()) // XXX DEBUG
		p.Mul(p, multiplyDigitsWithCache(r))
		//fmt.Printf("  p --> %s\n", p.String()) // XXX DEBUG
		if p.Cmp(big0) == 0 {
			return p
		}
		n.Set(q)
	}
	p.Mul(p, multiplyDigitsWithCache(n))
	return p
}

func multiplyDigitsWithCache(n *big.Int) *big.Int {
	s := n.String()
	if pCached, ok := productCache[s]; ok {
		cacheHits++
		//fmt.Printf("    cached: %s --> %s\n", n.String(), pCached.String()) // XXX DEBUG
		return pCached
	}

	p := new(big.Int)
	p.Rem(n, big10)
	if p.Cmp(big0) != 0 {
		//q := new(big.Int)
		r := new(big.Int)
		for i := 0; i < len(powers10) && n.Cmp(powers10[i]) > 0; i++ {
			//q.Quo(n, powers10[i])

			r.Rem(r.Quo(n, powers10[i]), big10)
			//fmt.Printf("    n = %s || q = %s || r = %s || p = %s --> ", n.String(), q.String(), r.String(), p.String()) // XXX DEBUG
			if r.Cmp(big0) == 0 {
				p.Set(big0)
				//fmt.Printf("0\n") // XXX DEBUG
				break
			}
			p.Mul(r, p)
			//fmt.Printf("%s\n", p.String()) // XXX DEBUG
		}
	}
	productCache[s] = p
	cacheMisses++
	return p
}

func search(size int) SearchResult {
	n := NewNumber(size)
	sr := NewSearchResult(size)
	start := time.Now()

	for {
		steps := n.Persistence()
		//fmt.Printf("%s -> %d\n", n.String(), steps) // XXX DEBUG
		if steps > sr.maxPersistence {
			sr.maxPersistence = steps
			sr.mostPersistent = []string{n.String()}
		} else if steps == sr.maxPersistence {
			sr.mostPersistent = append(sr.mostPersistent, n.String())
		}
		sr.numbersCount++

		if !n.Increment() {
			break
		}
	}
	sr.searchTime = time.Since(start)

	return sr
}

func main() {
	cpuProfile := flag.String("cpuprofile", "", "write cpu profile to file")
	dCache := flag.Int("dcache", 4, "number of digits for the multiplication cache")
	flag.Parse()

	if *cpuProfile != "" {
		f, err := os.Create(*cpuProfile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	args := flag.Args()
	if len(args) < 1 {
		log.Fatalf("Usage: %s [-cpuprofile=<file>] [-dCache=<int>] size\n", os.Args[0])
	}

	// Build caches
	productCache = make(map[string]*big.Int, int(math.Pow10(*dCache)+math.Pow10(*dCache-1)))
	powers10 = make([]*big.Int, *dCache-1)
	dCacheLimit = big.NewInt(int64(math.Pow10(*dCache)))
	for i := 1; i < *dCache; i++ {
		powers10[i-1] = big.NewInt(int64(math.Pow10(i)))
	}

	var start = 2
	var stop = 2
	var err error
	if len(args) == 1 {
		stop, err = strconv.Atoi(args[0])
		start = stop
	} else {
		start, err = strconv.Atoi(args[0])
		stop, err = strconv.Atoi(args[1])
	}
	if err != nil {
		log.Fatalf("Invalid argument: %s\n", err)
	}

	sr := SearchResults{
		results: make([]SearchResult, 0),
	}
	for i := start; i <= stop; i++ {
		sr.results = append(sr.results, search(i))
	}
	fmt.Print(sr.ToCSV())
	fmt.Printf("Cache: limit=%d hits=%d misses=%d\n", dCacheLimit, cacheHits, cacheMisses)
}
