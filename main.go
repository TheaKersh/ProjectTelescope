package main

import (
	//"bufio"
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"html/template"
	"net/url"
	"strings"

	//"io"
	"io/fs"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	//"strings"
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

/*Utility functions*/
func check(err error) {
	if err != nil {
		panic(err)
	}
}

//I was running into some weird errors where
//Random stuff appeared in my json files
//This trims them by searching them character by character
func trimJson(path string) {
	str, err := os.ReadFile(path)
	check(err)
	depth := 0
	for i := 0; i < len(str); i++ {
		if depth == 0 && str[i] != byte('{') && str[i] != byte('}') {
			if str[i+1] == byte('{') || str[i+1] == byte('}') {
				str = str[i+1:]
				fmt.Print("\n\n\n" + string(str) + "\n\n\n")
				break
			}
		}
		if str[i] == byte('{') {
			depth++
		}
		if str[i] == byte('}') {
			if depth == 0 {
				str = str[:i]
				break
			}
			depth--
		}
	}
	for i := 0; i < len(str); i++ {
		if str[i] == byte('{') {
			depth++
		}
		if str[i] == byte('}') {
			if depth == 0 {
				str = str[:i]
				break
			}
			depth--
		}
	}
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_TRUNC, fs.ModeAppend)
	check(err)
	f.Write(str)
}

func marshalSession(path string, toMarshal Session) error {
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_TRUNC, fs.ModeAppend)
	if err == nil {
		return err
	}
	currBinary, err := json.Marshal(toMarshal)
	if err != nil {
		return err
	}
	f.Write(currBinary)
	return nil
}

// Searches url path for specified string
func searchPath(url string, r *http.Request) string {
	title := r.URL.Path[len("/"+url+"/"):]
	return title
}

//Trims underscores for comparison
func trimUnderscores(toTrim string) string {
	retVal := toTrim
	for strings.Contains(retVal, "_") {
		ind := strings.Index(retVal, "_")
		retVal = retVal[ind:]
	}
	return retVal
}

/*Handler functions*/

//Default/Test get for undata
//Useful for demo/front page
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

//Just redirects to search(for now)
func titleDefault(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/search/", http.StatusFound)
}

//Gets the data set that the user searches for
//Will eventually feature a more user friendly
//way of searching
func searchForData(w http.ResponseWriter, r *http.Request) {
	//Asks the user for search term when get request received,
	//Gets possible mdata features and redirects user otherwise
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
		buf.WriteTo(f2)

		byteVal, err := ioutil.ReadAll(f2)
		check(err)
		xml.Unmarshal(byteVal, &structures)
		element := Dataflow{}
		for _, index := range structures.Structures.Flows.Dataflow {
			fmt.Println(index.Id == SearchTerm)
			if index.Name == SearchTerm {
				element = index
			}
		}

		fmt.Println("\nhttps://data.un.org/ws/rest/datastructure/" + element.DataStructure.RefID.AgencyID + "/" + element.DataStructure.RefID.Id)
		req, err = http.NewRequest("GET", "https://data.un.org/ws/rest/datastructure/"+element.DataStructure.RefID.AgencyID+"/"+element.DataStructure.RefID.Id, nil)
		check(err)
		resp, err = client.Do(req)
		check(err)
		reader := new(bytes.Buffer)
		reader.ReadFrom(resp.Body)
		f, err := os.Create("templateTest.html")
		check(err)
		rss := new(RSS)
		err = xml.Unmarshal(reader.Bytes(), rss)
		check(err)
		currBinary, err := os.ReadFile("current.json")
		check(err)
		current := EmptySess()
		err = json.Unmarshal(currBinary, &current)
		check(err)
		tmpl := template.Must(template.ParseFiles("innerTemplate.html"))
		tmpl.Execute(f, rss)
		if !current.FillY {
			current.X_vals.Name, current.X_vals.Id = element.Name, element.Id
			for _, d := range rss.Data.Set.Components.DimensionList {
				current.X_vals.Params = append(current.X_vals.Params, d.Id)
			}
		} else {
			current.Y_vals.Name, current.Y_vals.Id = element.Name, element.Id
			for _, d := range rss.Data.Set.Components.DimensionList {
				current.Y_vals.Params = append(current.X_vals.Params, d.Id)
			}
		}
		current.FillY = !current.FillY
		err = marshalSession("current.json", current)
		check(err)
		http.Redirect(w, r, "/mdata", http.StatusAccepted)
	}
}

