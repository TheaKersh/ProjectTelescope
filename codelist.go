package main

import(
  "encoding/xml"
)
type Top struct {
  XmlName xml.Name   `xml:"Structure"`
  Codelist Codelist  `xml:"Structures>Codelists>Codelist"`
}


type Codelist struct {
  XmlName xml.Name   `xml:"Codelist"`
  Id    string     `xml:"id,attr"`
  Name  string     `xml:"Name"`  
  Codes []Code     `xml:"Code"`
}
type Code struct {
  XmlName xml.Name   `xml:"Code"`
  Name  string       `xml:"Name"`
  Id    string       `xml:"id,attr"`
}