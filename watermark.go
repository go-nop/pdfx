package pdfx

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
)

// removeWatermarks is a function to remove watermarks from a PDF file
func (p *PDFProcessor) removeWatermarks() error {
	if err := pdfcpu.DetectWatermarks(p.pdfContext); err != nil {
		log.Println("No watermarks found")
	}

	if !p.pdfContext.Watermarked {
		lastPage := fmt.Sprintf("%d", p.pdfContext.PageCount)
		pages, err := api.PagesForPageSelection(p.pdfContext.PageCount, []string{lastPage}, true, true)
		if err != nil {
			return err
		}

		images, _, err := pdfcpu.Images(p.pdfContext, pages)
		if err != nil {
			return err
		}

		for _, img := range images {
			for _, i := range img {
				if i.Name == "I1" || i.Name == "XO1" || i.Name == "XO2" || i.Name == "XO3" || i.Name == "R72" {
					log.Println("Watermark found")

					// create image wihite with same dimensions
					rect := image.Rect(0, 0, i.Width, i.Height)
					white := color.RGBA{R: 255, G: 255, B: 255, A: 255}
					img := image.NewRGBA(rect)

					for y := 0; y < i.Height; y++ {
						for x := 0; x < i.Width; x++ {
							img.Set(x, y, white)
						}
					}

					var buf bytes.Buffer

					err := png.Encode(&buf, img)
					if err != nil {
						panic(err)
					}

					reader := bytes.NewReader(buf.Bytes())

					err = pdfcpu.UpdateImagesByObjNr(p.pdfContext, reader, i.ObjNr)
					if err != nil {
						return err
					}
					break
				}
			}
		}
	}

	pages, err := api.PagesForPageSelection(p.pdfContext.PageCount, nil, true, true)
	if err != nil {
		return err
	}

	if err = pdfcpu.RemoveWatermarks(p.pdfContext, pages); err != nil {
		return err
	}

	err = p.Optimize()
	if err != nil {
		return err
	}

	return nil

}
