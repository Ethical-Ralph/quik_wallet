package server

import "github.com/gin-gonic/gin"

func (s *Server) Routes() *gin.Engine {
	router := s.router

	r := router.Group("/api")
	{
		r.GET("/wallets", s.GetWallet())
	}

	return router
}
