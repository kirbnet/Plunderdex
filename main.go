package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"sort"
	"strings"

	"github.com/gorilla/mux"
)

//Global Variables
var plunderlings Plunderlings
var plunderclasses []string
var plunderfolks []string
var plundercolors []string
var plunderwaves []string
var plundertags []string

//Templates
var tpl = template.Must(template.ParseFiles("static/index.html"))
var indvtpl = template.Must(template.ParseFiles("static/plunderling_template.html"))

func main() {
	//Open and Parse file to object in memory
	GetPlunderlings()

	//If there is a preconfigured port
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	//Mux Http Handler
	router := mux.NewRouter()

	//Request handlers
	router.HandleFunc("/", namesHandler)
	router.HandleFunc("/class/", classesHandler)
	router.HandleFunc("/class/{class}", classHandler)
	router.HandleFunc("/name/", namesHandler)
	router.HandleFunc("/color/", colorsHandler)
	router.HandleFunc("/color/{color}", colorHandler)
	router.HandleFunc("/plunderfolk/", plunderfolksHandler)
	router.HandleFunc("/plunderfolk/{plunderfolk}", plunderfolkHandler)
	router.HandleFunc("/wave/", wavesHandler)
	router.HandleFunc("/wave/{wave}", waveHandler)
	router.HandleFunc("/tag/", tagsHandler)
	router.HandleFunc("/tag/{tag}", tagHandler)
	router.HandleFunc("/search", searchHandler)
	router.HandleFunc("/{figure}", figureHandler)
	router.HandleFunc("/accessory/{accessory}", accessoryHandler)

	//Static Resources
	fs := http.FileServer(http.Dir("./static"))
	router.PathPrefix("/static").Handler(http.StripPrefix("/static/", fs))
	//Start Port Listener/Web Server
	http.ListenAndServe(":"+port, router)
}

//Single figure display
func figureHandler(w http.ResponseWriter, r *http.Request) {
	var myfigure Plunderling

	//parse request data
	reqvars := mux.Vars(r)
	figure := reqvars["figure"]
	for _, individual := range plunderlings.Plunderlings {
		if individual.Name == figure {
			myfigure = individual
		}
	}
	//display new subset
	pagedata := &FigurePageData{myfigure, template.HTML(myfigure.Bio), plunderclasses, plunderfolks, plundercolors, plunderwaves, plundertags}
	indvtpl.Execute(w, pagedata)
}

//Main page and default handler. Displays everything.
func indexHandler(w http.ResponseWriter, r *http.Request) {
	pagedata := &PageData{plunderlings, plunderclasses, plunderfolks, plundercolors, plunderwaves, plundertags}
	if err := tpl.Execute(w, pagedata); err != nil {
		fmt.Println(err)
	}
}

//Nebulous catch-most for searching for things. Quite vague.
func searchHandler(w http.ResponseWriter, r *http.Request) {
	var foundlings Plunderlings

	//parse request data
	q := r.URL.Query().Get("search")
	for _, individual := range plunderlings.Plunderlings {
		if strings.Contains(individual.Class, q) ||
			strings.Contains(individual.Color, q) ||
			strings.Contains(individual.Name, q) ||
			strings.Contains(individual.Notes, q) ||
			strings.Contains(individual.Plunderfolk, q) ||
			strings.Contains(individual.Bio, q) {
			foundlings.AddItem(individual)
		}
	}
	//display the new subset
	pagedata := &PageData{foundlings, plunderclasses, plunderfolks, plundercolors, plunderwaves, plundertags}
	tpl.Execute(w, pagedata)
}

//Display all, sorting by Name
func namesHandler(w http.ResponseWriter, r *http.Request) {
	sort.Slice(plunderlings.Plunderlings, func(i, j int) bool {
		return plunderlings.Plunderlings[i].Name < plunderlings.Plunderlings[j].Name
	})
	pagedata := &PageData{plunderlings, plunderclasses, plunderfolks, plundercolors, plunderwaves, plundertags}
	tpl.Execute(w, pagedata)
}

