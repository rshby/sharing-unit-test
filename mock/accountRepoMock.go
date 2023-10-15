package mock

import (
	"context"
	"github.com/stretchr/testify/mock"
	entity "sharingunittest/account"
	"sync"
)

type AccountRepoMock struct {
	Mock mock.Mock
}

func (a *AccountRepoMock) Insert(ctx context.Context, input *entity.Account) (*entity.Account, error) {
	args := a.Mock.Called(ctx, input)
	acc := args.Get(0)

	if acc == nil {
		return nil, args.Error(1)
	}

	return acc.(*entity.Account), nil
}

func (a *AccountRepoMock) GetByid(ctx context.Context, wg *sync.WaitGroup, email string, chanRes chan entity.Account, chanErr chan error) {
	wg.Add(1)
	defer wg.Done()

	args := a.Mock.Called(ctx, wg, email, chanRes, chanErr)

	chanRes <- args.Get(0).(entity.Account)
	chanErr <- args.Error(1)

}
