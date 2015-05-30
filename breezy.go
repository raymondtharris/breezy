package main

import (
	"breezy/breezynlp"
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

//prefix br is from Client side
//prefix breezy is for data that goes to the database

const brImage, brAudio, brVideo = 0, 1, 2
const brPost = 0
const DB_URL = "45.55.192.173" //URL that points to MongoDB

type mgoSession struct {
	DB_Session *mgo.Session
	DB_Err     error
}

var mdbSession *mgo.Session

//var dbSession mgoSession

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
	Links  []string //an array of links found in a post
	Images []string //an array of image urls found within a post
	Audio  []string //an array of audio file urls found within a post
	Video  []string //an array of video file urls found within a post
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
	http.ServeFile(w, r, "index.html")
}

func breezyLoginHandler(w http.ResponseWriter, r *http.Request) {
	//Handler function to present the Breezy admin login HTML file
	http.ServeFile(w, r, "views/login.html")
}

func breezyLoginCredentrials(w http.ResponseWriter, r *http.Request) {
	//Handler function to take in user loginCredentials and verify info against database to determine if login is succesful
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	var userCred loginCredentials
	err = json.Unmarshal([]byte(string(body[:])), &userCred)
	if err != nil {
		panic(err)
	}
	//fmt.Println(string(body[:]), "\n", vls)
	fmt.Println(userCred.Password)
	hashedPassword, errP := bcrypt.GenerateFromPassword([]byte(userCred.Password), 10)
	fmt.Println(string(hashedPassword))
	_ = errP
	co := mdbSession.DB("test").C("Users")
	var res []breezyUser
	err = co.Find(bson.M{"username": userCred.Username}).All(&res)
	fmt.Println("Users Found:", res)
	passMatch := false
	for i := 0; i < len(res); i++ {
		errCheck := bcrypt.CompareHashAndPassword([]byte(res[i].Password), []byte(userCred.Password))
		if errCheck == nil {
			passMatch = true
			fmt.Println("found")
		}
	}
	//if res.Username == userCred.Username && res.Password == string(hashedPassword){
	if passMatch {
		w.Write([]byte("true"))
	} else {
		w.Write([]byte("false"))
	}
	//} else{
	//	w.Write([]byte("false"))
	//}
}

type brNewUser struct {
	Username string
	Password string
	Name     string
}

func breezyNewUserHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	var newUser brNewUser
	err = json.Unmarshal([]byte(string(body[:])), &newUser)
	if err != nil {
		panic(err)
	}
	cusers := mdbSession.DB("test").C("Users")
	var res []breezyUser
	err = cusers.Find(bson.M{"username": newUser.Username}).All(&res)
	if err != nil {
		panic(err)
	}
	fmt.Println("Users: ", res)
	if len(res) == 0 {
		fmt.Println("Create new user.")
		hashedPassword, errP := bcrypt.GenerateFromPassword([]byte(newUser.Password), 10)
		_ = errP
		newUserToSave := breezyUser{bson.NewObjectId(), newUser.Username, string(hashedPassword[:]), newUser.Name, "Editor", time.Now()}
		dbError := cusers.Insert(&newUserToSave)
		_ = dbError
		w.Write([]byte("New user created."))
	} else {
		w.Write([]byte("User with name already exists."))
	}
}

func breezyEditHandler(w http.ResponseWriter, r *http.Request) {
	//Handler function to present the Breezy Editor HTML file
	http.ServeFile(w, r, "views/edit.html")
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
	http.ServeFile(w, r, "views/dashboard.html")
}

func breezySettingsHandler(w http.ResponseWriter, r *http.Request) {
	//Handler function to present the settings HTML file
	http.ServeFile(w, r, "views/settings.html")
}

type breezySetupConfig struct {
	Username string
	Password string
	Name     string
	Blogname string
}

type breezySetupConfigDB struct {
	Username string
	Name     string
	Blogname string
}