//Display all, sorting by "color"
func colorsHandler(w http.ResponseWriter, r *http.Request) {
	sort.Slice(plunderlings.Plunderlings, func(i, j int) bool {
		//Logic to sort by NAME ASC as secondary sort
		if plunderlings.Plunderlings[i].Color == plunderlings.Plunderlings[j].Color {
			return plunderlings.Plunderlings[i].Name < plunderlings.Plunderlings[j].Name
		}
		return plunderlings.Plunderlings[i].Color < plunderlings.Plunderlings[j].Color
	})
	pagedata := &PageData{plunderlings, plunderclasses, plunderfolks, plundercolors, plunderwaves, plundertags}
	tpl.Execute(w, pagedata)
}

//Display a single "color"
func colorHandler(w http.ResponseWriter, r *http.Request) {
	var justOneColor Plunderlings

	//parse request data
	reqvars := mux.Vars(r)
	color := reqvars["color"]
	for _, individual := range plunderlings.Plunderlings {
		if individual.Color == color {
			justOneColor.AddItem(individual)
		}
	}
	//display new subset
	pagedata := &PageData{justOneColor, plunderclasses, plunderfolks, plundercolors, plunderwaves, plundertags}
	tpl.Execute(w, pagedata)
}

//Display all, sorting by "Class"
func classesHandler(w http.ResponseWriter, r *http.Request) {
	sort.Slice(plunderlings.Plunderlings, func(i, j int) bool {
		//Logic to sort by NAME ASC as secondary sort
		if plunderlings.Plunderlings[i].Class == plunderlings.Plunderlings[j].Class {
			return plunderlings.Plunderlings[i].Name < plunderlings.Plunderlings[j].Name
		}
		return plunderlings.Plunderlings[i].Class < plunderlings.Plunderlings[j].Class
	})
	pagedata := &PageData{plunderlings, plunderclasses, plunderfolks, plundercolors, plunderwaves, plundertags}
	tpl.Execute(w, pagedata)
}

//Display a single "class"
func classHandler(w http.ResponseWriter, r *http.Request) {
	var justOneClass Plunderlings

	//parse request data
	reqvars := mux.Vars(r)
	class := reqvars["class"]
	for _, individual := range plunderlings.Plunderlings {
		if individual.Class == class {
			justOneClass.AddItem(individual)
		}
	}
	//display new subset
	pagedata := &PageData{justOneClass, plunderclasses, plunderfolks, plundercolors, plunderwaves, plundertags}
	tpl.Execute(w, pagedata)
}

//Display all, sorting by "plunderfolk"
func plunderfolksHandler(w http.ResponseWriter, r *http.Request) {
	sort.Slice(plunderlings.Plunderlings, func(i, j int) bool {
		//Logic to sort by NAME ASC as secondary sort
		if plunderlings.Plunderlings[i].Plunderfolk == plunderlings.Plunderlings[j].Plunderfolk {
			return plunderlings.Plunderlings[i].Name < plunderlings.Plunderlings[j].Name
		}
		return plunderlings.Plunderlings[i].Plunderfolk < plunderlings.Plunderlings[j].Plunderfolk
	})
	pagedata := &PageData{plunderlings, plunderclasses, plunderfolks, plundercolors, plunderwaves, plundertags}
	tpl.Execute(w, pagedata)
}

//Display single "plunderfolk"
func plunderfolkHandler(w http.ResponseWriter, r *http.Request) {
	var justOnePlunderfolk Plunderlings

	//parse request data
	reqvars := mux.Vars(r)
	plunderfolk := reqvars["plunderfolk"]
	for _, individual := range plunderlings.Plunderlings {
		if individual.Plunderfolk == plunderfolk {
			justOnePlunderfolk.AddItem(individual)
		}
	}
	//display new subset
	pagedata := &PageData{justOnePlunderfolk, plunderclasses, plunderfolks, plundercolors, plunderwaves, plundertags}
	tpl.Execute(w, pagedata)
}

