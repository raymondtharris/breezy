package main

import (
	"fmt"
	//"log"
	"net/http"
)

const brImage, brAudio , brVideo = 0, 1, 2

type breezyImage struct{
	name, filename string
	fileSize float32
}

type breezyAudio struct{
	name, filename string
	fileSize float32
}

type breezyVideo struct{
	name, filename string
	fileSize float32
}

func (m breezyImage) String() string{
	return fmt.Sprintf("%v", m.name)
}
type breezyMedia struct{
	name filename string
	filesSize float32
	brType int
}

func main(){
	http.Handle("/string", String("I'm all good."))
	http.Handle("/struct", &Struct{"Hello", ":", "Gophers!"})
	http.ListenAndServe("localhost:4000", nil)
}