type breezyUser struct {
	//breezyUser stores Userdata to be used on the databse
	ID       bson.ObjectId `bson:"_id,omitempty"` //ID variable for MongoDB
	Username string        //Username of user stored on DB
	Password string        //Hashed password of the user stored on DB
	Name     string        //Name of the user stored on DB
	Access   string        //Level of access of the user stored on DB
	Created  time.Time     //The timestamp of when the user was created stored on DB
}
type breezyBlog struct {
	//breezyBlog stores the relevant blog data to the database
	ID      bson.ObjectId `bson: "_id,omitempty"` //ID variable for MongoDB
	Name    string        //Name of the blog store on DB
	Creator string        //Name of the creator of the blog stored on DB
	Users   []string      //An Array of usernames that have access to the admin section of the blog stored on DB
	Created time.Time     //The timestamp of when the blog was created stored on DB
	Posts   int           //The number of posts stored on the DB
}

func (br breezySetupConfig) String() string {
	return fmt.Sprintf("Username: %v\nName: %v\nBlog Name: %v\n", br.Username, br.Name, br.Blogname)
}

func breezySetupConfigHandler(w http.ResponseWriter, r *http.Request) {
	//breezySetupConfigHandler function takes configuration data from user and sets up
	//blog on backend
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	var userConfig breezySetupConfig
	err = json.Unmarshal([]byte(string(body[:])), &userConfig)
	fmt.Println(userConfig.Password)
	hashedPassword, errP := bcrypt.GenerateFromPassword([]byte(userConfig.Password), 10)
	_ = errP
	fmt.Println(string(hashedPassword))
	//Create User
	var user = breezyUser{bson.NewObjectId(), userConfig.Username, string(hashedPassword[:]), userConfig.Name, "Admin", time.Now()}
	var dberr error
	co := mdbSession.DB("test").C("Users")
	fmt.Println(user)
	dberr = co.Insert(&user)
	_ = dberr
	var temp []breezyUser
	dberr = co.Find(nil).All(&temp)
	fmt.Println("User:", temp)

	//Creating and saving Blog data to Database
	conBlog := mdbSession.DB("test").C("Blog")
	userList := []string{userConfig.Username}
	var blogInfo = breezyBlog{bson.NewObjectId(), userConfig.Blogname, userConfig.Name, userList, time.Now(), 0}
	errB := conBlog.Insert(&blogInfo)
	_ = errB

	// Setup App directories
	breezyAppDirectories()

	//Create config file
	f, err := os.Stat("app/config.json")
	_ = f

	if err != nil {
		// write configuration
		os.Create("app/config.josn")
		writeConfiguration(userConfig)
	}

	//Create log file
	if _, err2 := os.Stat("../app/user/setup_log.json"); err2 != nil {
		fmt.Println("Creating Log File.")
		logFile, err := os.Create("../app/user/setup_log.json")
		_ = err
		_ = logFile
		//write stuff to log like creation date
		currentTime := time.Now()
		fmt.Println(currentTime)
		writeToLog(currentTime.String()+"\n", 0)
		fmt.Println(userConfig)
		var jsonString string
		jsonString = "{username:" + userConfig.Username + ", name:" + userConfig.Name + ", blogname:" + userConfig.Blogname + "}"
		jsonToWrite, err := json.Marshal(jsonString)
		writeToLog(string(jsonToWrite[:]), 0)
	}

}

func breezyAppDirectories() {
	//Checks if App directories exists and if not creates them.
	_, err := os.Stat("../app/")
	if err != nil {
		os.Mkdir("../app/", 0777)
	}
	_, err = os.Stat("../app/user")
	if err != nil {
		os.Mkdir("../app/user", 0777)
	}
	_, err = os.Stat("../app/user/backup")
	if err != nil {
		os.Mkdir("../app/user/backup", 0777)
	}
	_, err = os.Stat("../app/user/logs")
	if err != nil {
		os.Mkdir("../app/user/logs", 0777)
	}
	_, err = os.Stat("../app/user/media")
	if err != nil {
		os.Mkdir("../app/user/media", 0777)
	}

}

func writeConfiguration(userConfig breezySetupConfig) {
	fmt.Println("Writing User Configuration\n")
	fmt.Println("Configuration write complete")
}

