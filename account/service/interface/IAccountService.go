package account

import (
	"context"
	entity "sharingunittest/account"
	"sharingunittest/dto"
)

type IAccountService interface {
	Insert(ctx context.Context, request *dto.InsertAccountRequest) (*entity.Account, error)
	GetById(ctx context.Context, request *dto.GetAccountRequest) ([]entity.Account, error)
}
