package main

import "encoding/xml"

type CodeStructure struct {
	XMLName        xml.Name `xml:"Structure"`
	Header         Header   `xml:"Header"`
	CodeStructures CodeStructures
}

type CodeStructures struct {
	XMLName   xml.Name  `xml:"Structures"`
	CodeLists CodeLists `xml:"Codelists"`
  OrderStructures OrderDatastructures `xml:"Datastructures"`
}

type OrderDatastructures struct {
  XmlName xml.Name                 `xml:"Datastructures"`
  Datastructure OrderDatastructure `xml:"DataStructure"`
}

type OrderDatastructure struct {
  XmlName xml.Name `xml:"DataStructure"`
  QueryOrder DatastructureComponents `xml:"DataStructureComponents"`
  Id string `xml:"id,attr"`
}

type DatastructureComponents struct {
  XmlName xml.Name `xml:"DataStructureComponents"`
  DimensionList DimensionList `xml:"DimensionList"`
}

type DimensionList struct {
  XmlName xml.Name `xml:"DimensionList"`
  Dimensions []Dimension `xml:"DimensionList>Dimension"`
}

type Dimension struct {
  XmlName xml.Name `xml:"Dimension"`
  Id string `xml:"id,attr"`
  Position string `xml:"position,attr"` 
}

type CodeLists struct {
	XMlName   xml.Name   `xml:"Codelists"`
	CodeLists []CodeList `xml:"Codelist"`
}

type CodeList struct {
	XMlName     xml.Name `xml:"Codelist"`
	Description string   `xml:"Description"`
	Name        string   `xml:"Name"`
	Codes       []Code   `xml:"Code"`
}

type Code struct {
	XmlName    xml.Name `xml:"Code"`
	Id         string   `xml:"id,attr"`
	CommonName string   `xml:"Name"`
}
