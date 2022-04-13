package server

import "github.com/gin-gonic/gin"

func (s *Server) Routes() *gin.Engine {
	router := s.router

	r := router.Group("/api/v1")
	{
		r.GET("/wallet/:wallet_id/balance", s.GetWallet())
		r.POST("/wallet/:wallet_id/credit", s.CreditWallet())
		r.POST("/wallet/:wallet_id/debit", s.DebitWallet())
		r.POST("/wallet/create", s.CreateWallet())
	}

	return router
}
