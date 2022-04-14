package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/Ethical-Ralph/quik_wallet/infrastructure/cache"
	"github.com/Ethical-Ralph/quik_wallet/infrastructure/database/models"
	"github.com/Ethical-Ralph/quik_wallet/repository"
	"github.com/Ethical-Ralph/quik_wallet/server"
	"github.com/Ethical-Ralph/quik_wallet/service"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Redis struct {
}

func (r *Redis) Set(key string, value interface{}) error {
	return nil
}

func (r *Redis) Get(key string) string {
	if key == "wallet_id" {
		return "123dd.43"
	}
	return "value"
}

func (r *Redis) Del(key string) {
}

func mockRedis() *Redis {
	return &Redis{}
}

type Repository struct {
}

func (r *Repository) CreateWallet(walletIdentifier string) error {
	return nil
}

func (r *Repository) FindWallet(walletId string) (*models.Wallet, error) {
	return &models.Wallet{}, nil
}

func (r *Repository) UpdateWalletBalance(walletId string, amount float64) error {
	return nil
}

func mockRepo() *Repository {
	return &Repository{}
}

func BootStrapServer() *gin.Engine {
	router := gin.Default()
	gin.SetMode(gin.TestMode)

	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetOutput(os.Stdout)

	cache := cache.NewCache(mockRedis())
	repository := repository.NewRepository(mockRepo())

	service := service.NewService(repository, cache)
	server := server.NewServer(service, router, logger)

	return server.Routes()
}

func MakeRequest(url string, method string, payload []byte) (*httptest.ResponseRecorder, error) {
	server := BootStrapServer()

	req, err := http.NewRequest(method, url, bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}

	w := httptest.NewRecorder()

	server.ServeHTTP(w, req)

	return w, nil
}
func TestGetWallet(t *testing.T) {
	request, err := MakeRequest("/api/v1/wallet/1/balance", http.MethodGet, nil)
	if err != nil {
		t.Fatalf("Couldn't make request: %v\n", err)
	}
	if request.Code == http.StatusOK {
		t.Logf("Expected to get status %d is same ast %d\n", http.StatusOK, request.Code)
	} else {
		t.Fatalf("Expected to get status %d but instead got %d\n", http.StatusOK, request.Code)
	}
}

func TestCreateWallet(t *testing.T) {
	request, err := MakeRequest("/api/v1/wallet/create", http.MethodPost, []byte(`{"walletIdentifier": "wallet"}`))
	if err != nil {
		t.Fatalf("Couldn't make request: %v\n", err)
	}

	if request.Code == http.StatusCreated {
		t.Logf("Expected to get status %d is same ast %d\n", http.StatusCreated, request.Code)
	} else {
		t.Fatalf("Expected to get status %d but instead got %d\n", http.StatusCreated, request.Code)
	}
}
