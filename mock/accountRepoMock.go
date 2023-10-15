package mock

import (
	"context"
	"github.com/stretchr/testify/mock"
	"sharingunittest/account"
	"sync"
)

type AccountRepoMock struct {
	Mock mock.Mock
}

func (a *AccountRepoMock) Insert(ctx context.Context, input *account.Account) (*account.Account, error) {
	args := a.Mock.Called(ctx, input)
	acc := args.Get(0)

	if acc == nil {
		return nil, args.Error(1)
	}

	return acc.(*account.Account), nil
}

func (a *AccountRepoMock) GetByid(ctx context.Context, wg *sync.WaitGroup, email string, chanRes chan account.Account, chanErr chan error) {
	wg.Add(1)
	defer wg.Done()

	arg := a.Mock.Called(ctx, wg, email, chanRes, chanErr)
	chanRes <- arg.Get(0).(account.Account)
	chanErr <- arg.Error(1)
}
