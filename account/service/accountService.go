package account

import (
	"context"
	"database/sql"
	"errors"
	"reflect"
	entity "sharingunittest/account"
	iaccount "sharingunittest/account/repository/interface"
	"sharingunittest/dto"
	"sync"
)

type AccountService struct {
	AccountRepo iaccount.IAccountRepository
}

func NewAccountService(accRepo iaccount.IAccountRepository) *AccountService {
	return &AccountService{accRepo}
}

// method insert
func (a *AccountService) Insert(ctx context.Context, request *dto.InsertAccountRequest) (*entity.Account, error) {
	// create entity
	input := entity.Account{
		Email:    request.Email,
		Username: request.Username,
		Password: request.Password,
		FullName: sql.NullString{String: request.FullName, Valid: true},
		Gender: sql.NullString{
			String: request.Gender,
			Valid:  true,
		},
	}

	// call procedure insert in repository
	res, err := a.AccountRepo.Insert(ctx, &input)
	if err != nil {
		return nil, err
	}

	// success insert
	return res, nil
}

func (a *AccountService) GetById(ctx context.Context, request *dto.GetAccountRequest) ([]entity.Account, error) {
	wg := &sync.WaitGroup{}

	lenEmail := len(request.Email) // lenEmail: 3
	chanRes := make(chan entity.Account, lenEmail)
	chanErr := make(chan error, lenEmail)

	var response []entity.Account

	// get account with async
	for _, email := range request.Email {
		go a.AccountRepo.GetByid(ctx, wg, email, chanRes, chanErr)
	}

	counter := 1
	for data := range chanRes {
		_ = <-chanErr
		if !reflect.DeepEqual(data, entity.Account{}) {
			response = append(response, data)
		}
		if counter == lenEmail {
			close(chanRes)
			close(chanErr)
		}
		counter++
	}
	wg.Wait()

	if len(response) == 0 {
		return nil, errors.New("record not found")
	}

	// success
	return response, nil
}
