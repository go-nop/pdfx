package main

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-nop/pdfx"
)

func main() {
	inputDir := "testdata"
	outputDir := "result"

	// Create the output directory if it doesn't exist
	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}

	// List all PDF files in the input directory
	files, err := os.ReadDir(inputDir)
	if err != nil {
		log.Fatalf("Failed to read input directory: %v", err)
	}

	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".pdf") {
			inputPath := filepath.Join(inputDir, file.Name())
			outputPath := filepath.Join(outputDir, file.Name())

			ctx := context.Background()

			log.Printf("Processing %s...", inputPath)

			// Create a new PDFProcessor
			processor, err := pdfx.New(ctx, inputPath, outputPath, pdfx.WithPassword("momoka"))
			if err != nil {
				log.Printf("Failed to create PDFProcessor for %s: %v", inputPath, err)
				continue
			}

			// save debug info to txt
			debugStr := processor.Debug()
			debugPath := strings.TrimSuffix(inputPath, ".pdf") + ".txt"
			if err := os.WriteFile(debugPath, []byte(debugStr), os.ModePerm); err != nil {
				log.Printf("Failed to write debug info to %s: %v", debugPath, err)
			}

			// Remove watermarks
			if err := processor.RemoveWatermarks(); err != nil {
				log.Printf("Failed to remove watermarks from %s: %v", inputPath, err)

			}

			// Remove signatures
			if err := processor.RemoveSignatures(); err != nil {
				log.Printf("Failed to remove signatures from %s: %v", inputPath, err)

			}

			// Remove QR codes
			if err := processor.RemoveQRCode(); err != nil {
				log.Printf("Failed to remove QR codes from %s: %v", inputPath, err)

			}

			// save debug info to txt
			debugStr = processor.Debug()
			debugPath = strings.TrimSuffix(outputPath, ".pdf") + ".txt"
			if err := os.WriteFile(debugPath, []byte(debugStr), os.ModePerm); err != nil {
				log.Printf("Failed to write debug info to %s: %v", debugPath, err)
			}

			// Write the output file
			if err := processor.WriteFile(); err != nil {
				log.Fatalf("Failed to write file %s: %v", outputPath, err)
			}

			log.Printf("Processed %s successfully", inputPath)
		}
	}
}
