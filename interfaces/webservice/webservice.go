package webservice

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"srv-video-info/usecases"
)

type videoInfoInteractor interface {
	Get(filename string) (videoInfo usecases.VideoInfo, err error)
}

type webService struct {
	videoInfoInteractor videoInfoInteractor
}

func New(interactor videoInfoInteractor) *webService {
	service := webService{
		videoInfoInteractor: interactor,
	}

	return &service
}

func uniq() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatal(err)
	}
	uuid := fmt.Sprintf("%x-%x-%x-%x-%x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
	return uuid
}

func (h webService) VideoInfo(w http.ResponseWriter, r *http.Request) {
	log.Println("New request", r.RemoteAddr)
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		h.sendResponse(w, err.Error(), http.StatusInternalServerError)
		fmt.Println(err)
		return
	}
	file, _, err := r.FormFile("file")
	if err != nil {
		h.sendResponse(w, err.Error(), http.StatusInternalServerError)
		fmt.Println(err)
		return
	}
	defer file.Close()

	filename := os.TempDir() + string(os.PathSeparator) + uniq()

	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		h.sendResponse(w, err.Error(), http.StatusInternalServerError)
		fmt.Println(err)
		return
	}
	defer f.Close()
	defer os.Remove(filename)

	_, err = io.Copy(f, file)
	if err != nil {
		h.sendResponse(w, err.Error(), http.StatusInternalServerError)
		fmt.Println(err)
		return
	}
	f.Close()

	info, err := h.videoInfoInteractor.Get(filename)
	if err != nil {
		h.sendResponse(w, err.Error(), http.StatusInternalServerError)
		fmt.Println(err)
		return
	}
	h.sendResponse(w, info, http.StatusOK)
}

func (h *webService) sendResponse(w http.ResponseWriter, response interface{}, code int) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(code)
	resByte, err := json.Marshal(response)
	if err != nil {
		log.Println(err)
		return
	}
	h.send(w, resByte)
}

func (h *webService) send(w http.ResponseWriter, response []byte) {
	if _, err := w.Write(response); err != nil {
		log.Println("error send response " + err.Error())
	}
}
