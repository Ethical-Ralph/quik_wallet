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

type redisMock struct {
}

func (r *redisMock) Set(key string, value interface{}) error {
	return nil
}

func (r *redisMock) Get(key string) string {
	if key == "wallet_id" {
		return "123.43"
	}
	return "value"
}

func (r *redisMock) Del(key string) {
}

func mockRedis() *redisMock {
	return &redisMock{}
}

type repositoryMock struct {
}

func (r *repositoryMock) CreateWallet(walletIdentifier string) error {
	return nil
}

func (r *repositoryMock) FindWallet(walletId string) (*models.Wallet, error) {
	return &models.Wallet{}, nil
}

func (r *repositoryMock) UpdateWalletBalance(walletId string, amount float64) error {
	return nil
}

func mockRepo() *repositoryMock {
	return &repositoryMock{}
}

func bootStrapServer() *gin.Engine {
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

func makeRequest(url string, method string, payload []byte) (*httptest.ResponseRecorder, error) {
	server := bootStrapServer()

	req, err := http.NewRequest(method, url, bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}

	w := httptest.NewRecorder()

	server.ServeHTTP(w, req)

	return w, nil
}
func TestGetWallet(t *testing.T) {
	request, err := makeRequest("/api/v1/wallet/wallet_id/balance", http.MethodGet, nil)
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
	request, err := makeRequest("/api/v1/wallet/create", http.MethodPost, []byte(`{"walletIdentifier": "wallet"}`))
	if err != nil {
		t.Fatalf("Couldn't make request: %v\n", err)
	}

	if request.Code == http.StatusCreated {
		t.Logf("Expected to get status %d is same ast %d\n", http.StatusCreated, request.Code)
	} else {
		t.Fatalf("Expected to get status %d but instead got %d\n", http.StatusCreated, request.Code)
	}
}
