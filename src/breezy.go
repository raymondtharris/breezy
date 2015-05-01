package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
)

const brImage, brAudio, brVideo = 0, 1, 2
const brPost = 0
const brBlogName = "Temp"

var editingPost brPostContent

type breezyMedia struct {
	name, filename string
	filesSize      float32
	brType         int
}

func (m breezyMedia) String() string {
	var typeDescription = ""
	switch m.brType {
	case brImage:
		typeDescription = "Image"
	case brAudio:
		typeDescription = "Audio"
	case brVideo:
		typeDescription = "Video"

	}
	return fmt.Sprintf("%v, %v", m.name, typeDescription)
}

type breezyActivity struct {
	name           string
	activityBody   string
	mediaStructure [5]breezyMedia
	brType         int
	//dateCreated, dateModified Time
}

type brPostData struct {
	// brPostData stores data about a post eccept for the post content
	Title        string //Title of the post
	DateCreated  string //Date the post was created and saved
	DateModified string //Last Modification date of post that was saved
}

type brPostMediaData struct {
	// brPostMediaData stores urls found within a post for fast referencing
	Links  [10]string //an array of links found in a post
	Images [10]string //an array of image urls found within a post
	Audio  [4]string  //an array of audio file urls found within a post
	Video  [4]string  //an array of video file urls found within a post
}

func (mediaData brPostMediaData) String() string {
	//formated string to print data within brPostMediaData
	return fmt.Sprintf("Links:%v\nImages:%v\n", mediaData.Links, mediaData.Images)
}

type brPostContent struct {
	// brPostContent stores all post data plus the actual markdown and markup content
	MarkdownContent string          //Markdown version of the post
	MarkupContent   string          //Converted Markup of the Markdown for the post
	PostData        brPostData      //Store other post data extracted from editor
	MediaData       brPostMediaData //Stores Links and Other Media URLs of post
}

type loginCredentials struct {
	// loginCredentials stores the data needed to login to server
	Username string //username provided at login from user
	Password string //hash of the users password from login
}

func (l loginCredentials) String() string {
	//Formated string to print string of loginCredentials
	return fmt.Sprintf("{username: %s, password: %s}", l.Username, l.Password)
}

func webBlogHandler(w http.ResponseWriter, r *http.Request) {
	//Handler function to present the Blog HTML file
	http.ServeFile(w, r, "../src/index.html")
}

func breezyLoginHandler(w http.ResponseWriter, r *http.Request) {
	//Handler function to present the Breezy admin login HTML file
	http.ServeFile(w, r, "../src/views/login.html")
}

func breezyLoginCredentrials(w http.ResponseWriter, r *http.Request) {
	//Handler function to take in user loginCredentials and verify info against database to determine if login is succesful
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	var vls loginCredentials
	err = json.Unmarshal([]byte(string(body[:])), &vls)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(body[:]), "\n", vls)

	w.Write([]byte("OK"))
}

func breezyEditHandler(w http.ResponseWriter, r *http.Request) {
	//Handler function to present the Breezy Editor HTML file
	http.ServeFile(w, r, "../src/views/edit.html")
}

func breezyMarkdownHandler(w http.ResponseWriter, r *http.Request) {
	//Handler function to initiate the markdown to markup conversion of what is in the Breezy Editor
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	var currentBlogContent brPostContent
	err = json.Unmarshal([]byte(string(body[:])), &currentBlogContent)
	if err != nil {
		panic(err)
	}
	//fmt.Println(string(body[:]) ,"\n", currentBlogContent)
	currentBlogContent = markdownConverter(currentBlogContent)
	//fmt.Println(currentBlogContent)
	// Send blog data back
	jsRes, err2 := json.Marshal(currentBlogContent)
	if err2 != nil {
		panic(err2)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsRes)
}

func breezyDashboardHandler(w http.ResponseWriter, r *http.Request) {
	//Handler function to present the Dashboard HTML file
	http.ServeFile(w, r, "../src/views/dashboard.html")
}

func breezySettingsHandler(w http.ResponseWriter, r *http.Request) {
	//Handler function to present the settings HTML file
	http.ServeFile(w, r, "../src/views/settings.html")
}

