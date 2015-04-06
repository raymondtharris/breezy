package main

import (
	"fmt"
	"encoding/json"
	"net/http"
	"regexp"
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
	fmt.Println(arr[0])
	switch arr[0]{
		case "#":
			currentLine = strings.Replace(currentLine, "#", "<h1>", 1)+"</h1>"
		case "##":
			currentLine = strings.Replace(currentLine, "##", "<h2>", 1)+"</h2>"
		case "###":
			currentLine = strings.Replace(currentLine, "###", "<h3>", 1)+"</h3>"
		case "####":
			currentLine = strings.Replace(currentLine, "####", "<h4>", 1)+"</h4>"
		case "#####":
			currentLine = strings.Replace(currentLine, "#####", "<h5>", 1)+"</h5>"
		case "######":
			currentLine = strings.Replace(currentLine, "######", "<h6>", 1)+"</h6>"
		default:
			if strings.Contains(currentLine, "!["){
				if strings.Index(currentLine, "![") > 0 {
					//url inside string
					var urlString = markdownHandleURL(currentLine)
					fmt.Println(urlString)
				} else {
					//url is first element
					var urlString = markdownHandleURL(currentLine)
					currentLine = urlString
				}
			} else{
			currentLine = "<p>"+currentLine+"</p>"
			}
	}
	return currentLine
}

func markdownHandleURL(currentLine string) string{
	parenRegex := regexp.MustCompile("\\((.*?)\\)")
	altTextRegex := regexp.MustCompile("\\[(.*?)\\]")
	
	parens := parenRegex.FindAllString(currentLine, -1)
	altText := altTextRegex.FindAllString(currentLine, -1)
	
	splitParens := strings.Split(parens[0], " ")
	urlString := splitParens[0]
	urlTitle := splitParens[1]
	
	// determine if link or media
	postDotResult := strings.Split(urlString, ".")
	
	var returnString =""
	if strings.Contains(postDotResult[1], "png"){
		returnString = "<div><img src='"+urlString+"' alt='"+altText[0]+"' title="+urlTitle+"/></div>"
		returnString = removeLeftoversInLink(returnString)
	} else if strings.Contains(postDotResult[1], "jpg"){
		returnString = "<div><img src='"+urlString+"' alt='"+altText[0]+"' title="+urlTitle+"/></div>"
		returnString = removeLeftoversInLink(returnString)
	} else if strings.Contains(postDotResult[1], "jpeg"){
		returnString = "<div><img src='"+urlString+"' alt='"+altText[0]+"' title="+urlTitle+"/></div>"
		returnString = removeLeftoversInLink(returnString)
	} else if strings.Contains(postDotResult[1], "gif"){
		returnString = "<div><img src='"+urlString+"' alt='"+altText[0]+"' title="+urlTitle+"/></div>"
		returnString = removeLeftoversInLink(returnString)
	} else{
		returnString = "<a href='"+urlString+"' alt='"+altText[0]+"'>"+urlTitle+"</a>"
		returnString =  removeLeftoversInLink(returnString)	
	}
	
	
	return returnString
}

func removeLeftoversInLink(linkUrl string) string{
	temp1 := strings.Split(linkUrl,"(")
	linkUrl = temp1[0]+temp1[1]
	temp1 = strings.Split(linkUrl,")")
	linkUrl = temp1[0]+temp1[1]
	temp1 = strings.Split(linkUrl,"[")
	linkUrl = temp1[0]+temp1[1]
	temp1 = strings.Split(linkUrl,"]")
	linkUrl = temp1[0]+temp1[1]

	fmt.Println(linkUrl)
	return linkUrl
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