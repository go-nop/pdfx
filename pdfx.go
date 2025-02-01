package pdfx

import (
	"context"
	"io"
	"log"
	"os"
	"sync"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
)

// PDFProcessor is a struct for processing PDF files
type PDFProcessor struct {
	// ctx is the context for the PDFProcessor
	// It is used to cancel the processing of the PDF file
	// It is also used to pass the context
	// to the underlying PDF processing functions
	ctx context.Context

	rs io.ReadSeeker

	inputFilePath  string
	outputFilePath string
	configuration  *model.Configuration
	pdfContext     *model.Context
}

var (
	once          sync.Once
	defaultConfig *model.Configuration
)

// defaultConfiguration returns the default configuration for PDFProcessor
func defaultConfiguration() *model.Configuration {
	once.Do(func() {
		defaultConfig = model.NewDefaultConfiguration()
	})
	return defaultConfig
}

// New creates a new PDFProcessor
func New(ctx context.Context, inputPath, outputPath string, opts ...Option) (*PDFProcessor, error) {
	conf := defaultConfiguration()
	p := &PDFProcessor{
		ctx:            ctx,
		inputFilePath:  inputPath,
		outputFilePath: outputPath,
		configuration:  conf,
	}

	// Apply options
	for _, opt := range opts {
		opt(p)
	}
	// Create a new PDFProcessor
	rs, err := os.Open(inputPath)
	if err != nil {
		return nil, err
	}
	defer rs.Close()

	p.rs = rs

	pdfCtx, err := api.ReadValidateAndOptimize(rs, conf)

	p.pdfContext = pdfCtx

	return p, nil
}

// WriteFile is a function to write the PDFProcessor's PDFContext to a file
func (p *PDFProcessor) WriteFile() error {
	err := pdfcpu.OptimizeXRefTable(p.pdfContext)
	if err != nil {
		log.Printf("Failed to optimize PDF: %v", err)
	}

	return api.WriteContextFile(p.pdfContext, p.outputFilePath)
}

// Debug is a function to print the PDFProcessor's PDFContext
func (p *PDFProcessor) Debug() string {
	return p.pdfContext.String()
}

// Optimize is a function to optimize a PDF file
func (p *PDFProcessor) Optimize() error {
	return api.OptimizeContext(p.pdfContext)
}

// RemoveSignatures is a function to remove signatures from a PDF file
func (p *PDFProcessor) RemoveSignatures() error {
	return p.removeSignatures()
}

// RemoveWatermarks is a function to remove watermarks from a PDF file
func (p *PDFProcessor) RemoveWatermarks() error {
	return p.removeWatermarks()
}

// RemoveQRCode is a function to remove QR codes from a PDF file
func (p *PDFProcessor) RemoveQRCode() error {
	return p.removeQRCode()
}
