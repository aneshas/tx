package example

import (
	"context"
	"github.com/aneshas/tx/v2"
	"github.com/aneshas/tx/v2/pgxtxv5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	var pool *pgxpool.Pool

	svc := NewAccountService(
		tx.New(pgxtxv5.NewDBFromPool(pool)),
		NewAccountRepo(pool),
	)

	_ = svc
}

type Account struct {
	// ...
}

type Repo interface {
	Save(ctx context.Context, account Account) error
	Find(ctx context.Context, id int) (*Account, error)
}

func NewAccountService(transactor tx.Transactor, repo Repo) *AccountService {
	return &AccountService{Transactor: transactor, repo: repo}
}

type AccountService struct {
	// Embedding transactional behavior in your service
	tx.Transactor

	repo Repo
}

type ProvisionAccountReq struct {
	// ...
}

func (s *AccountService) ProvisionAccount(ctx context.Context, r ProvisionAccountReq) error {
	return s.WithTransaction(ctx, func(ctx context.Context) error {
		// ctx contains an embedded transaction and as long as
		// we pass it to our repo methods, they will be able to unwrap it and use it

		// eg. multiple calls to different repos

		return s.repo.Save(ctx, Account{
			// ...
		})
	})
}

func NewAccountRepo(pool *pgxpool.Pool) *AccountRepo {
	return &AccountRepo{
		pool: pool,
	}
}

type AccountRepo struct {
	pool *pgxpool.Pool
}

func (r *AccountRepo) Save(ctx context.Context, account Account) error {
	_, err := r.conn(ctx).Exec(ctx, "...")

	return err
}

func (r *AccountRepo) Find(ctx context.Context, id int) (*Account, error) {
	rows, err := r.conn(ctx).Query(ctx, "...")
	if err != nil {
		return nil, err
	}

	_ = rows

	return nil, nil
}

type Conn interface {
	Exec(ctx context.Context, sql string, arguments ...any) (commandTag pgconn.CommandTag, err error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
}

func (r *AccountRepo) conn(ctx context.Context) Conn {
	if tx, ok := pgxtxv5.From(ctx); ok {
		return tx
	}

	return r.pool
}
