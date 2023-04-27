package controller

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
	"tix-id/config"
	"tix-id/models"
	"tix-id/tool"

	"github.com/gin-gonic/gin"
)

// GetMovies godoc
// @Summary Get Movies
// @Description Get All Movies sorted by realse date
// @Tags Guest
// @Param show_time query string false "Filter movies by show time"
// @Param branch query string false "Filter movies by branch"
// @Param rating query string false "Filter movies by rating"
// @Accept json
// @Produce json
// @Success 200 {object} models.MoviesResponse
// @Router /movies [get]

func GetMovies(c *gin.Context) {
	db := config.ConnectDB()
	defer db.Close()

	var movies []models.Movie
	params := []interface{}{}
	redisKey := "movies"
	redisClient := tool.NewRedisClient()

	query := "select m.id, m.title, m.description, m.duration, m.rating, m.release_date from movie m where 1 = 1"
	// check if there is params show_time
	if showTime := c.Query("show_time"); showTime != "" {
		redisKey += showTime + ":"
		query += " AND DATE(m.show_time) = ?"
		params = append(params, showTime)
	}

	// // check if there is params rating
	if rating := c.Query("rating"); rating != "" {
		redisKey += rating + ":"
		query += " AND m.rating > ?"
		params = append(params, rating)
	}

	moviesCache, err := tool.GetRedisValue(redisClient, redisKey)
	if err == nil {
		fmt.Println("Mengambil data dari redis")
		err := json.Unmarshal([]byte(moviesCache), &movies)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unmarshal cached data"})
			return
		}
	} else {

		rows, err := db.Query(query, params...)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
			return
		}

		noData := true
		for rows.Next() {
			noData = false
			var movie models.Movie
			if err := rows.Scan(&movie.ID, &movie.Title, &movie.Description, &movie.Duration, &movie.Rating, &movie.ReleaseDate); err != nil {
				log.Println(err)
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			movies = append(movies, movie)
		}

		if noData {
			response := models.Response{
				Status:  404,
				Message: "No movie found!",
			}
			c.JSON(http.StatusNotFound, response)
			return
		}

		moviesJSON, err := json.Marshal(movies)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal movie data"})
			return
		}
		err = tool.SetRedisValue(redisClient, redisKey, string(moviesJSON), 10*time.Second)
	}

	responseData := models.MoviesResponse{
		Response: models.Response{
			Status:  200,
			Message: "Movies retrieved successfully",
		},
		Movies: movies,
	}

	// Send response
	c.JSON(http.StatusOK, responseData)
}

// SearchMovies godoc
// @Summary Search movies by title and genre
// @Description Search movies by title and genre
// @Tags Guest
// @Accept json
// @Produce json
// @Param title query string false "Movie title to search"
// @Param genre query string false "Movie genre to search"
// @Success 200 {object} models.MoviesResponse
// @Router /movies/search [get]
func SearchMovies(c *gin.Context) {
	db := config.ConnectDB()
	defer db.Close()

	var movies []models.Movie
	// title := c.Query("title")
	// genre := c.Query("genre")
	// if title != "" && genre != "" {
	// 	models.DB.Where("title LIKE ? AND genre LIKE ?", "%"+title+"%", "%"+genre+"%").Find(&movies)
	// } else if title != "" {
	// 	models.DB.Where("title LIKE ?", "%"+title+"%").Find(&movies)
	// } else if genre != "" {
	// 	models.DB.Where("genre LIKE ?", "%"+genre+"%").Find(&movies)
	// } else {
	// 	c.JSON(http.StatusBadRequest, models.Response{
	// 		Status:  http.StatusBadRequest,
	// 		Message: "Please provide either title or genre to search for movies",
	// 	})
	// 	return
	// }

	responseData := models.MoviesResponse{
		Response: models.Response{
			Status:  http.StatusOK,
			Message: "Movies retrieved successfully",
		},
		Movies: movies,
	}

	c.JSON(http.StatusOK, responseData)
}

