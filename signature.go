package pdfx

import (
	"errors"
	"log"

	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/form"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/types"
)

// removeSignatures removes all signature fields from the PDF.
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

	objIDs := make([]string, 0)

	// each field is a dictionary
	for _, fieldObj := range fieldsArr {
		log.Print("fieldObj: ", fieldObj)
		fieldRef := fieldObj.(types.IndirectRef)

		fieldDict, err := p.pdfContext.DereferenceDict(fieldRef)
		if err != nil {
			return errors.New("can't dereference field dictionary")
		}

		// Check if the field is a signature field
		if _, found := fieldDict.Find("V"); !found {
			continue
		}

		objIDs = append(objIDs, fieldRef.ObjectNumber.String())
	}

	// get perms if sealed
	permsObj, found := rootDict.Find("Perms")
	if found {
		permsDict, err := p.pdfContext.DereferenceDict(permsObj)
		if err != nil {
			return errors.New("can't dereference perms dictionary")

		}

		docMdpObj, foundMDP := permsDict.Find("DocMDP")
		if foundMDP {
			docMdpRef := docMdpObj.(types.IndirectRef)
			docMdpDict, err := p.pdfContext.DereferenceDict(docMdpObj)
			if err != nil {
				return errors.New("can't dereference docmdp dictionary")
			}

			// Check if the field type is signature
			if docMDPType, found := docMdpDict.Find("Type"); found {
				if docMDPType.String() == "Sig" {
					// remove the object
					err := p.pdfContext.DeleteObject(docMdpRef)
					if err != nil {
						return errors.New("failed to remove object")
					}

					err = p.pdfContext.DeleteObject(docMdpObj)
					if err != nil {
						return errors.New("failed to remove object")
					}

					err = p.pdfContext.DeleteDictEntry(permsDict, "DocMDP")
					if err != nil {
						return errors.New("failed to remove object")
					}

				}
			}
		}
	}

	// Remove the signature fields
	ok, err = form.RemoveFormFields(p.pdfContext, objIDs)
	if err != nil {
		return errors.New("failed to remove signature fields")
	}

	return nil
}
