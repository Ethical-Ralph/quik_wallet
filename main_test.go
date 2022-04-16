package main

import (
	"bytes"
	"encoding/json"
	"fmt"
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
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
)

type redisMock struct {
}

func (r *redisMock) Set(key string, value interface{}) error {
	return nil
}

var walletBalanceFromCache string = "123.43"
var walletBalanceFromDB float64 = 232.34
var walletExistBalance float64 = 500

func (r *redisMock) Get(key string) string {
	if key == "wallet_from_cache" {
		return walletBalanceFromCache
	}
	return ""
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
	switch walletId {
	case "wallet_with_cash":
		return &models.Wallet{
			WalletIdentifier: "wallet_with_cash",
			Balance:          500.23,
		}, nil
	case "wallet_exist":
		return &models.Wallet{
			WalletIdentifier: "wallet_exist",
			Balance:          walletExistBalance,
		}, nil
	case "wallet_from_db":
		return &models.Wallet{
			WalletIdentifier: "wallet_from_db",
			Balance:          walletBalanceFromDB,
		}, nil
	default:
		return nil, nil
	}
}

func (r *repositoryMock) UpdateWalletBalance(walletId string, amount float64) error {
	return nil
}

func mockRepo() *repositoryMock {
	return &repositoryMock{}
}