//Gets the metadata features and presents them to the user in menu format
//Uses go template library
//See innerTemplate.html for details
func testParameterization(w http.ResponseWriter, r *http.Request) {
	//Conditional branch for request-response pairings
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
		sess := new(Session)
		vals, err := os.ReadFile("current.json")
		check(err)
		err = json.Unmarshal(vals, sess)
		check(err)
		query := "https://data.un.org/ws/rest/data/"
		if sess.FillY {
			query += sess.Y_vals.Id
		} else {
			query += sess.X_vals.Id
		}
		query += "/"
		fmt.Print(query)
		f, err := os.Open("termOrder.html")
		check(err)
		slice, err := ioutil.ReadAll(f)
		check(err)
		var RSS *RSS = new(RSS)
		xml.Unmarshal(slice, RSS)

		//Metadata features
		var features []string = make([]string, 0)
			for index := range RSS.Data.Set.Components.DimensionList {
				features = append(features, r.PostForm[listnames[index]][0]+".")
			}
		}

		if !sess.FillY {
			sess.X_vals.Params = features
		} else {
			sess.Y_vals.Params = features
		}
		for _, element := range features {
			fmt.Print(element)
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
		fmt.Print(buf.String())

		f, err = os.OpenFile("independent.json", os.O_WRONLY, fs.ModeAppend)
		check(err)
		f2, err := os.OpenFile("dependent.json", os.O_WRONLY, fs.ModeAppend)
		check(err)
		if sess.FillY {
			_, err = buf.WriteTo(f2)
			check(err)
		} else {
			_, err = buf.WriteTo(f)
			check(err)
		}
		if buf.String() != "NoRecordsFound" && buf.String() != "Could not find Dataflow and/or DSD related with this data request" {
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
		if sess.FillY {
			http.Redirect(w, r, "/output/", http.StatusFound)
		} else {
			http.Redirect(w, r, "/confirm/", http.StatusFound)
		}
		sess.FillY = !sess.FillY
		f, err = os.Create("current.json")
		check(err)
		vals, err = json.Marshal(sess)
		check(err)
		f.Write(vals)
	}

}

//Confirms to the user that the selected params are correct
//Not sure whether i'll keep this
func IntermediateStep(w http.ResponseWriter, r *http.Request) {
	//Conditional branch for request-response pairing
	if r.Method == "GET" {
		tmpl := template.Must(template.ParseFiles("dataReview.html"))
		tmpl.Execute(w, postForm)
	} else {
		http.Redirect(w, r, "/search", http.StatusFound)
	}
}

//Writes html for graph page to receiver (see graph.html for details)
func outputGraph(w http.ResponseWriter, r *http.Request) {
	bytes, err := os.ReadFile("graph.html")
	check(err)
	w.Write(bytes)
}

//Retrieves js
//Needed for webpage to fetch scripts
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

//Json file getter
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
	//Trims json at beginning and end of function
	//This for debug: So changes to current.json can be seen
	trimJson("current.json")
	f, err := os.OpenFile("current.json", os.O_RDWR, fs.ModeAppend)
	check(err)
	buf := new(bytes.Buffer)
	buf.ReadFrom(f)
	Sess := new(Session)
	strVal := strings.TrimSpace(buf.String())

	if strVal == "" {
		Sess := EmptySess()
		byteVal, _ := json.Marshal(Sess)
		f.Write(byteVal)

	} else {
		err = json.Unmarshal(buf.Bytes(), Sess)
		check(err)
		executed = Sess.FillY
	}
	trimJson("current.json")
	http.HandleFunc("/javascript/", retrieveJS)
	http.HandleFunc("/view/", getRequest)
	http.HandleFunc("/", titleDefault)
	http.HandleFunc("/search/", searchForData)
	http.HandleFunc("/mdata/", testParameterization)
	http.HandleFunc("/json/", retrieveJSON)
	http.HandleFunc("/output/", outputGraph)
	http.HandleFunc("/confirm/", IntermediateStep)
	s := &http.Server{
		Addr:           ":8080",
		MaxHeaderBytes: 1 << 20,
		ErrorLog:       log.Default(),
	}
	//Function will return an error if it encounters one
	check(s.ListenAndServe())
}