func writeToLog(dataToWrite string, logNum int) {
	//writeToLog function writes data to log.json file
	fmt.Println("Writing Log File")
	if logNum == 0 {
		//write to setup_log
		f, err := os.OpenFile("../app/user/setup_log.json", os.O_APPEND|os.O_WRONLY, 0666)
		_ = err
		_ = f
		defer f.Close()
		temp, err := f.WriteString(dataToWrite)
		_ = temp
	} else {
		//write to weekly log
		currentTime := time.Now()
		//check for year directory
		_, err := os.Stat("../app/user/logs/" + strconv.Itoa(currentTime.Year()) + "/")
		if err != nil {
			os.Mkdir("../app/user/logs/"+strconv.Itoa(currentTime.Year())+"/", 0777)
		}
		//check for month directory
		m, err := os.Stat("../app/user/logs/" + strconv.Itoa(currentTime.Year()) + "/" + currentTime.Month().String() + "/")
		_ = m
		if err != nil {
			os.Mkdir("../app/user/logs/"+strconv.Itoa(currentTime.Year())+"/"+currentTime.Month().String()+"/", 0777)
		}
		//check for log file
		_, err2 := os.Stat("../app/user/logs/" + strconv.Itoa(currentTime.Year()) + "/" + currentTime.Month().String() + "/" + strconv.Itoa(currentTime.Day()) + "_log.json")
		if err2 != nil {
			os.Create("../app/user/logs/" + strconv.Itoa(currentTime.Year()) + "/" + currentTime.Month().String() + "/" + strconv.Itoa(currentTime.Day()) + "_log.json")
		}
		f, err := os.OpenFile("../app/user/logs/"+strconv.Itoa(currentTime.Year())+"/"+currentTime.Month().String()+"/"+strconv.Itoa(currentTime.Day())+"_log.json", os.O_APPEND|os.O_WRONLY, 0666)
		defer f.Close()
		t, err := f.WriteString(dataToWrite)
		_ = t
	}
}

func breezyBackupHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Creating Backup.")
	//Write to Log

	//Get all PostData
	//Write PostData to a Backup File
	//Add Media to backup that is not currently available there
}

type backupScheduleType struct {
	ScheduleType string
}

func breezyBackupScheduleHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		panic(err)
	}
	var backupType backupScheduleType
	err = json.Unmarshal([]byte(string(body[:])), &backupType)
	switch backupType.ScheduleType {
	case "Daily":
		fmt.Println("Backup Everyday")
	case "Weekly":
		fmt.Println("Backup Weekly")
	case "Monthly":
		fmt.Println("Backup Monthly")
	default:

	}
}

func BackupBlog(scheduleOption string) {
	dir, err := os.Stat("../app/user/backup/")
	if err != nil {
		//backup directory does not exist
		os.Mkdir("../app/user/backup/", 0777)
	}
	_ = dir
	currentTime := time.Now()
	y, err := os.Stat("../app/user/backup/" + strconv.Itoa(currentTime.Year()) + "/")
	_ = y
	if err != nil {
		os.Mkdir("../app/user/backup/"+strconv.Itoa(currentTime.Year())+"/", 0777)
	}
	m, err := os.Stat("../app/user/backup/" + strconv.Itoa(currentTime.Year()) + "/" + currentTime.Month().String() + "/")
	_ = m
	if err != nil {
		os.Mkdir("../app/user/backup/"+strconv.Itoa(currentTime.Year())+"/"+currentTime.Month().String()+"/", 0777)
	}
	if scheduleOption == "Monthly" {
		//write to month log
	} else {
		if scheduleOption == "Weekly" {
			//look for the weeks directory remaining within the month
		} else {
			d, err := os.Stat("../app/user/backup/" + strconv.Itoa(currentTime.Year()) + "/" + currentTime.Month().String() + "/" + strconv.Itoa(currentTime.Day()) + "/")
			_ = d
			if err != nil {
				os.Mkdir("../app/user/backup/"+strconv.Itoa(currentTime.Year())+"/"+currentTime.Month().String()+"/"+strconv.Itoa(currentTime.Day())+"/", 0777)
			}
			// Write Daily Backup

		}
	}

}

func HandleDirs() {
	//HandlerDirs sets up the handling of the other directories need to make Breezy work
	http.Handle("/lib/", http.StripPrefix("/lib/", http.FileServer(http.Dir("lib/"))))
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("js/"))))
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css/"))))
	http.Handle("/views/", http.StripPrefix("/views/", http.FileServer(http.Dir("views/"))))
	http.Handle("/app/", http.StripPrefix("/app/", http.FileServer(http.Dir("app/"))))
}

