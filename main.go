package main

import (
	"fmt"

	"github.com/gokul1063/uuid-generator/testutils"
)

func main() {
	fmt.Println("=== Running Decode Tests ===")
	testutils.TestDecodeMultiple()

	fmt.Println("=== Running Generator Tests ===")
	testutils.TestGeneratorMultiple()

	fmt.Println("All tests completed. Check ../logs/decode.log and ../logs/generator.log")
}

