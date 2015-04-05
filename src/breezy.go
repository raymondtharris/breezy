package main

import (
	"fmt"
	"encoding/json"
	"net/http"
	//"html/template"
	"io/ioutil"
	"strings"
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


type brBlogContent struct{
	MarkdownContent string
	MarkupContent string
}


type loginCredentials struct{
	Username string
	Password string
}

func (l loginCredentials) String() string{
	return fmt.Sprintf("{username: %s, password: %s}", l.Username, l.Password)
}


func webBlogHandler(w http.ResponseWriter, r *http.Request){
	http.ServeFile(w, r, "../src/index.html")
}

func breezyLoginHandler(w http.ResponseWriter, r *http.Request){
	http.ServeFile(w, r, "../src/views/login.html")
}

func breezyLoginCredentrials(w http.ResponseWriter, r *http.Request){
	body, err := ioutil.ReadAll(r.Body)
	if err != nil{
		panic(err)
	}
	
	var vls loginCredentials
	err = json.Unmarshal([]byte(string(body[:])), &vls)
	if err != nil{
		panic(err)
	}
	fmt.Println(string(body[:]) ,"\n", vls)
	
	w.Write([]byte("OK"))
}

func breezyEditHandler(w http.ResponseWriter, r *http.Request){
	http.ServeFile(w, r, "../src/views/edit.html")
}

func breezyMarkdownHandler(w http.ResponseWriter, r*http.Request){
	body, err := ioutil.ReadAll(r.Body)
	if err != nil{
		panic(err)
	}
	var currentBlogContent brBlogContent
	err = json.Unmarshal([]byte(string(body[:])), &currentBlogContent)
	if err != nil{
		panic(err)
	}
	//fmt.Println(string(body[:]) ,"\n", currentBlogContent)
	currentBlogContent = markdownConverter(currentBlogContent)
	//fmt.Println(currentBlogContent)
	// Send blog data back
	jsRes, err2 := json.Marshal(currentBlogContent)
	if err2 != nil{
		panic(err2)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsRes)
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

func markdownConverter(br brBlogContent) brBlogContent{
	br.MarkupContent = ""
	//strings.Replace(br.MarkupContent, "<br>", "\n", -1)
	arr := strings.Split(br.MarkdownContent, "\n")
	for i := 0; i< len(arr); i++{
		if(len(arr[i])) >0{
			//fmt.Println(arr[i])
			br.MarkupContent = br.MarkupContent + markdownConvertLine(arr[i])
		}
	}
	
	//fmt.Println(br.MarkupContent)
	return br
}

func markdownConvertLine(currentLine string) string{
	arr := strings.Split(currentLine, " ")
	switch arr[0]{
		case "#":
			currentLine = "<h1>"+currentLine+"</h1>"
		case "##":
			currentLine = "<h2>"+currentLine+"</h2>"
		case "###":
			currentLine = "<h3>"+currentLine+"</h3>"
		default:
			currentLine = "<p>"+currentLine+"</p>"
	}
	return currentLine
}


func main(){
	HandleDirs()
	http.HandleFunc("/admin", breezyLoginHandler)
	http.HandleFunc("/checkcredentials", breezyLoginCredentrials)
	
	
	http.HandleFunc("/edit", breezyEditHandler)
	http.HandleFunc("/mdowntomup",breezyMarkdownHandler)
	
	http.HandleFunc("/dashboard", breezyDashboardHandler)
	http.HandleFunc("/settings", breezySettingsHandler)
	http.HandleFunc("/", webBlogHandler)
	http.ListenAndServe("localhost:4000", nil)
}