package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime/pprof"
	"strconv"
)

var productCache map[int64]int64

type SearchResult struct {
	maxSteps     int
	totalNumbers int
	mostSteps    []string
}

func NewSearchResult() SearchResult {
	sr := SearchResult{}
	sr.mostSteps = make([]string, 1)
	return sr
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
	size      int     // Number of digits of the number
	digits    []int   // Array of digits
	pProducts []int64 // Partial products of the digits
}

func NewNumber(size int) *Number {
	n := Number{
		size:      size,
		digits:    make([]int, size),
		pProducts: make([]int64, size),
	}

	// Starting point is 2666666....
	for i := 0; i < size-1; i++ {
		n.digits[i] = 6
	}
	n.digits[size-1] = 2

	// Array of the partial products of the digits from the highest to the
	// lowest digit
	n.pProducts[size-1] = 2
	for i := size - 2; i >= 0; i-- {
		n.pProducts[i] = n.pProducts[i+1] * int64(n.digits[i])
	}

	return &n
}

func (n *Number) Product() int64 {
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
			n.pProducts[i] = int64(n.digits[i])
		} else {
			n.pProducts[i] = n.pProducts[i+1] * int64(n.digits[i])
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

func (n *Number) ToString() string {
	s := ""
	for i := n.size - 1; i >= 0; i-- {
		s += strconv.Itoa(n.digits[i])
	}
	return s
}

func (n *Number) Details() string {
	return fmt.Sprintf(
		"%s (%d):\n - digits:    %v\n - pProducts: %v\n",
		n.ToString(), n.size, n.digits, n.pProducts,
	)
}

func persistRecursive(n int64, step int) int {
	p := multiply(n)
	if p < 10 {
		return step + 1
	}

	return persistence(p, step+1)
}

func multiply(n int64) int64 {
	var p int64
	p = 1

	for n >= 10000 {
		m := n / 10000
		p *= multiply_4d(n - m*10000)
		if p == 0 {
			return 0
		}
		n = m
	}
	return p * multiply_3d(n)
}

func multiply_4d(n int64) int64 {
	var p int64
	var ok bool
	if p, ok = productCache[n]; ok {
		return p
	}
	p = ((n / 1000) % 10) * ((n / 100) % 10) * ((n / 10) % 10) * (n % 10)
	productCache[n] = p
	return p
}

func multiply_3d(n int64) int64 {
	var p int64
	var ok bool
	if p, ok = productCache[n]; ok {
		return p
	}

	switch {
	case n < 10:
		p = n
	case n < 100:
		p = ((n / 10) % 10) * (n % 10)
	case n < 1000:
		p = ((n / 100) % 10) * ((n / 10) % 10) * (n % 10)
	default:
		p = ((n / 1000) % 10) * ((n / 100) % 10) * ((n / 10) % 10) * (n % 10)
	}
	productCache[n] = p
	return p
}

/// OLD STUFF

func increment(arr *[]int, i int) {
	switch (*arr)[i] {
	case 4:
		// Jump the 5
		(*arr)[i] += 2
	case 9:
		// If this is the last digit and it is already a 9, this is the end
		if i == len(*arr)-1 {
			*arr = nil
			return
		}

		// Edge case when we hit 2999..., the next one is 3466...
		if i == len(*arr)-3 && (*arr)[i+1] == 9 && (*arr)[i+2] == 2 {
			(*arr)[i+2] = 3
			(*arr)[i+1] = 4
			(*arr)[i] = 6
			return
		}

		// Increment the next digit
		increment(arr, i+1)

		// If the increment of the next digit is actually the end, stop now
		if *arr == nil {
			return
		}

		// Once the next digit is incremented, place the current to the same
		// value to avoid duplicate permutations
		(*arr)[i] = (*arr)[i+1]
	default:
		(*arr)[i]++
	}
	return
}

func persistence(n int64, steps int) int {
	p := multiply(n)
	if p < 10 {
		return steps + 1
	}

	return persistence(p, steps+1)
}

func multiplyArray(arr []int) int64 {
	var p int64
	p = 1
	for i := 0; i < len(arr); i++ {
		if arr[i] == 0 {
			return 0
		} else if arr[i] == 1 {
			continue
		}
		p *= int64(arr[i])
	}
	return p
}

func intArrayToStr(arr []int) string {
	s := ""
	for i := len(arr) - 1; i >= 0; i-- {
		s += strconv.Itoa(arr[i])
	}
	return s
}

func searchOLD(size int) SearchResult {
	// Starting point is 2666666....
	arr := make([]int, size)
	for i := 0; i < size-1; i++ {
		arr[i] = 6
	}
	arr[size-1] = 2

	sr := NewSearchResult()
	for arr != nil {
		steps := persistence(multiplyArray(arr), 1)
		if steps > sr.maxSteps {
			sr.maxSteps = steps
			/*
					sr.mostSteps = []string{intArrayToStr(arr)}
				} else if steps == sr.maxSteps {
					sr.mostSteps = append(sr.mostSteps, intArrayToStr(arr))
			*/
		}
		sr.totalNumbers++

		increment(&arr, 0)
	}

	return sr
}

/// END OF OLD STUFF

func search(size int) SearchResult {
	n := NewNumber(size)
	sr := NewSearchResult()

	for {
		steps := n.Persistence()
		if steps > sr.maxSteps {
			sr.maxSteps = steps
			/*
					sr.mostSteps = []string{intArrayToStr(arr)}
				} else if steps == sr.maxSteps {
					sr.mostSteps = append(sr.mostSteps, intArrayToStr(arr))
			*/
		}
		sr.totalNumbers++

		if !n.Increment() {
			break
		}
	}

	return sr
}

func main() {
	cpuprofile := flag.String("cpuprofile", "", "write cpu profile to file")
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	args := flag.Args()
	if len(args) < 1 {
		log.Fatalf("Usage: %s [-cpuprofile=<file>] size\n", os.Args[0])
	}
	productCache = make(map[int64]int64, 11000)

	size, err := strconv.Atoi(args[0])
	if err != nil {
		log.Fatalf("Invalid argument: %s\n", err)
	}

	result := search(size)
	fmt.Printf("Results: %+v\n", result)

	//n := NewNumber(size)
	//for i := 0; i < 100; i++ {
	//	fmt.Print(n.Details())
	//	n.Increment()
	//}
	//fmt.Print(n.Details())
}
