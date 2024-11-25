package handlers

import (
	"encoding/json"
	"net/http"
	"testeff/internal/config"
	"testeff/internal/db"
)

func Init() start {
	localhost, dbAdress, api := config.Flags()
	db, _ := db.NewDataBase(dbAdress)
	return start{URL: localhost, database: db, api: api}
}

type start struct {
	URL      string
	database db.Storage
	api      string
}

func (st start) Delete(res http.ResponseWriter, req *http.Request) {
	var song db.Song
	if err := json.NewDecoder(req.Body).Decode(&song); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
	err := st.database.DeleteSong(song)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	res.WriteHeader(http.StatusOK)
}

func (st start) Update(res http.ResponseWriter, req *http.Request) {
	var song db.SongData
	if err := json.NewDecoder(req.Body).Decode(&song); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
	err := st.database.UpdateSong(song)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	res.WriteHeader(http.StatusOK)
}

func (st start) TextSong(res http.ResponseWriter, req *http.Request) {
	var song db.Song
	if err := json.NewDecoder(req.Body).Decode(&song); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
	text, err := st.database.GetTextSong(song)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	res.Header().Set("content-type", "application/json")
	if err := json.NewEncoder(res).Encode(text); err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	res.WriteHeader(http.StatusOK)
}

func (st start) Add(res http.ResponseWriter, req *http.Request) {
	var song db.Song

	if err := json.NewDecoder(req.Body).Decode(&song); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
	resp, err := http.Get(st.api + "/info?group=" + song.Group + "&song=" + song.Song)

	if resp.StatusCode == http.StatusInternalServerError {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	if resp.StatusCode == http.StatusBadRequest {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
	defer resp.Body.Close()
	var songData db.SongData
	if err := json.NewDecoder(resp.Body).Decode(&songData); err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	err1 := st.database.AddSong(songData)
	if err1 != nil {
		http.Error(res, err1.Error(), http.StatusInternalServerError)
		return
	}
	res.WriteHeader(http.StatusOK)
}

func (st start) Data(res http.ResponseWriter, req *http.Request) {
	var typeSort db.Sort
	if err := json.NewDecoder(req.Body).Decode(&typeSort); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
	list, err := st.database.Info(typeSort)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	res.Header().Set("content-type", "application/json")
	if err := json.NewEncoder(res).Encode(list); err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	res.WriteHeader(http.StatusOK)
}
