package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"strconv"
	"strings"
)

type Field struct {
	IsAllowed bool
	PathCount int64
}

func ParseFile(filename string) (n int64, k int64, fields [][]Field, err error) {
	// Read file
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return 0, 0, [][]Field{}, err
	}
	// Parse data
	splits := strings.Split(string(file), "\n")
	// Parse header
	header := strings.Split(splits[0], " ")
	if len(header) < 2 {
		return 0, 0, [][]Field{}, errors.New("invalid header")
	}
	n, err = strconv.ParseInt(header[0], 10, 64)
	if err != nil {
		return 0, 0, [][]Field{}, err
	}
	k, err = strconv.ParseInt(header[1], 10, 64)
	if err != nil {
		return 0, 0, [][]Field{}, err
	}
	fields = make([][]Field, n)
	for i, split := range splits[1:] {
		s := strings.Split(split, "")
		for v := range s {
			value, err := strconv.ParseInt(s[v], 10, 64)
			if err != nil {
				return 0, 0, nil, err
			}
			fields[i] = append(fields[i], Field{
				IsAllowed: value == 0,
				PathCount: 0,
			})
		}
	}
	return
}

func CalculatePathCount(fields [][]Field, n int64, k int64) [][]Field {
	odd := make([][]Field, n+2)
	even := make([][]Field, n+2)
	for i := int64(0); i < n+2; i++ {
		odd[i] = make([]Field, n+2)
		even[i] = make([]Field, n+2)
	}
	for i := int64(0); i < n+2; i++ {
		odd[0][i] = Field{}
		odd[i][0] = Field{}
		odd[n+1][i] = Field{}
		odd[i][n+1] = Field{}
		even[0][i] = Field{}
		even[i][0] = Field{}
		even[n+1][i] = Field{}
		even[i][n+1] = Field{}
	}
	for i := int64(1); i < n+1; i++ {
		for j := int64(1); j < n+1; j++ {
			odd[i][j] = fields[i-1][j-1]
			even[i][j] = fields[i-1][j-1]
		}
	}
	odd[1][1].PathCount = 1
	even[1][1].PathCount = 1

	var isEven bool
	for step := int64(0); step < k; step++ {
		isEven = math.Mod(float64(step), 2) == 0
		for i := int64(1); i < n+1; i++ {
			for j := int64(1); j < n+1; j++ {
				if odd[i][j].IsAllowed {
					if isEven {
						even[i][j].PathCount += odd[i-1][j].PathCount +
							odd[i+1][j].PathCount +
							odd[i][j+1].PathCount +
							odd[i][j-1].PathCount
						continue
					}
					odd[i][j].PathCount += even[i-1][j].PathCount +
						even[i+1][j].PathCount +
						even[i][j+1].PathCount +
						even[i][j-1].PathCount
				}
			}
		}
	}

	if isEven {
		return even
	}
	return odd
}

func main() {
	n, k, fields, err := ParseFile("input")
	if err != nil {
		log.Fatalf("Failed to parse file: %v", err)
	}
	r := CalculatePathCount(fields, n, k)

	printedPathField := ""
	for i := int64(0); i < n+2; i++ {
		for j := int64(0); j < n+2; j++ {
			printedPathField += fmt.Sprintf("%v ", r[i][j].PathCount)
		}
		printedPathField += "\n"
	}
	log.Printf("\n%s\nResult is: %v", printedPathField, r[n][n].PathCount)
}
