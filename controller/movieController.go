package controller

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"
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

	// Menambahkan kondisi pencarian berdasarkan parameter show_time
	// if showTime := c.Query("show_time"); showTime != "" {
	// 	query += " AND show_time = ?"
	// 	params = append(params, showTime)
	// }

	// // Menambahkan kondisi pencarian berdasarkan parameter branch
	// if branch := c.Query("branch"); branch != "" {
	// 	query += " AND branch = ?"
	// 	params = append(params, branch)
	// }

	// // Menambahkan kondisi pencarian berdasarkan parameter rating
	// if rating := c.Query("rating"); rating != "" {
	// 	query += " AND rating = ?"
	// 	params = append(params, rating)
	// }

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

	// Ensure the database connection is closed when the function returns
	defer db.Close()

	// Execute a SELECT query to retrieve all movies from the database
	rows, err := db.Query("SELECT * FROM movie")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve movies from database"})
		return
	}

	// Iterate over the rows returned from the query and store them in a slice of movie structs
	var movie []models.Movie
	for rows.Next() {
		var branch models.Movie
		err := rows.Scan(&movie.ID, &movie.Name, &movie.Address)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve movie from database"})
			return
		}

		// Execute a SELECT query to retrieve all theatres for the current movie
		theatreRows, err := db.Query("SELECT id, name FROM theatre WHERE movie_id = ?", movie.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error1": err})
			return
		}

		// Create an empty slice of Theatre structs to hold the retrieved theatres
		var theatres []models.Theatre

		// Iterate over the theatre rows returned from the query and store them in the theatres slice
		for theatreRows.Next() {
			var theatre models.Theatre
			err := theatreRows.Scan(&theatre.ID, &theatre.Name)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error2": fmt.Sprintf("%v", err)})
				return
			}
			theatres = append(theatres, theatre)
		}

		// Add the retrieved theatres to the current branch
		branch.Theatres = &theatres

		// Add the current branch to the branches slice
		branches = append(branches, branch)
	

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
