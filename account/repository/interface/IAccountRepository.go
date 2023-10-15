package account

import (
	"context"
	"sharingunittest/account"
	"sync"
)

type IAccountRepository interface {
	Insert(ctx context.Context, input *account.Account) (*account.Account, error)
	GetByid(ctx context.Context, wg *sync.WaitGroup, email string, chanRes chan account.Account, chanErr chan error)
}