// GetMovieById godoc
// @Summary Get a movie by ID
// @Description Get a movie by ID
// @Tags Guest
// @Accept json
// @Produce json
// @Param movie_id path string true "Movie ID"
// @Success 200 {object} models.Movie
// @Router /movies/{movie_id} [get]
func GetMovieById(c *gin.Context) {
	// Connect to database
	db := config.ConnectDB()

	// Ensure the database connection is closed when the function returns
	defer db.Close()
	// movieId := c.Query("movieId")
	movieId, err := strconv.Atoi(c.Param("movieId"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve branches from database"})
		return
	}
	var movie models.Movie
	err = db.QueryRow("Select title,description,duration,rating,release_date from movie where id =?", movieId).Scan(&movie.Title, &movie.Description, &movie.Duration, &movie.Rating, &movie.ReleaseDate)
	if err != nil {
		if err == sql.ErrNoRows {
			response := models.Response{
				Status:  404,
				Message: "the movie is not found!",
			}
			c.JSON(http.StatusNotFound, response)
			return
		}
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	responseData := models.MovieResponse{
		Response: models.Response{
			Status:  200,
			Message: "Movie retrieved successfully",
		},
		Movie: movie,
	}
	c.JSON(http.StatusOK, responseData)
}

// CreateMovie godoc
// @Summary Create Movie
// @Description Create a new movie
// @Tags Admin
// @Accept json
// @Produce json
// @Param body body models.Movie true "Movie details"
// @Success 201 {object} models.MovieResponse
// @Router /movies [post]
func CreateMovie(c *gin.Context) {
	// Parse request body to Movie struct
	var movie models.Movie
	if err := c.ShouldBindJSON(&movie); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Connect to database
	db := config.ConnectDB()

	// Ensure the database connection is closed when the function returns
	defer db.Close()

	// Insert the movie into the database
	result, err := db.Exec("INSERT INTO movie (title, description, duration, rating, release_date) VALUES (?, ?, ?, ?, ?)", movie.Title, movie.Description, movie.Duration, movie.Rating, movie.ReleaseDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	// Get the ID of the inserted movie
	id, err := result.LastInsertId()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get ID of inserted movie"})
		return
	}

	// Set the ID of the movie to the inserted ID
	movie.ID = int(id)

	// Create a MovieResponse struct with the inserted movie and send it as a JSON response
	responseData := models.MovieResponse{
		Response: models.Response{
			Status:  200,
			Message: "Movie inserted successfully",
		},
		Movie: movie,
	}
	c.JSON(http.StatusOK, responseData)
}

// UpdateMovie godoc
// @Summary Update Movie
// @Description Update an existing movie
// @Tags Admin
// @Accept json
// @Produce json
// @Param movieId path string true "Movie ID"
// @Param body body models.Movie true "Movie details"
// @Success 200 {object} models.MovieResponse
// @Router /movies/{movieId} [put]
func UpdateMovie(c *gin.Context) {
	// Connect to database
	db := config.ConnectDB()

	// Ensure the database connection is closed when the function returns
	defer db.Close()
	// Get the movie ID from the request URL parameter
	movieID, err := strconv.Atoi(c.Param("movieId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid movie ID"})
		return
	}

	// Bind the updated movie data from the request body to the movie variable
	var movie models.Movie
	if err := c.ShouldBindJSON(&movie); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	movie.ID = movieID
	// Update movie in the database
	result, err := db.Exec("UPDATE movie SET title=?, description=?, duration=?, rating=?, release_date=? WHERE id=?", movie.Title, movie.Description, movie.Duration, movie.Rating, movie.ReleaseDate, movie.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Movie not found"})
		return
	}

	responseData := models.MovieResponse{
		Response: models.Response{
			Status:  200,
			Message: "Movie updated successfully",
		},
		Movie: movie,
	}
	c.JSON(http.StatusOK, responseData)
}

// DeleteMovie godoc
// @Summary Delete a movie
// @Description Delete a movie by ID
// @Tags Admin
// @Param movieId path string true "Movie ID"
// @Success 204 {object} models.Response
// @Router /movies/{movieId} [delete]
func DeleteMovie(c *gin.Context) {
	db := config.ConnectDB()
	defer db.Close()

	id := c.Param("movieId")
	log.Println("id: ", id)
	// Check if the movie exists
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM movie WHERE id=?", id).Scan(&count)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if count == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Movie not found"})
		return
	}

	// Delete the movie
	result, err := db.Exec("DELETE FROM movie WHERE id=?", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if rowsAffected == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete movie"})
		return
	}

	log.Printf("Movie with id %s has been deleted", id)

	responseData := models.Response{
		Status:  200,
		Message: "Movie deleted successfully",
	}
	c.JSON(http.StatusOK, responseData)
}