//Display all, sorting by "wave"
func wavesHandler(w http.ResponseWriter, r *http.Request) {
	sort.Slice(plunderlings.Plunderlings, func(i, j int) bool {
		//Logic to sort by NAME ASC as secondary sort
		if plunderlings.Plunderlings[i].Wave == plunderlings.Plunderlings[j].Wave {
			return plunderlings.Plunderlings[i].Name < plunderlings.Plunderlings[j].Name
		}
		return plunderlings.Plunderlings[i].Wave < plunderlings.Plunderlings[j].Wave
	})
	pagedata := &PageData{plunderlings, plunderclasses, plunderfolks, plundercolors, plunderwaves, plundertags}
	tpl.Execute(w, pagedata)
}

//Display single "wave"
func waveHandler(w http.ResponseWriter, r *http.Request) {
	var justOneWave Plunderlings

	//parse request data
	reqvars := mux.Vars(r)
	wave := reqvars["wave"]
	for _, individual := range plunderlings.Plunderlings {
		if individual.Wave == wave {
			justOneWave.AddItem(individual)
		}
	}
	//display new subset
	pagedata := &PageData{justOneWave, plunderclasses, plunderfolks, plundercolors, plunderwaves, plundertags}
	tpl.Execute(w, pagedata)
}

//Display all, sorting by "wave"
func tagsHandler(w http.ResponseWriter, r *http.Request) {
	var tagged Plunderlings
	for _, individual := range plunderlings.Plunderlings {
		if individual.Tag != "" {
			tagged.AddItem(individual)
		}
	}
	sort.Slice(tagged.Plunderlings, func(i, j int) bool {
		//Logic to sort by NAME ASC as secondary sort
		if tagged.Plunderlings[i].Tag == tagged.Plunderlings[j].Tag {
			return tagged.Plunderlings[i].Name < tagged.Plunderlings[j].Name
		}
		return tagged.Plunderlings[i].Tag < tagged.Plunderlings[j].Tag
	})
	pagedata := &PageData{tagged, plunderclasses, plunderfolks, plundercolors, plunderwaves, plundertags}
	tpl.Execute(w, pagedata)
}

//Search for an accessory
func accessoryHandler(w http.ResponseWriter, r *http.Request) {
	var foundlings Plunderlings

	//parse request data
	reqvars := mux.Vars(r)
	accessory := reqvars["accessory"]
	for _, individual := range plunderlings.Plunderlings {
		for _, acc := range individual.Accessories {
			if strings.Contains(acc, accessory) {
				foundlings.AddItem(individual)
			}
		}
	}
	//display the new subset
	pagedata := &PageData{foundlings, plunderclasses, plunderfolks, plundercolors, plunderwaves, plundertags}
	tpl.Execute(w, pagedata)
}

//Display single "tag"
func tagHandler(w http.ResponseWriter, r *http.Request) {
	var justOneTag Plunderlings

	//parse request data
	reqvars := mux.Vars(r)
	tag := reqvars["tag"]
	for _, individual := range plunderlings.Plunderlings {
		if individual.Tag == tag {
			justOneTag.AddItem(individual)
		}
	}
	//display new subset
	pagedata := &PageData{justOneTag, plunderclasses, plunderfolks, plundercolors, plunderwaves, plundertags}
	tpl.Execute(w, pagedata)
}

//DATA FUNCTIONS
//Load and parse JSON database into object
func GetPlunderlings() {
	jsonFile, err := os.Open("static/Plunderlings.json")
	if err != nil {
		fmt.Println(err)
	}
	//fmt.Println("Opened Plunderlings.json!")
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	//var plunderlings Plunderlings
	json.Unmarshal(byteValue, &plunderlings)

	//Now that the database is put into structs, create groups of data for soring
	GroupPlunderfolks()
	GroupColors()
	GroupClasses()
	GroupWaves()
	GroupTags()
}

