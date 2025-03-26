package main

import (
	"airista.com/ConvertedTempProcessor/pkg" // Import the pkg package
)

func main() {
	// Initialize the processor
	processor := pkg.NewXpertTempCheckProcessor()

	// Start the processor
	processor.Start(false)

	// Keep the main function running
	select {}
}
