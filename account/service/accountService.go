package account

import (
	"context"
	"database/sql"
	"errors"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
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
	span, ctxTracing := opentracing.StartSpanFromContext(ctx, "Service Insert")
	defer span.Finish()

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
	res, err := a.AccountRepo.Insert(ctxTracing, &input)
	if err != nil {
		return nil, err
	}

	// success insert
	span.LogFields(
		log.Object("request", *request),
		log.Object("response", *res))
	return res, nil
}

func (a *AccountService) GetById(ctx context.Context, request *dto.GetAccountRequest) ([]entity.Account, error) {
	span, ctxTracing := opentracing.StartSpanFromContext(ctx, "Service GetByEmail")
	defer span.Finish()

	wg := &sync.WaitGroup{}

	lenEmail := len(request.Email) // lenEmail: 3
	chanRes := make(chan entity.Account, lenEmail)
	chanErr := make(chan error, lenEmail)

	var response []entity.Account

	// get account with async
	for _, email := range request.Email {
		go a.AccountRepo.GetByid(ctxTracing, wg, email, chanRes, chanErr)
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
	span.LogFields(
		log.Int("request-email-len", len(request.Email)),
		log.Object("request-email", request.Email),
		log.Object("response-object", response))

	return response, nil
}
