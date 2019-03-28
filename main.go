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
		r := new(big.Int)
		for i := 0; i < len(powers10) && n.Cmp(powers10[i]) > 0; i++ {
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
		//fmt.Printf("%s -> %d\n", n.String(), steps) // XXX DEBUG
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
			log.Printf("Computation for %d digits in %.2fms\n", size, float64(sr.searchTime)/float64(time.Millisecond))
			size++
			n = NewNumber(size)
			sr = NewSearchResult(size)
			start = time.Now()
		}
	}
	log.Println("Maximum number of digits to look for attained.")

	return r
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
	dCacheSize := math.Pow10(*dCache)
	dCacheLimit = big.NewInt(int64(dCacheSize))
	for i := 1; i < *dCache; i++ {
		powers10[i-1] = big.NewInt(int64(math.Pow10(i)))
	}
	log.Printf("Cache size: %d entries\n", int64(dCacheSize))

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

	if start < 3 {
		log.Fatalf("Cannot be used with integers with less than 3 digits\n")
	}

	sr := search(start, stop+1)
	fmt.Print(sr.CSV())

	cacheUsage := 100 * float64(cacheMisses) / dCacheSize
	log.Printf("Cache results: %d hits, %d misses, %.2f%% cache filled\n", cacheHits, cacheMisses, cacheUsage)
}
