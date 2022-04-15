package server

import (
	"net/http"

	"github.com/Ethical-Ralph/quik_wallet/infrastructure/database/models"
	"github.com/Ethical-Ralph/quik_wallet/utils"
	"github.com/gin-gonic/gin"
)

func (s *Server) CreateWallet() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		var wallet models.Wallet

		err := c.BindJSON(&wallet)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid JSON Object",
				"success": false,
			})
			return
		}

		err = s.service.CreateWallet(&wallet)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   err.Error(),
				"success": false,
			})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"success": true,
			"message": "Wallet created successfully",
		})
	}
}

func (s *Server) GetWallet() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		walletId := c.Param("wallet_id")

		wallet, err := s.service.GetWalletBalance(walletId)
		if err != nil {
			status := utils.GetStatus(err)
			c.JSON(status, gin.H{
				"error":   err.Error(),
				"success": false,
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"data":    wallet,
			"message": "Wallet balance fetched successfully",
		})
	}
}

type WalletUpdatePayload struct {
	Amount float64 `json:"amount"`
}

func (s *Server) CreditWallet() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		walletId := c.Param("wallet_id")
		var payload WalletUpdatePayload

		err := c.BindJSON(&payload)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid JSON Object",
				"success": false,
			})
			return
		}

		newBalance, err := s.service.CreditWallet(walletId, payload.Amount)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   err.Error(),
				"success": false,
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Wallet credited successfully",
			"data": gin.H{
				"balance": newBalance,
			},
		})
	}
}

func (s *Server) DebitWallet() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		walletId := c.Param("wallet_id")
		var payload WalletUpdatePayload

		err := c.BindJSON(&payload)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid JSON Object",
				"success": false,
			})
			return
		}

		newBalance, err := s.service.DebitWallet(walletId, payload.Amount)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   err.Error(),
				"success": false,
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Wallet debited successfully",
			"data": gin.H{
				"balance": newBalance,
			},
		})
	}
}
