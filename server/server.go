package server

import (
	"goflix/db"
	"goflix/middleware"
	"goflix/models"
	"net/http"
	"strconv"

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

	s.router.Use(middleware.JwtMiddleware())

	s.router.GET("/users/:userID", s.handelGetUsers)
	s.router.DELETE("/users/:userID", s.handelDeleteUsers)
	s.router.PUT("/users/:userID", s.handelUpdateUsers)

	s.router.GET("/series/", s.handelGetListSeries)           //    - GET /series : Récupérer la liste des séries disponibles.
	s.router.GET("/movies", s.handelGetListMovies)            //	  - GET /movies : Récupérer la liste des films disponibles.
	s.router.GET("/movies/:movieID", s.handelGetmovie)        //    - GET /movies/{movieID} : Obtenir les détails d'un film spécifique.
	s.router.POST("/movies/", s.handelAddMovies)              //    - POST /movies : Ajouter un nouveau film au catalogue.
	s.router.DELETE("/movies/:movieID", s.handelDeleteMovies) //    - DELETE /movies/{movieID} : Supprimer un film du catalogue.

}

func (s *Serve) handelHello(c *gin.Context) {
	c.String(200, "Hello Goflix")
}

// * * *  USER * * *

func (s *Serve) handelAddUsers(c *gin.Context) {
	if user := s.decodeUserJSON(c); user != nil {
		err := s.db.SaveUser(user)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
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
		token, err := middleware.GenerateToken(user.Id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		} else {
			c.JSON(http.StatusOK, gin.H{"token": token})
		}
	}
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
	movieID, err := strconv.Atoi(c.Param("userID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return 0, err
	}
	return movieID, nil
}

// * * * MOVIE * * *

func (s *Serve) handelGetListSeries(c *gin.Context) {
	if user := s.decodeUserJSON(c); user != nil {
		err := s.db.SaveUser(user)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "user saved"})
	}
}
func (s *Serve) handelGetListMovies(c *gin.Context) {

}
func (s *Serve) handelGetmovie(c *gin.Context) {

}
func (s *Serve) handelAddMovies(c *gin.Context) {

}
func (s *Serve) handelDeleteMovies(c *gin.Context) {

}

// * * * *

func (s *Serve) decodeMovieJSON(c *gin.Context) *models.User {
	var user models.User
	err := c.ShouldBindJSON(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return nil
	}
	return &user
}
func (s *Serve) getMovieID(c *gin.Context) (int, error) {
	movieID, err := strconv.Atoi(c.Param("userID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return 0, err
	}
	return movieID, nil
}
