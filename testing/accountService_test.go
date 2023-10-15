package testing

import (
	"context"
	"database/sql"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	entity "sharingunittest/account"
	account "sharingunittest/account/service"
	"sharingunittest/dto"
	mck "sharingunittest/mock"
	"testing"
)

func TestInsert(t *testing.T) {
	t.Run("test success insert", func(t *testing.T) {
		accRepo := &mck.AccountRepoMock{mock.Mock{}}
		accService := account.NewAccountService(accRepo)

		input := entity.Account{
			Id:       0,
			Email:    "reo@gmail.com",
			Username: "rshby",
			Password: "1234",
			FullName: sql.NullString{"Reo Sahobby", true},
			Gender:   sql.NullString{"M", true},
		}

		accRepo.Mock.On("Insert", mock.Anything, &input).Return(&input, nil)

		res, err := accService.Insert(context.Background(), &dto.InsertAccountRequest{
			Email:    input.Email,
			Username: input.Username,
			Password: input.Password,
			FullName: input.FullName.String,
			Gender:   input.Gender.String,
		})

		assert.Nil(t, err)
		assert.NotNil(t, res)
	})
	t.Run("test error insert email already exist", func(t *testing.T) {
		accRepo := &mck.AccountRepoMock{mock.Mock{}}
		accService := account.NewAccountService(accRepo)

		accRepo.Mock.On("Insert", mock.Anything, mock.Anything).Return(nil, errors.New("email already exist"))

		res, err := accService.Insert(context.Background(), &dto.InsertAccountRequest{
			Email:    "reo@gmail.com",
			Username: "rshby",
			Password: "1234",
		})

		assert.Error(t, err)
		assert.NotNil(t, err)
		assert.Nil(t, res)
		assert.Equal(t, "email already exist", err.Error())
	})
}

// func Test Get By Email
func TestGetByEmail(t *testing.T) {
	t.Run("test get one email", func(t *testing.T) {
		accRepo := &mck.AccountRepoMock{mock.Mock{}}
		accService := account.NewAccountService(accRepo)

		email := "reo@gmail.com"
		accRepo.Mock.On("GetByid", mock.Anything, mock.Anything, email, mock.Anything, mock.Anything).Return(entity.Account{
			Id:       1,
			Email:    email,
			Password: "1234",
		}, nil).Times(1)

		res, err := accService.GetById(context.Background(), &dto.GetAccountRequest{
			Email: []string{email},
		})

		assert.Nil(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, 1, len(res))
		accRepo.Mock.AssertExpectations(t)
	})
	t.Run("test 2 email get 1 data", func(t *testing.T) {
		accRepo := &mck.AccountRepoMock{mock.Mock{}}
		accService := account.NewAccountService(accRepo)

		email := []string{"reo@gmail.com", "sahobby@gmail.com"}
		accRepo.Mock.On("GetByid", mock.Anything, mock.Anything, email[0], mock.Anything, mock.Anything).Return(entity.Account{
			Id:       1,
			Email:    email[0],
			Username: "rshby",
			Password: "1234",
		}, nil).Times(1)
		accRepo.Mock.On("GetByid", mock.Anything, mock.Anything, email[1], mock.Anything, mock.Anything).Return(entity.Account{}, nil).Times(1)

		res, err := accService.GetById(context.Background(), &dto.GetAccountRequest{Email: email})

		assert.Nil(t, err)
		assert.Equal(t, 1, len(res))
		accRepo.Mock.AssertExpectations(t)
	})
	t.Run("test not found from 2 email", func(t *testing.T) {
		accRepo := &mck.AccountRepoMock{mock.Mock{}}
		accService := account.NewAccountService(accRepo)

		email := []string{"satu@gmail.com", "dua@gmail.com"}
		accRepo.Mock.On("GetByid", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(entity.Account{}, nil).Times(2)

		res, err := accService.GetById(context.Background(), &dto.GetAccountRequest{email})

		// validate result
		assert.NotNil(t, err)
		assert.Equal(t, "record not found", err.Error())
		assert.Nil(t, res)
		accRepo.Mock.AssertExpectations(t)
	})
}
