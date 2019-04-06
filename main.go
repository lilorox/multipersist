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
)

// Global variables
var (
	//Build management
	Version string
	Build   string

	// Cache management
	dCache           int
	dCacheLimit      *big.Int
	productCache     map[string]*big.Int
	productCacheSize int
	powers10         []*big.Int

	// Cache statistics
	cacheHits   = 0
	cacheMisses = 0
)

func initCache(size int) {
	dCache = size
	dCacheSize := math.Pow10(size)
	dCacheLimit = big.NewInt(int64(dCacheSize))

	productCacheSize = int(dCacheSize + math.Pow10(dCache-1))
	productCache = make(map[string]*big.Int, productCacheSize)

	powers10 = make([]*big.Int, dCache-1)
	for i := 1; i < dCache; i++ {
		powers10[i-1] = big.NewInt(int64(math.Pow10(i)))
	}
}

func main() {
	version := flag.Bool("version", false, "version information")
	cpuProfile := flag.String("cpuprofile", "", "write cpu profile to file")
	dCacheFlag := flag.Int("dcache", 4, "number of digits for the multiplication cache")
	flag.Parse()

	if *version {
		fmt.Printf("Version: %s, Build: %s\n", Version, Build)
		os.Exit(0)
	}

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
		fmt.Fprintf(
			os.Stderr,
			"Usage: %s [-version] [-cpuprofile=<file>] [-dCache=<int>] size\n",
			os.Args[0],
		)
		os.Exit(1)
	}

	// Initialize caches
	initCache(*dCacheFlag)

	var (
		start = 2
		stop  = 2
		err   error
	)
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

	results := NewResults()
	for size := start; size <= stop; size++ {
		s := NewSearcher(size)
		*results = append(*results, s.Search())
	}
	results.Print()
	log.Printf("Cache stats: %d hits, %d misses, %d entries\n", cacheHits, cacheMisses, productCacheSize)
}
