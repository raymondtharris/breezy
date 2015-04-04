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



type loginCredentials struct{
	Username string
	Password string
}

func (l loginCredentials) String() string{
	return fmt.Sprintf("{username: %s, password: %s}", l.Username, l.Password)
}

func loadPage(pageName string) (*Page, error){
	pageBody, err := ioutil.ReadFile(pageName+".html")
	if err != nil {
	        return nil, err
	    }
	return &Page{Title:pageName, Body: pageBody }, nil
}

func webBlogHandler(w http.ResponseWriter, r *http.Request){
	p,_ := loadPage("index")
	t, _ := template.ParseFiles("../src/index.html")
	t.Execute(w, p)
}

func breezyLoginHandler(w http.ResponseWriter, r *http.Request){
	http.ServeFile(w, r, "../src/views/login.html")
}
type test struct{
	what string
}

func breezyLoginCredentrials(w http.ResponseWriter, r *http.Request){
	//defer r.Body.Close()
	body, err2 := ioutil.ReadAll(r.Body)
	if err2 != nil{
		panic(err2)
	}
	
	var vls loginCredentials
	err2 = json.Unmarshal([]byte(string(body[:])), &vls)
	if err2 != nil{
		panic(err2)
	}
	fmt.Println(string(body[:]) ,"\n", vls)
	
	w.Write([]byte("OK"))
}

func breezyEditHandler(w http.ResponseWriter, r *http.Request){
	http.ServeFile(w, r, "../src/views/edit.html")
}

func breezyDashboardHandler(w http.ResponseWriter, r *http.Request){
	http.ServeFile(w, r, "../src/views/dashboard.html")
}

func breezySettingsHandler(w http.ResponseWriter, r *http.Request){
	http.ServeFile(w, r, "../src/views/settings.html")
}



func HandleDirs(){
	http.Handle("/lib/", http.StripPrefix("/lib/", http.FileServer(http.Dir("../src/lib/"))))
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("../src/js/"))))
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("../src/css/"))))
	http.Handle("/views/", http.StripPrefix("/views/", http.FileServer(http.Dir("../src/views/"))))
}



func main(){
	HandleDirs()
	http.HandleFunc("/admin", breezyLoginHandler)
	http.HandleFunc("/checkcredentials", breezyLoginCredentrials)
	
	
	http.HandleFunc("/edit", breezyEditHandler)
	http.HandleFunc("/dashboard", breezyDashboardHandler)
	http.HandleFunc("/settings", breezySettingsHandler)
	http.HandleFunc("/", webBlogHandler)
	http.ListenAndServe("localhost:4000", nil)
}