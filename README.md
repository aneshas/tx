# tx
[![Go Test](https://github.com/aneshas/tx/actions/workflows/test.yml/badge.svg)](https://github.com/aneshas/tx/actions/workflows/test.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/aneshas/tx)](https://goreportcard.com/report/github.com/aneshas/tx)
[![Coverage Status](https://coveralls.io/repos/github/aneshas/tx/badge.svg)](https://coveralls.io/github/aneshas/tx)
[![Go Reference](https://pkg.go.dev/badge/github.com/aneshas/tx.svg)](https://pkg.go.dev/github.com/aneshas/tx)

`go get github.com/aneshas/tx/v2`

Package tx provides a simple abstraction which leverages `context.Context` in order to provide a transactional behavior
which one could use in their use case orchestrator (eg. application service, command handler, etc...).  You might think of it
as closest thing in `Go` to `@Transactional` annotation in Java or the way you could scope a transaction in `C#`.

Many people tend to implement this pattern in one way or another (I have seen it and did it quite a few times), and
more often then not, the implementations still tend to couple your use case orchestrator with your database adapters (eg. repositories) or
on the other hand, rely to heavily on `context.Context` and end up using it as a dependency injection mechanism.

This package relies on `context.Context` in order to simply pass the database transaction down the stack in a safe and clean way which
still does not violate the reasoning behind context package - which is to carry `request scoped` data across api boundaries - which is
a database transaction in this case.

## Drivers
Library currently supports `pgx` and stdlib `sql` out of the box although it is very easy to implement any additional ones
you might need.

## Example
Let's assume we have the following very common setup of an example account service which has a dependency to account repository. 

```go
type Repo interface {
    Save(ctx context.Context, account Account) error
    Find(ctx context.Context, id int) (*Account, error)
}

func NewAccountService(transactor tx.Transactor, repo Repo) *AccountService {
    return &AccountService{
        Transactor: transactor, 
        repo: repo,
    }
}

type AccountService struct {
    // Embedding Transactor interface in order to decorate the service with transactional behavior,
    // although you can choose how and when you use it freely
    tx.Transactor

    repo Repo 
}

type ProvisionAccountReq struct {
    // ...
}

func (s *AccountService) ProvisionAccount(ctx context.Context, r ProvisionAccountReq) error {
    return s.WithTransaction(ctx, func (ctx context.Context) error {
        // ctx contains an embedded transaction and as long as
        // we pass it to our repo methods, they will be able to unwrap it and use it

        // eg. multiple calls to the same or different repos

        return s.repo.Save(ctx, Account{
            // ...
        })
    })
}
```

You will notice that the service looks mostly the same as it would normally apart from embedding `Transactor` interface
and wrapping the use case execution using `WithTransaction`, both of which say nothing of the way the mechanism is implemented (no infrastructure dependencies).

If the function wrapped via `WithTransaction` errors out or panics the transaction itself will be rolled back and if nil error is
returned the transaction will be committed. (this behavior can be changed by providing `WithIgnoredErrors(...)` option to `tx.New`)

### Repo implementation
Then, your repo might use postgres with pgx and have the following example implementation:

```go
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
```

Again, you may freely choose how you implement this and whether or not you actually do use the wrapped
transaction or not.

### main
Then your main function would simply tie everything together like this for example:

```go
func main() {
    var pool *pgxpool.Pool

    svc := NewAccountService(
        tx.New(pgxtxv5.NewDBFromPool(pool)),
        NewAccountRepo(pool),
    )

    _ = svc
}
```

This way, your infrastructural concerns stay in the infrastructure layer where they really belong.

*Please note that this is only one way of using the abstraction*

## Next up
- [ ] Add a way to configure transaction isolation levels for individual drivers eg. `pgxtx.NewDBFromPool(pool, ...opts)`