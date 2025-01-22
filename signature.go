package pdfx

import (
	"errors"
	"regexp"
	"strings"
	"unicode"

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

	var newFieldsArr types.Array

	// each field is a dictionary
	for _, fieldObj := range fieldsArr {
		fieldDict, err := p.pdfContext.DereferenceDict(fieldObj)
		if err != nil {
			return errors.New("can't dereference field dictionary")
		}

		// check if field has a signature
		if _, found := fieldDict.Find("V"); found {

			err = p.pdfContext.DeleteObject(fieldObj)
			if err != nil {
				return errors.New("can't delete field object")
			}

			err = p.pdfContext.DeleteObject(fieldObj)
			if err != nil {
				return errors.New("can't delete field object")
			}

		} else {
			newFieldsArr = append(newFieldsArr, fieldObj)
		}
	}

	// Update the Fields array in the AcroForm dictionary
	acroFormDict.Update("Fields", newFieldsArr)

	// Remove SigFlags if present
	if _, found := acroFormDict.Find("SigFlags"); found {
		acroFormDict.Delete("SigFlags")
	}

	err = p.Optimize()
	if err != nil {
		return err
	}
	return nil
}

var reSignatureAssistPriv = []*regexp.Regexp{
	regexp.MustCompile(`(?mi)(.*?)([0-9a-fA-F]{8}\b-[0-9a-fA-F]{4}\b-[0-9a-fA-F]{4}\b-[0-9a-fA-F]{4}\b-[0-9a-fA-F]{12})_(:?qrcode|assistpriv)`),
	regexp.MustCompile(`(?mi)(.*?)(:?qrcode|assistpriv)`),
}

func isAssistPrivy(str string) bool {
	str = strings.Map(func(r rune) rune {
		if r > unicode.MaxASCII || (r < 32 && r != '\t' && r != '\n' && r != '\r') {
			return -1
		}
		if unicode.IsPrint(r) {
			return r
		}
		return -1
	}, str)
	for _, re := range reSignatureAssistPriv {
		if matches := re.FindAllString(str, -1); len(matches) != 0 {
			return true
		}
	}
	return false
}
