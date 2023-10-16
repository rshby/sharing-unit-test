package mock

import (
	"context"
	"github.com/stretchr/testify/mock"
	"sharingunittest/account"
	"sharingunittest/dto"
)

type AccountServiceMock struct {
	Mock mock.Mock
}

func (a *AccountServiceMock) Insert(ctx context.Context, request *dto.InsertAccountRequest) (*account.Account, error) {
	args := a.Mock.Called(ctx, request)
	acc := args.Get(0)
	if acc == nil {
		return nil, args.Error(1)
	}

	return acc.(*account.Account), nil
}

func (a *AccountServiceMock) GetById(ctx context.Context, request *dto.GetAccountRequest) ([]account.Account, error) {
	args := a.Mock.Called(ctx, request)
	acc := args.Get(0)
	if acc == nil {
		return nil, args.Error(1)
	}

	return acc.([]account.Account), nil
}
