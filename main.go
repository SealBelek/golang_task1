package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	server := Server{startStorageManager()}
	http.HandleFunc("/get", server.get)
	http.HandleFunc("/put", server.set)
	http.HandleFunc("/delete", server.delete)
	log.Fatal(http.ListenAndServe(":8000", nil))
}

func (s *Server) get(w http.ResponseWriter, r *http.Request) {
	r.Context()
	val, _ := s.storageService.Get(r.Context(), r.URL.Query().Get("key"))
	fmt.Println(val)
}

func (s *Server) set(w http.ResponseWriter, r *http.Request) {
	r.Context()
	s.storageService.Put(
		r.Context(),
		r.URL.Query().Get("key"),
		r.URL.Query().Get("val"))
	fmt.Println("ok")
}

func (s *Server) delete(w http.ResponseWriter, r *http.Request) {
	r.Context()
	s.storageService.Delete(r.Context(), r.URL.Query().Get("key"))
	fmt.Println("ok")
}

type Server struct {
	storageService StorageService
}
