import "encoding/xml"


type Structure struct {
	XMLName xml.Name `xml:"Structure"`
  Header Header `xml:"Header"`
  Structures Structures `xml:"Structures"`
}


type Structures struct {
  XMLName xml.Name `xml:"Structures"`
  Flows Dataflows `xml:"Dataflows"`
}

type Dataflows struct {
  XmlName xml.Name `xml:"Dataflows"`
  Dataflow []Dataflow `xml:"Dataflow"`
}
type Dataflow struct {
  XMLName xml.Name `xml:"Dataflow"`
  Id string `xml:"id,attr"`
  Name string `xml:"Name"`
  DataStructure DataStructure `xml:"Structure"`
}

type DataStructure struct {
  XMlName xml.Name `xml:"Structure"`
  RefID refID `xml:"Ref"`
}

type refID struct {
  XMLName xml.Name `xml:"Ref"`
  Id string `xml:"id,attr"`
  AgencyID string `xml:"agencyID,attr"`
}