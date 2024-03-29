package server

import (
	"fmt"
	"goflix/db"
	"goflix/middleware"
	"goflix/models"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type Server interface {
	Run()
}

type Serve struct {
	router *gin.Engine
	db     db.Storage
}

func New(db db.Storage) Server {
	gin.SetMode(gin.ReleaseMode)
	return &Serve{
		router: gin.Default(),
		db:     db,
	}
}

func (s *Serve) Run() {
	s.routes()
	s.router.Run(":4123")
}

func (s *Serve) routes() {

	s.router.GET("/", s.handelHello)
	s.router.POST("/login", s.handelLogin)
	s.router.POST("/users", s.handelAddUsers)

	// Routes for connected user
	s.router.Use(middleware.JwtMiddleware())

	s.router.GET("/users/:userID", s.handelGetUsers)
	s.router.DELETE("/users/:userID", s.handelDeleteUsers)
	s.router.PUT("/users/:userID", s.handelUpdateUsers)

	s.router.GET("/series/", s.handelGetListSeries)
	s.router.GET("/movies", s.handelGetListMovies)
	s.router.GET("/movies/:movieID", s.handelGetmovie)

	s.router.POST("/ratings", s.handelSaveRatingsUsers)
	s.router.GET("/ratings/:userID", s.handelGetRatingsUsers)

	s.router.POST("//favorites", s.handelSaveFavoriteUsers)
	s.router.GET("/favorites/:userID", s.handelGetFavoriteUsers)
	s.router.DELETE("/favorites/:userID/:favoriteID", s.handelDeleteFavoriteUsers)

	// Routes for admin user only
	s.router.Use(middleware.AdminOnly())

	s.router.POST("/movies/", s.handelAddMovies)
	s.router.DELETE("/movies/:movieID", s.handelDeleteMovies)

}

func (s *Serve) handelHello(c *gin.Context) {
	c.String(200, "Hello Goflix")
}

// * * *  USER * * *

func (s *Serve) handelAddUsers(c *gin.Context) {
	if user := s.decodeUserJSON(c); user != nil {
		err := s.db.SaveUser(user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "user saved"})
	}
}

func (s *Serve) handelGetUsers(c *gin.Context) {
	if id, err := s.getUserID(c); err == nil {
		user, err := s.db.GetUser(id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, user)
	}
}

func (s *Serve) handelDeleteUsers(c *gin.Context) {
	if id, err := s.getUserID(c); err == nil {
		err := s.db.DeleteUser(id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "user deleted"})
	}
}

func (s *Serve) handelUpdateUsers(c *gin.Context) {
	if user := s.decodeUserJSON(c); user != nil {
		if id, err := s.getUserID(c); err == nil {
			user.Id = id
		}
		err := s.db.UpdateUser(user)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "user updated"})
	}
}

func (s *Serve) handelLogin(c *gin.Context) {
	if user := s.decodeUserJSON(c); user != nil {
		err := s.db.GetID(user)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		user, err := s.db.GetUser(user.Id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		token, err := middleware.GenerateToken(user.Id, user.Account)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		} else {
			c.JSON(http.StatusOK, gin.H{"token": token})
		}
	}
}

// * * * RANTING * * *

func (s *Serve) handelGetRatingsUsers(c *gin.Context) {
	if userId, err := s.getUserID(c); err == nil {
		result, err := s.db.GetRatingByUser(userId)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, result)
	}
}

func (s *Serve) handelSaveRatingsUsers(c *gin.Context) {

	if ranting := s.decodeRatingJSON(c); ranting != nil {
		err := s.db.SaveRating(ranting)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "ranting saved"})
	}
}
func (s *Serve) decodeRatingJSON(c *gin.Context) *models.Rating {
	var rating models.Rating
	err := c.ShouldBindJSON(&rating)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return nil
	}
	return &rating
}

// * * * FAVORITE * * *

func (s *Serve) handelGetFavoriteUsers(c *gin.Context) {
	var favorite models.Favorite
	var err error
	if favorite.UserId, err = s.getUserID(c); err == nil {
		err := s.db.GetFavoriteByUser(&favorite)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, favorite)
	}
}

func (s *Serve) handelSaveFavoriteUsers(c *gin.Context) {

	if favorite := s.decodeFavoriteJSON(c); favorite != nil {
		newFavor := favorite.MoviesID
		err := s.db.GetFavoriteByUser(favorite)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		// delete if exist
		favorite.MoviesID = strings.Replace(favorite.MoviesID, fmt.Sprintf("#%s|", newFavor), "", -1)

		// add new
		favorite.MoviesID += fmt.Sprintf("#%s|", newFavor)

		err = s.db.SaveFavorite(favorite)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "favorite saved"})
	}
}

func (s *Serve) handelDeleteFavoriteUsers(c *gin.Context) {
	userId, err := s.getUserID(c)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	}
	favoriteId := s.getFavoritesID(c)

	var favorite models.Favorite
	favorite.UserId = userId
	err = s.db.GetFavoriteByUser(&favorite)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	// delete
	favorite.MoviesID = strings.Replace(favorite.MoviesID, fmt.Sprintf("#%s|", favoriteId), "", -1)

	err = s.db.SaveFavorite(&favorite)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "favorite deleted"})

}

func (s *Serve) decodeFavoriteJSON(c *gin.Context) *models.Favorite {
	var favorite models.Favorite
	err := c.ShouldBindJSON(&favorite)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return nil
	}
	return &favorite
}

// * * * *

func (s *Serve) decodeUserJSON(c *gin.Context) *models.User {
	var user models.User
	err := c.ShouldBindJSON(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return nil
	}
	return &user
}
func (s *Serve) getUserID(c *gin.Context) (int, error) {
	id, err := strconv.Atoi(c.Param("userID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return 0, err
	}
	return id, nil
}
func (s *Serve) getFavoritesID(c *gin.Context) string {
	return c.Param("favoriteID")
}

// * * * MOVIE * * *

func (s *Serve) handelGetListSeries(c *gin.Context) {
	series, err := s.db.GetSeries()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"series": series})
}
func (s *Serve) handelGetListMovies(c *gin.Context) {
	movies, err := s.db.GetMovies()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"movie": movies})
}
func (s *Serve) handelGetmovie(c *gin.Context) {
	if id, err := s.getMovieID(c); err == nil {
		user, err := s.db.GetMoviesById(id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, user)
	}
}
func (s *Serve) handelAddMovies(c *gin.Context) {
	if movie := s.decodeMovieJSON(c); movie != nil {
		err := s.db.AddMovie(movie)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "movie saved"})
	}
}
func (s *Serve) handelDeleteMovies(c *gin.Context) {
	if id, err := s.getMovieID(c); err == nil {
		err := s.db.DeleteMovieByID(id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "movie deleted"})
	}
}

// * * * *

func (s *Serve) decodeMovieJSON(c *gin.Context) *models.Movies {
	var movie models.Movies
	err := c.ShouldBindJSON(&movie)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return nil
	}
	return &movie
}
func (s *Serve) getMovieID(c *gin.Context) (int, error) {
	movieID, err := strconv.Atoi(c.Param("movieID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid movie ID"})
		return 0, err
	}
	return movieID, nil
}
