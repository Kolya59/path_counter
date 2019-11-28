package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
)

type Field struct {
	IsAllowed bool
	PathCount int64
}

func ParseFile(filename string) (n int64, k int64, fields [][]*Field, err error) {
	// Read file
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return 0, 0, [][]*Field{}, err
	}

	// Parse data
	splits := strings.Split(string(file), "\n")

	// Parse header
	header := strings.Split(splits[0], " ")
	if len(header) < 2 {
		return 0, 0, [][]*Field{}, errors.New("invalid header")
	}
	n, err = strconv.ParseInt(header[0], 10, 64)
	if err != nil {
		return 0, 0, [][]*Field{}, err
	}
	k, err = strconv.ParseInt(header[1], 10, 64)
	if err != nil {
		return 0, 0, [][]*Field{}, err
	}

	fields = make([][]*Field, n)
	for i, split := range splits[1:] {
		s := strings.Split(split, "")
		for v := range s {
			value, err := strconv.ParseInt(s[v], 10, 64)
			if err != nil {
				return 0, 0, nil, err
			}
			fields[i] = append(fields[i], &Field{
				IsAllowed: value == 0,
				PathCount: 0,
			})
		}
	}

	return
}

func CalculatePathCount(x int64, y int64, fields [][]*Field, n int64, k int64) int64 {
	current := fields[x][y]
	current.PathCount++
	if current.PathCount >= k {
		return current.PathCount
	}

	switch true {
	case x < n - 1 && fields[x+1][y].IsAllowed:
		CalculatePathCount(x+1, y, fields, n, k)
	case y < n - 1 && fields[x][y+1].IsAllowed:
		CalculatePathCount(x, y+1, fields, n, k)
	case x > 0 && fields[x-1][y].IsAllowed:
		CalculatePathCount(x-1, y, fields, n, k)
	case y > 0 && fields[x][y-1].IsAllowed:
		CalculatePathCount(x, y-1, fields, n, k)
	}
	return current.PathCount
}

func main() {
	n, k, fields, err := ParseFile("input")
	if err != nil {
		log.Fatalf("Failed to parse file: %v", err)
	}

	count := CalculatePathCount(0, 0, fields, n, k)

	printedField := ""
	printedPathField := ""
	for i := int64(0); i < n; i++ {
		for j := int64(0); j < n; j++ {
			printedField += fmt.Sprintf("%v ", fields[i][j].IsAllowed)
			printedPathField += fmt.Sprintf("%v ", fields[i][j].PathCount)
		}
		printedField += "\n"
		printedPathField += "\n"
	}
	log.Printf("\n%s\n%s\n%v\n", printedField, printedPathField, count)

}
