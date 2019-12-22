package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
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

func CalculatePathCount(fields [][]Field, n int64, k int64) int64 {
	wrappedFields := make([][][]Field, k+1)
	for step := int64(0); step < k+1; step++ {
		wrappedFields[step] = make([][]Field, n+2)
		for i := int64(0); i < n+2; i++ {
			wrappedFields[step][i] = make([]Field, n+2)
		}
	}
	for step := int64(0); step < k+1; step++ {
		for i := int64(0); i < n+2; i++ {
			wrappedFields[step][0][i] = Field{}
			wrappedFields[step][i][0] = Field{}
			wrappedFields[step][n+1][i] = Field{}
			wrappedFields[step][i][n+1] = Field{}
		}
		for i := int64(1); i < n+1; i++ {
			for j := int64(1); j < n+1; j++ {
				wrappedFields[step][i][j] = fields[i-1][j-1]
			}
		}
	}
	wrappedFields[0][1][1].PathCount = 1

	for step := int64(1); step < k+1; step++ {
		for i := int64(1); i < n+1; i++ {
			for j := int64(1); j < n+1; j++ {
				if wrappedFields[step][i][j].IsAllowed {
					wrappedFields[step][i][j].PathCount +=
						wrappedFields[step-1][i-1][j].PathCount +
							wrappedFields[step-1][i+1][j].PathCount +
							wrappedFields[step-1][i][j+1].PathCount +
							wrappedFields[step-1][i][j-1].PathCount
				}
			}
		}
		printedPathField := fmt.Sprintf("Step: %v\n", step)
		for i := int64(0); i < n+2; i++ {
			for j := int64(0); j < n+2; j++ {
				printedPathField += fmt.Sprintf("%v ", wrappedFields[step][i][j].PathCount)
			}
			printedPathField += "\n"
		}
		log.Printf("%s\n", printedPathField)

	}

	return wrappedFields[k][n][n].PathCount
}

func main() {
	n, k, fields, err := ParseFile(filepath.Base("./input"))
	if err != nil {
		log.Fatalf("Failed to parse file: %v", err)
	}
	r := CalculatePathCount(fields, n, k)
	log.Printf("Result is: %v", r)
}
