package testing

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"io"
	"net/http"
	"net/http/httptest"
	entity "sharingunittest/account"
	account "sharingunittest/account/handler"
	"sharingunittest/dto"
	mck "sharingunittest/mock"
	"strings"
	"testing"
)

func TestHandlerInsert(t *testing.T) {
	t.Run("test success insert", func(t *testing.T) {
		accService := &mck.AccountServiceMock{mock.Mock{}}
		accHandler := account.NewAccountHandler(accService)

		reqBody := &dto.InsertAccountRequest{
			Email:    "reoshby@gmail.com",
			Username: "rshby",
			Password: "1234",
		}

		accService.Mock.On("Insert", mock.Anything, mock.Anything).Return(&entity.Account{
			Id:       1,
			Email:    reqBody.Email,
			Username: reqBody.Username,
			Password: reqBody.Password,
		}, nil).Times(1)

		router := gin.Default()
		router.Handle(http.MethodPost, "/account", accHandler.Insert)

		reqJson, _ := json.Marshal(&reqBody)
		requestBody := strings.NewReader(string(reqJson))
		req := httptest.NewRequest(http.MethodPost, "/account", requestBody)
		req.Header.Add("content-type", "application/json")
		recorder := httptest.NewRecorder()

		router.ServeHTTP(recorder, req)

		// get response
		response := recorder.Result()
		assert.Equal(t, http.StatusOK, response.StatusCode)
		accService.Mock.AssertExpectations(t)
	})
	t.Run("test failed insert", func(t *testing.T) {
		accService := &mck.AccountServiceMock{mock.Mock{}}
		accHandler := account.NewAccountHandler(accService)

		input := dto.InsertAccountRequest{
			Email:    "reo@gmail.com",
			Username: "rshby",
			Password: "1234",
		}
		accService.Mock.On("Insert", mock.Anything, &input).Return(nil, errors.New("gagal insert")).Times(1)

		router := gin.Default()
		router.Handle(http.MethodPost, "/", accHandler.Insert)

		inputJson, _ := json.Marshal(&input)
		requestBody := strings.NewReader(string(inputJson))
		req := httptest.NewRequest(http.MethodPost, "/", requestBody)
		req.Header.Add("content-type", "application/json")
		recorder := httptest.NewRecorder()

		router.ServeHTTP(recorder, req)
		response := recorder.Result()
		body, _ := io.ReadAll(response.Body)
		responseBody := map[string]any{}

		json.Unmarshal(body, &responseBody)

		// validate
		assert.Equal(t, http.StatusInternalServerError, response.StatusCode)
		assert.Equal(t, "gagal insert", responseBody["message"].(string))
		accService.Mock.AssertExpectations(t)
	})
	t.Run("test failed validation", func(t *testing.T) {
		accService := &mck.AccountServiceMock{mock.Mock{}}
		accHandler := account.NewAccountHandler(accService)

		router := gin.Default()
		router.Handle(http.MethodPost, "/", accHandler.Insert)

		input := dto.InsertAccountRequest{
			Email: "aa",
		}
		inputJson, _ := json.Marshal(&input)
		requestBody := strings.NewReader(string(inputJson))
		req := httptest.NewRequest(http.MethodPost, "/", requestBody)
		req.Header.Add("content-type", "application/json")
		recorder := httptest.NewRecorder()

		router.ServeHTTP(recorder, req)

		response := recorder.Result()
		body, _ := io.ReadAll(response.Body)
		responseBody := map[string]any{}
		json.Unmarshal(body, &responseBody)

		// validate
		assert.Equal(t, http.StatusBadRequest, response.StatusCode)
		assert.Equal(t, http.StatusBadRequest, int(responseBody["status_code"].(float64)))
	})
	t.Run("test failed not validation", func(t *testing.T) {
		accService := &mck.AccountServiceMock{mock.Mock{}}
		accHandler := account.NewAccountHandler(accService)

		input := dto.InsertAccountRequest{}
		inputJson, _ := json.Marshal(&input)
		requestBody := strings.NewReader(string(inputJson))

		router := gin.Default()
		router.POST("/", accHandler.Insert)

		req := httptest.NewRequest(http.MethodPost, "/", requestBody)
		req.Header.Add("content-type", "application/json")
		recorder := httptest.NewRecorder()

		router.ServeHTTP(recorder, req)

		response := recorder.Result()

		assert.Equal(t, http.StatusBadRequest, response.StatusCode)
	})
}

// test get by id
func TestGetDataByEmail(t *testing.T) {
	t.Run("test success get one data", func(t *testing.T) {
		accService := &mck.AccountServiceMock{mock.Mock{}}
		accHandler := account.NewAccountHandler(accService)

		accService.Mock.On("GetById", mock.Anything, mock.Anything).Return([]entity.Account{
			{1, "reo@gmail.com", "rshby", "1234", sql.NullString{"Reo Sahobby", true}, sql.NullString{"M", true}},
		}, nil).Times(1)

		input := dto.GetAccountRequest{[]string{"reo@gmail.com"}}
		inputJson, _ := json.Marshal(&input)
		requestBody := strings.NewReader(string(inputJson))

		router := gin.Default()
		router.POST("/", accHandler.GetById)

		req := httptest.NewRequest(http.MethodPost, "/", requestBody)
		req.Header.Add("content-type", "application/json")
		recorder := httptest.NewRecorder()

		router.ServeHTTP(recorder, req)

		response := recorder.Result()
		body, _ := io.ReadAll(response.Body)
		responseBody := map[string]any{}
		json.Unmarshal(body, &responseBody)

		// validate
		assert.Equal(t, http.StatusOK, response.StatusCode)
		accService.Mock.AssertExpectations(t)
	})

	t.Run("test success get one data from two emails", func(t *testing.T) {
		accService := &mck.AccountServiceMock{mock.Mock{}}
		accHandler := account.NewAccountHandler(accService)

		input := dto.GetAccountRequest{
			Email: []string{"reo@gmail.com", "reo1@gmail.com"},
		}
		accService.Mock.On("GetById", mock.Anything, &input).Return([]entity.Account{
			{
				Id:       1,
				Email:    input.Email[0],
				Username: "rshby",
				Password: "1234",
			},
		}, nil).Times(1)

		router := gin.Default()
		router.POST("/", accHandler.GetById)

		inputJson, _ := json.Marshal(&input)
		requestBody := strings.NewReader(string(inputJson))

		req := httptest.NewRequest(http.MethodPost, "/", requestBody)
		req.Header.Add("content-type", "application/json")
		recorder := httptest.NewRecorder()

		router.ServeHTTP(recorder, req)

		response := recorder.Result()
		body, _ := io.ReadAll(response.Body)
		responseBody := dto.ApiResponse{}
		json.Unmarshal(body, &responseBody)

		// validate
		assert.Equal(t, http.StatusOK, response.StatusCode)
		accService.Mock.AssertExpectations(t)
	})

	t.Run("test get two data", func(t *testing.T) {

	})

	t.Run("test not found", func(t *testing.T) {

	})

	t.Run("test bad request", func(t *testing.T) {

	})
}