//Data Class and Functions
type Plunderlings struct {
	Plunderlings []Plunderling `json:"plunderlings"`
}

func (plunderlings *Plunderlings) AddItem(plunderling Plunderling) {
	plunderlings.Plunderlings = append(plunderlings.Plunderlings, plunderling)
}

type Plunderling struct {
	Name        string   `json:"name"`
	Class       string   `json:"class"`
	Color       string   `json:"color"`
	Notes       string   `json:"notes"`
	Plunderfolk string   `json:"plunderfolk"`
	Wave        string   `json:"wave"`
	Tag         string   `json:"tag"`
	Accessories []string `json:"accessories"`
	Bio         string   `json:"bio"`
}

//Create a slice of "plunderfolk"s based on data set
func GroupPlunderfolks() {
	tmpplunderfolks := make(map[string]int)
	for i := range plunderlings.Plunderlings {
		_, exists := tmpplunderfolks[plunderlings.Plunderlings[i].Plunderfolk]
		if !exists {
			tmpplunderfolks[plunderlings.Plunderlings[i].Plunderfolk] = i
		}
	}
	//Take map keys and convert to slice and sort
	plunderfolks = MapKeysToSliceSorted(tmpplunderfolks)
}

//Create a slice of "color"s based on data set
func GroupColors() {
	colors := make(map[string]int)
	for i := range plunderlings.Plunderlings {
		_, exists := colors[plunderlings.Plunderlings[i].Color]
		if !exists {
			colors[plunderlings.Plunderlings[i].Color] = i
		}
	}
	//Take map keys and convert to slice and sort
	plundercolors = MapKeysToSliceSorted(colors)
}

//Create a slice of "class"es based on data set
func GroupClasses() {
	classes := make(map[string]int)
	for i := range plunderlings.Plunderlings {
		_, exists := classes[plunderlings.Plunderlings[i].Class]
		if !exists {
			classes[plunderlings.Plunderlings[i].Class] = i
		}
	}
	//Take map keys and convert to slice and sort
	plunderclasses = MapKeysToSliceSorted(classes)
}

//Create a slice of "wave"s based on data set
func GroupWaves() {
	classes := make(map[string]int)
	for i := range plunderlings.Plunderlings {
		_, exists := classes[plunderlings.Plunderlings[i].Wave]
		if !exists {
			classes[plunderlings.Plunderlings[i].Wave] = i
		}
	}
	//Take map keys and convert to slice and sort
	plunderwaves = MapKeysToSliceSorted(classes)
}

//Create a slice of "tag"s based on data set
func GroupTags() {
	tags := make(map[string]int)
	for i := range plunderlings.Plunderlings {
		_, exists := tags[plunderlings.Plunderlings[i].Tag]
		if !exists {
			tags[plunderlings.Plunderlings[i].Tag] = i
		}
	}
	//Take map keys and convert to slice and sort
	plundertags = MapKeysToSliceSorted(tags)
}

//Generic support functions
func MapKeysToSliceSorted(m map[string]int) []string {
	keys := make([]string, len(m))

	i := 0
	for k := range m {
		keys[i] = k
		i++
	}
	//Sort
	sort.Slice(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})
	return keys
}

//Struct to handle Page Data
//First element is the individual Plunderlings/items
//Remainder are lists aggregated from data
type PageData struct {
	Lings        Plunderlings
	Classes      []string
	Plunderfolks []string
	Colors       []string
	Waves        []string
	Tags         []string
}
type FigurePageData struct {
	Figure       Plunderling
	BioData      template.HTML
	Classes      []string
	Plunderfolks []string
	Colors       []string
	Waves        []string
	Tags         []string
}
