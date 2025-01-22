package main

import (
	"context"
	"log"
	"os"

	"github.com/go-nop/pdfx"
)

func main() {
	if len(os.Args) < 3 {
		log.Fatal("Usage: pdfx <input.pdf> <output.pdf>")
	}

	inputPath := os.Args[1]
	outputPath := os.Args[2]

	ctx := context.Background()

	// Create a new PDFProcessor
	processor, err := pdfx.New(ctx, inputPath, outputPath)
	if err != nil {
		log.Fatalf("Failed to create PDFProcessor: %v", err)
	}

	// Remove watermarks
	err = processor.RemoveWatermarks()
	if err != nil {
		log.Printf("Failed to remove watermarks: %v", err)
	}

	// // Remove signatures
	err = processor.RemoveSignatures()
	if err != nil {
		log.Printf("Failed to remove signatures: %v", err)
	}

	processor.Debug()

	err = processor.WriteFile()
	if err != nil {
		log.Fatalf("Failed to write file: %v", err)
	}

	log.Println("Watermarks removed successfully")
}
