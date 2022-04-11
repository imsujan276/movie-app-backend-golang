package handlers

import (
	"backend/models"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
)

type jsonResp struct {
	OK      bool   `json:"ok"`
	Message string `json:"message"`
}

func (m *Repository) GetAllMovies(w http.ResponseWriter, r *http.Request) {
	movies, err := m.App.Models.DB.All()
	if err != nil {
		m.errorJSON(w, err)
		return
	}
	err = m.writeJSON(w, http.StatusOK, movies, "movies")
	if err != nil {
		m.errorJSON(w, err)
		return
	}
}

func (m *Repository) GetOneMovie(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		m.errorJSON(w, err)
		return
	}
	movie, err := m.App.Models.DB.Get(id)
	err = m.writeJSON(w, http.StatusOK, movie, "movie")
	if err != nil {
		m.errorJSON(w, err)
		return
	}
}

func (m *Repository) GetAllGenres(w http.ResponseWriter, r *http.Request) {
	genres, err := m.App.Models.DB.GenresAll()
	if err != nil {
		m.errorJSON(w, err)
		return
	}
	err = m.writeJSON(w, http.StatusOK, genres, "genres")
	if err != nil {
		m.errorJSON(w, err)
		return
	}
}

func (m *Repository) GetAllMoviesByGenre(w http.ResponseWriter, r *http.Request) {
	genreID, err := strconv.Atoi(chi.URLParam(r, "genre_id"))
	if err != nil {
		m.errorJSON(w, err)
		return
	}
	movies, err := m.App.Models.DB.All(genreID)
	if err != nil {
		m.errorJSON(w, err)
		return
	}
	err = m.writeJSON(w, http.StatusOK, movies, "movies")
	if err != nil {
		m.errorJSON(w, err)
		return
	}
}

type MoviePayload struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Year        string `json:"year"`
	ReleaseDate string `json:"release_date"`
	RunTime     string `json:"run_time"`
	Rating      string `json:"rating"`
	MPAARating  string `json:"mpaa_rating"`
}

func (m *Repository) AddUpdateMovie(w http.ResponseWriter, r *http.Request) {
	var payload MoviePayload
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		m.errorJSON(w, err)
		return
	}

	var movie models.Movie

	if payload.ID != "0" {
		id, _ := strconv.Atoi(payload.ID)
		mov, _ := m.App.Models.DB.Get(id)
		movie = *mov
		movie.UpdatedAt = time.Now()
	}

	movie.ID, _ = strconv.Atoi(payload.ID)
	movie.Title = payload.Title
	movie.Description = payload.Description
	movie.RunTime, _ = strconv.Atoi(payload.RunTime)
	movie.ReleaseDate, _ = time.Parse("2006-01-02", payload.ReleaseDate)
	movie.Year = movie.ReleaseDate.Year()
	movie.Rating, _ = strconv.Atoi(payload.Rating)
	movie.MPAARating = payload.MPAARating
	movie.CreatedAt = time.Now()
	movie.UpdatedAt = time.Now()

	if movie.Poster == "" {
		movie = getPoster(movie)
	}

	if movie.ID == 0 {
		err = m.App.Models.DB.InsertMovie(movie)
	} else {
		err = m.App.Models.DB.UpdateMovie(movie)
	}
	if err != nil {
		m.errorJSON(w, err)
		return
	}

	operationType := "created"
	if movie.ID > 0 {
		operationType = "updated"
	}
	ok := jsonResp{
		OK:      true,
		Message: fmt.Sprintf("Movie %s successfully", operationType),
	}

	err = m.writeJSON(w, http.StatusOK, ok, "response")
	if err != nil {
		m.errorJSON(w, err)
		return
	}
}

func (m *Repository) DeleteMovie(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		m.errorJSON(w, err)
		return
	}
	err = m.App.Models.DB.DeleteMovie(id)
	if err != nil {
		m.errorJSON(w, err)
		return
	}
	ok := jsonResp{
		OK:      true,
		Message: "Movie deleted successfully",
	}

	err = m.writeJSON(w, http.StatusOK, ok, "response")
	if err != nil {
		m.errorJSON(w, err)
		return
	}
}

func (m *Repository) SearchMovies(w http.ResponseWriter, r *http.Request) {

}

func getPoster(movie models.Movie) models.Movie {
	type TheMovieDB struct {
		Page    int `json:"page"`
		Results []struct {
			Adult            bool    `json:"adult"`
			BackdropPath     string  `json:"backdrop_path"`
			GenreIds         []int   `json:"genre_ids"`
			ID               int     `json:"id"`
			OriginalLanguage string  `json:"original_language"`
			OriginalTitle    string  `json:"original_title"`
			Overview         string  `json:"overview"`
			Popularity       float64 `json:"popularity"`
			PosterPath       string  `json:"poster_path"`
			ReleaseDate      string  `json:"release_date"`
			Title            string  `json:"title"`
			Video            bool    `json:"video"`
			VoteAverage      float64 `json:"vote_average"`
			VoteCount        int     `json:"vote_count"`
		} `json:"results"`
		TotalPages   int `json:"total_pages"`
		TotalResults int `json:"total_results"`
	}

	client := &http.Client{}
	key := os.Getenv("TMDBAPiKey")
	theUrl := "https://api.themoviedb.org/3/search/movie?api_key="
	theFullUrl := theUrl + key + "&query=" + url.QueryEscape(movie.Title)

	req, err := http.NewRequest("GET", theFullUrl, nil)
	if err != nil {
		log.Println(err)
		return movie
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return movie
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return movie
	}

	var responseObject TheMovieDB
	json.Unmarshal(bodyBytes, &responseObject)

	if len(responseObject.Results) > 0 {
		movie.Poster = responseObject.Results[0].PosterPath
	}

	return movie
}
