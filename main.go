package main

import (
	"log"
	"net/http"
	"testeff/internal/handlers"

	"github.com/gorilla/mux"
)

func main() {
	st := handlers.Init()
	r := mux.NewRouter()
	r.HandleFunc("/delete", st.Delete).Methods(http.MethodPost)
	r.HandleFunc("/update", st.Update).Methods(http.MethodPost)
	r.HandleFunc("/add", st.Add).Methods(http.MethodPost)
	r.HandleFunc("/data", st.Data).Methods(http.MethodPost)
	r.HandleFunc("/text", st.TextSong).Methods(http.MethodPost)

	log.Println("server is running")
	err := http.ListenAndServe(st.URL, r)

	if err != nil {
		panic(err)
	}
}
