package account

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"net/http"
	account "sharingunittest/account/service/interface"
	"sharingunittest/dto"
	"sharingunittest/helper"
	"strings"
)

type AccountHandler struct {
	AccountService account.IAccountService
}

func NewAccountHandler(accSerice account.IAccountService) *AccountHandler {
	return &AccountHandler{
		accSerice,
	}
}

// handler insert
func (a *AccountHandler) Insert(c *gin.Context) {
	span, ctxTracing := opentracing.StartSpanFromContext(c, "Handler Insert")
	defer span.Finish()

	var request dto.InsertAccountRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		statusCode := http.StatusBadRequest
		if errMsg, ok := err.(validator.ValidationErrors); ok {
			var errorMessages []string
			for _, msg := range errMsg {
				errorMessages = append(errorMessages, fmt.Sprintf("error field %v with tag %v", msg.Field(), msg.Tag()))
			}

			c.JSON(statusCode, &dto.ApiResponse{
				StatusCode: statusCode,
				Status:     helper.GenerateStatusFromCode(statusCode),
				Message:    strings.Join(errorMessages, ". "),
			})
			return
		}
	}

	// call insert in service
	res, err := a.AccountService.Insert(ctxTracing, &request)
	if err != nil {
		statusCode := http.StatusInternalServerError
		c.JSON(statusCode, &dto.ApiResponse{
			StatusCode: statusCode,
			Status:     helper.GenerateStatusFromCode(statusCode),
			Message:    err.Error(),
		})
		return
	}

	// success
	statusCode := http.StatusOK
	response := &dto.ApiResponse{
		StatusCode: statusCode,
		Status:     helper.GenerateStatusFromCode(statusCode),
		Message:    "success insert data",
		Data:       res,
	}

	span.LogFields(
		log.Object("request", request),
		log.Int("response-status-code", statusCode),
		log.Object("response-body", *response))

	c.JSON(statusCode, response)
	return
}

// handler get by email
func (a *AccountHandler) GetById(c *gin.Context) {
	span, ctxTracing := opentracing.StartSpanFromContext(c, "Handler GetByEmail")
	defer span.Finish()

	var request dto.GetAccountRequest
	if err := c.BindJSON(&request); err != nil {
		statusCode := http.StatusBadRequest
		c.JSON(statusCode, &dto.ApiResponse{
			StatusCode: statusCode,
			Status:     helper.GenerateStatusFromCode(statusCode),
			Message:    err.Error(),
		})
		return
	}

	span.LogFields(
		log.Object("request-body", request),
		log.Object("request-email", request.Email))

	// call get by id in service
	res, err := a.AccountService.GetById(ctxTracing, &request)
	if err != nil {
		statusCode := http.StatusNotFound
		c.JSON(statusCode, &dto.ApiResponse{
			StatusCode: statusCode,
			Status:     helper.GenerateStatusFromCode(statusCode),
			Message:    err.Error(),
		})
		return
	}

	// success get data
	statusCode := http.StatusOK

	span.LogFields(
		log.Int("response-status-code", statusCode),
		log.Object("response-data", res))

	c.JSON(statusCode, &dto.ApiResponse{
		StatusCode: statusCode,
		Status:     helper.GenerateStatusFromCode(statusCode),
		Message:    "success get data",
		Data:       res,
	})
	return
}
