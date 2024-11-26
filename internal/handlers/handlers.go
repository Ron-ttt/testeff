package handlers

import (
	"encoding/json"
	"net/http"
	"testeff/internal/config"
	"testeff/internal/db"

	"github.com/labstack/gommon/log"
)

func Init() start {
	localhost, dbAdress, api := config.Flags()
	db, err := db.NewDataBase(dbAdress)
	if err != nil {
		log.Fatal(err)
	}
	log.Info(localhost, dbAdress, api)
	return start{URL: localhost, database: db, api: api}
}

type start struct {
	URL      string
	database db.Storage
	api      string
}

func (st start) DeleteSong(res http.ResponseWriter, req *http.Request) {
	var song db.Song
	if err := json.NewDecoder(req.Body).Decode(&song); err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		log.Debug("deletesong err: ", err)
		return
	}
	log.Info(song)
	err := st.database.DeleteSong(song)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		log.Debug("deletesong err: ", err)
		return
	}
	res.WriteHeader(http.StatusOK)
}

func (st start) Update(res http.ResponseWriter, req *http.Request) {
	var song db.SongData
	if err := json.NewDecoder(req.Body).Decode(&song); err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		log.Debug("updatesong err: ", err)
		return
	}
	log.Info(song)
	err := st.database.UpdateSong(song)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		log.Debug("updatesong err: ", err)
		return
	}
	res.WriteHeader(http.StatusOK)
}

func (st start) TextSong(res http.ResponseWriter, req *http.Request) {
	var song db.Song
	if err := json.NewDecoder(req.Body).Decode(&song); err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		log.Debug("textsong err: ", err)
		return
	}
	log.Info(song)
	text, err := st.database.GetTextSong(song)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		log.Debug("textsong err: ", err)
		return
	}
	res.Header().Set("content-type", "application/json")
	if err := json.NewEncoder(res).Encode(text); err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		log.Debug("textsong err: ", err)
		return
	}
	res.WriteHeader(http.StatusOK)
}

func (st start) Add(res http.ResponseWriter, req *http.Request) {
	var song db.Song

	if err := json.NewDecoder(req.Body).Decode(&song); err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		log.Debug("addsong err: ", err)
		return
	}
	log.Info(song)
	resp, err := http.Get(st.api + "/info?group=" + song.Group + "&song=" + song.Song)

	if resp.StatusCode != http.StatusOK {
		http.Error(res, err.Error(), resp.StatusCode)
		log.Debug("addsong err: ", err)
		return
	}

	defer resp.Body.Close()
	var songData db.SongData
	if err := json.NewDecoder(resp.Body).Decode(&songData); err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		log.Debug("addsong err: ", err)
		return
	}
	log.Info(songData)
	err1 := st.database.AddSong(songData)
	if err1 != nil {
		http.Error(res, err1.Error(), http.StatusBadRequest)
		log.Debug("addsong err: ", err1)
		return
	}
	res.WriteHeader(http.StatusOK)
}

func (st start) GetSongs(res http.ResponseWriter, req *http.Request) {
	var typeSort db.Sort
	if err := json.NewDecoder(req.Body).Decode(&typeSort); err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		log.Debug("getsongs err: ", err)
		return
	}
	log.Info(typeSort)
	list, err := st.database.Info(typeSort)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		log.Debug("getsongs err: ", err)
		return
	}
	res.Header().Set("content-type", "application/json")
	if err := json.NewEncoder(res).Encode(list); err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		log.Debug("getsongs err: ", err)
		return
	}
	res.WriteHeader(http.StatusOK)
}
