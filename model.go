package main

import (
	// remember to add encoding/xml to your list of imports
	"encoding/xml"
)

type QvdFieldHeader struct {
	XMLName     xml.Name `xml:"QvdFieldHeader"`
	Name        string   `xml:"FieldName"`
	BitOffset   int      `xml:"BitOffset"`
	BitWidth    int      `xml:"BitWidth"`
	Offset      int      `xml:"Offset"`
	Length      int      `xml:"Length"`
	Comment     string   `xml:"Comment"`
	SymbolCount int      `xml:"NoOfSymbols"`
	Symbols     []string
}

type QvdHeader struct {
	XMLName xml.Name         `xml:"QvdTableHeader"`
	Name    string           `xml:"TableName"`
	Comment string           `xml:"Comment"`
	Length  int              `xml:"Length"`
	Fields  []QvdFieldHeader `xml:"Fields>QvdFieldHeader"`
	Offset  int              `xml:"Offset"`
}
