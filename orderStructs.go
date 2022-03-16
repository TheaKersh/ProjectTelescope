package main

import "encoding/xml"

type RSS struct {
  XmlName xml.Name `xml:"Structure"`
  OrderData OrderData `xml:"Structures"`
}

type OrderData struct {
  XmlName xml.Name `xml:"Structures"`
  OrderSet OrderSet  `xml:"DataStructures>DataStructure"`
}

type OrderSet struct {
  XmlName xml.Name `xml:"DataStructures>DataStructure"`
  Components Components `xml:"DataStructureComponents"`
}

type Components struct {
  XmlName xml.Name `xml:"DataStructureComponents"`
  DimensionList []Dimension `xml:"DimensionList>Dimension"`
}

type Dimension struct {
  XmlName xml.Name `xml:"DimensionList>Dimension"`
  Postion string `xml:"position,attr"`
  Id string `xml:"id,attr"`
}

