package pdfx

import (
	"context"
	"fmt"
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
func New(ctx context.Context, inputPath, outputPath string, opts ...Option) *PDFProcessor {
	// Create a new PDFProcessor
	rs, err := os.Open(inputPath)
	if err != nil {
		log.Fatal(err)
	}
	defer rs.Close()

	conf := defaultConfiguration()

	// read context from the input file
	pdfCtx, err := pdfcpu.ReadWithContext(ctx, rs, conf)
	if err != nil {
		log.Fatal(err)
	}

	err = pdfcpu.OptimizeXRefTable(pdfCtx)
	if err != nil {
		log.Fatal(err)
	}

	p := &PDFProcessor{
		ctx:            ctx,
		rs:             rs,
		inputFilePath:  inputPath,
		outputFilePath: outputPath,
		configuration:  conf,
		pdfContext:     pdfCtx,
	}

	// Apply options
	for _, opt := range opts {
		opt(p)
	}

	return p
}

func (p *PDFProcessor) WriteFile() error {
	f, err := os.Create(p.outputFilePath)
	if err != nil {
		return err
	}
	defer f.Close()

	return api.WriteContextFile(p.pdfContext, p.outputFilePath)
}

func (p *PDFProcessor) Debug() {
	fmt.Print(p.pdfContext.String())
}

func (p *PDFProcessor) Optimize() error {
	return api.OptimizeContext(p.pdfContext)
}

func (p *PDFProcessor) Images() {
	pages, err := api.PagesForPageSelection(p.pdfContext.PageCount, nil, true, true)
	if err != nil {
		log.Fatal(err)
	}

	images, _, err := pdfcpu.Images(p.pdfContext, pages)
	if err != nil {
		log.Fatal(err)
	}

	for _, img := range images {
		log.Println(img)
	}
}
