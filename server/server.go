package server

import (
	"github.com/Ethical-Ralph/quik_wallet/service"
	"github.com/gin-gonic/gin"
)

type Server struct {
	service *service.Service
	router  *gin.Engine
}

func NewServer(service *service.Service, router *gin.Engine) *Server {
	return &Server{
		service: service,
		router:  router,
	}
}

func (s *Server) Run() error {
	r := s.Routes()

	err := r.Run()

	if err != nil {
		return err
	}

	return nil
}
