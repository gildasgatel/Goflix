package server

import (
	"goflix/db"

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

}

func (s *Serve) handelHello(c *gin.Context) {
	c.String(200, "Hello Goflix")
}

func (s *Serve) handelLogin(c *gin.Context) {
	c.String(200, "Hello Goflix")
}
