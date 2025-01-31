package pdfx

import (
	"errors"
	"log"

	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/types"
)

func (p *PDFProcessor) removeSignatures() error {
	rootDict, err := p.pdfContext.Catalog()
	if err != nil {
		return err
	}

	acroFormObj, ok := rootDict.Find("AcroForm")
	if !ok {
		return errors.New("acroform dictionary not found")
	}

	acroFormDict, err := p.pdfContext.DereferenceDict(acroFormObj)
	if err != nil {
		return err
	}

	fields, found := acroFormDict.Find("Fields")

	if !found {
		return errors.New("fields not found in acroform dictionary")
	}

	fieldsArr, err := p.pdfContext.DereferenceArray(fields)
	if err != nil {
		return errors.New("can't dereference fields array")
	}

	// each field is a dictionary
	for _, fieldObj := range fieldsArr {
		annotationRef, ok := fieldObj.(types.IndirectRef)
		if !ok {
			log.Println("field is not an indirect reference")
		}

		annotationDict, err := p.pdfContext.DereferenceDict(annotationRef)
		if err != nil {
			return err
		}

		// Remove the annotation dictionary if it is a signature
		if v, found := annotationDict.Find("V"); found {
			_, err := p.pdfContext.DereferenceDict(v)
			if err != nil {
				return errors.New("can't dereference field value")
			}

			err = p.pdfContext.DeleteObjectGraph(v)
			if err != nil {
				return errors.New("can't delete field value")
			}

			err = p.pdfContext.DeleteObjectGraph(annotationRef)
			if err != nil {
				return errors.New("can't delete field object")
			}

			err = p.pdfContext.DeleteObjectGraph(fieldObj)
			if err != nil {
				return errors.New("can't delete field object")
			}

		}

	}

	// Update the Fields array in the AcroForm dictionary
	acroFormDict.Delete("Fields")

	// Remove SigFlags if present
	if _, found := acroFormDict.Find("SigFlags"); found {
		acroFormDict.Delete("SigFlags")
	}

	return nil
}
