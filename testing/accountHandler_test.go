package testing

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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
		}, nil)

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
	})
}
