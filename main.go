package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
)

func persistence(nbStr string, stepsIn int) (stepsOut int, err error) {
	res := 1
	stepsOut = stepsIn + 1

	for _, c := range nbStr {
		if c == '0' {
			fmt.Println("   0")
			return
		}

		// Trick to convert a rune digit into an int
		res *= int(c - '0')
	}
	fmt.Printf("   %d\n", res)

	if res < 10 {
		return
	}

	return persistence(strconv.Itoa(res), stepsOut)
}

func generate(size int) (c chan string, err error) {
	if size < 3 {
		err = errors.New("Won't generate integers that have less than 3 digits.")
		return
	}

	c = make(chan string)

	go func() {
		// Iterate over [2,4,6-9][6-9]*, 3[4,6-9]*

		s := make([]string, size)
		for _, n0 := range [7]string{"2", "3", "4", "6", "7", "8", "9"} {
			s[0] = n0
		}

		close(c)
	}()

	return
}

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s number [number]\n", os.Args[0])
		return
	}

	for _, nbStr := range os.Args[1:] {
		fmt.Printf("-> %s\n", nbStr)
		p, err := persistence(nbStr, 0)
		if err != nil {
			fmt.Printf("%s raised an error: %s\n", nbStr, err)
		} else {
			fmt.Printf("%s => %d\n", nbStr, p)
		}
	}
}