func markdownConverter(br brPostContent) brPostContent {
	//markdownConverter function runs through brPostContent MarkdownContent line by line to convert the data to usable HTML markup
	br.MarkupContent = ""

	arr := strings.Split(br.MarkdownContent, "\n")

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
				//br.MediaData.Images = append(br.MediaData.Images, brNewLine.convertedString)
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
				var urlString = markdownHandleURL(currentLine, false)
				fmt.Println(urlString)
			} else {
				//url is first element
				var urlString = markdownHandleURL(currentLine, true)
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

func markdownHandleURL(currentLine string, isSingle bool) markdownConvertedLine {
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
		if isSingle == true {
			returnString = returnString + "<br/><br/>"
		}
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

type breezyMarkdownMarkup struct {
	//Markdown and Markup storage to be used for the DB
	Markdown string // Markdown to be saved to the DB
	Markup   string // Converted markup to be saved to the DB
}
type breezyMediaCollection struct {
	//breezyMediaCollection stores all the media object in the DB
	Links  []string //An Array of links for in post to be stored on DB
	Images []string //An Array of Images in post to be stored on DB
	Audio  []string //An Array of Audio urls to be stored in the DB
	Video  []string //An Array of Video urls to be stored in the DB
}
type breezyPost struct {
	//breezyPost stores data needed for a post on the DB
	ID      bson.ObjectId         `bson:"_id,omitempty"` //ID variable for MongoDB
	Title   string                //Title of post to be stored on DB
	Created time.Time             //Created timestamp of the post stored on DB
	Updated time.Time             //Upated timestamp of the post to be stored on DB
	Creator string                //Creator of post stored on the DB
	Content breezyMarkdownMarkup  //Markdown and Markup of post to be stored on DB
	Media   breezyMediaCollection //Collection of links for posts to be stored on the DB
}
type brPostToSave struct {
	Title        string
	Content      breezyMarkdownMarkup
	Media        breezyMediaCollection
	ContentDirty bool
}

func breezySavePostHandler(w http.ResponseWriter, r *http.Request) {
	//make post variable
	//add to database
	//if save to db succesful send ok to client
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	var currentBlogContent brPostToSave
	err = json.Unmarshal([]byte(string(body[:])), &currentBlogContent)
	if err != nil {
		panic(err)
	}
	fmt.Println(currentBlogContent)
	coPosts := mdbSession.DB("test").C("Posts")
	currentTime := time.Now()
	newPost := breezyPost{bson.NewObjectId(), currentBlogContent.Title, currentTime, currentTime, "Tim", currentBlogContent.Content, currentBlogContent.Media}
	dbErr := coPosts.Insert(&newPost)
	_ = dbErr
	w.Write([]byte("Post Saved"))

}

func breezyPostListHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "views/posts.html")
}

func breezyPostDeleteHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.RequestURI)
	temp := strings.Split(r.RequestURI, "/")
	fmt.Println(temp[2])
	coPost := mdbSession.DB("test").C("Posts")
	err := coPost.Remove(bson.M{"_id": bson.ObjectIdHex(temp[2])})
	_ = err
	//fmt.Println(res)
	//coPost.RemoveId(id)
}

func breezyAllPostsHandler(w http.ResponseWriter, r *http.Request) {
	coPosts := mdbSession.DB("test").C("Posts")
	var res []breezyPost
	dbErr := coPosts.Find(nil).Sort("-created").All(&res)
	_ = dbErr
	//fmt.Println("Posts:", res)

	jsRes, err2 := json.Marshal(res)
	if err2 != nil {
		panic(err2)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsRes)
}

func breezyMediaListHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "views/media.html")
}

func breezyAllMediaHandler(w http.ResponseWriter, r *http.Request) {
	//breezyAllMediaHandler function gets all the media objects from the database and
	// then sends it back to the client as json array.
	coMedia := mdbSession.DB("test").C("Media")
	var res []breezyMediaObject
	dbErr := coMedia.Find(nil).Sort("-added").All(&res)
	_ = dbErr

	jsRes, err2 := json.Marshal(res)
	if err2 != nil {
		panic(err2)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsRes)
}

