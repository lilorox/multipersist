package main

import (
	"fmt"
	"math/big"
	"strconv"
)

// Frequently used big constants
var (
	big0  = big.NewInt(0)
	big10 = big.NewInt(10)
)

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
	n := &Number{
		size:      size,
		digits:    make([]int, size),
		pProducts: make([]*big.Int, size),
	}

	n.initDigits()
	n.initProducts()

	return n
}

// initDigits fills the digits slice with the starting Number of this size.
// Starting point is 2666666....
func (n *Number) initDigits() {
	for i := 0; i < n.size-1; i++ {
		n.digits[i] = 6
	}
	n.digits[n.size-1] = 2
}

// initProducts compute the pProducts slice from the digits
func (n *Number) initProducts() {
	n.pProducts[n.size-1] = big.NewInt(int64(n.digits[n.size-1]))
	for i := n.size - 2; i >= 0; i-- {
		if n.pProducts[i] == nil {
			n.pProducts[i] = new(big.Int)
		}
		n.pProducts[i] = new(big.Int).Mul(
			n.pProducts[i+1],
			big.NewInt(int64(n.digits[i])),
		)
	}
}

func (n *Number) Size() int {
	return n.Size()
}

func (n *Number) Resize(size int) {
	if size > n.size {
		n.digits = append(n.digits, make([]int, size-n.size)...)
		n.pProducts = append(n.pProducts, make([]*big.Int, size-n.size)...)
	} else if size < n.size {
		n.digits = make([]int, size)
		n.pProducts = make([]*big.Int, size)
	}

	n.size = size
	n.initDigits()
	n.initProducts()
}

func (n *Number) Increment() *big.Int {
	highest := n.incRecursive(0)
	if highest == -1 {
		return nil
	}

	// Update all the partial products down from the highest updated digit
	for i := highest; i >= 0; i-- {
		if i == n.size-1 {
			n.pProducts[i] = big.NewInt(int64(n.digits[i]))
		} else {
			n.pProducts[i].Mul(n.pProducts[i+1], big.NewInt(int64(n.digits[i])))
		}
	}

	return n.pProducts[0]
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

func (n *Number) Product() *big.Int {
	return n.pProducts[0]
}

func persistence(n *big.Int, step int) int {
	p := multiplyDigits(n)

	if p.Cmp(big10) == -1 {
		return step + 1
	}

	return persistence(p, step+1)
}

func multiplyDigits(n *big.Int) *big.Int {
	p := big.NewInt(1)
	s := n.String()

	for len(s) > dCache {
		p.Mul(p, multiplyDigitsWithCache(s[len(s)-dCache:]))
		if p.Cmp(big0) == 0 {
			return p
		}
		s = s[:len(s)-dCache]
	}
	if len(s) > 0 {
		p.Mul(p, multiplyDigitsWithCache(s))
	}
	return p
}

func multiplyDigitsWithCache(s string) *big.Int {
	if pCached, ok := productCache[s]; ok {
		cacheHits++
		return pCached
	}

	p := new(big.Int)
	if s[0] != '0' {
		p.SetInt64(int64(s[0] - '0'))
		for i := 1; i < len(s); i++ {
			if s[i] == '0' {
				p.Set(big0)
				break
			}
			p.Mul(p, big.NewInt(int64(s[i]-'0')))
		}
	}
	productCache[s] = p
	cacheMisses++
	return p
}
