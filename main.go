package main

import (
	//"bufio"
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"html/template"
	"net/url"
<<<<<<< HEAD
	"strings"
=======
>>>>>>> 9edba4da8db13b36d3d6acd7e84dbf76240a28af

	//"io"
	"io/fs"
	"io/ioutil"
	"log"
	"net/http"
	"os"
<<<<<<< HEAD
	//"strings"
=======
	"strconv"
	"strings"
>>>>>>> 9edba4da8db13b36d3d6acd7e84dbf76240a28af
)

var listnames []string
var listIds []string
var postForm url.Values
var executed bool

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

func checkList(slice []Dimension, id string, testmap map[string]string) bool {
	for _, element := range slice {
		if element.Id == id {
			return true
		} else if element.Id == testmap[id] {
			return true
		}
	}
	return false
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
	http.Redirect(w, r, "/search/", http.StatusFound)
<<<<<<< HEAD
=======
}

func GetFromJson(path string) []float64 {
	m := new(map[string]interface{})
	dat, err := os.ReadFile(path)
	check(err)
	json.Unmarshal(dat, m)
	items := (*m)["dataSets"].([]interface{})
	points := items[0].(map[string]interface{})["series"].(map[string]interface{})["0:0:0:0"].(map[string]interface{})["observations"].(map[string]interface{})
	inter := make(map[string]float64)
	for key, value := range points {
		inter[key] = value.([]interface{})[0].(float64)
	}
	list := make([]float64, 31)
	for key, value := range inter {
		ind, err := strconv.Atoi(key)
		check(err)
		list[ind] = value
	}
	return list
}

func Correl(w http.ResponseWriter, r *http.Request) {
	toUse := correlationOptions[searchPath("correlation", r)]
	fmt.Fprint(w, toUse(GetFromJson("independent.json"), GetFromJson("dependent.json")))

}

func Pearsons(x []float64, y []float64) float64 {
	pearson := constrPearson(x, y)
	return pearson.Correl()
}

var correlationOptions map[string]func(x []float64, b []float64) float64 = map[string]func(x []float64, b []float64) float64{
	"Pearsons": Pearsons,
>>>>>>> 9edba4da8db13b36d3d6acd7e84dbf76240a28af
}

func CreateDataflow() Dataflow {
	return *new(Dataflow)
}

func searchForData(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		tmpl := template.Must(template.ParseFiles("askForSubject.html"))
		tmpl.Execute(w, nil)
	} else {
		SearchTerm := r.FormValue("SearchTerm")
		var structures Structure
		client := http.Client{}
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
		fmt.Printf("%#v\n\n\n", codestructures)
		tmpl := template.Must(template.ParseFiles("innerTemplate.html"))
		tmpl.Execute(f, codestructures)
		listnames = make([]string, 0)
		listIds = make([]string, 0)
		for _, element := range codestructures.CodeStructures.CodeLists.CodeLists {
			listnames = append(listnames, element.Name)
			listIds = append(listIds, element.Id)
			fmt.Print(element.Name + "\n")
		}
		resp, err = http.Get("https://data.un.org/ws/rest/datastructure/" + element.DataStructure.RefID.AgencyID + "/" + element.DataStructure.RefID.Id)
		check(err)
		reader = new(bytes.Buffer)
		reader.ReadFrom(resp.Body)
		f, err = os.Create("termOrder.html")
		check(err)
		f.Write(reader.Bytes())
		http.Redirect(w, r, "/mdata", http.StatusFound)
	}
}

