package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Adhiana46/echo-boilerplate/config"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
)

type Server struct {
	e   *echo.Echo
	cfg *config.Config
	db  *sqlx.DB
}

func NewServer(cfg *config.Config, db *sqlx.DB) *Server {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	// TODO: init

	return &Server{
		e:   e,
		cfg: cfg,
		db:  db,
	}
}

func (s *Server) Run() error {
	log.Println("[Server]:", fmt.Sprintf("Running server on %s:%s", s.cfg.Http.Host, s.cfg.Http.Port))
	return s.e.Start(fmt.Sprintf("%s:%s", s.cfg.Http.Host, s.cfg.Http.Port))
}
