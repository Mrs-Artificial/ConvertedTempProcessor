package main

import (
	"ConvertedTempProcessor/pkg" // Import the temp package
	"fmt"
)

func main() {
	processor := temp.NewXpertTempCheckProcessor()
	processor.Start(false)
	fmt.Println("test")
}
