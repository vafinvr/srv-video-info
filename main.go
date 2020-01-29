package main

import (
	"log"
	"srv-video-info/infrastructure"
	"srv-video-info/interfaces/webservice"
	"srv-video-info/usecases"
)

func main() {
	port := ":4000"

	mediaInfo := infrastructure.NewMediaInfo()
	interactor := usecases.NewVideoInfoInteractor(mediaInfo)

	ws := webservice.New(interactor)

	hs := infrastructure.NewHttpServer(&port)
	hs.AddPost("/api/uploadVideo", ws.VideoInfo)

	log.Fatal(hs.ListenAndServe())

}
