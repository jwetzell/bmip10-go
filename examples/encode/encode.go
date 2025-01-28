package main

import (
	"fmt"

	bmip10 "github.com/jwetzell/bmip10-go"
)

func main() {
	var rawData []int32 = []int32{1, 1, 2, 3, 5, 8, 13, 21, 34, 55, 89, 144, 233, 377, 610, 987}
	encoder := bmip10.NewEncoder(10, 8)

	fmt.Printf("index\traw\tenc\n")
	for i := 0; i < len(rawData); i++ {
		encoded := encoder.Encode(rawData[i])
		fmt.Printf("%d\t%d\t%d\n", i, rawData[i], encoded)
	}
}
