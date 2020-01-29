package infrastructure

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func NewHttpServer(addr *string) *httpServer {
	server := &httpServer{
		router:  mux.NewRouter(),
		address: ":8080",
	}

	if addr != nil {
		server.address = *addr
	}

	return server
}

type httpServer struct {
	router  *mux.Router
	address string
}

// ServeHttp Вызов методов для тестов
func (h *httpServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.router.ServeHTTP(w, r)
}

// AddGet Добавления маршрута на GET метод
func (h *httpServer) AddGet(path string, handler func(w http.ResponseWriter, r *http.Request)) {
	h.router.HandleFunc(path, handler).Methods(http.MethodGet)
}

// AddPost Добавления маршрута на POST метод
func (h *httpServer) AddPost(path string, handler func(w http.ResponseWriter, r *http.Request)) {
	h.router.HandleFunc(path, handler).Methods(http.MethodPost)
}

// AddPost Добавления маршрута на OPTIONS метод
func (h *httpServer) AddOptions(path string, handler func(w http.ResponseWriter, r *http.Request)) {
	h.router.HandleFunc(path, handler).Methods(http.MethodOptions)
}

// ListenAndServe Запуск сервера
func (h httpServer) ListenAndServe() error {
	h.router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		log.Printf("call unknown url %s from %s\n", r.URL.Path, r.Host)
	})

	h.router.MethodNotAllowedHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)
		log.Printf("call unallowed method %s url %s from %s", r.Method, r.URL.Path, r.Host)
	})

	return http.ListenAndServe(h.address, h.router)
}
