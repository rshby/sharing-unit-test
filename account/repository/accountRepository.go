package account

import (
	"context"
	"database/sql"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"sharingunittest/account"
	iaccount "sharingunittest/account/repository/interface"
	"sync"
)

type AccountRepository struct {
	DB *sql.DB
}

// function Provider
func NewAccountRepository(db *sql.DB) iaccount.IAccountRepository {
	return &AccountRepository{db}
}

// method insert
func (a *AccountRepository) Insert(ctx context.Context, input *account.Account) (*account.Account, error) {
	span, ctxTracing := opentracing.StartSpanFromContext(ctx, "Repository Insert")
	defer span.Finish()

	if input.FullName.String == "" {
		input.FullName.Valid = false
	}

	if input.Gender.String == "" {
		input.Gender.Valid = false
	}

	statement, err := a.DB.PrepareContext(ctxTracing, "INSERT INTO accounts(email, username, password, full_name, gender) VALUES(?, ?, ?, ? ,?)")
	if err != nil {
		return nil, err
	}

	// execute
	res, err := statement.ExecContext(ctxTracing, input.Email, input.Username, input.Password, input.FullName, input.Gender)
	if err != nil {
		return nil, err
	}

	// get id
	if id, err := res.LastInsertId(); err != nil {
		input.Id = int(id)
	}

	// success insert
	span.LogFields(
		log.Object("input", *input),
		log.Object("response", *input))
	return input, nil
}

// method get by id
func (a *AccountRepository) GetByid(ctx context.Context, wg *sync.WaitGroup, email string, chanRes chan account.Account, chanErr chan error) {
	span, ctxTracing := opentracing.StartSpanFromContext(ctx, "Repository GetByEmail")
	defer span.Finish()

	wg.Add(1)
	defer wg.Done()

	statement, err := a.DB.PrepareContext(ctxTracing, "SELECT id, email, username, password, full_name, gender FROM accounts WHERE email=?")
	defer statement.Close()
	if err != nil {
		chanRes <- account.Account{}
		chanErr <- err
		return
	}

	// execute
	row := statement.QueryRowContext(ctxTracing, email)
	if row.Err() != nil {
		chanRes <- account.Account{}
		chanErr <- err
		return
	}

	var acc account.Account
	if err := row.Scan(&acc.Id, &acc.Email, &acc.Username, &acc.Password, &acc.FullName, &acc.Gender); err != nil {
		chanRes <- account.Account{}
		chanErr <- err
		return
	}

	span.LogFields(
		log.String("request-email", email),
		log.Object("response-object", acc))

	// success
	chanRes <- acc
	chanErr <- nil
}
