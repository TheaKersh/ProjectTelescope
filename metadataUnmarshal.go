package main

import "encoding/xml"

type CodeStructure struct {
	XMLName        xml.Name       `xml:"Structure"`
	Header         Header         `xml:"Header"`
	CodeStructures CodeStructures `xml:"Structures"`
}

type CodeStructures struct {
	XMLName   xml.Name  `xml:"Structures"`
	CodeLists CodeLists `xml:"Codelists"`
}
type CodeLists struct {
	XMlName   xml.Name   `xml:"Codelists"`
	CodeLists []CodeList `xml:"Codelist"`
}

type CodeList struct {
	XMlName     xml.Name `xml:"Codelist"`
	Description string   `xml:"Description"`
	Name        string   `xml:"Name"`
	Id          string   `xml:"id,attr"`
	Codes       []Code   `xml:"Code"`
}

type Code struct {
	XmlName    xml.Name `xml:"Code"`
	Id         string   `xml:"id,attr"`
	CommonName string   `xml:"Name"`
}
