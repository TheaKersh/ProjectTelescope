package main

import "encoding/xml"

type RSS struct {
	XmlName xml.Name `xml:"Structure"`
	Data    Data     `xml:"Structures"`

}

type Data struct {
	XmlName xml.Name `xml:"Structures"`
	Set     Set      `xml:"DataStructures>DataStructure"`
}

type Set struct {
	XmlName    xml.Name   `xml:"DataStructures>DataStructure"`
	Components Components `xml:"DataStructureComponents"`
}

type Components struct {
	XmlName       xml.Name    `xml:"DataStructureComponents"`
	DimensionList []Dimension `xml:"DimensionList>Dimension"`
}

type Dimension struct {
	XmlName xml.Name `xml:"DimensionList>Dimension"`
	Postion string   `xml:"position,attr"`
	Id      string   `xml:"id,attr"`
  RefID   RefID  `xml:"LocalRepresentation>Enumeration>Ref"`
}

type RefID struct {
  XmlName xml.Name `xml:"Ref"`
  Id string        `xml:"id,attr"`
  AgencyID string  `xml:"agencyID,attr"`
}