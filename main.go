package main

import (
	//"bufio"
	"bytes"
	"encoding/xml"
	"fmt"
	"html/template"

	//"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	//"strings"
)

//Types for dealing with common elements in undata xml
type Header struct {
	XMLName xml.Name `xml:"Header"`
	IDRef   string   `xml:"ID"`
	Test    string   `xml:"Test"`
}

//Types to unmarshal and search xml when looking for a specific dataset
func (c CodeStructure) toString() [][]string {
	var codelists CodeLists = c.CodeStructures.CodeLists
	var retVal [][]string = make([][]string, 30)
	for i := range retVal {
		retVal[i] = make([]string, 1000)
	}
	for index, codelist := range codelists.CodeLists {
		retVal[index][0] += codelist.Description + "\n"
		for j, code := range codelist.Codes {
			retVal[index][j] += code.CommonName + "  " + code.Id + "\n"
		}
	}
	return retVal
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func initNameQueryMap() map[string]string {
	var m map[string]string = make(map[string]string)
	m["Greenhouse Gases"] = "DF_UNData_UNFCC"
	m["Carbon"] = m["Greenhouse Gases"] + "/A.EN_ATM_PFCE.AUS.Gg_CO2"
	return m
}

func getRequest(w http.ResponseWriter, r *http.Request) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://data.un.org/ws/rest/data/DF_UNData_UNFCC/A.EN_ATM_PFCE.AUS.Gg_CO2", nil)
	check(err)
	req.Header.Set("Accept", "text/json")
	resp, err := client.Do(req)
	check(err)

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	defer resp.Body.Close()
	fmt.Fprintf(w, "%s", string(buf.Bytes()))

}

func titleDefault(w http.ResponseWriter, r *http.Request) {
	buf, err := os.ReadFile("index.html")
	check(err)
	w.Write(buf)
}

func CreateDataflow() Dataflow {
	return *new(Dataflow)
}

func searchForData(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("askForSubject.html"))
	client := &http.Client{}
	tmpl.Execute(w, nil)
	SearchTerm := r.FormValue("SearchTerm")
	var structures Structure
	//m := initNameQueryMap()
	req, err := http.NewRequest("GET", "https://data.un.org/ws/rest/dataflow/", nil)
	check(err)

	resp, err := client.Do(req)
	check(err)
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)

	f2, err := os.Open("test.html")
	check(err)
	byteVal, err := ioutil.ReadAll(f2)
	check(err)
	xml.Unmarshal(byteVal, &structures)
	element := CreateDataflow()
	for _, index := range structures.Structures.Flows.Dataflow {
		fmt.Println(index.Id == SearchTerm)
		if index.Id == SearchTerm {
			element = index
		}
	}
	fmt.Println("\nhttps://data.un.org/ws/rest/datastructure/" + element.DataStructure.RefID.AgencyID + "/" + element.DataStructure.RefID.Id + "/?references=children\n")
	req, err = http.NewRequest("GET", "https://data.un.org/ws/rest/datastructure/"+element.DataStructure.RefID.AgencyID+"/"+element.DataStructure.RefID.Id+"/?references=children", nil)
	check(err)
	resp, err = client.Do(req)
	check(err)
	reader := new(bytes.Buffer)
	reader.ReadFrom(resp.Body)
	f, err := os.Create("templateTest.html")
	check(err)
	var codestructures *CodeStructure = new(CodeStructure)
	xml.Unmarshal(reader.Bytes(), codestructures)
	tmpl = template.Must(template.ParseFiles("innerTemplate.html"))
	tmpl.Execute(f, codestructures)

}

func testParameterization(w http.ResponseWriter, r *http.Request) {
	f, err := os.Open("templateTest.html")
	check(err)
	bytes, err := ioutil.ReadAll(f)
	check(err)
	w.Write(bytes)

}

func main() {
	http.HandleFunc("/view/", getRequest)
	http.HandleFunc("/", titleDefault)
	http.HandleFunc("/search/", searchForData)
	http.HandleFunc("/mdata/", testParameterization)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