func testParameterization(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		f, err := os.Open("templateTest.html")
		check(err)
		slice, err := ioutil.ReadAll(f)
		check(err)
		buf := bytes.NewBuffer(slice)
		buf.WriteTo(w)
	} else {
		r.ParseForm()
		fmt.Print(r.PostForm)
		postForm = r.PostForm
		query := "https://data.un.org/ws/rest/data/DF_UNDATA_ENERGY/"
		fmt.Print(query)
		f, err := os.Open("termOrder.html")
		check(err)
		slice, err := ioutil.ReadAll(f)
		check(err)
		var RSS *RSS = new(RSS)
		xml.Unmarshal(slice, RSS)
		//Metadata features
		var features []string = make([]string, 0)
		for ind, element := range listIds {
			fmt.Println(element + "\n\n")
			for _, index := range RSS.OrderData.OrderSet.Components.DimensionList {
				if strings.Index(index.Id, "_") != -1 {
					index.Id = index.Id[strings.Index(index.Id, "_"):]
				}
				fmt.Println(index.Id)
				if strings.Contains(element, index.Id) {
					features = append(features, r.PostForm[listnames[ind]][0]+".")
				}
			}
		}

		for _, element := range features {
			query = query + element
		}

		query = strings.TrimRight(query, ".")
		fmt.Println("\n\n\n\n" + query)
		client := http.Client{}
		req, err := http.NewRequest("GET", query, nil)
		check(err)
		req.Header.Set("Accept", "text/json")
		resp, err := client.Do(req)
		check(err)
		buf := new(bytes.Buffer)
		buf.ReadFrom(resp.Body)
<<<<<<< HEAD
		fmt.Print(buf.String())

=======
>>>>>>> 9edba4da8db13b36d3d6acd7e84dbf76240a28af
		f, err = os.OpenFile("independent.json", os.O_WRONLY, fs.ModeAppend)
		check(err)
		f2, err := os.OpenFile("dependent.json", os.O_WRONLY, fs.ModeAppend)
		check(err)
		if executed {
<<<<<<< HEAD
			_, err = buf.WriteTo(f2)
			check(err)
		} else {
			_, err = buf.WriteTo(f)
			check(err)
=======
			buf.WriteTo(f2)
		} else {
			buf.WriteTo(f)
>>>>>>> 9edba4da8db13b36d3d6acd7e84dbf76240a28af
		}
		if buf.String() != "NoRecordsFound" {
			_, err = f.Write(buf.Bytes())
			check(err)
		} else {
			f, err := os.Open("templateTest.html")
			check(err)
			slice, err := ioutil.ReadAll(f)
			check(err)
			buf := bytes.NewBuffer(slice)
			buf.WriteTo(w)
		}
		if executed {
			http.Redirect(w, r, "/output/", http.StatusFound)
		} else {
			http.Redirect(w, r, "/confirm/", http.StatusFound)
		}
		executed = !executed
<<<<<<< HEAD
		if 
=======
	}
}
func trimJson(path string) {
	if !strings.Contains(path, ".json") {
		path = path + ".json"
	}
	slice, err := os.ReadFile(path)
	check(err)
	depth := 0
	prevDepth := depth
	for ind, element := range slice {
		if rune(element) == '{' {
			depth++
		}
		if rune(element) == '}' {
			depth--
		}
		if prevDepth > depth && depth == 0 {
			if ind < len(slice) {
				slice = slice[:ind+1]
				os.Truncate(path, 0)
				os.WriteFile(path, slice, fs.ModeAppend)
			}
			break
		}
		prevDepth = depth
>>>>>>> 9edba4da8db13b36d3d6acd7e84dbf76240a28af
	}
}

func IntermediateStep(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		tmpl := template.Must(template.ParseFiles("dataReview.html"))
		tmpl.Execute(w, postForm)
	} else {
		http.Redirect(w, r, "/search", http.StatusFound)
	}
}

func outputGraph(w http.ResponseWriter, r *http.Request) {
	bytes, err := os.ReadFile("graph.html")
	check(err)
	w.Write(bytes)
}

func searchPath(path string, r *http.Request) string {
	title := r.URL.Path[len("/"+path+"/"):]
	return title
}

func retrieveJS(w http.ResponseWriter, r *http.Request) {
	title := searchPath("javascript", r)
	buf := new(bytes.Buffer)
	file, err := os.Open("javascript/" + title)
	check(err)
	_, err = buf.ReadFrom(file)
	check(err)
	fmt.Print(buf.String())
	w.Write(buf.Bytes())
}

func retrieveJSON(w http.ResponseWriter, r *http.Request) {
	title := searchPath("json", r)
	buf := new(bytes.Buffer)
	trimJson(title + ".json")
	file, err := os.Open(title + ".json")
	check(err)
	_, err = buf.ReadFrom(file)
	check(err)
	w.Write(buf.Bytes())
}

func main() {
	http.HandleFunc("/javascript/", retrieveJS)
	http.HandleFunc("/view/", getRequest)
	http.HandleFunc("/", titleDefault)
	http.HandleFunc("/search/", searchForData)
	http.HandleFunc("/mdata/", testParameterization)
	http.HandleFunc("/json/", retrieveJSON)
	http.HandleFunc("/output/", outputGraph)
	http.HandleFunc("/confirm/", IntermediateStep)
<<<<<<< HEAD
=======
	http.HandleFunc("/correlation/", Correl)
>>>>>>> 9edba4da8db13b36d3d6acd7e84dbf76240a28af
	s := &http.Server{
		Addr:           ":8080",
		MaxHeaderBytes: 1 << 20,
		ErrorLog:       log.New(os.Stdout, "err:", log.Ldate|log.Ltime|log.Lshortfile),
	}
	log.Println(s.ListenAndServe())
}
