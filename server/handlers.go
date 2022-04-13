package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) GetWallet() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, map[string]interface{}{
			"sucess": true,
		})
	}
}

func (s *Server) CreditWallet() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, map[string]interface{}{
			"sucess": true,
		})
	}
}

func (s *Server) DebitWallet() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, map[string]interface{}{
			"sucess": true,
		})
	}
}
