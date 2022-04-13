package server

import (
	"github.com/Ethical-Ralph/quik_wallet/service"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Server struct {
	service *service.Service
	router  *gin.Engine
	log     *logrus.Logger
}

func NewServer(service *service.Service, router *gin.Engine, logger *logrus.Logger) *Server {
	return &Server{
		service: service,
		router:  router,
		log:     logger,
	}
}

func (s *Server) Run() error {
	s.router.Use(s.LoggerMiddleware())
	r := s.Routes()

	err := r.Run()
	if err != nil {
		return err
	}

	return nil
}

func (s *Server) LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		log := s.log

		log.WithFields(logrus.Fields{
			"method": c.Request.Method,
			"path":   c.Request.URL.Path,
		}).Info("REQUEST INFO")

		c.Next()
	}
}
