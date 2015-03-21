package main

import (
	"fmt"
	//"time"
	"net/http"
	"html/template"
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
	Body []btye
	
}

func loadPage(pageName String) (*Page, error){
	pageBody, err := ioutil.ReadFile(pageName+".html")
	if err != nil {
	        return nil, err
	    }
	return &Page{Title:pageName, Body: pageBody }, nil
}

func webHandler(w http.ResponseWriter, r *http.Request){
	t, _ := template.ParseFiles("index.html")
	t.Execute(w, nil)
}


func main(){
	//http.Handle("/string", String("I'm all good."))
	http.HandleFunc("/", webHandler)
	http.ListenAndServe("localhost:4000", nil)
}