func breezyNewestMediaHandler(w http.ResponseWriter, r *http.Request) {
	// breezyNewestMediaHandler function gets the latest piece of media added to the database
	coMedia := mdbSession.DB("test").C("Media")
	var res breezyMediaObject
	dbErr := coMedia.Find(nil).Sort("-added").One(&res)
	_ = dbErr
	jsRes, err2 := json.Marshal(res)
	if err2 != nil {
		panic(err2)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsRes)
}

func breezyMediaDeleteHandler(w http.ResponseWriter, r *http.Request) {
	//breezyMediaDeleteHandler function looks at the request and finds the _id of a piece
	// of media and uses that to remove it from the database.
	fmt.Println(r.RequestURI)
	temp := strings.Split(r.RequestURI, "/")
	fmt.Println(temp[2])
	coMedia := mdbSession.DB("test").C("Media")

	// Find media file to remove on DB
	var mediTotFind breezyMediaObject
	dbErr := coMedia.Find(bson.M{"_id": bson.ObjectIdHex(temp[2])}).One(&mediTotFind)
	if dbErr != nil {
		panic(dbErr)
	}
	// Remove file on Server
	fiErr := os.Remove(".." + mediTotFind.FileWithPath)
	if fiErr != nil {
		panic(fiErr)
	}
	// Remove file record off the database
	err := coMedia.Remove(bson.M{"_id": bson.ObjectIdHex(temp[2])})
	_ = err
}

func breezySetupHandler(w http.ResponseWriter, r *http.Request) {
	d, err := os.Stat("app/user/setup_log.json")
	_ = d
	if err == nil {
		http.Redirect(w, r, "/admin", http.StatusFound)
	} else {
		http.ServeFile(w, r, "views/setup.html")
	}
}

type blogSettings struct {
	Name string // Name of the blog
	//SearchEnabled bool   //Stores if search is available to users
}

func breezyBlogInfoHandler(w http.ResponseWriter, r *http.Request) {
	coBlog := mdbSession.DB("test").C("Blog")
	var blog breezyBlog
	dbErr := coBlog.Find(nil).One(&blog)
	_ = dbErr
	blogData := blogSettings{blog.Name}
	jsRes, err := json.Marshal(blogData)
	if err != nil {
		panic(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsRes)
}

func breezyBlogInfoUpdateHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	var blogUpdateData blogSettings
	err = json.Unmarshal([]byte(string(body[:])), &blogUpdateData)
	if err != nil {
		panic(err)
	}
	coBlog := mdbSession.DB("test").C("Blog")
	dbErr := coBlog.Update(nil, bson.M{"name": blogUpdateData.Name})
	if dbErr != nil {
		panic(dbErr)
	}

}

type blogDisplayInfo struct {
	Title     string
	UserCount int
}

func breezyBlogDisplayInfoHandler(w http.ResponseWriter, r *http.Request) {
	coBlog := mdbSession.DB("test").C("Blog")
	var blog breezyBlog
	dbErr := coBlog.Find(nil).One(&blog)
	_ = dbErr
	displayInfo := blogDisplayInfo{blog.Name, len(blog.Users)}
	jsRes, err := json.Marshal(displayInfo)
	if err != nil {
		panic(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsRes)
}

type breezyMediaObject struct {
	ID           bson.ObjectId `bson:"_id,omitempty"` //ID variable for MongoDB
	Filename     string        // Filename string of the the file
	FileWithPath string        //Path and filename to the media object
	MediaType    string        //The type of media of the media object
	FileType     string        //The file extension of the media object
	Added        time.Time     //The timestamp of when the media object was added
}

func breezyFileUploadHandler(w http.ResponseWriter, r *http.Request) {
	//breezyFileUploadHandler function recieves file data from client and stores it in the appropriate location

	file, handler, err := r.FormFile("file")
	if err != nil {
		panic(err)
	}
	data, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}
	fileType := handler.Header.Get("Content-Type")
	_ = fileType
	currentTime := time.Now()
	path := mediaDirectoryCheck(currentTime)
	err = ioutil.WriteFile(path+handler.Filename, data, 0777)
	if err != nil {
		panic(err)
	}
	//Save file information to the database
	coMedia := mdbSession.DB("test").C("Media")
	objectTypes := strings.Split(fileType, "/")
	// Might need to change the saved path of the directory to the file.
	newMediaObject := breezyMediaObject{bson.NewObjectId(), handler.Filename, path[2:len(path)] + handler.Filename, objectTypes[0], objectTypes[1], currentTime}
	dbErr := coMedia.Insert(&newMediaObject)
	_ = dbErr
}

