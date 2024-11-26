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
	r.HandleFunc("/delete-song", st.DeleteSong).Methods(http.MethodPost)
	r.HandleFunc("/update-song", st.Update).Methods(http.MethodPost)
	r.HandleFunc("/add-song", st.Add).Methods(http.MethodPost)
	r.HandleFunc("/get-songs", st.GetSong).Methods(http.MethodPost)
	r.HandleFunc("/get-song-text", st.TextSong).Methods(http.MethodPost)

	log.Println("server is running")
	err := http.ListenAndServe(st.URL, r)

	if err != nil {
		panic(err)
	}
}
