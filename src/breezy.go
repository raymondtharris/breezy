package main

import (
	"fmt"
	"encoding/json"
	"net/http"
	"html/template"
	"io/ioutil"
)

const brImage, brAudio , brVideo = 0, 1, 2
const brPost = 0

type breezyMedia struct{
	name, filename string
	filesSize float32
	brType int
}

func (m breezyMedia) String() string{
	var typeDescription = ""
	switch m.brType{
	case brImage:
		typeDescription = "Image"
	case brAudio:
		typeDescription = "Audio"
	case brVideo:
		typeDescription = "Video"
	
	}
	return fmt.Sprintf("%v, %v", m.name , typeDescription)
}

type breezyActivity struct{
	name string
	activityBody string
	mediaStructure [5]breezyMedia
	brType int
	//dateCreated, dateModified Time
}


type Page struct{
	Title string
	Body []byte
	
}

func loadPage(pageName string) (*Page, error){
	pageBody, err := ioutil.ReadFile(pageName+".html")
	if err != nil {
	        return nil, err
	    }
	return &Page{Title:pageName, Body: pageBody }, nil
}

func webHandler(w http.ResponseWriter, r *http.Request){
	p,_ := loadPage("index")
	t, _ := template.ParseFiles("../src/index.html")
	t.Execute(w, p)
}

func breezyLoginHandler(w http.ResponseWriter, r *http.Request){
	var temp = ""
	
	temp = r.FormValue("username")
	if len(temp) > 0 {
		print(temp)
	}
	
	p,_ := loadPage("login")
	t, _ := template.ParseFiles("../src/views/login.html")
	t.Execute(w, p)
}
type test struct{
	what string
}

func breezyLoginCredentrials(w http.ResponseWriter, r *http.Request){
	decoder := json.NewDecoder(r.Body)
	var test1 test 
	err:= decoder.Decode(&test1)
	if err != nil{
		panic(err)
	}
	fmt.Println(test1.what)
	
}

func breezyEditHandler(w http.ResponseWriter, r *http.Request){
	p,_ := loadPage("editor")
	t,_ := template.ParseFiles("../src/views/edit.html")
	t.Execute(w,p)
}

func breezyDashboardHandler(w http.ResponseWriter, r *http.Request){
	p,_ := loadPage("dashboard")
	t, _ := template.ParseFiles("../src/views/dashboard.html")
	t.Execute(w, p)
}



func HandleDirs(){
	http.Handle("/lib/", http.StripPrefix("/lib/", http.FileServer(http.Dir("../src/lib/"))))
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("../src/js/"))))
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("../src/css/"))))
	http.Handle("/views/", http.StripPrefix("/views/", http.FileServer(http.Dir("../src/views/"))))
}



func main(){
	//http.HandleFunc("/admin", breezyLoginHandler)
	HandleDirs()
	http.HandleFunc("/admin", breezyLoginHandler)
	http.HandleFunc("/checkcredentials", breezyLoginCredentrials)
	
	
	http.HandleFunc("/edit", breezyEditHandler)
	http.HandleFunc("/dashboard", breezyDashboardHandler)
	http.HandleFunc("/", webHandler)
	http.ListenAndServe("localhost:4000", nil)
}