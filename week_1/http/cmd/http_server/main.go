package main

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/go-chi/chi"
)

const (
	baseUrl       = "localhost:8081"
	createPostfix = "/notes"
	getPostfix    = "/notes/%id"
)

type NoteInfo struct {
	Title    string `json:"title"`
	Context  string `json:"context"`
	Author   string `json:"author"`
	IsPublic bool   `json:"is_public"`
}

type Note struct {
	ID        int64     `json:"id"`
	Info      NoteInfo  `json:"info"`
	CreatedAT time.Time `json:"created_at"`
	UpdatedAT time.Time `json:"updated_at"`
}

type SyncMap struct {
	elems map[int64]*Note
	m     sync.RWMutex
}

var notes = SyncMap{
	elems: make(map[int64]*Note),
}

func createNoteHandler(w http.ResponseWriter, r *http.Request) {
	info := &NoteInfo{}
	if err := json.NewDecoder(r.Body).Decode(info); err != nil {
		http.Error(w, "Failed to decode note data", http.StatusBadRequest)
		return
	}
	now := time.Now()
	nBig, err := rand.Int(rand.Reader, big.NewInt(time.Now().UnixNano()))
	if err != nil {
		panic(err)
	}
	note := &Note{
		ID:        nBig.Int64(),
		Info:      *info,
		CreatedAT: now,
		UpdatedAT: now,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(note); err != nil {
		http.Error(w, "Failed to encode note data", http.StatusInternalServerError)
		return
	}
	notes.m.Lock()
	defer notes.m.Unlock()
	notes.elems[note.ID] = note
}

// nBig, err := rand.Int(rand.Reader, big.NewInt(time.Now().UnixNano()))
func getNoteHandler(w http.ResponseWriter, r *http.Request) {
	noteID := chi.URLParam(r, "id")
	id, err := parseNoteID(noteID)
	if err != nil {
		http.Error(w, "Invalid note ID", http.StatusBadRequest)
		return
	}
	notes.m.RLock()
	defer notes.m.RUnlock()
	note, ok := notes.elems[id]
	if !ok {
		http.Error(w, "Note not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(note); err != nil {
		http.Error(w, "Failed to encode note data", http.StatusInternalServerError)
		return
	}
}

func parseNoteID(idStr string) (int64, error) {
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func main() {
	r := chi.NewRouter()
	r.Post(createPostfix, createNoteHandler)
	r.Get(getPostfix, getNoteHandler)

	fmt.Println("server starts on port 8081")

	err := http.ListenAndServe(baseUrl, r)
	if err != nil {
		log.Fatal(err)
	}

}
