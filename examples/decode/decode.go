package main

import (
	"fmt"

	bmip10 "github.com/jwetzell/bmip10-go"
)

func main() {

	decoder := bmip10.NewDecoder(10, 8)

	encodedData := []int32{128, 1, 2, 3, 5, 8, 13, 21, 34, 55, 89, 116, 143, 163, 196, 250}

	fmt.Printf("index\tenc\tdec\n")

	for i := 0; i < len(encodedData); i++ {
		decoded := decoder.Decode(encodedData[i])
		fmt.Printf("%d\t%d\t%d\n", i, encodedData[i], decoded)
	}
}