func HandleDirs() {
	//HandlerDirs sets up the handling of the other directories need to make Breezy work
	http.Handle("/lib/", http.StripPrefix("/lib/", http.FileServer(http.Dir("../src/lib/"))))
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("../src/js/"))))
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("../src/css/"))))
	http.Handle("/views/", http.StripPrefix("/views/", http.FileServer(http.Dir("../src/views/"))))
}

func markdownConverter(br brPostContent) brPostContent {
	//markdownConverter function runs through brPostContent MarkdownContent line by line to convert the data to usable HTML markup
	br.MarkupContent = ""

	arr := strings.Split(br.MarkdownContent, "\n")
	var mediaDB brPostMediaData
	for i := 0; i < len(arr); i++ {
		if (len(arr[i])) > 0 {
			var brNewLine markdownConvertedLine
			brNewLine = markdownConvertLine(arr[i])
			br.MarkupContent = br.MarkupContent + brNewLine.convertedString
			//Check line conversionType to see if it needs to be added to br.PostData

			switch brNewLine.conversionType {
			case "Title":
				br.PostData.Title = brNewLine.convertedString
			case "Link":
				br.MediaData.Links = append(br.MediaData.Links, brNewLine.convertedString)
			case "Image":
				br.MediaData.Images = append(br.MediaData.Images, brNewLine.convertedString)
			default:

			}
			//br.MarkupContent = br.MarkupContent + markdownConvertLine(arr[i])
			fmt.Println(br.MediaData)
		}
	}

	//fmt.Println(br.MarkupContent)
	return br
}

type markdownConvertedLine struct {
	//markdownConvertedLine stores the data for eachline and what type of conversion was made
	convertedString string //String variable storing the line that was converted
	conversionType  string //String variable storing the type of conversion that was done on the line
}

func markdownConvertLine(currentLine string) markdownConvertedLine {
	//markdownConvertLine function takes in a string of the currently presented line and then determines the actions needed
	//to be done convert it from markdown to markup and finally returns the conerted string and the type of conversion that was
	//done on the line
	arr := strings.Split(currentLine, " ")
	fmt.Println(arr[0]) //printout of the first object in the array
	var convertedLine markdownConvertedLine

	switch arr[0] { // switch to determine what needs to be done to line based on first element in array
	//Convert for Headers H1-H6
	case "#":
		currentLine = strings.Replace(currentLine, "#", "<h1>", 1) + "</h1>"
		convertedLine.conversionType = "Title"
	case "##":
		currentLine = strings.Replace(currentLine, "##", "<h2>", 1) + "</h2>"
		convertedLine.conversionType = "H2"
	case "###":
		currentLine = strings.Replace(currentLine, "###", "<h3>", 1) + "</h3>"
		convertedLine.conversionType = "H3"
	case "####":
		currentLine = strings.Replace(currentLine, "####", "<h4>", 1) + "</h4>"
		convertedLine.conversionType = "H4"
	case "#####":
		currentLine = strings.Replace(currentLine, "#####", "<h5>", 1) + "</h5>"
		convertedLine.conversionType = "H5"
	case "######":
		currentLine = strings.Replace(currentLine, "######", "<h6>", 1) + "</h6>"
		convertedLine.conversionType = "H6"
	default:
		if strings.Contains(currentLine, "![") {
			//Check for link or media URL present
			if strings.Index(currentLine, "![") > 0 {
				//url inside string
				var urlString = markdownHandleURL(currentLine)
				fmt.Println(urlString)
			} else {
				//url is first element
				var urlString = markdownHandleURL(currentLine)
				currentLine = urlString.convertedString
				convertedLine.conversionType = urlString.conversionType
			}
		} else {
			//Create a paragraph
			currentLine = "<p>" + currentLine + "</p>"
			convertedLine.conversionType = "Text"
		}
	}

	convertedLine.convertedString = currentLine
	return convertedLine
}

