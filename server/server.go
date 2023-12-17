package server

import (
	"fmt"
	"goflix/db"
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

	s.router.POST("/users", s.handelSaveUsers)
	s.router.GET("/users/:userID", s.handelGetUsers)
	s.router.DELETE("/users/:userID", s.handelDeleteUsers)
	s.router.PUT("/users/:userID", s.handelUpdateUsers)

	s.router.POST("/getid", s.handelGetId)

}

func (s *Serve) handelHello(c *gin.Context) {
	c.String(200, "Hello Goflix")
}

func (s *Serve) handelSaveUsers(c *gin.Context) {
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

func (s *Serve) handelGetId(c *gin.Context) {
	if user := s.decodeUserJSON(c); user != nil {
		err := s.db.GetID(user)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("user id: %d", user.Id)})
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
