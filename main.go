package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
)

// User represents a user in the system.
type User struct {
	ID         string     `json:"id"`
	SecretCode string     `json:"secretCode"`
	Name       string     `json:"name"`
	Email      string     `json:"email"`
	Playlists  []Playlist `json:"playlists"`
}

// Playlist represents a playlist.
type Playlist struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Songs []Song `json:"songs"`
}

// Song represents a song.
type Song struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Composers string `json:"composers"`
	MusicURL  string `json:"musicURL"`
}

var (
	users    []User
	usersMu  sync.Mutex
	nextID   = 1
	nextCode = 1000
)

// implementing all the Handlers

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var requestData struct {
		SecretCode string `json:"secretCode"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	usersMu.Lock()
	defer usersMu.Unlock()

	for _, user := range users {
		if user.SecretCode == requestData.SecretCode {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(user)
			return
		}
	}

	http.Error(w, "User not found", http.StatusNotFound)
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var newUser User

	if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	usersMu.Lock()
	defer usersMu.Unlock()

	// Generating unique ID and Secret Code as required
	newUser.ID = fmt.Sprintf("user%d", nextID)
	nextID++
	newUser.SecretCode = fmt.Sprintf("code%d", nextCode)
	nextCode++

	users = append(users, newUser)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newUser)
}

// ViewProfileHandler for returnig  all current playlists of a user.
func viewProfileHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Extracting user ID from request headers or tokens (implement authentication)
	userID := "user1" //actual user authentication logic required for multiple users

	usersMu.Lock()
	defer usersMu.Unlock()

	for _, user := range users {
		if user.ID == userID {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(user.Playlists)
			return
		}
	}

	http.Error(w, "User not found", http.StatusNotFound)
}

func getAllSongsOfPlaylistHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Extracting playlist ID from request (validate input)
	playlistID := r.URL.Query().Get("playlistID")
	if playlistID == "" {
		http.Error(w, "Missing playlistID parameter", http.StatusBadRequest)
		return
	}

	// Finding the user and playlist (implement authentication)
	userID := "user1" // actual user authentication logic to be implemented for multiple users

	usersMu.Lock()
	defer usersMu.Unlock()

	for _, user := range users {
		if user.ID == userID {
			for _, playlist := range user.Playlists {
				if playlist.ID == playlistID {
					w.Header().Set("Content-Type", "application/json")
					json.NewEncoder(w).Encode(playlist.Songs)
					return
				}
			}
		}
	}

	http.Error(w, "Playlist not found", http.StatusNotFound)
}

func createPlaylistHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Extracting user ID from request headers or tokens (implement authentication)
	userID := "user1" // actual user authentication logic to implemented for multiple users

	usersMu.Lock()
	defer usersMu.Unlock()

	for i, user := range users {
		if user.ID == userID {
			var newPlaylist Playlist
			if err := json.NewDecoder(r.Body).Decode(&newPlaylist); err != nil {
				http.Error(w, "Invalid request body", http.StatusBadRequest)
				return
			}

			// Generate a unique playlist ID
			newPlaylist.ID = fmt.Sprintf("playlist%d", len(user.Playlists)+1)

			// Add the new playlist to the user
			users[i].Playlists = append(users[i].Playlists, newPlaylist)

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(newPlaylist)
			return
		}
	}

	http.Error(w, "User not found", http.StatusNotFound)
}

// Adding a new  song to a playlist
func addSongToPlaylistHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Parseing request body to get playlist ID and song data
	var requestData struct {
		PlaylistID string `json:"playlistID"`
		Song       Song   `json:"song"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	usersMu.Lock()
	defer usersMu.Unlock()

	// Finding the user and playlist to which the song should be added
	for i, user := range users {
		for j, playlist := range user.Playlists {
			if playlist.ID == requestData.PlaylistID {
				// Appending the song to the playlist
				users[i].Playlists[j].Songs = append(playlist.Songs, requestData.Song)

				// Returning the updated playlist
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(users[i].Playlists[j])
				return
			}
		}
	}

	http.Error(w, "Playlist not found", http.StatusNotFound)
}

// implementing Delete a song Logic from a playlist
func deleteSongFromPlaylistHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Parseing the request body to get playlist ID and song ID
	var requestData struct {
		PlaylistID string `json:"playlistID"`
		SongID     string `json:"songID"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	usersMu.Lock()
	defer usersMu.Unlock()

	// Finding the user and playlist from which the song should be deleted
	for i, user := range users {
		for j, playlist := range user.Playlists {
			if playlist.ID == requestData.PlaylistID {
				// Finding the song index in the playlist
				for k, song := range playlist.Songs {
					if song.ID == requestData.SongID {
						// Removeing the song from the playlist
						users[i].Playlists[j].Songs = append(playlist.Songs[:k], playlist.Songs[k+1:]...)

						// Returning success message
						w.WriteHeader(http.StatusNoContent)
						return
					}
				}

				http.Error(w, "Song not found in playlist", http.StatusNotFound)
				return
			}
		}
	}

	http.Error(w, "Playlist not found", http.StatusNotFound)
}

// Deleteing a playlist
func deletePlaylistHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Parseing the request body to get playlist ID
	var requestData struct {
		PlaylistID string `json:"playlistID"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	usersMu.Lock()
	defer usersMu.Unlock()

	// Finding the user and playlist to delete
	for i, user := range users {
		for j, playlist := range user.Playlists {
			if playlist.ID == requestData.PlaylistID {
				// Removeing the playlist from the user's playlists
				users[i].Playlists = append(user.Playlists[:j], user.Playlists[j+1:]...)

				// Returning success message
				w.WriteHeader(http.StatusNoContent)
				return
			}
		}
	}

	http.Error(w, "Playlist not found", http.StatusNotFound)
}

// Getting the song details
func getSongDetailHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Parseing the query parameter to get song ID
	songID := r.URL.Query().Get("songID")

	usersMu.Lock()
	defer usersMu.Unlock()

	// Finding the user and playlist to delete
	for _, user := range users {
		for _, playlist := range user.Playlists {
			for _, song := range playlist.Songs {
				if song.ID == songID {
					// Returning the song details
					w.Header().Set("Content-Type", "application/json")
					json.NewEncoder(w).Encode(song)
					return
				}
			}
		}
	}

	http.Error(w, "Song not found", http.StatusNotFound)
}

func main() {
	// ...
	http.HandleFunc("/Login", loginHandler)
	http.HandleFunc("/Register", registerHandler)
	http.HandleFunc("/ViewProfile", viewProfileHandler)
	http.HandleFunc("/GetAllSongsOfPlaylist", getAllSongsOfPlaylistHandler)
	http.HandleFunc("/CreatePlaylist", createPlaylistHandler)
	http.HandleFunc("/AddSongToPlaylist", addSongToPlaylistHandler)
	http.HandleFunc("/DeleteSongFromPlaylist", deleteSongFromPlaylistHandler)
	http.HandleFunc("/DeletePlaylist", deletePlaylistHandler)
	http.HandleFunc("/GetSongDetail", getSongDetailHandler)

	fmt.Println("Server is running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
