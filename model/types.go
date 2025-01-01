package model

import (
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"

	"github.com/go-nop/pdfx/utils"
)

// PDFObject is the interface that represents a PDF object.
type PDFObject interface {
	// String returns the string representation of the PDF object.
	String() string
	// WriteString writes the string representation of the PDF object to the writer.
	WriteString() string
}

////////////////////////
// ObjNull
////////////////////////

// ObjNull represents a null object.
type ObjNull struct{}

var _ PDFObject = ObjNull{}

// MakeNull returns a new null object.
func MakeNull() ObjNull {
	return ObjNull{}
}

// String returns the string representation of the PDF object.
func (o ObjNull) String() string {
	return "null"
}

// WriteString writes the string representation of the PDF object to the writer.
func (o ObjNull) WriteString() string {
	return "null"
}

////////////////////////
// ObjBool
////////////////////////

// ObjBool represents a boolean object.
type ObjBool bool

var _ PDFObject = ObjBool(false)

// MakeBool returns a new boolean object.
func MakeBool(b bool) ObjBool {
	return ObjBool(b)
}

// String returns the string representation of the PDF object.
func (o ObjBool) String() string {
	if o {
		return "true"
	}
	return "false"
}

// WriteString writes the string representation of the PDF object to the writer.
func (o ObjBool) WriteString() string {
	return o.String()
}

////////////////////////
// ObjInt
////////////////////////

// ObjInt represents an integer object.
type ObjInt int64

var _ PDFObject = ObjInt(0)

// MakeInt returns a new integer object.
func MakeInt(i int64) ObjInt {
	return ObjInt(i)
}

// String returns the string representation of the PDF object.
func (o ObjInt) String() string {
	return fmt.Sprintf("%d", o)
}

// WriteString writes the string representation of the PDF object to the writer.
func (o ObjInt) WriteString() string {
	return o.String()
}

////////////////////////
// ObjFloat
////////////////////////

// ObjFloat represents a float object.
type ObjFloat float64

var _ PDFObject = ObjFloat(0)

// MakeFloat returns a new float object.
func MakeFloat(f float64) ObjFloat {
	return ObjFloat(f)
}

// String returns the string representation of the PDF object.
func (o ObjFloat) String() string {
	return fmt.Sprintf("%f", o)
}

// WriteString writes the string representation of the PDF object to the writer.
func (o ObjFloat) WriteString() string {
	return o.String()
}

////////////////////////
// ObjString
////////////////////////

// ObjString represents a string object.
type ObjString struct {
	value string
	isHex bool
}

var _ PDFObject = ObjString{}

// MakeString returns a new string object.
func MakeString(value string) ObjString {
	return ObjString{value: value}
}

// MakeHexString returns a new hex string object.
func MakeHexString(value string) ObjString {
	return ObjString{value: value, isHex: true}
}

// String returns the string representation of the PDF object.
func (o ObjString) String() string {
	return o.value
}

// WriteString writes the string representation of the PDF object to the writer.
func (o ObjString) WriteString() string {
	var sb strings.Builder
	if o.isHex {
		sHex := hex.EncodeToString([]byte(o.value))
		sb.WriteString("<")
		sb.WriteString(sHex)
		sb.WriteString(">")
		return sb.String()
	}

	// Otherwise regular string.
	escapeSequences := map[byte]string{
		'\n': "\\n",
		'\r': "\\r",
		'\t': "\\t",
		'\b': "\\b",
		'\f': "\\f",
		'(':  "\\(",
		')':  "\\)",
		'\\': "\\\\",
	}

	sb.WriteString("(")
	for i := 0; i < len(o.value); i++ {
		b := o.value[i]
		if es, found := escapeSequences[b]; found {
			sb.WriteString(es)
		} else {
			sb.WriteByte(b)
		}
	}
	sb.WriteString(")")

	return sb.String()
}

////////////////////////
// ObjName
////////////////////////

// ObjName represents a name object.
type ObjName string

var _ PDFObject = ObjName("")

// MakeName returns a new name object.
func MakeName(value string) ObjName {
	return ObjName(value)
}

// String returns the string representation of the PDF object.
func (o ObjName) String() string {
	return "/" + string(o)
}

// WriteString writes the string representation of the PDF object to the writer.
func (o ObjName) WriteString() string {
	var sb strings.Builder

	sb.WriteString("/")
	for i := 0; i < len(o); i++ {
		ch := (o)[i]
		if !utils.IsPrintable(ch) || ch == '#' || utils.IsDelimiter(ch) {
			sb.WriteString(fmt.Sprintf("#%.2x", ch))
		} else {
			sb.WriteByte(ch)
		}
	}

	return sb.String()
}

////////////////////////
// ObjArray
////////////////////////

// ObjArray represents an array object.
type ObjArray struct {
	values []PDFObject
}

var _ PDFObject = ObjArray{}

// MakeArray returns a new array object.
func MakeArray(values ...PDFObject) ObjArray {
	return ObjArray{values: values}
}

// String returns the string representation of the PDF object.
func (o ObjArray) String() string {
	var sb strings.Builder

	sb.WriteString("[")
	for i, v := range o.values {
		if i > 0 {
			sb.WriteString(" ")
		}
		sb.WriteString(v.String())
	}
	sb.WriteString("]")

	return sb.String()
}

