package account

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
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

		c.JSON(statusCode, &dto.ApiResponse{
			StatusCode: statusCode,
			Status:     helper.GenerateStatusFromCode(statusCode),
			Message:    err.Error(),
		})
		return
	}

	// call insert in service
	res, err := a.AccountService.Insert(c, &request)
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
	c.JSON(statusCode, &dto.ApiResponse{
		StatusCode: statusCode,
		Status:     helper.GenerateStatusFromCode(statusCode),
		Message:    "success insert data",
		Data:       res,
	})
	return
}

// handler get by email
func (a *AccountHandler) GetById(c *gin.Context) {
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

	// call get by id in service
	res, err := a.AccountService.GetById(c, &request)
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
	c.JSON(statusCode, &dto.ApiResponse{
		StatusCode: statusCode,
		Status:     helper.GenerateStatusFromCode(statusCode),
		Message:    "success get data",
		Data:       res,
	})
	return
}
