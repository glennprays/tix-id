package controller

import (
	"net/http"
	"tix-id/config"
	"tix-id/models"

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
	var movies []models.Movie

	// showTime := c.Query("show_time")
	// branch := c.Query("branch")
	// rating := c.Query("rating")

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
	// TODO: Implement logic to search movies by title and genre
	var movies []models.Movie

	responseData := models.MoviesResponse{
		Response: models.Response{
			Status:  200,
			Message: "Searched movies retrieved successfully",
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
	// movieId := c.Query("movieId")
	var movie models.Movie

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
	// movieId := c.Query("movieId")
	var movie models.Movie
	if err := c.ShouldBindJSON(&movie); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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
	// movieId := c.Query("movieId")
	responseData := models.Response{
		Status:  200,
		Message: "Movie deleted successfully",
	}
	c.JSON(http.StatusOK, responseData)
}
