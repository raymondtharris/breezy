package main

import (
	"fmt"
	"time"
	"net/http"
)

const brImage, brAudio , brVideo = 0, 1, 2
const brPost = 0

type breezyMedia struct{
	name, filename string
	filesSize float32
	brType int
}

func (m breezyMedia) String() string{
	var typeDescription
	switch m.brType{
	case brImage:
		typeDescription = "Image"
	case brAudio:
		typeDescription = "Audio"
	case brVideo:
		typeDescription = "Video"
		
	case default:
		typeDescription = ""
	}
	return fmt.Sprintf("%v, %v", m.name , typeDescription)
}

type breezyActivity struct{
	name string
	activityBody string
	mediaStructure [5]breezyMedia
	brType int
	dateCreated Time
}


func main(){
	//http.Handle("/string", String("I'm all good."))
	//http.Handle("/struct", &Struct{"Hello", ":", "Gophers!"})
	http.ListenAndServe("localhost:4000", nil)
}