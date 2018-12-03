# tx
Package tx provides a simple transaction abstraction I created to be used
mostly by data-mappers/repositories in order to decouple persistence from application/domain logic.
(Something like `@Transactional` annotation in Java, without an annotation)

## Usage
In order to use tx package, your persistence abstractions (be it repositories or anything else)
need to implement `tx.Transactional` interface containing three
relatively simple methods:

```Go
type Transactional interface {
	// RunTx should create new client side transaction object that can
	// be used to Commit or Rollback the transaction at a later point, and
	// wrap into *Tx by using tx.Wrap and then it should call tx.Run
	//
	// Example body:
	// {
	// 	t, _ := sql.BeginTx(ctx) // create client side transaction object
	// 	transaction := tx.Wrap(t) // wrap it inside of a tx transaction

	// TODO - Add any meta data you might need later to the transaction

	// 	return tx.Run(ctx, ur, transaction, f)
	// }
	RunTx(context.Context, func(context.Context) error) error

	// Commit is called every time RunTx transaction func returns nil error
	// indicating that transaction shuld be commited
	//
	// If Commit returns an error it is propagated and returned by RunTx
	//
	// Commit should never be called directly by client code
	Commit(*Tx) error

	// Rollback is called every time RunTx transaction func
	// returns an error indicating that transaction should be aborted
	//
	// If Rollback returns an error it wraps the error returned from
	// transaction func and both are propagated and returned by RunT
	//
	// Rollback should never be called directly by client code
	Rollback(*Tx) error
}
```

This packages also provides simple default `Transactional` implementations
which you can use by simply embedding them inside of your database access objects.

Example:
```Go
func NewAccountDB(db *sql.DB) *AccountDB {
    return &AccountDB{
        SQL: tx.SQL{
            DB: db,
        }
    }
}

type AccountDB struct {
    tx.SQL
}
```

Available transactional implementations:
- [SQL](sql.go) 

See [example](example/) package for a full example implementation.

## Application service
Ultimately you want to use it inside of your application service use case, eg:

```Go
svc.RunTx(
	ctx,
	// This func will be run in a transaction
	// It is crucial for you to use the context provided by the
	// call to this function, otherwise statements will not be run in a transaction
	func(ctx context.Context) error {

		// Use your repo methods here passing in ctx	

		return nil
	},
)
```

See [example](example/) package for a full example implementation.

### Stuff I might implement:
- Transaction options
- Transaction meta data 