func mediaDirectoryCheck(currentTime time.Time) string {
	// Check if directories needed to write media exists and if not makes them.
	path := "../app/user/media/"
	_, err := os.Stat(path + strconv.Itoa(currentTime.Year()) + "/")
	if err != nil {
		os.Mkdir(path+strconv.Itoa(currentTime.Year())+"/", 0777)
	}
	path = path + strconv.Itoa(currentTime.Year()) + "/"
	_, err = os.Stat(path + currentTime.Month().String() + "/")
	if err != nil {
		os.Mkdir(path+currentTime.Month().String()+"/", 0777)
	}
	path = path + currentTime.Month().String() + "/"
	_, err = os.Stat(path + strconv.Itoa(currentTime.Day()) + "/")
	if err != nil {
		os.Mkdir(path+strconv.Itoa(currentTime.Day())+"/", 0777)
	}
	path = path + strconv.Itoa(currentTime.Day()) + "/"
	return path

}

type SearchInput struct {
	Searchtext string
}

func breezySearchHandler(w http.ResponseWriter, r *http.Request) {
	//fmt.Println(r.RequestURI)
	//temp := strings.Split(r.RequestURI, "/")
	//fmt.Println(temp[2])
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	var searchData SearchInput
	err = json.Unmarshal([]byte(string(body[:])), &searchData)
	if err != nil {
		panic(err)
	}
	fmt.Println(searchData.Searchtext)
}

func main() {
	var sErr error
	mdbSession, sErr = mgo.Dial(DB_URL)
	if sErr != nil {
		panic(sErr)
		//fmt.Println("Cannot connect to DB")
	}
	defer mdbSession.Close()

	//var tempNode breezynlp.BreezyNode
	tempNode := breezynlp.BreezyNode{1, "What", nil}
	fmt.Println(tempNode.Payload)

	HandleDirs()
	http.HandleFunc("/admin", breezyLoginHandler)
	http.HandleFunc("/checkcredentials", breezyLoginCredentrials)
	http.HandleFunc("/newuser", breezyNewUserHandler)

	http.HandleFunc("/edit", breezyEditHandler)
	http.HandleFunc("/mdowntomup", breezyMarkdownHandler)
	http.HandleFunc("/savepost", breezySavePostHandler)

	http.HandleFunc("/uploadfile", breezyFileUploadHandler)

	http.HandleFunc("/dashboard", breezyDashboardHandler)
	http.HandleFunc("/postlist", breezyPostListHandler)
	http.HandleFunc("/deletepost/", breezyPostDeleteHandler)
	http.HandleFunc("/get_all_posts", breezyAllPostsHandler)

	http.HandleFunc("/medialist", breezyMediaListHandler)
	http.HandleFunc("/deletemedia/", breezyMediaDeleteHandler)
	http.HandleFunc("/get_all_media", breezyAllMediaHandler)
	http.HandleFunc("/newest_media", breezyNewestMediaHandler)

	http.HandleFunc("/settings", breezySettingsHandler)
	http.HandleFunc("/blog_info", breezyBlogInfoHandler)
	http.HandleFunc("/blog_info_update", breezyBlogInfoUpdateHandler)

	http.HandleFunc("/backup", breezyBackupHandler)
	http.HandleFunc("/scheduledbackup", breezyBackupScheduleHandler)

	http.HandleFunc("/setup", breezySetupHandler)
	http.HandleFunc("/setup_config", breezySetupConfigHandler)
	http.HandleFunc("/", webBlogHandler)
	http.HandleFunc("/getsearch/", breezySearchHandler)
	http.HandleFunc("/get_blog_display", breezyBlogDisplayInfoHandler)
	http.ListenAndServe("localhost:4000", nil)
}
