package main

import (
	"net/http"
	"testeff/internal/handlers"

	"github.com/gorilla/mux"
	"github.com/labstack/gommon/log"
)

func main() {
	st := handlers.Init()
	r := mux.NewRouter()
	r.HandleFunc("/delete-song", st.DeleteSong).Methods(http.MethodPost)
	r.HandleFunc("/update-song", st.Update).Methods(http.MethodPost)
	r.HandleFunc("/add-song", st.Add).Methods(http.MethodPost)
	r.HandleFunc("/get-songs", st.GetSongs).Methods(http.MethodPost)
	r.HandleFunc("/get-song-text", st.TextSong).Methods(http.MethodPost)

	log.Info("server is running")
	log.Debug(st.URL)
	err := http.ListenAndServe(st.URL, r)

	if err != nil {
		log.Fatal(err)
	}
}