func markdownHandleURL(currentLine string) markdownConvertedLine {
	//markdownHandleURL function pulls data out of the currentLine string and construct the appropriate
	//type of markup neeeded to display the URL
	parenRegex := regexp.MustCompile("\\((.*?)\\)")   //regexp to isolate the string within the parens
	altTextRegex := regexp.MustCompile("\\[(.*?)\\]") //regexp to isolate the string within the doublequotes

	parens := parenRegex.FindAllString(currentLine, -1)
	altText := altTextRegex.FindAllString(currentLine, -1)

	splitParens := strings.Split(parens[0], " ")
	urlString := splitParens[0]
	urlTitle := splitParens[1]

	// determine if link or media
	postDotResult := strings.Split(urlString, ".")
	var convertedURL markdownConvertedLine
	var returnString = ""
	// conditional if statements to determine if Image, Video, Audio, or Link
	if strings.Contains(postDotResult[1], "png") {
		returnString = "<div><img src='" + urlString + "' alt='" + altText[0] + "' title=" + urlTitle + "/></div>"
		returnString = removeLeftoversInLink(returnString)
		convertedURL.conversionType = "Image"
	} else if strings.Contains(postDotResult[1], "jpg") {
		returnString = "<div><img src='" + urlString + "' alt='" + altText[0] + "' title=" + urlTitle + "/></div>"
		returnString = removeLeftoversInLink(returnString)
		convertedURL.conversionType = "Image"
	} else if strings.Contains(postDotResult[1], "jpeg") {
		returnString = "<div><img src='" + urlString + "' alt='" + altText[0] + "' title=" + urlTitle + "/></div>"
		returnString = removeLeftoversInLink(returnString)
		convertedURL.conversionType = "Image"
	} else if strings.Contains(postDotResult[1], "gif") {
		returnString = "<div><img src='" + urlString + "' alt='" + altText[0] + "' title=" + urlTitle + "/></div>"
		returnString = removeLeftoversInLink(returnString)
		convertedURL.conversionType = "Image"
	} else {
		urlTitle = urlTitle[1:len(urlTitle)-2] + ")" // Gotta take out the doublequotes but leave in the parent to be removed later
		returnString = "<a href='" + urlString + "' alt='" + altText[0] + "'>" + urlTitle + "</a>"
		returnString = removeLeftoversInLink(returnString)
		convertedURL.conversionType = "Link"
	}

	convertedURL.convertedString = returnString
	return convertedURL

}

func removeLeftoversInLink(linkUrl string) string {
	//removeLeftoversInLink function finishes making appropriate url Markup string by taking out parens and square brakets
	temp1 := strings.Split(linkUrl, "(")
	linkUrl = temp1[0] + temp1[1]
	temp1 = strings.Split(linkUrl, ")")
	linkUrl = temp1[0] + temp1[1]
	temp1 = strings.Split(linkUrl, "[")
	linkUrl = temp1[0] + temp1[1]
	temp1 = strings.Split(linkUrl, "]")
	linkUrl = temp1[0] + temp1[1]

	//fmt.Println(linkUrl)
	return linkUrl
}
func breezySavePostHandler(w http.ResponseWriter, r *http.Request) {
	//make post variable
	//add to database
	//if save to db succesful send ok to client
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	var currentBlogContent brPostContent
	err = json.Unmarshal([]byte(string(body[:])), &currentBlogContent)
	if err != nil {
		panic(err)
	}
	fmt.Println(currentBlogContent)
}

func breezyPostListHandle(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "../src/views/posts.html")
}

func breezyFileUploadHandler(w http.ResponseWriter, r *http.Request) {
	//breezyFileUploadHandler function recieves file data from client and stores it in the appropriate location

}

func main() {
	HandleDirs()
	http.HandleFunc("/admin", breezyLoginHandler)
	http.HandleFunc("/checkcredentials", breezyLoginCredentrials)

	http.HandleFunc("/edit", breezyEditHandler)
	http.HandleFunc("/mdowntomup", breezyMarkdownHandler)
	http.HandleFunc("/savepost", breezySavePostHandler)

	http.HandleFunc("/uploadfile", breezyFileUploadHandler)

	http.HandleFunc("/dashboard", breezyDashboardHandler)
	http.HandleFunc("/settings", breezySettingsHandler)
	http.HandleFunc("/", webBlogHandler)
	http.ListenAndServe("localhost:4000", nil)
}