// WriteString writes the string representation of the PDF object to the writer.
func (o ObjArray) WriteString() string {
	var sb strings.Builder

	sb.WriteString("[")
	for i, v := range o.values {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(v.WriteString())
	}
	sb.WriteString("]")

	return sb.String()
}

////////////////////////
// ObjDict
////////////////////////

// ObjDict represents a dictionary object.
type ObjDict struct {
	values map[ObjName]PDFObject
	keys   []ObjName
}

var _ PDFObject = &ObjDict{}

// MakeDict returns a new dictionary object.
func MakeDict() *ObjDict {
	d := &ObjDict{}
	d.values = map[ObjName]PDFObject{}
	d.keys = []ObjName{}
	return d
}

// String returns the string representation of the PDF object.
func (o ObjDict) String() string {
	var sb strings.Builder

	sb.WriteString("Dict(")
	for _, k := range o.keys {
		v := o.values[k]
		sb.WriteString(`"` + k.String() + `": `)
		sb.WriteString(v.String())
		sb.WriteString(`, `)
	}
	sb.WriteString(")")
	return sb.String()
}

// WriteString writes the string representation of the PDF object to the writer.
func (o ObjDict) WriteString() string {
	var sb strings.Builder

	sb.WriteString("<<")
	for _, k := range o.keys {
		v := o.values[k]
		sb.WriteString(k.WriteString())
		sb.WriteString(" ")
		sb.WriteString(v.WriteString())
	}

	sb.WriteString(">>")
	return sb.String()
}

// Set sets the value of the key in the dictionary.
func (o *ObjDict) Set(key ObjName, value PDFObject) {
	if _, found := o.values[key]; !found {
		o.keys = append(o.keys, key)
	}
	o.values[key] = value
}

////////////////////////
// ObjRef
////////////////////////

// ObjRef represents a reference object.
type ObjRef struct {
	ObjNumber        int64
	GenerationNumber int64
}

var _ PDFObject = ObjRef{}

// MakeRef returns a new reference object.
func MakeRef(objNumber, genNumber int64) ObjRef {
	return ObjRef{ObjNumber: objNumber, GenerationNumber: genNumber}
}

// String returns the string representation of the PDF object.
func (o ObjRef) String() string {
	return fmt.Sprintf("Ref(%d %d)", o.ObjNumber, o.GenerationNumber)
}

// WriteString writes the string representation of the PDF object to the writer.
func (o ObjRef) WriteString() string {
	var sb strings.Builder
	sb.WriteString(strconv.FormatInt(o.ObjNumber, 10))
	sb.WriteString(" ")
	sb.WriteString(strconv.FormatInt(o.GenerationNumber, 10))
	sb.WriteString(" R")
	return sb.String()
}

////////////////////////
// ObjStream
////////////////////////

// ObjStream represents a stream object.
type ObjStream struct {
	*ObjDict
	ObjRef
	Stream []byte
}

var _ PDFObject = ObjStream{}

// MakeStream returns a new stream object.
func MakeStream(contents []byte) (*ObjStream, error) {
	stream := &ObjStream{}

	stream.ObjDict.Set("Length", MakeInt(int64(len(contents))))
	stream.Stream = contents

	return stream, nil
}

// String returns the string representation of the PDF object.
func (o ObjStream) String() string {
	return fmt.Sprintf("Stream(%s)", string(o.Stream))
}

// WriteString writes the string representation of the PDF object to the writer.
func (o ObjStream) WriteString() string {
	var sb strings.Builder
	sb.WriteString(strconv.FormatInt(o.ObjNumber, 10))
	sb.WriteString(" 0 R")
	return sb.String()
}

// ObjStreams represents a collection of stream objects.
type ObjStreams struct {
	values []PDFObject
	ObjRef
}

var _ PDFObject = ObjStreams{}

// MakeStreams returns a new stream object.
func MakeStreams(objs ...PDFObject) *ObjStreams {
	streams := &ObjStreams{}
	streams.values = []PDFObject{}
	streams.values = append(streams.values, objs...)

	return streams
}

// String returns the string representation of the PDF object.
func (o ObjStreams) String() string {
	return fmt.Sprintf("Streams(%d)", len(o.values))
}

// WriteString writes the string representation of the PDF object to the writer.
func (o ObjStreams) WriteString() string {
	var sb strings.Builder
	sb.WriteString(strconv.FormatInt(o.ObjNumber, 10))
	sb.WriteString(" 0 R")
	return sb.String()
}

// ObjIndirect represents an indirect object.
type ObjIndirect struct {
	PDFObject
	ObjRef
}

var _ PDFObject = ObjIndirect{}

// MakeIndirect returns a new indirect object.
func MakeIndirect(obj PDFObject) *ObjIndirect {
	ind := ObjIndirect{}
	ind.PDFObject = obj
	return &ind
}

// String returns the string representation of the PDF object.
func (o ObjIndirect) String() string {
	return fmt.Sprintf("Indirect(%s)", o.PDFObject.String())
}

// WriteString writes the string representation of the PDF object to the writer.
func (o ObjIndirect) WriteString() string {
	var sb strings.Builder
	sb.WriteString(strconv.FormatInt(o.ObjNumber, 10))
	sb.WriteString(" 0 R")
	return sb.String()
}