func to2sf(num decimal.Decimal) string {
	n, _ := num.Float64()
	return fmt.Sprintf("%.2f", n)
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

type Data struct {
	Balance float64 `json:"balance"`
}
type walletBalanceResponseStuct struct {
	Data    Data
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func TestGetWalletGetFromCache(t *testing.T) {
	wallet_id := "wallet_from_cache"

	request, err := makeRequest(fmt.Sprintf("/api/v1/wallet/%s/balance", wallet_id), http.MethodGet, nil)
	if err != nil {
		t.Fatalf("Couldn't make request: %v\n", err)
	}

	res := walletBalanceResponseStuct{}
	err = json.Unmarshal(request.Body.Bytes(), &res)
	if err != nil {
		t.Fatal("Couldn't unmarshal JSON")
	}

	if request.Code != http.StatusOK {
		t.Fatalf("Expected to get status %d but instead got %d\n", http.StatusOK, request.Code)
	}

	expected, _ := decimal.NewFromString(walletBalanceFromCache)
	value := decimal.NewFromFloat(res.Data.Balance)

	if !expected.Equal(value) {
		t.Fatalf("Expected to get balance %s but instead got %s\n", to2sf(expected), to2sf(value))
	}
}

func TestGetWalletGetFromDB(t *testing.T) {
	wallet_id := "wallet_from_db"

	request, err := makeRequest(fmt.Sprintf("/api/v1/wallet/%s/balance", wallet_id), http.MethodGet, nil)
	if err != nil {
		t.Fatalf("Couldn't make request: %v\n", err)
	}

	res := walletBalanceResponseStuct{}
	err = json.Unmarshal(request.Body.Bytes(), &res)
	if err != nil {
		t.Fatal("Couldn't unmarshal JSON")
	}

	if request.Code != http.StatusOK {
		t.Fatalf("Expected to get status %d but instead got %d\n", http.StatusOK, request.Code)
	}

	expected := decimal.NewFromFloat(walletBalanceFromDB)
	value := decimal.NewFromFloat(res.Data.Balance)

	if !expected.Equal(value) {
		t.Fatalf("Expected to get balance %s but instead got %s\n", to2sf(expected), to2sf(value))
	}
}

func TestCreateWallet(t *testing.T) {
	request, err := makeRequest("/api/v1/wallet/create", http.MethodPost, []byte(`{"walletIdentifier": "wallet"}`))
	if err != nil {
		t.Fatalf("Couldn't make request: %v\n", err)
	}

	if request.Code != http.StatusCreated {
		t.Fatalf("Expected to get status %d but instead got %d\n", http.StatusCreated, request.Code)
	}
}

type responseStuct struct {
	Message string `json:"message"`
	Success bool   `json:"success"`
}

type errorResponseStuct struct {
	Error   string `json:"error"`
	Success bool   `json:"success"`
}

func TestCreditWallet(t *testing.T) {
	wallet_id := "wallet_exist"
	amount_to_credit := 100

	request, err := makeRequest(fmt.Sprintf("/api/v1/wallet/%s/credit", wallet_id), http.MethodPost, []byte(fmt.Sprintf(`{"amount": %d}`, amount_to_credit)))
	if err != nil {
		t.Fatalf("Couldn't make request: %v\n", err)
	}

	res := walletBalanceResponseStuct{}
	err = json.Unmarshal(request.Body.Bytes(), &res)
	if err != nil {
		t.Fatal("Couldn't unmarshal JSON")
	}

	if request.Code != http.StatusOK {
		t.Fatalf("Expected to get status %d but instead got %d\n", http.StatusOK, request.Code)
	}

	if !res.Success {
		t.Fatal("Credit Wallet Failed")
	}

	oldBal := decimal.NewFromFloat(walletExistBalance)
	amountToCredit := decimal.NewFromInt(int64(amount_to_credit))

	expected := oldBal.Add(amountToCredit)
	value := decimal.NewFromFloat(res.Data.Balance)

	if !expected.Equal(value) {
		t.Fatalf("Expected to get balance %s but instead got %s\n", to2sf(expected), to2sf(value))
	}
}

func TestDebitWallet(t *testing.T) {
	wallet_id := "wallet_exist"
	amount_to_debit := 100

	request, err := makeRequest(fmt.Sprintf("/api/v1/wallet/%s/debit", wallet_id), http.MethodPost, []byte(fmt.Sprintf(`{"amount": %d}`, amount_to_debit)))
	if err != nil {
		t.Fatalf("Couldn't make request: %v\n", err)
	}

	res := walletBalanceResponseStuct{}
	err = json.Unmarshal(request.Body.Bytes(), &res)
	if err != nil {
		t.Fatal("Couldn't unmarshal JSON")
	}

	if request.Code != http.StatusOK {
		t.Fatalf("Expected to get status %d but instead got %d\n", http.StatusOK, request.Code)
	}

	if !res.Success {
		t.Fatal("Credit Wallet Failed")
	}

	oldBal := decimal.NewFromFloat(walletExistBalance)
	amountToDebit := decimal.NewFromInt(int64(amount_to_debit))

	expected := oldBal.Sub(amountToDebit)
	value := decimal.NewFromFloat(res.Data.Balance)

	if !expected.Equal(value) {
		t.Fatalf("Expected to get balance %s but instead got %s\n", to2sf(expected), to2sf(value))
	}
}

func TestCreditWalletInvalidAmount(t *testing.T) {
	wallet_id := "wallet_with_cash"

	request, err := makeRequest(fmt.Sprintf("/api/v1/wallet/%s/credit", wallet_id), http.MethodPost, []byte(`{"amount": -232}`))
	if err != nil {
		t.Fatalf("Couldn't make request: %v\n", err)
	}

	res := responseStuct{}
	err = json.Unmarshal(request.Body.Bytes(), &res)
	if err != nil {
		t.Fatal("Couldn't unmarshal JSON")
	}

	if request.Code != http.StatusBadRequest {
		t.Fatalf("Expected to get status %d but instead got %d\n", http.StatusOK, request.Code)
	}

	if res.Success {
		t.Fatal("Failed")
	}
}

func TestCreditWalletFailWhenWalletDoesntExist(t *testing.T) {
	wallet_id := "wallet_doesnt_exist"

	request, err := makeRequest(fmt.Sprintf("/api/v1/wallet/%s/credit", wallet_id), http.MethodPost, []byte(`{"amount": 232}`))
	if err != nil {
		t.Fatalf("Couldn't make request: %v\n", err)
	}

	res := errorResponseStuct{}
	err = json.Unmarshal(request.Body.Bytes(), &res)
	if err != nil {
		t.Fatal("Couldn't unmarshal JSON")
	}

	if request.Code != http.StatusBadRequest {
		t.Fatalf("Expected to get status %d but instead got %d\n", http.StatusOK, request.Code)
	}

	if res.Success {
		t.Fatal("Failed")
	}

	if res.Error != "Wallet doesn't exist" {
		t.Fatal("Failed")
	}
}

func TestCreditWalletFailWhenWalletExist(t *testing.T) {
	wallet_id := "wallet_exist"

	request, err := makeRequest("/api/v1/wallet/create", http.MethodPost, []byte(fmt.Sprintf(`{"walletIdentifier": "%s"}`, wallet_id)))
	if err != nil {
		t.Fatalf("Couldn't make request: %v\n", err)
	}

	res := errorResponseStuct{}
	err = json.Unmarshal(request.Body.Bytes(), &res)
	if err != nil {
		t.Fatal("Couldn't unmarshal JSON")
	}

	if request.Code != http.StatusBadRequest {
		t.Fatalf("Expected to get status %d but instead got %d\n", http.StatusOK, request.Code)
	}

	if res.Success {
		t.Fatal("Failed")
	}
}

func TestDebitWalletInsufficientBalance(t *testing.T) {
	wallet_id := "wallet_with_cash"

	request, err := makeRequest(fmt.Sprintf("/api/v1/wallet/%s/debit", wallet_id), http.MethodPost, []byte(`{"amount": 1000000}`))
	if err != nil {
		t.Fatalf("Couldn't make request: %v\n", err)
	}

	res := errorResponseStuct{}
	err = json.Unmarshal(request.Body.Bytes(), &res)
	if err != nil {
		t.Fatal("Couldn't unmarshal JSON")
	}

	if request.Code != http.StatusBadRequest {
		t.Fatalf("Expected to get status %d but instead got %d\n", http.StatusBadRequest, request.Code)
	}

	if res.Success {
		t.Fatal("Failed")
	}

	if res.Error != "Insufficent Balance" {
		t.Fatal("Failed")
	}
}

func TestDebitWalletInvalidAmount(t *testing.T) {
	wallet_id := "wallet_with_cash"

	request, err := makeRequest(fmt.Sprintf("/api/v1/wallet/%s/debit", wallet_id), http.MethodPost, []byte(`{"amount": -232}`))
	if err != nil {
		t.Fatalf("Couldn't make request: %v\n", err)
	}

	res := errorResponseStuct{}
	err = json.Unmarshal(request.Body.Bytes(), &res)
	if err != nil {
		t.Fatal("Couldn't unmarshal JSON")
	}

	if request.Code != http.StatusBadRequest {
		t.Fatalf("Expected to get status %d but instead got %d\n", http.StatusBadRequest, request.Code)
	}

	if res.Success {
		t.Fatal("Failed")
	}

	if res.Error != "Invalid Amount" {
		t.Fatal("Failed")
	}
}
