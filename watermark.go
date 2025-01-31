package pdfx

import (
	"log"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/types"
)

// removeWatermarks is a function to remove watermarks from a PDF file
func (p *PDFProcessor) removeWatermarks() error {
	if err := pdfcpu.DetectWatermarks(p.pdfContext); err != nil {
		log.Println("No watermarks found")
	}

	if !p.pdfContext.Watermarked {
		return nil
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

// removeQRCode is a function to remove QR codes from a PDF file
func (p *PDFProcessor) removeQRCode() error {
	rootDict, err := p.pdfContext.Catalog()
	if err != nil {
		return err
	}

	pagesObj, found := rootDict.Find("Pages")
	if !found {
		return err
	}

	pagesDict, _ := p.pdfContext.DereferenceDict(pagesObj)
	pagesKidsObj, _ := pagesDict.Find("Kids")
	pagesKidsArray, _ := p.pdfContext.DereferenceArray(pagesKidsObj)

	if pagesKidsArray == nil {
		return nil
	}

	for _, kidObj := range pagesKidsArray {
		kidRef, ok := kidObj.(types.IndirectRef)
		if !ok {
			log.Println("page kid is not an indirect reference")
			continue
		}

		kidDict, err := p.pdfContext.DereferenceDict(kidRef)
		if err != nil {
			log.Println("can't dereference page kid")
			continue
		}

		resourceObj, found := kidDict.Find("Resources")
		if !found {
			continue
		}

		resourceDict, err := p.pdfContext.DereferenceDict(resourceObj)
		if err != nil {
			log.Println("can't dereference resource dictionary")
			continue
		}

		// fetch XObject
		xObjectObj, found := resourceDict.Find("XObject")
		if !found {
			continue
		}

		xObjectDict, err := p.pdfContext.DereferenceDict(xObjectObj)
		if err != nil {
			log.Println("can't dereference xObject dictionary")
			continue
		}

		// delete X0 object if found
		if _, found := xObjectDict.Find("X0"); found {
			log.Println("found X0 object")

			err = p.pdfContext.DeleteDictEntry(xObjectDict, "X0")
			if err != nil {
				log.Println("can't delete X0 object")
			}

			continue
		}

		// delete X1 object if found
		if _, found := xObjectDict.Find("X1"); found {
			log.Println("found X1 object")
			err := p.pdfContext.DeleteDictEntry(xObjectDict, "X1")
			if err != nil {
				log.Println("can't delete X1 object")
			}

			continue
		}

		// delete X3 object if found
		if _, found := xObjectDict.Find("X3"); found {
			log.Println("found X3 object")
			err := p.pdfContext.DeleteDictEntry(xObjectDict, "X3")
			if err != nil {
				log.Println("can't delete X3 object")
			}

			continue
		}

		// delete R19 object if found
		if _, found := xObjectDict.Find("R19"); found {
			log.Println("found R19 object")
			err := p.pdfContext.DeleteDictEntry(xObjectDict, "R19")
			if err != nil {
				log.Println("can't delete R19 object")
			}

			continue
		}

		// delete I1 object if found
		if _, found := xObjectDict.Find("I1"); found {
			log.Println("found I1 object")
			err := p.pdfContext.DeleteDictEntry(xObjectDict, "I1")
			if err != nil {
				log.Println("can't delete I1 object")
			}

			continue
		}

		// delete XO1 object if found
		if _, found := xObjectDict.Find("XO1"); found {
			log.Println("found XO1 object")
			err := p.pdfContext.DeleteDictEntry(xObjectDict, "XO1")
			if err != nil {
				log.Println("can't delete XO1 object")
			}

			continue
		}

		// delete XO2 object if found
		if _, found := xObjectDict.Find("XO2"); found {
			log.Println("found XO2 object")
			err := p.pdfContext.DeleteDictEntry(xObjectDict, "XO2")
			if err != nil {
				log.Println("can't delete XO2 object")
			}

			continue
		}

		// delete XO3 object if found
		if _, found := xObjectDict.Find("XO3"); found {
			log.Println("found XO3 object")
			err := p.pdfContext.DeleteDictEntry(xObjectDict, "XO3")
			if err != nil {
				log.Println("can't delete XO3 object")
			}

			continue
		}

		// delete R72 object if found
		if _, found := xObjectDict.Find("R72"); found {
			log.Println("found R72 object")
			err := p.pdfContext.DeleteDictEntry(xObjectDict, "R72")
			if err != nil {
				log.Println("can't delete R72 object")
			}

			continue
		}

		// delete Fm0 object if found
		if _, found := xObjectDict.Find("Fm0"); found {
			log.Println("found Fm0 object")
			err := p.pdfContext.DeleteDictEntry(xObjectDict, "Fm0")
			if err != nil {
				log.Println("can't delete Fm0 object")
			}

			continue
		}

		// delete X5 object if found
		if _, found := xObjectDict.Find("X5"); found {
			log.Println("found X5 object")
			err := p.pdfContext.DeleteDictEntry(xObjectDict, "X5")
			if err != nil {
				log.Println("can't delete X5 object")
			}

			continue
		}
	}

	err = p.Optimize()
	if err != nil {
		return err
	}

	return nil
}
