package main

import (
	"flag"
	"log"
	"math"
	"math/big"
	"os"
	"runtime/pprof"
	"strconv"
)

// Global variables
var dCache int
var dCacheLimit *big.Int
var productCache map[string]*big.Int
var powers10 []*big.Int

var cacheHits = 0
var cacheMisses = 0

func main() {
	cpuProfile := *flag.String("cpuprofile", "", "write cpu profile to file")
	dCache = *flag.Int("dcache", 4, "number of digits for the multiplication cache")
	flag.Parse()

	if cpuProfile != "" {
		f, err := os.Create(cpuProfile)
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

	// Initialize caches
	productCache = make(map[string]*big.Int, int(math.Pow10(dCache)+math.Pow10(dCache-1)))
	powers10 = make([]*big.Int, dCache-1)
	dCacheSize := math.Pow10(dCache)
	dCacheLimit = big.NewInt(int64(dCacheSize))
	for i := 1; i < dCache; i++ {
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

	sr := search(start, stop)
	//fmt.Print(sr.CSV())
	sr.Print()

	cacheUsage := 100 * float64(cacheMisses) / dCacheSize
	log.Printf("Cache results: %d hits, %d misses, %.2f%% cache filled\n", cacheHits, cacheMisses, cacheUsage)
}
